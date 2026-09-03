package rules

import (
	"errors"
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ tflint.Rule = new(AzapiResourceTagRule)

type tagCapabilityResolver func(resourceType string) (tagcapability.Status, error)

// AzapiResourceTagRule checks that managed AzAPI resources propagate the standard tags input.
type AzapiResourceTagRule struct {
	tflint.DefaultRule
	resolveTagCapability tagCapabilityResolver
}

// NewAzapiResourceTagRule returns an AzAPI tags propagation rule.
func NewAzapiResourceTagRule() *AzapiResourceTagRule {
	return newAzapiResourceTagRule(embeddedTagCapability)
}

func newAzapiResourceTagRule(resolve tagCapabilityResolver) *AzapiResourceTagRule {
	return &AzapiResourceTagRule{resolveTagCapability: resolve}
}

func (r *AzapiResourceTagRule) Name() string {
	return "avm_azapi_resource_tags_required"
}

func (r *AzapiResourceTagRule) Link() string {
	return "https://aka.ms/avm/spec/TFFR9"
}

func (r *AzapiResourceTagRule) Enabled() bool {
	return true
}

func (r *AzapiResourceTagRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AzapiResourceTagRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(azapiResourceTagBodySchema, &tflint.GetModuleContentOption{
		ExpandMode: tflint.ExpandModeNone,
	})
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) != 2 || block.Labels[0] != "azapi_resource" {
			continue
		}
		typeAttribute, ok := block.Body.Attributes["type"]
		if !ok {
			continue
		}

		var resourceType string
		wantType := cty.String
		if err := runner.EvaluateExpr(typeAttribute.Expr, &resourceType, &tflint.EvaluateExprOption{WantType: &wantType}); err != nil {
			logger.Debug("skip azapi_resource_tag because the resource type cannot be evaluated: %s", err)
			continue
		}

		status, err := r.resolveTagCapability(resourceType)
		if err != nil {
			if knownCapabilityLookupError(err) {
				logger.Debug("skip azapi_resource_tag because the tag capability is unavailable: %s", err)
				continue
			}
			return fmt.Errorf("resolve tags support for %q: %w", resourceType, err)
		}

		tagsAttribute, hasTags := block.Body.Attributes["tags"]
		switch status {
		case tagcapability.StatusWritable:
			if !hasTags {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("AzAPI resource type `%s` supports tags and must set `tags`", resourceType),
					block.DefRange,
				); err != nil {
					return err
				}
			}
		case tagcapability.StatusReadOnly, tagcapability.StatusUnsupported:
			if hasTags {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("AzAPI resource type `%s` does not support writable tags and must not set `tags`", resourceType),
					tagsAttribute.Range,
				); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("resolve tags support for %q: unexpected property status %q", resourceType, status)
		}
	}
	return nil
}

var azapiResourceTagBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "type"},
					{Name: "tags"},
				},
			},
		},
	},
}

func knownCapabilityLookupError(err error) bool {
	return errors.Is(err, tagcapability.ErrUnknownResource)
}

func embeddedTagCapability(resourceType string) (tagcapability.Status, error) {
	snapshot, err := tagcapability.Default()
	if err != nil {
		return "", err
	}
	return snapshot.Lookup(resourceType)
}
