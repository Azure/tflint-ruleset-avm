package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// DisallowedProviderRule checks that a specific provider is not used in the module.
// It performs two checks when ProviderName is "azurerm":
//  1. The `required_providers` block must not declare a provider whose source is hashicorp/azurerm.
//  2. No resource / data / ephemeral blocks whose first label starts with "azurerm_" should exist.
//
// If any are found, an issue is emitted.
// The rule is generic for potential future expansion by changing ProviderName/ProviderSource/Prefix.
// Currently only azurerm is required per user request.
var _ tflint.Rule = new(DisallowedProviderRule)

type DisallowedProviderRule struct {
	tflint.DefaultRule
	ProviderName   string
	ProviderSource string
}

func NewDisallowedProviderRule(providerName, providerSource string) *DisallowedProviderRule {
	return &DisallowedProviderRule{ProviderName: providerName, ProviderSource: providerSource}
}

func (r *DisallowedProviderRule) Name() string {
	// naming convention: provider_<name>_disallowed
	return fmt.Sprintf("provider_%s_disallowed", r.ProviderName)
}

func (r *DisallowedProviderRule) Enabled() bool { return true }

func (r *DisallowedProviderRule) Severity() tflint.Severity { return tflint.ERROR }

// requiredProvidersBodySchema dynamically requests the attribute for the provider name inside required_providers.
func (r *DisallowedProviderRule) requiredProvidersBodySchema() *hclext.BodySchema {
	return &hclext.BodySchema{Blocks: []hclext.BlockSchema{{Type: "terraform", Body: &hclext.BodySchema{Blocks: []hclext.BlockSchema{{Type: "required_providers", Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: r.ProviderName}}}}}}}}}
}

var disallowedProviderUsageSchema = &hclext.BodySchema{Blocks: []hclext.BlockSchema{
	{Type: "resource", LabelNames: []string{"type", "name"}},
	{Type: "data", LabelNames: []string{"type", "name"}},
	{Type: "ephemeral", LabelNames: []string{"type", "name"}},
}}

func (r *DisallowedProviderRule) Check(run tflint.Runner) error {
	path, err := run.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() { // only root module enforcement for consistency with other rules
		return nil
	}

	// 1. Check required_providers for disallowed provider with matching source
	content, err := run.GetModuleContent(r.requiredProvidersBodySchema(), &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}
	for _, tb := range content.Blocks { // terraform blocks
		for _, rp := range tb.Body.Blocks { // required_providers blocks
			if attr, ok := rp.Body.Attributes[r.ProviderName]; ok {
				provider := struct {
					Source string `cty:"source"`
				}{}
				wantType := cty.Object(map[string]cty.Type{"source": cty.String})
				if evalErr := run.EvaluateExpr(attr.Expr, &provider, &tflint.EvaluateExprOption{WantType: &wantType}); evalErr == nil {
					if strings.EqualFold(provider.Source, r.ProviderSource) {
						if issueErr := run.EmitIssue(r, fmt.Sprintf("provider '%s' (source %s) is disallowed", r.ProviderName, r.ProviderSource), attr.Range); issueErr != nil {
							return issueErr
						}
					}
				}
			}
		}
	}

	// 2. Scan for any usage blocks with provider name prefix
	usage, err := run.GetModuleContent(disallowedProviderUsageSchema, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	var errs error
	for _, blk := range usage.Blocks {
		if len(blk.Labels) == 0 { // defensive
			continue
		}
		first := blk.Labels[0]
		if strings.HasPrefix(first, r.ProviderName+"_") {
			if e := run.EmitIssue(r, fmt.Sprintf("%s block type '%s' is disallowed (provider '%s')", blk.Type, first, r.ProviderName), blk.DefRange); e != nil {
				errs = multierror.Append(errs, e)
			}
		}
	}
	return errs
}

func (r *DisallowedProviderRule) Link() string { return "" }
