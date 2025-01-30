package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	goverison "github.com/hashicorp/go-version"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ tflint.Rule = new(ProviderVersionRule)

type ProviderVersionRule struct {
	tflint.DefaultRule
	ProviderName          string
	ProviderSource        string
	Version               string
	RecommendedConstraint string
	MustExist             bool
}

func NewProviderVersionRule(providerName, providerSource, ver, recConstr string, mustExist bool) *ProviderVersionRule {
	return &ProviderVersionRule{
		ProviderName:          providerName,
		ProviderSource:        providerSource,
		Version:               ver,
		MustExist:             mustExist,
		RecommendedConstraint: recConstr,
	}
}

func (m *ProviderVersionRule) Name() string {
	return fmt.Sprintf("provider_%s_version_constraint", m.ProviderName)
}

func (m *ProviderVersionRule) Enabled() bool {
	return true
}

func (m *ProviderVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (m *ProviderVersionRule) Check(r tflint.Runner) error {
	ver, err := version.NewVersion(m.Version)
	if err != nil {
		return fmt.Errorf("invalid version constraint: %s", err)
	}
	content, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: "required_providers",
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{
										Name: m.ProviderName,
									},
								},
							},
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}
	if len(content.Blocks) == 0 {
		return nil
	}
	providerFound := false
	requiredProviderFound := false
	for _, tb := range content.Blocks {
		for _, rpb := range tb.Body.Blocks {
			requiredProviderFound = true
			providerAttr, ok := rpb.Body.Attributes[m.ProviderName]
			if !ok {
				continue
			}
			providerFound = true
			provider := struct {
				Source  string `cty:"source"`
				Version string `cty:"version"`
			}{}
			wantType := cty.Object(map[string]cty.Type{
				"source":  cty.String,
				"version": cty.String,
			})
			if err = r.EvaluateExpr(providerAttr.Expr, &provider, &tflint.EvaluateExprOption{WantType: &wantType}); err != nil {
				return err
			}
			if !strings.EqualFold(provider.Source, m.ProviderSource) {
				return r.EmitIssue(m, fmt.Sprintf("provider `%s`'s source should be %s, got %s", m.ProviderName, m.ProviderSource, provider.Source), providerAttr.Range)
			}
			constraint, err := goverison.NewConstraint(provider.Version)
			if err != nil {
				return fmt.Errorf("invalid version constraint: %s", err)
			}
			if !constraint.Check(ver) {
				r.EmitIssue(m, fmt.Sprintf("provider `%s`'s version should satisfy %s, got %s. Recommended version constraint `%s`", m.ProviderName, m.Version, provider.Version, m.RecommendedConstraint), providerAttr.Range)
			}
		}
	}
	if !requiredProviderFound {
		return nil
	}
	if !providerFound && m.MustExist {
		return r.EmitIssue(m, fmt.Sprintf("`%s` provider should be declared in the `required_providers` block", m.ProviderName), content.Blocks[0].DefRange)
	}
	return nil
}
