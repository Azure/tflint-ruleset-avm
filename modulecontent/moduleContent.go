// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package modulecontent

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/terraform"
	"github.com/terraform-linters/tflint/terraform/addrs"
)

// AppFs is the virtual filesystem we use to that we can mock testing with the real Terraform evaluator,
// bypassing the tflint test runner.
var AppFs = afero.Afero{
	Fs: afero.NewOsFs(),
}

// BlockFetcher is an interface that permits partial content to be evaluated within a Terraform module.
// If a tflint rule satisfies this interface it can use the FetchResources and FetchAttributes functions to
// retrieve the resources and attributes of a given resource type.
type BlockFetcher interface {
	BlockType() string    // The type of block to fetch, e.g. `resource`.
	LabelOne() string     // The value of the first label of the block to fetch, e.g. `azapi_resource`.
	LabelNames() []string // The labels of the block to fetch, e.g. `["type", "name"]` for Terraform resources.
	Attributes() []string // The attributes to fetch from the block.
}

// FetchResources fetches the attributes of given resource type and the attribute if they exist.
func FetchAttributes(f BlockFetcher, runner tflint.Runner) (*terraform.Evaluator, []*hclext.Attribute, hcl.Diagnostics) {
	config, ctx, diags := initEvaluator(runner)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	attrs, diags := getAttributesFilterByLabelOne(ctx, config.Module, f)
	return ctx, attrs, diags
}

// FetchBlocks returns a slice of resources with the given resource type and the attribute if they exist.
func FetchBlocks(f BlockFetcher, runner tflint.Runner) (*terraform.Evaluator, []*hclext.Block, hcl.Diagnostics) {
	config, ctx, diags := initEvaluator(runner)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	blocks, diags := blocksFilterByLabelOne(ctx, config.Module, f)

	return ctx, blocks, diags
}

// blocksFilterByLabelOne returns a slice of resources with the given resource type and the attribute if they exist.
func blocksFilterByLabelOne(ctx *terraform.Evaluator, module *terraform.Module, bf BlockFetcher) ([]*hclext.Block, hcl.Diagnostics) {
	resources, diags := blocksWithPartialContent(ctx, module, bf)
	if diags.HasErrors() {
		return nil, diags
	}
	filteredResources := make([]*hclext.Block, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != bf.LabelOne() {
			continue
		}
		filteredResources = append(filteredResources, resource)
	}
	return filteredResources, nil
}

// getAttributesFilterByLabelOne returns a slice of attributes with the given attribute name from the resources of the given resource type.
func getAttributesFilterByLabelOne(ctx *terraform.Evaluator, module *terraform.Module, bf BlockFetcher) ([]*hclext.Attribute, hcl.Diagnostics) {
	resources, diags := blocksWithPartialContent(ctx, module, bf)
	if diags.HasErrors() {
		return nil, diags
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != bf.LabelOne() {
			continue
		}
		for _, attribute := range bf.Attributes() {
			if attribute := attrFromBlock(resource, attribute); attribute != nil {
				attrs = append(attrs, attribute)
			}
		}
	}
	return attrs, nil
}

// attrFromBlock returns the attribute with the given attribute name from the block.
func attrFromBlock(block *hclext.Block, attributeName string) *hclext.Attribute {
	attribute, exists := block.Body.Attributes[attributeName]
	if !exists {
		return nil
	}
	return attribute
}

// initEvaluator initializes the evaluator with the given runner.
// This uses a virtual filesystem to load the Terraform configuration so we can use it in prod and testing.
// It dows not use the tflint test runner as this limits the tests we can run.
// e.g. using this we have support for `optional()` evaluation, etc.
func initEvaluator(runner tflint.Runner) (*terraform.Config, *terraform.Evaluator, hcl.Diagnostics) {
	wd, _ := runner.GetOriginalwd()
	loader, err := terraform.NewLoader(AppFs, wd)
	if err != nil {
		return nil, nil, hcl.Diagnostics{{
			Summary: err.Error(),
		}}
	}
	config, diags := loader.LoadConfig(".", terraform.CallLocalModule)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	vvals, diags := terraform.VariableValues(config)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	ctx := &terraform.Evaluator{
		Meta: &terraform.ContextMeta{
			Env:                "",
			OriginalWorkingDir: wd,
		},
		Config:         config,
		VariableValues: vvals,
		ModulePath:     addrs.RootModuleInstance,
	}
	return config, ctx, nil
}

// blocksWithPartialContent returns the blocks with the given resource type and the attribute if they exist.
func blocksWithPartialContent(ctx *terraform.Evaluator, module *terraform.Module, bf BlockFetcher) (*hclext.BodyContent, hcl.Diagnostics) {
	attrSchema := make([]hclext.AttributeSchema, 0, len(bf.Attributes()))
	for _, attr := range bf.Attributes() {
		attrSchema = append(attrSchema, hclext.AttributeSchema{
			Name:     attr,
			Required: false,
		})
	}
	resources, diags := module.PartialContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       bf.BlockType(),
				LabelNames: bf.LabelNames(),
				Body: &hclext.BodySchema{
					Attributes: attrSchema,
				},
			},
		},
	}, ctx)
	return resources, diags
}
