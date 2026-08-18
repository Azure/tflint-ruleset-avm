package interfaces

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestRegisteredInterfaceSchemasAreCanonicalHCL(t *testing.T) {
	for _, registration := range interfaceRuleRegistrations {
		for i, variant := range registration.variants {
			t.Run(registration.name+"/"+string(rune('0'+i)), func(t *testing.T) {
				if strings.ContainsRune(variant.VarTypeString, '\t') {
					t.Fatalf("schema contains a tab:\n%s", variant.VarTypeString)
				}

				source := "schema = " + variant.VarTypeString + "\n"
				formatted := string(hclwrite.Format([]byte(source)))
				if source != formatted {
					t.Fatalf("schema is not canonical HCL formatting:\n%s\nformatted:\n%s", source, formatted)
				}
			})
		}
	}
}
