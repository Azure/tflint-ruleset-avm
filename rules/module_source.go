package rules

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(ModuleSourceRule)

type ModuleSourceRule struct {
	tflint.DefaultRule
}

func NewModuleSourceRule() *ModuleSourceRule {
	return &ModuleSourceRule{}
}

func (t *ModuleSourceRule) Name() string {
	return "required_module_source_tffr1"
}

func (t *ModuleSourceRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr1---category-composition---cross-referencing-modules"
}

func (t *ModuleSourceRule) Enabled() bool {
	return true
}

func (t *ModuleSourceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *ModuleSourceRule) Check(r tflint.Runner) error {
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

func (t *ModuleSourceRule) checkBlock(r tflint.Runner, block *hclsyntax.Block) error {
	source, exists := block.Body.Attributes["source"]
	if !exists {
		return r.EmitIssue(
			t,
			"The `source` field should be declared in the `module` block",
			block.Range(),
		)
	}

	if err := r.EvaluateExpr(source.Expr, t.isAVMModule(r, source.NameRange), &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}); err != nil {
		return err
	}

	return nil
}

func (t *ModuleSourceRule) isAVMModule(r tflint.Runner, issueRange hcl.Range) func(string) error {
	return func(source string) error {
		if !(strings.HasPrefix(source, "Azure/") && strings.Contains(source, "avm-")) {
			return r.EmitIssue(
				t,
				"The `source` property constraint should start with `Azure/` and contain `avm-` to only involve AVM Module",
				issueRange,
			)
		}

		return nil
	}
}
