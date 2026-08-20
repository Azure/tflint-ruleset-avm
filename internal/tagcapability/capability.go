// Package tagcapability loads the embedded AzAPI resource tag capability snapshot.
package tagcapability

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const schemaVersion = 1

var (
	// ErrUnknownResource indicates that a resource type and API version are not in the snapshot.
	ErrUnknownResource = errors.New("unknown resource type")
	// ErrInvalidSnapshot indicates that snapshot content does not satisfy the v1 contract.
	ErrInvalidSnapshot = errors.New("invalid tag capability snapshot")
)

// Status describes whether an AzAPI resource type accepts a top-level tags property.
type Status string

const (
	// StatusWritable requires tags to be set.
	StatusWritable Status = "writable"
	// StatusReadOnly prohibits tags because the API exposes it as read-only.
	StatusReadOnly Status = "read_only"
	// StatusUnsupported prohibits tags because the API does not expose it.
	StatusUnsupported Status = "unsupported"
)

// Source contains the immutable upstream revisions used to generate a snapshot.
type Source struct {
	AzureRESTAPISpecs string `json:"azure_rest_api_specs"`
	BicepTypesAZ      string `json:"bicep_types_az"`
	BicepTypes        string `json:"bicep_types"`
}

// Snapshot is the compact, standalone tag capability data consumed by the rule.
type Snapshot struct {
	SchemaVersion int               `json:"schema_version"`
	GeneratedAt   string            `json:"generated_at"`
	Source        Source            `json:"source"`
	Resources     map[string]Status `json:"resources"`
}

//go:embed data/azapi_tags_v1.json
var embeddedSnapshot []byte

var (
	defaultSnapshot     *Snapshot
	defaultSnapshotErr  error
	defaultSnapshotOnce sync.Once
)

// Default loads and validates the baseline embedded in the ruleset binary.
func Default() (*Snapshot, error) {
	defaultSnapshotOnce.Do(func() {
		defaultSnapshot, defaultSnapshotErr = Parse(embeddedSnapshot)
	})
	return defaultSnapshot, defaultSnapshotErr
}

// Parse validates and loads snapshot data.
func Parse(data []byte) (*Snapshot, error) {
	if err := rejectDuplicateKeys(data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSnapshot, err)
	}

	var snapshot Snapshot
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snapshot); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", ErrInvalidSnapshot, err)
	}
	if err := ensureEOF(decoder); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSnapshot, err)
	}
	if snapshot.SchemaVersion != schemaVersion {
		return nil, fmt.Errorf("%w: unsupported schema version %d", ErrInvalidSnapshot, snapshot.SchemaVersion)
	}
	if snapshot.GeneratedAt == "" {
		return nil, fmt.Errorf("%w: generated_at is required", ErrInvalidSnapshot)
	}
	if _, err := time.Parse(time.RFC3339, snapshot.GeneratedAt); err != nil {
		return nil, fmt.Errorf("%w: generated_at must be RFC3339: %v", ErrInvalidSnapshot, err)
	}
	if !snapshot.Source.Valid() {
		return nil, fmt.Errorf("%w: all source revisions must be full commit SHAs", ErrInvalidSnapshot)
	}
	if len(snapshot.Resources) == 0 {
		return nil, fmt.Errorf("%w: resources is required", ErrInvalidSnapshot)
	}

	normalized := make(map[string]Status, len(snapshot.Resources))
	for resourceType, status := range snapshot.Resources {
		key := strings.ToLower(resourceType)
		if key != resourceType || !validResourceKey(key) {
			return nil, fmt.Errorf("%w: invalid resource key %q", ErrInvalidSnapshot, resourceType)
		}
		if !validStatus(status) {
			return nil, fmt.Errorf("%w: invalid status %q for %q", ErrInvalidSnapshot, status, resourceType)
		}
		if _, exists := normalized[key]; exists {
			return nil, fmt.Errorf("%w: duplicate resource key %q", ErrInvalidSnapshot, resourceType)
		}
		normalized[key] = status
	}
	snapshot.Resources = normalized
	return &snapshot, nil
}

// Lookup returns the tag capability for an exact AzAPI resource type and API version.
func (s *Snapshot) Lookup(resourceType string) (Status, error) {
	status, ok := s.Resources[strings.ToLower(resourceType)]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownResource, resourceType)
	}
	return status, nil
}

// Valid reports whether Source contains full, lowercase Git commit SHAs.
func (s Source) Valid() bool {
	return validSHA(s.AzureRESTAPISpecs) && validSHA(s.BicepTypesAZ) && validSHA(s.BicepTypes)
}

func validStatus(status Status) bool {
	return status == StatusWritable || status == StatusReadOnly || status == StatusUnsupported
}

func validResourceKey(key string) bool {
	resourceType, apiVersion, found := strings.Cut(key, "@")
	return found && resourceType != "" && apiVersion != "" && !strings.Contains(apiVersion, "@")
}

func validSHA(value string) bool {
	if len(value) != 40 {
		return false
	}
	for _, character := range value {
		if character < '0' || (character > '9' && character < 'a') || character > 'f' {
			return false
		}
	}
	return true
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	err := decoder.Decode(&extra)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err == nil {
		return errors.New("multiple JSON values")
	}
	return err
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := walkJSON(decoder); err != nil {
		return err
	}
	return ensureEOF(decoder)
}

func walkJSON(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch delimiter := token.(type) {
	case json.Delim:
		switch delimiter {
		case '{':
			keys := map[string]struct{}{}
			for decoder.More() {
				token, err := decoder.Token()
				if err != nil {
					return err
				}
				key, ok := token.(string)
				if !ok {
					return errors.New("object key is not a string")
				}
				if _, exists := keys[key]; exists {
					return fmt.Errorf("duplicate key %q", key)
				}
				keys[key] = struct{}{}
				if err := walkJSON(decoder); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walkJSON(decoder); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected delimiter %q", delimiter)
		}
	default:
		return nil
	}
}
