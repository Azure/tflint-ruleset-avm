package requiredversion

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// variableBodySchema is the schema for the variable block that we want to extract from the config.
var requiredVersionBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type: "terraform",
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "required_version"},
				},
				Blocks: []hclext.BlockSchema{
					{
						Type: "required_providers",
						Body: &hclext.BodySchema{
							Attributes: []hclext.AttributeSchema{
								{Name: "azurerm"},
							},
						},
					},
				},
			},
		},
	},
}

// Check interface compliance with the tflint.Rule.
var _ tflint.Rule = new(RequiredVersionRule)

// RequiredVersionRule is the struct that represents a rule that
// check for the correct usage of an interface.
type RequiredVersionRule struct {
	tflint.DefaultRule
	link     string
	ruleName string
}

// NewRequiredVersionRule returns a new rule with the given variable.
func NewRequiredVersionRule(ruleName, link string) *RequiredVersionRule {
	return &RequiredVersionRule{
		ruleName: ruleName,
		link:     link,
	}
}

// Name returns the rule name.
func (or *RequiredVersionRule) Name() string {
	return or.ruleName
}

// Link returns the link to the rule documentation.
func (or *RequiredVersionRule) Link() string {
	return or.link
}

// Enabled returns whether the rule is enabled.
func (or *RequiredVersionRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (or *RequiredVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check checks whether the module satisfies the interface.
// It will search for a variable with the same name as the interface.
// It will check the type, default value and nullable attributes.
func (vcr *RequiredVersionRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	// Define the schema that we want to pull out of the module content.
	body, err := r.GetModuleContent(
		requiredVersionBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Iterate over the outputs and check for the name we are interested in.
	for _, b := range body.Blocks {
		if b.Labels[0] == "abc" {
			return nil
		}
	}
	return r.EmitIssue(
		vcr,
		fmt.Sprintf("module owners MUST output the `%s` in their modules", "aaa"),
		hcl.Range{
			Filename: "outputs.tf",
		},
	)
}
