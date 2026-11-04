package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformPatternModuleSourceRule)

// TerraformPatternModuleSourceRule checks that resource modules do not call AVM pattern modules.
type TerraformPatternModuleSourceRule struct {
	tflint.DefaultRule
}

// NewTerraformPatternModuleSourceRule returns a new rule.
func NewTerraformPatternModuleSourceRule() *TerraformPatternModuleSourceRule {
	return new(TerraformPatternModuleSourceRule)
}

func (r *TerraformPatternModuleSourceRule) Name() string {
	return "avm_terraform_pattern_module_source_disallowed"
}

func (r *TerraformPatternModuleSourceRule) Enabled() bool {
	return true
}

func (r *TerraformPatternModuleSourceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformPatternModuleSourceRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/RMFR8/"
}

// Check decodes the module class when the rule is used without the ruleset wrapper.
func (r *TerraformPatternModuleSourceRule) Check(runner tflint.Runner) error {
	moduleClass, err := decodeModuleClass(runner, r.Name())
	if err != nil {
		return err
	}
	return r.CheckWithModuleClass(runner, moduleClass)
}

// CheckWithModuleClass checks resource module roots and skips pattern and utility modules.
func (r *TerraformPatternModuleSourceRule) CheckWithModuleClass(runner tflint.Runner, moduleClass string) error {
	return checkResourceModuleRoot(runner, r.Name(), moduleClass, r.checkModule)
}

func (r *TerraformPatternModuleSourceRule) checkModule(runner tflint.Runner) error {
	blocks, err := moduleBlocks(runner)
	if err != nil {
		return err
	}
	var errList error
	for _, block := range blocks {
		if block.Type != "module" || len(block.Labels) == 0 {
			continue
		}
		source, ok := staticString(block, "source")
		if !ok || !isAVMPatternModuleSource(source) {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("Module `%s` must not call the AVM pattern module `%s`; resource modules may only call AVM resource or utility modules", block.Labels[0], source),
			block.Body.Attributes["source"].SrcRange,
		)
		if subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}
	return errList
}

// isAVMPatternModuleSource reports whether source is a registry address such as `Azure/avm-ptn-<name>/<provider>`, with an optional hostname and `//subdir`.
func isAVMPatternModuleSource(source string) bool {
	if strings.Contains(source, "::") || strings.Contains(source, "://") {
		return false
	}
	address, _, _ := strings.Cut(source, "//")
	parts := strings.Split(address, "/")
	if len(parts) == 4 && strings.Contains(parts[0], ".") {
		parts = parts[1:]
	}
	if len(parts) != 3 {
		return false
	}
	return strings.EqualFold(parts[0], "Azure") && strings.HasPrefix(strings.ToLower(parts[1]), "avm-ptn-")
}
