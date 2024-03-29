package attrvalue

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ AttrValueRule = baseValue{}

type baseValue struct {
	resourceType    string // e.g. "azurerm_storage_account"
	nestedBlockType *string
	attributeName   string // e.g. "account_replication_type"
	enabled         bool
	link            string
	severity        tflint.Severity
}

func (b baseValue) GetNestedBlockType() *string {
	return b.nestedBlockType
}

func newBaseValue(
	resourceType string,
	nestedBlockType *string,
	attributeName string,
	enabled bool,
	link string,
	severity tflint.Severity,
) baseValue {
	return baseValue{
		resourceType:    resourceType,
		nestedBlockType: nestedBlockType,
		attributeName:   attributeName,
		enabled:         enabled,
		link:            link,
		severity:        severity,
	}
}

func (b baseValue) GetResourceType() string {
	return b.resourceType
}

func (b baseValue) GetAttributeName() string {
	return b.attributeName
}

func (b baseValue) Enabled() bool {
	return b.enabled
}

func (b baseValue) Severity() tflint.Severity {
	return b.severity
}

func (b baseValue) checkAttributes(r tflint.Runner, ct cty.Type, c func(*hclext.Attribute, cty.Value) error) error {
	ctx, attrs, diags := fetchAttrsAndContext(b, r)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}
	for _, attr := range attrs {
		val, diags := ctx.EvaluateExpr(attr.Expr, ct)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate expression: %s", diags)
		}

		if err := c(attr, val); err != nil {
			return err
		}
	}
	return nil
}
