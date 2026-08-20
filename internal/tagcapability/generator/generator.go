// Package generator extracts top-level tag capabilities from public Bicep type data.
package generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
)

// Options identifies the generated Bicep corpus and source metadata to include in a snapshot.
type Options struct {
	TypesDirectory string
	GeneratedAt    time.Time
	Source         tagcapability.Source
}

type indexFile struct {
	Resources map[string]reference `json:"resources"`
}

type reference struct {
	Reference string `json:"$ref"`
}

type resourceType struct {
	Type string    `json:"$type"`
	Body reference `json:"body"`
}

type objectType struct {
	Type           string              `json:"$type"`
	Properties     map[string]property `json:"properties"`
	BaseProperties map[string]property `json:"baseProperties"`
}

type property struct {
	Flags int `json:"flags"`
}

// Extract creates a compact tag capability snapshot from a Bicep generated directory.
func Extract(options Options) (snapshot *tagcapability.Snapshot, err error) {
	if options.TypesDirectory == "" {
		return nil, errors.New("types directory is required")
	}
	if options.GeneratedAt.IsZero() {
		return nil, errors.New("generated_at is required")
	}
	if !options.Source.Valid() {
		return nil, errors.New("all source revisions must be full commit SHAs")
	}

	root, err := os.OpenRoot(options.TypesDirectory)
	if err != nil {
		return nil, fmt.Errorf("open generated directory: %w", err)
	}
	defer func() {
		if closeErr := root.Close(); closeErr != nil && err == nil {
			snapshot = nil
			err = fmt.Errorf("close generated directory: %w", closeErr)
		}
	}()

	var index indexFile
	if err := decodeFile(root, "index.json", &index); err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}
	if len(index.Resources) == 0 {
		return nil, errors.New("read index: resources is required")
	}

	cache := map[string][]json.RawMessage{}
	resources := make(map[string]tagcapability.Status, len(index.Resources))
	resourceTypes := make([]string, 0, len(index.Resources))
	for resourceType := range index.Resources {
		resourceTypes = append(resourceTypes, resourceType)
	}
	sort.Strings(resourceTypes)

	for _, resourceKey := range resourceTypes {
		key := strings.ToLower(resourceKey)
		if !validResourceType(key) {
			return nil, fmt.Errorf("resource %q: invalid resource type", resourceKey)
		}
		if _, exists := resources[key]; exists {
			return nil, fmt.Errorf("resource %q: duplicate normalized resource type", resourceKey)
		}

		file, indexPosition, err := parseReference(index.Resources[resourceKey].Reference, false)
		if err != nil {
			return nil, fmt.Errorf("resource %q: %w", resourceKey, err)
		}
		definitions, err := loadDefinitions(root, file, cache)
		if err != nil {
			return nil, fmt.Errorf("resource %q: %w", resourceKey, err)
		}
		if indexPosition >= len(definitions) {
			return nil, fmt.Errorf("resource %q: reference index %d is out of range", resourceKey, indexPosition)
		}

		var resource resourceType
		if err := json.Unmarshal(definitions[indexPosition], &resource); err != nil {
			return nil, fmt.Errorf("resource %q: decode resource type: %w", resourceKey, err)
		}
		if resource.Type != "ResourceType" {
			return nil, fmt.Errorf("resource %q: expected ResourceType, got %q", resourceKey, resource.Type)
		}
		_, bodyIndex, err := parseReference(resource.Body.Reference, true)
		if err != nil {
			return nil, fmt.Errorf("resource %q: body: %w", resourceKey, err)
		}
		if bodyIndex >= len(definitions) {
			return nil, fmt.Errorf("resource %q: body reference index %d is out of range", resourceKey, bodyIndex)
		}

		status, known, err := tagStatus(definitions[bodyIndex])
		if err != nil {
			return nil, fmt.Errorf("resource %q: body: %w", resourceKey, err)
		}
		if known {
			resources[key] = status
		}
	}

	return &tagcapability.Snapshot{
		SchemaVersion: 1,
		GeneratedAt:   options.GeneratedAt.UTC().Format(time.RFC3339),
		Source:        options.Source,
		Resources:     resources,
	}, nil
}

// Write serializes a snapshot using stable, compact JSON formatting.
func Write(path string, snapshot *tagcapability.Snapshot) error {
	if snapshot == nil {
		return errors.New("snapshot is required")
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	return nil
}

func tagStatus(raw json.RawMessage) (tagcapability.Status, bool, error) {
	var object objectType
	if err := json.Unmarshal(raw, &object); err != nil {
		return "", false, fmt.Errorf("decode body type: %w", err)
	}
	switch object.Type {
	case "ObjectType":
		return propertyStatus(object.Properties), true, nil
	case "DiscriminatedObjectType":
		return propertyStatus(object.BaseProperties), true, nil
	default:
		return "", false, nil
	}
}

func propertyStatus(properties map[string]property) tagcapability.Status {
	tag, found := properties["tags"]
	if !found {
		return tagcapability.StatusUnsupported
	}
	if tag.Flags&2 != 0 {
		return tagcapability.StatusReadOnly
	}
	return tagcapability.StatusWritable
}

func loadDefinitions(root *os.Root, relativePath string, cache map[string][]json.RawMessage) ([]json.RawMessage, error) {
	if definitions, found := cache[relativePath]; found {
		return definitions, nil
	}
	var definitions []json.RawMessage
	if err := decodeFile(root, relativePath, &definitions); err != nil {
		return nil, fmt.Errorf("read types file %q: %w", relativePath, err)
	}
	cache[relativePath] = definitions
	return definitions, nil
}

func parseReference(value string, localOnly bool) (string, int, error) {
	file, fragment, found := strings.Cut(value, "#/")
	if !found || fragment == "" || strings.Contains(fragment, "/") {
		return "", 0, fmt.Errorf("invalid reference %q", value)
	}
	if localOnly && file != "" {
		return "", 0, fmt.Errorf("body reference must be local, got %q", value)
	}
	if !localOnly && file == "" {
		return "", 0, fmt.Errorf("resource reference must include a file, got %q", value)
	}
	if !localOnly && (strings.Contains(file, "\\") || strings.Contains(file, "..")) {
		return "", 0, fmt.Errorf("invalid reference path %q", file)
	}
	position, err := strconv.Atoi(fragment)
	if err != nil || position < 0 {
		return "", 0, fmt.Errorf("invalid reference index %q", fragment)
	}
	return file, position, nil
}

func decodeFile(root *os.Root, path string, target any) error {
	data, err := root.ReadFile(path)
	if err != nil {
		return err
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func validResourceType(resourceType string) bool {
	name, apiVersion, found := strings.Cut(resourceType, "@")
	return found && name != "" && apiVersion != "" && !strings.Contains(apiVersion, "@")
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := walkJSON(decoder); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func walkJSON(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return nil
	}
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
}
