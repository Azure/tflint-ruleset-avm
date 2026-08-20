// Command azapi-tagsnapshot creates the embedded AzAPI tag capability snapshot.
package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability/generator"
)

func main() {
	var (
		typesDirectory = flag.String("types-dir", "", "path to bicep-types-az/generated")
		output         = flag.String("output", "", "snapshot JSON output path")
		previous       = flag.String("previous", "", "existing snapshot whose generated_at is retained when unchanged")
		generatedAt    = flag.String("generated-at", time.Now().UTC().Format(time.RFC3339), "RFC3339 generation time")
		azureSpecsSHA  = flag.String("azure-rest-api-specs-sha", "", "Azure/azure-rest-api-specs commit SHA")
		bicepTypesAZ   = flag.String("bicep-types-az-sha", "", "Azure/bicep-types-az commit SHA")
		bicepTypes     = flag.String("bicep-types-sha", "", "Azure/bicep-types submodule commit SHA")
	)
	flag.Parse()

	if *output == "" {
		fatalf("-output is required")
	}
	timestamp, err := time.Parse(time.RFC3339, *generatedAt)
	if err != nil {
		fatalf("parse -generated-at: %v", err)
	}

	snapshot, err := generator.Extract(generator.Options{
		TypesDirectory: *typesDirectory,
		GeneratedAt:    timestamp,
		Source: tagcapability.Source{
			AzureRESTAPISpecs: *azureSpecsSHA,
			BicepTypesAZ:      *bicepTypesAZ,
			BicepTypes:        *bicepTypes,
		},
	})
	if err != nil {
		fatalf("extract snapshot: %v", err)
	}
	if *previous != "" {
		data, err := os.ReadFile(*previous)
		if err != nil {
			fatalf("read -previous: %v", err)
		}
		existing, err := tagcapability.Parse(data)
		if err != nil {
			fatalf("parse -previous: %v", err)
		}
		retainSnapshotMetadata(snapshot, existing)
	}
	if err := generator.Write(*output, snapshot); err != nil {
		fatalf("write snapshot: %v", err)
	}
}

func fatalf(format string, arguments ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", arguments...)
	os.Exit(1)
}

func retainSnapshotMetadata(snapshot, existing *tagcapability.Snapshot) {
	if reflect.DeepEqual(existing.Resources, snapshot.Resources) {
		snapshot.GeneratedAt = existing.GeneratedAt
		snapshot.Source = existing.Source
	}
}
