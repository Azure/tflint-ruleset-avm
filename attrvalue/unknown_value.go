package attrvalue

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/terraform"
	"github.com/terraform-linters/tflint/terraform/addrs"
	"github.com/zclconf/go-cty/cty"
)

// UnknownValueRule checks whether an attribute value is null or part of a variable with no default value.
type UnknownValueRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType  string // e.g. "azurerm_storage_account"
	attributeName string // e.g. "account_replication_type"
}

var _ tflint.Rule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewUnknownValueRule(resourceType string, attributeName string) *UnknownValueRule {
	return &UnknownValueRule{
		resourceType:  resourceType,
		attributeName: attributeName,
	}
}

func (r *UnknownValueRule) Name() string {
	return fmt.Sprintf("%s.%s must be null", r.resourceType, r.attributeName)
}

func (r *UnknownValueRule) Enabled() bool {
	return true
}

func (r *UnknownValueRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *UnknownValueRule) Check(runner tflint.Runner) error {
	content, ctx, diags := prepareEvaluator(r, runner)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}
	// expr := content.Blocks[0].Body.Attributes[r.attributeName].Expr
	for _, block := range content.Blocks {
		if block.Labels[0] != r.resourceType {
			continue
		}
		for _, attr := range block.Body.Attributes {
			if attr.Name != r.attributeName {
				continue
			}
			expr := attr.Expr
			val, diags := ctx.EvaluateExpr(expr, cty.DynamicPseudoType)
			if diags.HasErrors() {
				return fmt.Errorf("could not evaluate expression: %s", diags)
			}
			if val.IsKnown() {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("invalid attribute value of `%s` - expecting unknown", r.attributeName),
					expr.Range(),
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func prepareEvaluator(r *UnknownValueRule, runner tflint.Runner) (*hclext.BodyContent, *terraform.Evaluator, hcl.Diagnostics) {
	var appFs afero.Afero
	// If we are using the tflint test runner then we need to create a new memory file system
	wd, _ := runner.GetOriginalwd()
	if _, ok := runner.(*helper.Runner); ok {
		appFs = afero.Afero{Fs: afero.NewMemMapFs()}
		fileName := "main.tf"
		mainTf, _ := runner.GetFile(fileName)
		file, _ := appFs.Create(fileName)
		file.Write(mainTf.Bytes)
	} else {
		appFs = afero.Afero{
			Fs: afero.NewBasePathFs(afero.NewOsFs(), wd),
		}
	}
	loader, _ := terraform.NewLoader(appFs, wd)
	config, _ := loader.LoadConfig(".", terraform.CallLocalModule)
	vvals, _ := terraform.VariableValues(config)
	ctx := &terraform.Evaluator{
		Meta: &terraform.ContextMeta{
			Env:                "",
			OriginalWorkingDir: wd,
		},
		Config:         config,
		VariableValues: vvals,
		CallStack:      terraform.NewCallStack(),
		ModulePath:     addrs.RootModuleInstance,
	}
	schema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     r.attributeName,
							Required: false,
						},
					},
				},
			},
		},
	}

	content, diags := config.Module.PartialContent(schema, ctx)
	return content, ctx, diags
}
