package tagcapability

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_validSnapshotLookup(t *testing.T) {
	snapshot, err := Parse([]byte(validSnapshot(`"microsoft.example/widgets@2026-01-01": "writable"`)))
	require.NoError(t, err)

	status, err := snapshot.Lookup("Microsoft.Example/Widgets@2026-01-01")
	require.NoError(t, err)
	assert.Equal(t, StatusWritable, status)

	_, err = snapshot.Lookup("Microsoft.Example/Widgets@2026-02-01")
	assert.ErrorIs(t, err, ErrUnknownResource)
}

func TestParse_rejectsInvalidSnapshots(t *testing.T) {
	tests := map[string]struct {
		data string
		want string
	}{
		"malformed JSON": {
			data: `{"schema_version":`,
			want: "EOF",
		},
		"unsupported schema version": {
			data: `{"schema_version": 2, "generated_at": "2026-01-01T00:00:00Z", "source": {"azure_rest_api_specs": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bicep_types_az": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "bicep_types": "cccccccccccccccccccccccccccccccccccccccc"}, "resources": {"microsoft.example/widgets@2026-01-01": "writable"}}`,
			want: "unsupported schema version",
		},
		"duplicate resource": {
			data: `{"schema_version": 1, "generated_at": "2026-01-01T00:00:00Z", "source": {"azure_rest_api_specs": "a", "bicep_types_az": "b", "bicep_types": "c"}, "resources": {"microsoft.example/widgets@2026-01-01": "writable", "microsoft.example/widgets@2026-01-01": "unsupported"}}`,
			want: "duplicate key",
		},
		"invalid status": {
			data: validSnapshot(`"microsoft.example/widgets@2026-01-01": "unknown"`),
			want: "invalid status",
		},
		"invalid source SHA": {
			data: `{"schema_version": 1, "generated_at": "2026-01-01T00:00:00Z", "source": {"azure_rest_api_specs": "short", "bicep_types_az": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "bicep_types": "cccccccccccccccccccccccccccccccccccccccc"}, "resources": {"microsoft.example/widgets@2026-01-01": "writable"}}`,
			want: "full commit SHAs",
		},
		"uppercase resource": {
			data: validSnapshot(`"Microsoft.Example/widgets@2026-01-01": "writable"`),
			want: "invalid resource key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Parse([]byte(test.data))
			assert.ErrorIs(t, err, ErrInvalidSnapshot)
			assert.ErrorContains(t, err, test.want)
		})
	}
}

func TestDefault_concurrentLoad(t *testing.T) {
	const workers = 20
	snapshots := make(chan *Snapshot, workers)
	errors := make(chan error, workers)
	var group sync.WaitGroup
	for range workers {
		group.Go(func() {
			snapshot, err := Default()
			snapshots <- snapshot
			errors <- err
		})
	}
	group.Wait()
	close(snapshots)
	close(errors)

	var first *Snapshot
	for snapshot := range snapshots {
		require.NotNil(t, snapshot)
		if first == nil {
			first = snapshot
			continue
		}
		assert.Same(t, first, snapshot)
	}
	for err := range errors {
		assert.NoError(t, err)
	}
}

func TestLookup_unknownResourceWrapsSentinel(t *testing.T) {
	snapshot, err := Parse([]byte(validSnapshot(`"microsoft.example/widgets@2026-01-01": "unsupported"`)))
	require.NoError(t, err)
	_, err = snapshot.Lookup("missing@2026-01-01")
	assert.True(t, errors.Is(err, ErrUnknownResource))
}

func validSnapshot(resources string) string {
	return `{
  "schema_version": 1,
  "generated_at": "2026-01-01T00:00:00Z",
  "source": {
    "azure_rest_api_specs": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "bicep_types_az": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
    "bicep_types": "cccccccccccccccccccccccccccccccccccccccc"
  },
  "resources": {` + resources + `}
}`
}
