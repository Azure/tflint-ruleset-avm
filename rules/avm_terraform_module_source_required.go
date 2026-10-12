package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var moduleSourceBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "module",
			LabelNames: []string{"name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{
						Name: "source",
					},
				},
			},
		},
	},
}

var _ tflint.Rule = new(ModuleSourceRule)

type ModuleSourceRule struct {
	tflint.DefaultRule
}

type moduleSourceRuleConfig struct {
	Severity    string `hclext:"severity,optional"`
	ModuleScope string `hclext:"module_scope,optional"`
}

func NewModuleSourceRule() *ModuleSourceRule {
	return &ModuleSourceRule{}
}

func (t *ModuleSourceRule) Name() string {
	return "avm_terraform_module_source_required"
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
	config := moduleSourceRuleConfig{ModuleScope: "root"}
	if err := r.DecodeRuleConfig(t.Name(), &config); err != nil {
		return err
	}
	return t.CheckWithModuleScope(r, config.ModuleScope)
}

// CheckWithModuleScope allows only direct local sources for the configured AVM scope.
func (t *ModuleSourceRule) CheckWithModuleScope(r tflint.Runner, moduleScope string) error {
	var localPrefix string
	switch moduleScope {
	case "root":
		localPrefix = "./modules/"
	case "submodule":
		localPrefix = "../"
	default:
		return fmt.Errorf(
			"invalid configuration for rule %q: module_scope must be one of %q or %q, got %q",
			t.Name(), "root", "submodule", moduleScope,
		)
	}

	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	body, err := r.GetModuleContent(
		moduleSourceBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	var errList error
	for _, block := range body.Blocks {
		if block.Type != "module" {
			continue
		}

		if subErr := t.checkBlock(r, block, localPrefix); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *ModuleSourceRule) checkBlock(r tflint.Runner, block *hclext.Block, localPrefix string) error {
	source, exists := block.Body.Attributes["source"]
	if !exists {
		return r.EmitIssue(
			t,
			"The `source` field should be declared in the `module` block",
			block.DefRange,
		)
	}

	if err := r.EvaluateExpr(source.Expr, t.isAVMModule(r, source.NameRange, localPrefix), &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}); err != nil {
		return err
	}

	return nil
}

func (t *ModuleSourceRule) isAVMModule(r tflint.Runner, issueRange hcl.Range, localPrefix string) func(string) error {
	return func(source string) error {
		if (strings.HasPrefix(source, "Azure/") && strings.Contains(source, "avm-")) || isDirectLocalModuleSource(source, localPrefix) {
			return nil
		}
		return r.EmitIssue(
			t,
			fmt.Sprintf("The `source` property must reference an AVM registry module (`Azure/` containing `avm-`) or a direct local module (`%s<name>`)", localPrefix),
			issueRange,
		)
	}
}

func isDirectLocalModuleSource(source, prefix string) bool {
	name, ok := strings.CutPrefix(source, prefix)
	name = strings.TrimSuffix(name, "/")
	return ok && name != "" &&
		!strings.ContainsAny(name, "/\\:\x00") &&
		!strings.HasSuffix(name, ".") && !strings.HasSuffix(name, " ")
}
