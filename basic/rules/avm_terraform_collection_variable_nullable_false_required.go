package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformCollectionVariableNullableRule)

// TerraformCollectionVariableNullableRule checks that collection-shaped variables declare `nullable = false`.
type TerraformCollectionVariableNullableRule struct {
	tflint.DefaultRule
}

// NewTerraformCollectionVariableNullableRule returns a new rule.
func NewTerraformCollectionVariableNullableRule() *TerraformCollectionVariableNullableRule {
	return new(TerraformCollectionVariableNullableRule)
}

func (r *TerraformCollectionVariableNullableRule) Name() string {
	return "avm_terraform_collection_variable_nullable_false_required"
}

func (r *TerraformCollectionVariableNullableRule) Enabled() bool {
	return true
}

func (r *TerraformCollectionVariableNullableRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformCollectionVariableNullableRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/"
}

func (r *TerraformCollectionVariableNullableRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformCollectionVariableNullableRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_collection_variable_nullable_false_required since it is not an HCL file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		name := block.Labels[0]
		// The AVM tags interface requires `tags` to stay a nullable map(string) with a null default.
		if name == "tags" {
			continue
		}
		if typeName, _ := variableType(block); !collectionTypeNames[typeName] {
			continue
		}
		// A literal `nullable = true` is reported by avm_terraform_variable_nullable_true_disallowed, and non-literal values can't be checked.
		if _, set := block.Body.Attributes["nullable"]; set {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("Collection variable `%s` must set `nullable = false`", name),
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
