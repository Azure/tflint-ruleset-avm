package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformParentIDVariableRule)

// parentScopeVariableNames lists parent-scope-specific variables that resource modules must not expose.
var parentScopeVariableNames = map[string]bool{
	"resource_group_name":        true,
	"resource_group_resource_id": true,
}

// TerraformParentIDVariableRule checks that AzAPI resource modules expose their parent scope as a required, non-nullable `parent_id` string variable.
type TerraformParentIDVariableRule struct {
	tflint.DefaultRule
}

// NewTerraformParentIDVariableRule returns a new rule.
func NewTerraformParentIDVariableRule() *TerraformParentIDVariableRule {
	return new(TerraformParentIDVariableRule)
}

func (r *TerraformParentIDVariableRule) Name() string {
	return "avm_terraform_parent_id_variable_required"
}

func (r *TerraformParentIDVariableRule) Enabled() bool {
	return true
}

func (r *TerraformParentIDVariableRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformParentIDVariableRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFRMFR1/"
}

// Check decodes the module class when the rule is used without the ruleset wrapper.
func (r *TerraformParentIDVariableRule) Check(runner tflint.Runner) error {
	moduleClass, err := decodeModuleClass(runner, r.Name())
	if err != nil {
		return err
	}
	return r.CheckWithModuleClass(runner, moduleClass)
}

// CheckWithModuleClass checks resource module roots and skips pattern and utility modules.
func (r *TerraformParentIDVariableRule) CheckWithModuleClass(runner tflint.Runner, moduleClass string) error {
	return checkResourceModuleRoot(runner, r.Name(), moduleClass, r.checkModule)
}

func (r *TerraformParentIDVariableRule) checkModule(runner tflint.Runner) error {
	blocks, err := moduleBlocks(runner)
	if err != nil {
		return err
	}
	resources := azapiManagedResources(runner, blocks)
	if len(resources) == 0 {
		return nil
	}

	var parentID *hclsyntax.Block
	var errList error
	emit := func(message string, issueRange hcl.Range) {
		if subErr := runner.EmitIssue(r, message, issueRange); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}
	for _, block := range blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		name := block.Labels[0]
		if name == "parent_id" && parentID == nil {
			parentID = block
		}
		if parentScopeVariableNames[name] {
			emit(fmt.Sprintf("Variable `%s` must not be declared; expose the parent scope through the `parent_id` variable instead", name), block.DefRange())
		}
	}

	if parentID == nil {
		emit("Variable `parent_id` must be declared to expose the parent scope of the primary AzAPI resource", hcl.Range{Filename: "variables.tf"})
	} else {
		if typeAttr, ok := parentID.Body.Attributes["type"]; !ok || hcl.ExprAsKeyword(typeAttr.Expr) != "string" {
			emit("Variable `parent_id` must be of type `string`", parentID.DefRange())
		}
		if defaultAttr, ok := parentID.Body.Attributes["default"]; ok {
			emit("Variable `parent_id` must be required and must not declare a default value", defaultAttr.SrcRange)
		}
		if nullableAttr, ok := parentID.Body.Attributes["nullable"]; !ok {
			emit("Variable `parent_id` must set `nullable = false`", parentID.DefRange())
		} else if nullable, ok := staticBool(parentID, "nullable"); !ok || nullable {
			emit("Variable `parent_id` must set `nullable = false`", nullableAttr.SrcRange)
		}
	}

	for _, resource := range resources {
		if resource.Labels[0] != "azapi_resource" || resource.Labels[1] != "this" {
			continue
		}
		attr, ok := resource.Body.Attributes["parent_id"]
		if !ok {
			emit("`azapi_resource.this` must set `parent_id = var.parent_id`", resource.DefRange())
		} else if !isVarParentIDReference(attr.Expr) {
			emit("`azapi_resource.this` must set `parent_id = var.parent_id`", attr.SrcRange)
		}
	}
	return errList
}

func isVarParentIDReference(expr hcl.Expression) bool {
	traversal, diags := hcl.AbsTraversalForExpr(expr)
	if diags.HasErrors() || len(traversal) != 2 {
		return false
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok || root.Name != "var" {
		return false
	}
	attr, ok := traversal[1].(hcl.TraverseAttr)
	return ok && attr.Name == "parent_id"
}
