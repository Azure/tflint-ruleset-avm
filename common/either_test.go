package common_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &mockRule{}

type mockRule struct {
	tflint.DefaultRule
	success bool
}

func (m *mockRule) Check(r tflint.Runner) error {
	if !m.success {
		_ = r.EmitIssue(m, "mock issue", hcl.Range{})
	}
	return nil
}

func (m *mockRule) Name() string {
	return "mock"
}

func (m *mockRule) Enabled() bool {
	return true
}

func (m *mockRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (m *mockRule) Link() string {
	return ""
}

func TestEitherPrivateEndpoints(t *testing.T) {
	cases := []struct {
		name          string
		rules         []tflint.Rule
		expectedIssue bool
	}{
		{
			name:          "correct",
			rules:         []tflint.Rule{&mockRule{success: true}, &mockRule{success: true}},
			expectedIssue: false,
		},
		{
			name:          "correct2",
			rules:         []tflint.Rule{&mockRule{success: true}, &mockRule{success: false}},
			expectedIssue: false,
		},
		{
			name:          "correct3",
			rules:         []tflint.Rule{&mockRule{success: false}, &mockRule{success: true}},
			expectedIssue: false,
		},
		{
			name:          "incorrect",
			rules:         []tflint.Rule{&mockRule{success: false}, &mockRule{success: false}},
			expectedIssue: true,
		},
		{
			name:          "threeRulesCorrect",
			rules:         []tflint.Rule{&mockRule{success: false}, &mockRule{success: false}, &mockRule{success: false}},
			expectedIssue: true,
		},
		{
			name:          "threeRulesIncorrect",
			rules:         []tflint.Rule{&mockRule{success: false}, &mockRule{success: false}, &mockRule{success: true}},
			expectedIssue: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filename := "variables.tf"

			runner := helper.TestRunner(t, map[string]string{filename: ""})

			sut := common.NewEitherCheckRule("either", true, tflint.ERROR, tc.rules...)
			err := sut.Check(runner)
			require.NoError(t, err)

			if tc.expectedIssue {
				assert.NotEmpty(t, runner.Issues)
			} else {
				assert.Empty(t, runner.Issues)
			}
		})
	}
}
