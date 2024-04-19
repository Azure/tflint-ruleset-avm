package rules

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(ModulesRule)

type ModulesRule struct {
	tflint.DefaultRule
}

func NewModulesRule() *ModulesRule {
	return &ModulesRule{}
}

func (t *ModulesRule) Name() string {
	return "required_providers_tffr1"
}

func (t *ModulesRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr1---category-composition---cross-referencing-modules"
}

func (t *ModulesRule) Enabled() bool {
	return true
}

func (t *ModulesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *ModulesRule) Check(r tflint.Runner) error {
	tFile, err := r.GetFile("terraform.tf")
	if err != nil {
		return err
	}

	body, ok := tFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	var errList error
	for _, block := range body.Blocks {
		if block.Type != "module" {
			continue
		}

		if subErr := t.checkBlock(r, block); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *ModulesRule) checkBlock(r tflint.Runner, block *hclsyntax.Block) error {
	var errList error
	source, exists := block.Body.Attributes["source"]
	if !exists {
		errList = multierror.Append(errList, r.EmitIssue(
			t,
			"The `source` field should be declared in the `module` block",
			block.DefRange(),
		))

		return errList
	}

	errList = multierror.Append(r.EvaluateExpr(source.Expr, t.isAVMModule(r, source.NameRange), &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}))

	return errList
}

func (t *ModulesRule) isAVMModule(r tflint.Runner, issueRange hcl.Range) func(string) error {
	return func(source string) error {
		if !strings.Contains(source, "terraform-azurerm-") {
			return r.EmitIssue(
				t,
				"The `source` property constraint should contain `terraform-azurerm-` to only involve AVM Module",
				issueRange,
			)
		}

		return nil
	}
}
