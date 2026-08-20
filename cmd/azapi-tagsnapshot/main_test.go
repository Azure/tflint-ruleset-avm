package main

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
	"github.com/stretchr/testify/assert"
)

func TestRetainSnapshotMetadata_unchangedCapabilities(t *testing.T) {
	snapshot := &tagcapability.Snapshot{
		GeneratedAt: "2026-02-01T00:00:00Z",
		Source: tagcapability.Source{
			AzureRESTAPISpecs: "dddddddddddddddddddddddddddddddddddddddd",
			BicepTypesAZ:      "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			BicepTypes:        "ffffffffffffffffffffffffffffffffffffffff",
		},
		Resources: map[string]tagcapability.Status{
			"microsoft.example/widgets@2026-01-01": tagcapability.StatusWritable,
		},
	}
	existing := &tagcapability.Snapshot{
		GeneratedAt: "2026-01-01T00:00:00Z",
		Source: tagcapability.Source{
			AzureRESTAPISpecs: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			BicepTypesAZ:      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			BicepTypes:        "cccccccccccccccccccccccccccccccccccccccc",
		},
		Resources: map[string]tagcapability.Status{
			"microsoft.example/widgets@2026-01-01": tagcapability.StatusWritable,
		},
	}

	retainSnapshotMetadata(snapshot, existing)

	assert.Equal(t, existing.GeneratedAt, snapshot.GeneratedAt)
	assert.Equal(t, existing.Source, snapshot.Source)
}

func TestRetainSnapshotMetadata_changedCapabilities(t *testing.T) {
	snapshot := &tagcapability.Snapshot{
		GeneratedAt: "2026-02-01T00:00:00Z",
		Source: tagcapability.Source{
			AzureRESTAPISpecs: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		Resources: map[string]tagcapability.Status{
			"microsoft.example/widgets@2026-01-01": tagcapability.StatusWritable,
		},
	}
	existing := &tagcapability.Snapshot{
		GeneratedAt: "2026-01-01T00:00:00Z",
		Source: tagcapability.Source{
			AzureRESTAPISpecs: "dddddddddddddddddddddddddddddddddddddddd",
		},
		Resources: map[string]tagcapability.Status{
			"microsoft.example/widgets@2026-01-01": tagcapability.StatusUnsupported,
		},
	}

	retainSnapshotMetadata(snapshot, existing)

	assert.Equal(t, "2026-02-01T00:00:00Z", snapshot.GeneratedAt)
	assert.NotEqual(t, existing.Source, snapshot.Source)
}
