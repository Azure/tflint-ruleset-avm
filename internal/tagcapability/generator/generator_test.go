package generator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract_capabilityStatuses(t *testing.T) {
	directory := t.TempDir()
	writeFixture(t, directory, `{
  "resources": {
    "Microsoft.Example/widgets@2026-01-01": {"$ref": "types.json#/0"},
    "Microsoft.Example/widgets@2026-02-01": {"$ref": "types.json#/10"},
    "Microsoft.Example/readOnly@2026-01-01": {"$ref": "types.json#/2"},
    "Microsoft.Example/unsupported@2026-01-01": {"$ref": "types.json#/4"},
    "Microsoft.Example/discriminated@2026-01-01": {"$ref": "types.json#/6"},
    "Microsoft.Example/unknown@2026-01-01": {"$ref": "types.json#/8"}
  }
}`, `[
  {"$type": "ResourceType", "body": {"$ref": "#/1"}},
  {"$type": "ObjectType", "properties": {"tags": {"flags": 0}}},
  {"$type": "ResourceType", "body": {"$ref": "#/3"}},
  {"$type": "ObjectType", "properties": {"tags": {"flags": 2}}},
  {"$type": "ResourceType", "body": {"$ref": "#/5"}},
  {"$type": "ObjectType", "properties": {}},
  {"$type": "ResourceType", "body": {"$ref": "#/7"}},
  {"$type": "DiscriminatedObjectType", "baseProperties": {"tags": {"flags": 0}}},
  {"$type": "ResourceType", "body": {"$ref": "#/9"}},
  {"$type": "UnionType"},
  {"$type": "ResourceType", "body": {"$ref": "#/11"}},
  {"$type": "ObjectType", "properties": {}}
]`)

	snapshot, err := Extract(testOptions(directory))
	require.NoError(t, err)
	assert.Equal(t, map[string]tagcapability.Status{
		"microsoft.example/widgets@2026-01-01":       tagcapability.StatusWritable,
		"microsoft.example/widgets@2026-02-01":       tagcapability.StatusUnsupported,
		"microsoft.example/readonly@2026-01-01":      tagcapability.StatusReadOnly,
		"microsoft.example/unsupported@2026-01-01":   tagcapability.StatusUnsupported,
		"microsoft.example/discriminated@2026-01-01": tagcapability.StatusWritable,
	}, snapshot.Resources)
}

func TestExtract_rejectsMalformedReferencesAndDefinitions(t *testing.T) {
	tests := map[string]struct {
		index string
		types string
		want  string
	}{
		"cross-file body reference": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"}}}`,
			types: `[{"$type":"ResourceType","body":{"$ref":"other.json#/1"}}]`,
			want:  "body reference must be local",
		},
		"local resource reference": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"#/0"}}}`,
			types: `[{"$type":"ResourceType","body":{"$ref":"#/0"}}]`,
			want:  "resource reference must include a file",
		},
		"out of range reference": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/9"}}}`,
			types: `[]`,
			want:  "out of range",
		},
		"unexpected resource kind": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"}}}`,
			types: `[{"$type":"ObjectType"}]`,
			want:  "expected ResourceType",
		},
		"duplicate index key": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"},"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"}}}`,
			types: `[{"$type":"ResourceType","body":{"$ref":"#/0"}}]`,
			want:  "duplicate key",
		},
		"duplicate normalized key": {
			index: `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"},"microsoft.example/widgets@2026-01-01":{"$ref":"types.json#/0"}}}`,
			types: `[{"$type":"ResourceType","body":{"$ref":"#/1"}},{"$type":"ObjectType","properties":{}}]`,
			want:  "duplicate normalized",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			directory := t.TempDir()
			writeFixture(t, directory, test.index, test.types)
			_, err := Extract(testOptions(directory))
			require.ErrorContains(t, err, test.want)
		})
	}
}

func TestWrite_isDeterministic(t *testing.T) {
	directory := t.TempDir()
	writeFixture(t, directory, `{"resources":{"Microsoft.Example/widgets@2026-01-01":{"$ref":"types.json#/0"}}}`, `[
  {"$type":"ResourceType","body":{"$ref":"#/1"}},
  {"$type":"ObjectType","properties":{"tags":{"flags":0}}}
]`)
	snapshot, err := Extract(testOptions(directory))
	require.NoError(t, err)

	first := filepath.Join(directory, "first.json")
	second := filepath.Join(directory, "second.json")
	require.NoError(t, Write(first, snapshot))
	require.NoError(t, Write(second, snapshot))
	firstData, err := os.ReadFile(first)
	require.NoError(t, err)
	secondData, err := os.ReadFile(second)
	require.NoError(t, err)
	assert.Equal(t, firstData, secondData)

	parsed, err := tagcapability.Parse(firstData)
	require.NoError(t, err)
	assert.Equal(t, snapshot.Resources, parsed.Resources)
}

func testOptions(directory string) Options {
	return Options{
		TypesDirectory: directory,
		GeneratedAt:    time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		Source: tagcapability.Source{
			AzureRESTAPISpecs: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			BicepTypesAZ:      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			BicepTypes:        "cccccccccccccccccccccccccccccccccccccccc",
		},
	}
}

func writeFixture(t *testing.T, directory, index, types string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(directory, "index.json"), []byte(index), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(directory, "types.json"), []byte(types), 0o600))
}
