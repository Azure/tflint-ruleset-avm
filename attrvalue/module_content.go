package attrvalue

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/terraform"
)

type AttrValueRule interface {
	GetResourceType() string
	GetNestedBlockType() *string
	GetAttributeName() string
}

func fetchAttrsAndContext(r AttrValueRule, runner tflint.Runner) (*terraform.Evaluator, []*hclext.Attribute, hcl.Diagnostics) {
	moduleEvaluator, diag := newModuleEvaluator(runner)
	if diag.HasErrors() {
		return nil, nil, diag
	}

	if r.GetNestedBlockType() != nil {
		attrs, diags := moduleEvaluator.getNestedBlockAttrs(r.GetResourceType(), *r.GetNestedBlockType(), r.GetAttributeName())
		return moduleEvaluator.Evaluator, attrs, diags
	}

	attrs, diags := moduleEvaluator.getSimpleAttrs(r.GetResourceType(), r.GetAttributeName())
	return moduleEvaluator.Evaluator, attrs, diags
}

func fetchResourcesAndContext(r AttrValueRule, runner tflint.Runner) ([]*hclext.Block, *terraform.Evaluator, hcl.Diagnostics) {
	moduleEvaluator, diag := newModuleEvaluator(runner)
	if diag.HasErrors() {
		return nil, nil, diag
	}

	if r.GetNestedBlockType() != nil {
		resources, diags := moduleEvaluator.getNestedResourcesWithBlockAttributes(r.GetResourceType(), *r.GetNestedBlockType(), r.GetAttributeName())
		return resources, moduleEvaluator.Evaluator, diags
	}

	resources, diags := moduleEvaluator.getSimpleResourcesWithAttributes(r.GetResourceType(), r.GetAttributeName())

	return resources, moduleEvaluator.Evaluator, diags
}
