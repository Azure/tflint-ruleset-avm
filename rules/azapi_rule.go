// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rules

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/blockquery"
	"github.com/Azure/tflint-ruleset-avm/modulecontent"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AzApiRule runs the specified gjson query on the `body` attribute of `azapi_resource` resources and checks if the result is as expected.
type AzApiRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation
	blockquery.BlockQuery
	expected          []cty.Value
	maximumApiVersion string
	minimumApiVersion string
	link              string
	resourceType      string
	ruleName          string
	mustExist         bool
}

var _ tflint.Rule = &AzApiRule{}
var _ modulecontent.BlockFetcher = &AzApiRule{}

// AzApiRule creates a rule to check the `body` attribute of `azapi_resource` resources.
// The `query` parameter is a gjson query string to run against the `body` attribute.
// The `compareFunc` parameter is a function to compare the result of the query with the expected results. E.g. `blockquery.IsOneOf`.
// The `expectedResults` parameter is a list of expected results, use the `blockquery.New*Results` functions to create them.
// The resource type is the first part of the `type` attribute of the resource, e.g. "Microsoft.Compute/virtualMachines" for VMs.
// Use the `minimumApiVersion` and `maximumApiVersion` parameters to filter resources based on their API version.
func NewAzApiRuleQueryMustExist(
	ruleName, link, resourceType, minimumApiVersion, maximumApiVersion, query string,
	compareFunc blockquery.ResultCompareFunc,
	expectedResults ...cty.Value,
) *AzApiRule {
	return &AzApiRule{
		BlockQuery: blockquery.NewBlockQuery(
			"resource",
			"azapi_resource",
			[]string{"type", "name"},
			"body",
			query,
			compareFunc,
		),
		expected:          expectedResults,
		link:              link,
		maximumApiVersion: maximumApiVersion,
		minimumApiVersion: minimumApiVersion,
		resourceType:      resourceType,
		ruleName:          ruleName,
		mustExist:         true,
	}
}

func NewAzApiRuleQueryOptionalExist(
	ruleName, link, resourceType, minimumApiVersion, maximumApiVersion, query string,
	compareFunc blockquery.ResultCompareFunc,
	expectedResults ...cty.Value,
) *AzApiRule {
	return &AzApiRule{
		BlockQuery: blockquery.NewBlockQuery(
			"resource",
			"azapi_resource",
			[]string{"type", "name"},
			"body",
			query,
			compareFunc,
		),
		expected:          expectedResults,
		link:              link,
		maximumApiVersion: maximumApiVersion,
		minimumApiVersion: minimumApiVersion,
		resourceType:      resourceType,
		ruleName:          ruleName,
		mustExist:         false,
	}
}

func (r *AzApiRule) Link() string {
	return r.link
}

func (r *AzApiRule) Enabled() bool {
	return true
}

func (r *AzApiRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AzApiRule) Name() string {
	return r.ruleName
}

func (r *AzApiRule) LabelOne() string {
	return "azapi_resource"
}

func (r *AzApiRule) LabelNames() []string {
	return []string{"type", "name"}
}

func (r *AzApiRule) BlockType() string {
	return "resource"
}

func (r *AzApiRule) Attributes() []string {
	return []string{"name", "type", "body"}
}

func (r *AzApiRule) Check(runner tflint.Runner) error {
	return r.queryResource(runner, cty.DynamicPseudoType)
}

func (r *AzApiRule) queryResource(runner tflint.Runner, ct cty.Type) error {
	ctx, resources, diags := modulecontent.FetchBlocks(r, runner)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}
	for _, resource := range resources {
		typeAttr, typeAttrExists := resource.Body.Attributes["type"]
		if !typeAttrExists {
			if err := runner.EmitIssue(
				r,
				"Resource does not have a `type` attribute",
				resource.DefRange,
			); err != nil {
				return err
			}
			continue
		}
		typeVal, diags := ctx.EvaluateExpr(typeAttr.Expr, cty.String)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate type expression: %s", diags)
		}
		typeStr := typeVal.AsString()
		if !checkAzApiType(typeStr, r.resourceType, r.minimumApiVersion, r.maximumApiVersion) {
			continue
		}
		bodyAttr, bodyAttrExists := resource.Body.Attributes["body"]
		if !bodyAttrExists {
			if err := runner.EmitIssue(
				r,
				"Resource does not have a `body` attribute",
				resource.DefRange,
			); err != nil {
				return err
			}
			continue
		}
		val, diags := ctx.EvaluateExpr(bodyAttr.Expr, ct)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate body expression: %s", diags)
		}
		qr, err := blockquery.QueryCty(val, r.Query)
		if err != nil {
			notExistsErr := &blockquery.QueryErrorNotFound{Query: r.Query}
			if errors.As(err, &notExistsErr) {
				if r.mustExist {
					if err := runner.EmitIssue(
						r,
						err.Error(),
						bodyAttr.Range,
					); err != nil {
						return err
					}
				}
				continue
			}
			return fmt.Errorf("could not query value: %w", err)
		}
		ok, msg, err := r.CompareFunc(qr, r.expected...)
		if err != nil {
			return fmt.Errorf("could not compare values: %w", err)
		}
		if !ok {
			if err := runner.EmitIssue(
				r,
				msg,
				bodyAttr.Range,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkAzApiType(gotType, wantType, minimumApiVersion, maximumApiVersion string) bool {
	gotSplit := strings.Split(gotType, "@")
	if len(gotSplit) != 2 {
		return false
	}
	if !strings.EqualFold(gotSplit[0], wantType) {
		return false
	}
	if minimumApiVersion != "" {
		if gotSplit[1] < minimumApiVersion {
			return false
		}
	}
	if maximumApiVersion != "" {
		if gotSplit[1] > maximumApiVersion {
			return false
		}
	}
	return true
}
