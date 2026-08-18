package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"reflect"
	"sort"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformVariableOrderRule checks whether the variables are sorted in expected order
type TerraformVariableOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableOrderRule returns a new rule
func NewTerraformVariableOrderRule() *TerraformVariableOrderRule {
	return &TerraformVariableOrderRule{}
}

// Name returns the rule name
func (r *TerraformVariableOrderRule) Name() string {
	return "terraform_variable_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariableOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVariableOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariableOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in expected order
func (r *TerraformVariableOrderRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkVariableOrder(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformVariableOrderRule) checkVariableOrder(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_variable_order check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks

	requiredVars := r.getSortedVariableNames(blocks, false)
	optionalVars := r.getSortedVariableNames(blocks, true)
	sortedVariableNames := append(requiredVars, optionalVars...)

	variableNames := r.getVariableNames(blocks)
	if reflect.DeepEqual(variableNames, sortedVariableNames) {
		return nil
	}

	firstRange := r.firstVariableRange(blocks)
	sortedVariableHclTxts := r.sortedVariableCodeTxts(blocks, file, sortedVariableNames)
	sortedVariableHclBytes := hclwrite.Format([]byte(strings.Join(sortedVariableHclTxts, "\n\n")))

	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended variable order:\n%s", sortedVariableHclBytes),
		*firstRange,
	)
}

func (r *TerraformVariableOrderRule) sortedVariableCodeTxts(blocks hclsyntax.Blocks, file *hcl.File, sortedVariableNames []string) []string {
	variableHclTxts := r.variableCodeTxts(blocks, file)
	var sortedVariableHclTxts []string
	for _, name := range sortedVariableNames {
		sortedVariableHclTxts = append(sortedVariableHclTxts, variableHclTxts[name])
	}
	return sortedVariableHclTxts
}

func (r *TerraformVariableOrderRule) variableCodeTxts(blocks hclsyntax.Blocks, file *hcl.File) map[string]string {
	variableHclTxts := make(map[string]string)
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		name := v.Labels[0]
		variableHclTxts[name] = string(v.Range().SliceBytes(file.Bytes))
	})
	return variableHclTxts
}

func (r *TerraformVariableOrderRule) firstVariableRange(blocks hclsyntax.Blocks) *hcl.Range {
	var firstRange *hcl.Range
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		if firstRange == nil {
			firstRange = ref(v.DefRange())
		}
	})
	return firstRange
}

func (r *TerraformVariableOrderRule) getVariableNames(blocks hclsyntax.Blocks) []string {
	var variableNames []string
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		variableNames = append(variableNames, v.Labels[0])
	})
	return variableNames
}

func (r *TerraformVariableOrderRule) getSortedVariableNames(blocks hclsyntax.Blocks, defaultWanted bool) []string {
	var sortedVariableNames []string
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		if _, hasDefault := v.Body.Attributes["default"]; hasDefault == defaultWanted {
			sortedVariableNames = append(sortedVariableNames, v.Labels[0])
		}
	})
	sort.Strings(sortedVariableNames)
	return sortedVariableNames
}

func (r *TerraformVariableOrderRule) forVariables(blocks hclsyntax.Blocks, action func(v *hclsyntax.Block)) {
	for _, block := range blocks {
		if block.Type == "variable" {
			action(block)
		}
	}
}
