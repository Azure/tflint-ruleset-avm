package rules_test

import (
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"strings"
	"testing"
)

func TestTerraformTf(t *testing.T) {
	cases := []struct {
		desc            string
		files           map[string]string
		expectIssue     bool
		expectedMessage string
	}{
		{
			desc:            "NoTerraformDotTfFileShouldEmitIssue",
			files:           map[string]string{},
			expectIssue:     true,
			expectedMessage: "must contain `terraform.tf` file",
		},
		{
			desc: "NoTerraformDotTfFileShouldEmitIssue2",
			files: map[string]string{
				"main.tf": "",
			},
			expectIssue:     true,
			expectedMessage: "must contain `terraform.tf` file",
		},
		{
			desc: "NoTerraformBlockInTerraformDotTfFileShouldEmitError",
			files: map[string]string{
				"terraform.tf": "",
			},
			expectIssue:     true,
			expectedMessage: "must contain `terraform` block",
		},
		{
			desc: "TerraformDotTfFileContainsBlockOtherThanTerraformBlockShouldEmitError",
			files: map[string]string{
				"terraform.tf": `locals {}
								 terraform {}`,
			},
			expectIssue:     true,
			expectedMessage: "must contain `terraform` block only",
		},
		{
			desc: "TerraformDotTfFileContainsBlockOtherThanTerraformBlockShouldEmitError2",
			files: map[string]string{
				"terraform.tf": `terraform {}
								 locals {}`,
			},
			expectIssue:     true,
			expectedMessage: "must contain `terraform` block only",
		},
		{
			desc: "TerraformDotTfFileContainsTerraformBlockOnly",
			files: map[string]string{
				"terraform.tf": `terraform {}`,
			},
			expectIssue: false,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.desc, func(t *testing.T) {
			r := helper.TestRunner(t, cc.files)
			sut := rules.NewTerraformDotTfRule()
			err := sut.Check(r)
			require.NoError(t, err)
			if cc.expectIssue {
				assert.True(t, issuesContainMessage(r.Issues, cc.expectedMessage))
			} else {
				assert.Empty(t, r.Issues)
			}
		})
	}
}

func issuesContainMessage(issues helper.Issues, msg string) bool {
	for _, issue := range issues {
		if strings.Contains(issue.Message, msg) {
			return true
		}
	}
	return false
}
