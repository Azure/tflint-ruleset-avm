package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformLockInterfaceRule)

type TerraformLockInterfaceRule struct {
	tflint.DefaultRule
}

func NewTerraformLockInterfaceRule() *TerraformLockInterfaceRule {
	return new(TerraformLockInterfaceRule)
}

func (t *TerraformLockInterfaceRule) Name() string {
	return "tfintlock"
}

func (t *TerraformLockInterfaceRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#resource-locks"
}

func (t *TerraformLockInterfaceRule) Enabled() bool {
	return false
}

func (t *TerraformLockInterfaceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *TerraformLockInterfaceRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}
	body, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "type"},
						{Name: "default"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		if variable.Labels[0] != "lock" {
			continue
		}
		// Check if the variable has a type attribute and that it is correct.
		typeattr, exists := variable.Body.Attributes["type"]
		if !exists {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`%s` variable type not declared", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}
		if interfaceErr := checkLockInterface(typeattr); interfaceErr != nil {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: %s", variable.Labels[0], interfaceErr.Error()),
				variable.DefRange,
			); err != nil {
				return err
			}
		}

		defaultattr, exists := variable.Body.Attributes["default"]
		if !exists {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: default not declared", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}
		defaultval, _ := defaultattr.Expr.Value(nil)
		if !defaultval.IsNull() {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: default value is not `null`", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkLockInterface(attr *hclext.Attribute) error {
	expr, ok := attr.Expr.(*hclsyntax.FunctionCallExpr)
	if !ok {
		return fmt.Errorf("expression is not a function call")
	}

	if expr.Name != "object" {
		return fmt.Errorf("expression function is not object()")
	}

	if len(expr.Args) != 1 {
		return fmt.Errorf("expression function object() does not have exactly one argument")
	}

	objExpr, ok := expr.Args[0].(*hclsyntax.ObjectConsExpr)
	if !ok {
		return fmt.Errorf("expression function object() argument is not a object")
	}

	em := objExpr.ExprMap()
	if len(em) != 2 {
		return fmt.Errorf("expression function object() argument does not have exactly two attributes")
	}

	for _, kvp := range em {
		key, _ := kvp.Key.Value(nil)
		switch key.AsString() {
		case "kind":
			v := hcl.ExprAsKeyword(kvp.Value)
			if v != "string" {
				return fmt.Errorf("expression function object() argument attribute `kind` value is not `string`")
			}
		case "name":
			namefn, ok := kvp.Value.(*hclsyntax.FunctionCallExpr)
			if !ok {
				return fmt.Errorf("expression function object() argument attribute `name` value is not a function call to `optional()`")
			}
			if namefn.Name != "optional" {
				return fmt.Errorf("expression function object() argument attribute `name` value is not `optional()` function")
			}
			if len(namefn.Args) != 2 {
				return fmt.Errorf("expression function object() argument attribute `name` value `optional()` does not have exactly two arguments")
			}
			if namefn.Args[0].Variables()[0].RootName() != "string" {
				return fmt.Errorf("expression function object() argument attribute `name` value `optional()` first argument is not `string`")
			}
			if v, _ := namefn.Args[1].Value(nil); !v.IsNull() {
				return fmt.Errorf("expression function object() argument attribute `name` value `optional()` second argument is not `null`")
			}

		default:
			return fmt.Errorf("expression function object() argument attribute `%s` is not valid", key.AsString())
		}
	}
	return nil
}
