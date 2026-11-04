package rules

import (
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformCountForEachUsageRule)

// TerraformCountForEachUsageRule checks that resources and module calls don't key instances by list position.
type TerraformCountForEachUsageRule struct {
	tflint.DefaultRule
}

// NewTerraformCountForEachUsageRule returns a new rule.
func NewTerraformCountForEachUsageRule() *TerraformCountForEachUsageRule {
	return new(TerraformCountForEachUsageRule)
}

func (r *TerraformCountForEachUsageRule) Name() string {
	return "avm_terraform_count_for_each_usage_valid"
}

func (r *TerraformCountForEachUsageRule) Enabled() bool {
	return true
}

func (r *TerraformCountForEachUsageRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformCountForEachUsageRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR7/"
}

func (r *TerraformCountForEachUsageRule) Check(runner tflint.Runner) error {
	blocks, err := moduleBlocks(runner)
	if err != nil {
		return err
	}

	listVariables := make(map[string]bool)
	for _, block := range blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		if typeName, _ := variableType(block); typeName == "list" || typeName == "tuple" {
			listVariables[block.Labels[0]] = true
		}
	}

	var errList error
	emit := func(message string, issueRange hcl.Range) {
		if subErr := runner.EmitIssue(r, message, issueRange); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}
	for _, block := range blocks {
		if block.Type != "resource" && block.Type != "module" {
			continue
		}
		if common.IsIgnoredResourceBlock(runner, block.Type, block.Labels) {
			continue
		}
		if count, ok := block.Body.Attributes["count"]; ok {
			for _, call := range countLengthCalls(count.Expr) {
				emit("`count` must not be set from `length(...)`; use `for_each` with a map or set so instances keep stable keys", call.Range())
			}
		}
		if forEach, ok := block.Body.Attributes["for_each"]; ok {
			r.checkForEach(forEach, listVariables, emit)
		}
	}
	return errList
}

func (r *TerraformCountForEachUsageRule) checkForEach(attr *hclsyntax.Attribute, listVariables map[string]bool, emit func(string, hcl.Range)) {
	expr := unwrapParens(attr.Expr)
	if name, ok := variableReference(expr); ok && listVariables[name] {
		emit(fmt.Sprintf("`for_each` must not iterate list variable `%s`; use a map or set variable instead", name), attr.Expr.Range())
		return
	}
	forExpr, ok := expr.(*hclsyntax.ForExpr)
	if !ok || forExpr.KeyExpr == nil || forExpr.KeyVar == "" {
		return
	}
	name, ok := variableReference(forExpr.CollExpr)
	if !ok || !listVariables[name] || !isLocalReference(forExpr.KeyExpr, forExpr.KeyVar) {
		return
	}
	emit(fmt.Sprintf("`for_each` must not key instances by the position of list variable `%s`; key them by a stable value instead", name), forExpr.KeyExpr.Range())
}

// countLengthCalls returns the `length(...)` calls that produce a count value directly or as a conditional result. Calls inside a condition are allowed.
func countLengthCalls(expr hclsyntax.Expression) []*hclsyntax.FunctionCallExpr {
	switch e := unwrapParens(expr).(type) {
	case *hclsyntax.FunctionCallExpr:
		if e.Name == "length" {
			return []*hclsyntax.FunctionCallExpr{e}
		}
	case *hclsyntax.ConditionalExpr:
		return append(countLengthCalls(e.TrueResult), countLengthCalls(e.FalseResult)...)
	}
	return nil
}

// isLocalReference reports whether expr is a bare reference to the local symbol name, such as a for expression key.
func isLocalReference(expr hclsyntax.Expression, name string) bool {
	traversal, ok := unwrapParens(expr).(*hclsyntax.ScopeTraversalExpr)
	return ok && len(traversal.Traversal) == 1 && traversal.Traversal.RootName() == name
}
