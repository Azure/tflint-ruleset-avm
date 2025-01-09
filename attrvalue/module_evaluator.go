package attrvalue

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/terraform"
	"github.com/terraform-linters/tflint/terraform/addrs"
)

var AppFs = afero.Afero{
	Fs: afero.NewOsFs(),
}

type moduleEvaluator struct {
	*terraform.Evaluator
}

func newModuleEvaluator(runner tflint.Runner) (*moduleEvaluator, hcl.Diagnostics) {
	wd, _ := runner.GetOriginalwd()
	loader, err := terraform.NewLoader(AppFs, wd)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Summary: err.Error(),
		}}
	}
	config, diags := loader.LoadConfig(".", terraform.CallLocalModule)
	if diags.HasErrors() {
		return nil, diags
	}
	vvals, diags := terraform.VariableValues(config)
	if diags.HasErrors() {
		return nil, diags
	}
	return &moduleEvaluator{
		Evaluator: &terraform.Evaluator{
			Meta: &terraform.ContextMeta{
				Env:                "",
				OriginalWorkingDir: wd,
			},
			Config:         config,
			VariableValues: vvals,
			ModulePath:     addrs.RootModuleInstance,
		},
	}, nil
}

// getSimpleResourcesWithAttributes returns a slice of resources with the given resource type and the attribute if it exists.
func (e *moduleEvaluator) getSimpleResourcesWithAttributes(resourceType string, attributeName string) ([]*hclext.Block, hcl.Diagnostics) {
	resources, diags := e.getResourcesOfResourceTypeIncludingSpecifiedAttribute(attributeName)
	if diags.HasErrors() {
		return nil, diags
	}
	filteredResources := make([]*hclext.Block, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		filteredResources = append(filteredResources, resource)
	}
	return filteredResources, nil
}

// getSimpleAttrs returns a slice of attributes with the given attribute name from the resources of the given resource type.
func (e *moduleEvaluator) getSimpleAttrs(resourceType string, attributeName string) ([]*hclext.Attribute, hcl.Diagnostics) {
	resources, diags := e.getResourcesOfResourceTypeIncludingSpecifiedAttribute(attributeName)
	if diags.HasErrors() {
		return nil, diags
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		if attribute := e.getAttrFromBlock(resource, attributeName); attribute != nil {
			attrs = append(attrs, attribute)
		}
	}
	return attrs, nil
}

// getNestedResourcesWithBlockAttributes returns a slice of resources with the given resource type and the attribute if it exists.
func (e *moduleEvaluator) getNestedResourcesWithBlockAttributes(resourceType, nestedBlockType, attributeName string) ([]*hclext.Block, hcl.Diagnostics) {
	resources, diags := e.getResourcesOfResourceTypeIncludingBlocksWithSpecifiedAttribute(nestedBlockType, attributeName)
	if diags.HasErrors() {
		return nil, diags
	}
	filteredResources := make([]*hclext.Block, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		filteredResources = append(filteredResources, resource)
	}
	return filteredResources, nil
}

// getNestedBlockAttrs returns a slice of attributes with the given attribute name from the nested blocks of the given resource type.
func (e *moduleEvaluator) getNestedBlockAttrs(resourceType, nestedBlockType, attributeName string) ([]*hclext.Attribute, hcl.Diagnostics) {
	resources, diags := e.getResourcesOfResourceTypeIncludingBlocksWithSpecifiedAttribute(nestedBlockType, attributeName)
	if diags.HasErrors() {
		return nil, diags
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		for _, block := range resource.Body.Blocks {
			if attr := e.getAttrFromBlock(block, attributeName); attr != nil {
				attrs = append(attrs, attr)
			}
		}
	}
	return attrs, nil
}

// getAttrFromBlock returns the attribute with the given attribute name from the block.
func (e *moduleEvaluator) getAttrFromBlock(block *hclext.Block, attributeName string) *hclext.Attribute {
	attribute, exists := block.Body.Attributes[attributeName]
	if !exists {
		return nil
	}
	return attribute
}

func (e *moduleEvaluator) getResourcesOfResourceTypeIncludingSpecifiedAttribute(attributeName string) (*hclext.BodyContent, hcl.Diagnostics) {
	resources, diags := e.Config.Module.PartialContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     attributeName,
							Required: false,
						},
					},
				},
			},
		},
	}, e.Evaluator)

	return resources, diags
}

func (e *moduleEvaluator) getResourcesOfResourceTypeIncludingBlocksWithSpecifiedAttribute(nestedBlockType string, attributeName string) (*hclext.BodyContent, hcl.Diagnostics) {
	resources, diags := e.Config.Module.PartialContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: nestedBlockType,
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{
										Name:     attributeName,
										Required: false,
									},
								},
							},
						},
					},
				},
			},
		},
	}, e.Evaluator)

	return resources, diags
}
