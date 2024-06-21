package outputs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// variableBodySchema is the schema for the variable block that we want to extract from the config.
var outputBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "output",
			LabelNames: []string{"name"},
			Body:       &hclext.BodySchema{},
		},
	},
}

// Check interface compliance with the tflint.Rule.
var _ tflint.Rule = new(RequiredOutputRule)

// RequiredOutputRule is the struct that represents a rule that
// check for the correct usage of an interface.
type RequiredOutputRule struct {
	tflint.DefaultRule
	outputName string
	link       string
	ruleName   string
}

// NewRequiredOutputRule returns a new rule with the given variable.
func NewRequiredOutputRule(ruleName, requiredOutputName, link string) *RequiredOutputRule {
	return &RequiredOutputRule{
		ruleName:   ruleName,
		outputName: requiredOutputName,
		link:       link,
	}
}

// Name returns the rule name.
func (or *RequiredOutputRule) Name() string {
	return or.ruleName
}

// Link returns the link to the rule documentation.
func (or *RequiredOutputRule) Link() string {
	return or.link
}

// Enabled returns whether the rule is enabled.
func (or *RequiredOutputRule) Enabled() bool {
	// Removed as per spec change 20240621 https://github.com/Azure/Azure-Verified-Modules/pull/1028
	return false
}

// Severity returns the severity of the rule.
func (or *RequiredOutputRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check checks whether the module satisfies the interface.
// It will search for a variable with the same name as the interface.
// It will check the type, default value and nullable attributes.
func (vcr *RequiredOutputRule) Check(r tflint.Runner) error {
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
		outputBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Iterate over the outputs and check for the name we are interested in.
	for _, b := range body.Blocks {
		if b.Labels[0] == vcr.outputName {
			return nil
		}
	}
	return r.EmitIssue(
		vcr,
		fmt.Sprintf("module owners MUST output the `%s` in their modules", vcr.outputName),
		hcl.Range{
			Filename: "outputs.tf",
		},
	)
}
