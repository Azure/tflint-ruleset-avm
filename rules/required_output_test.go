package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"os"
	"path/filepath"
	"testing"
)

func TestRequiredOutput(t *testing.T) {
	cases := []struct {
		desc         string
		config       string
		requiredName string
		issue        *helper.Issue
	}{
		{
			desc: "require resource_id, ok",
			config: `output "resource_id" {
  value = azurerm_kubernetes_cluster.this.id
}`,
			requiredName: "resource_id",
			issue:        nil,
		},
		{
			desc:         "require resource_id, not ok",
			config:       ``,
			requiredName: "resource_id",
			issue: &helper.Issue{
				Rule:    NewRequiredOutputRule("required_output", "resource_id", ""),
				Message: "module owners MUST output the `resource_id` in their modules",
				Range: hcl.Range{
					Filename: "outputs.tf",
				},
			},
		},
		{
			desc: "require resource, ok",
			config: `output "resource" {
  value = azurerm_kubernetes_cluster.this
}`,
			requiredName: "resource",
			issue:        nil,
		},
		{
			desc:         "require resource, not ok",
			config:       ``,
			requiredName: "resource",
			issue: &helper.Issue{
				Rule:    NewRequiredOutputRule("required_output", "resource", ""),
				Message: "module owners MUST output the `resource` in their modules",
				Range: hcl.Range{
					Filename: "outputs.tf",
				},
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.desc, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			stub := gostub.Stub(&Fs, &fsAdapter{
				Fs: fs,
			})
			defer stub.Reset()
			pwd, err := os.Getwd()
			require.NoError(t, err)
			path := filepath.Join(pwd, "outputs.tf")
			_ = afero.WriteFile(fs, path, []byte(cc.config), 0644)
			r := helper.TestRunner(t, map[string]string{
				"outputs.tf": cc.config,
			})
			sut := NewRequiredOutputRule("required_output", cc.requiredName, "")
			err = sut.Check(r)
			require.NoError(t, err)
			if cc.issue == nil {
				assert.Empty(t, r.Issues)
				return
			}
			helper.AssertIssues(t, []*helper.Issue{
				cc.issue,
			}, r.Issues)
		})
	}
}
