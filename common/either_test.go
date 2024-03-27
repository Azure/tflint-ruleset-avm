package common_test

import (
	"fmt"
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestEitherPrivateEndpoints(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  toTerraformVarType(interfaces.PrivateEndpoints),
			Expected: helper.Issues{},
		},
		{
			Name:     "also_correct",
			Content:  toTerraformVarType(interfaces.PrivateEndpointsWithSubresourceName),
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect private_endpoints",
			Content: `
variable "private_endpoints" {
    type = map(object({
      name               = optional(string, null)
      role_assignments   = optional(map(object({})), {})
      lock               = optional(object({}), {})
      tags               = optional(map(any), null)
      subnet_resource_id = string
      subresource_name                        = string
    }))
    default     = {}
    nullable    = false
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.PrivateEndpointsRule,
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.PrivateEndpointTypeString),
					Range:   hcl.Range{Filename: "variables.tf", Start: hcl.Pos{Line: 3, Column: 5}, End: hcl.Pos{Line: 10, Column: 8}},
				},
			},
		},
		{
			Name: "also incorrect private_endpoints with subresource_name",
			Content: `
variable "private_endpoints" {
    type = map(object({
      name               = optional(string, null)
      role_assignments   = optional(map(object({})), {})
      lock               = optional(object({}), {})
      tags               = optional(map(any), null)
      subnet_resource_id = string
    }))
    default     = {}
    nullable    = false
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.PrivateEndpointsRule,
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.PrivateEndpointTypeString),
					Range:   hcl.Range{Filename: "variables.tf", Start: hcl.Pos{Line: 3, Column: 5}, End: hcl.Pos{Line: 9, Column: 8}},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rules.PrivateEndpointsRule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

func toTerraformVarType(i interfaces.AvmInterface) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	varBlock := rootBody.AppendNewBlock("variable", []string{i.RuleName})
	varBody := varBlock.Body()

	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenStringLit,
			Bytes:        []byte(i.VarTypeString),
			SpacesBefore: 1,
		},
	})
	varBody.SetAttributeValue("default", i.Default)
	if !i.Nullable {
		varBody.SetAttributeValue("nullable", cty.False)
	}
	return string(f.Bytes())
}
