package rules_test

import (
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"strings"
	"testing"
)

func TestTerraformTfFileRule(t *testing.T) {
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
			expectedMessage: "must contain a `terraform.tf` file",
		},
		{
			desc: "NoTerraformDotTfFileShouldEmitIssue2",
			files: map[string]string{
				"main.tf": "",
			},
			expectIssue:     true,
			expectedMessage: "must contain a `terraform.tf` file",
		},
		{
			desc: "NoTerraformBlockInTerraformDotTfFileShouldEmitError",
			files: map[string]string{
				"terraform.tf": "",
			},
			expectIssue:     true,
			expectedMessage: "must contain exactly one `terraform` block",
		},
		{
			desc: "TerraformDotTfFileContainsBlockOtherThanTerraformBlockShouldEmitError",
			files: map[string]string{
				"terraform.tf": `locals {}
								 terraform {}`,
			},
			expectIssue:     true,
			expectedMessage: "must contain exactly one `terraform` block",
		},
		{
			desc: "TerraformDotTfFileContainsBlockOtherThanTerraformBlockShouldEmitError2",
			files: map[string]string{
				"terraform.tf": `terraform {}
								 locals {}`,
			},
			expectIssue:     true,
			expectedMessage: "must contain exactly one `terraform` block",
		},
		{
			desc: "TerraformDotTfFileContainsMultipleTerraformBlocks",
			files: map[string]string{
				"terraform.tf": `terraform {}
								 terraform {}`,
			},
			expectIssue:     true,
			expectedMessage: "must contain exactly one `terraform` block",
		},
		{
			desc: "TerraformDotTfFileContainsTopLevelAttribute",
			files: map[string]string{
				"terraform.tf": `unexpected = true
								 terraform {}`,
			},
			expectIssue:     true,
			expectedMessage: "no other top-level content",
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
			sut := rules.NewTerraformTfFileRule()
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
