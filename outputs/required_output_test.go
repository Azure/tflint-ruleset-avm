package outputs_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/outputs"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRequiredOutput(t *testing.T) {
	cases := []struct {
		desc         string
		config       string
		requiredName string
		issues       helper.Issues
	}{
		{
			desc: "require resource_id, ok",
			config: `output "resource_id" {
  value = azurerm_kubernetes_cluster.this.id
}`,
			requiredName: "resource_id",
			issues:       helper.Issues{},
		},
		{
			desc:         "require resource_id, not ok",
			config:       ``,
			requiredName: "resource_id",
			issues: helper.Issues{
				{
					Rule:    outputs.NewRequiredOutputRule("required_output", "resource_id", ""),
					Message: "module owners MUST output the `resource_id` in their modules",
					Range: hcl.Range{
						Filename: "outputs.tf",
					},
				},
			},
		},
		{
			desc: "require resource, ok",
			config: `output "resource" {
  value = azurerm_kubernetes_cluster.this
}`,
			requiredName: "resource",
			issues:       helper.Issues{},
		},
		{
			desc:         "require resource, not ok",
			config:       ``,
			requiredName: "resource",
			issues: helper.Issues{
				{
					Rule:    outputs.NewRequiredOutputRule("required_output", "resource", ""),
					Message: "module owners MUST output the `resource` in their modules",
					Range: hcl.Range{
						Filename: "outputs.tf",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := outputs.NewRequiredOutputRule("required_output", tc.requiredName, "")
			filename := "variables.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}
