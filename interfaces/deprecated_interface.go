package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type deprecatedInterfaceVariantRule struct {
	tflint.DefaultRule
	name         string
	variableName string
	message      string
	link         string
	variants     []AvmInterface
}

func newDeprecatedInterfaceVariantRule(name, variableName, message, link string, variants ...AvmInterface) tflint.Rule {
	return &deprecatedInterfaceVariantRule{
		name:         name,
		variableName: variableName,
		message:      message,
		link:         link,
		variants:     variants,
	}
}

func (r *deprecatedInterfaceVariantRule) Name() string {
	return r.name
}

func (r *deprecatedInterfaceVariantRule) Enabled() bool {
	return true
}

func (r *deprecatedInterfaceVariantRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *deprecatedInterfaceVariantRule) Link() string {
	return r.link
}

func (r *deprecatedInterfaceVariantRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	body, err := runner.GetModuleContent(
		variableBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone},
	)
	if err != nil {
		return err
	}

	for _, block := range body.Blocks {
		if block.Labels[0] != r.variableName {
			continue
		}

		for _, variant := range r.variants {
			passes, err := common.RulePasses(runner, NewVarCheckRuleFromAvmInterface(variant))
			if err != nil {
				return err
			}
			if passes {
				return runner.EmitIssue(r, r.message, block.DefRange)
			}
		}
		return nil
	}

	return nil
}
