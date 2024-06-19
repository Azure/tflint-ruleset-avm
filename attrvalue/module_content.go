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

type AttrValueRule interface {
	GetResourceType() string
	GetNestedBlockType() *string
	GetAttributeName() string
}

// getSimpleResources returns a slice of resources with the given resource type and the attribute if it exists.
func getSimpleResourcesWithAttribute(module *terraform.Module, resourceType string, attributeName string,  ctx *terraform.Evaluator) ([]*hclext.Block, hcl.Diagnostics) {
	resources, diags := module.PartialContent(&hclext.BodySchema{
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
	}, ctx)
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
func getSimpleAttrs(module *terraform.Module, resourceType string, attributeName string, ctx *terraform.Evaluator) ([]*hclext.Attribute, hcl.Diagnostics) {
	resources, diags := module.PartialContent(&hclext.BodySchema{
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
	}, ctx)
	if diags.HasErrors() {
		return nil, diags
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		if attribute := getAttrFromBlock(resource, attributeName); attribute != nil {
			attrs = append(attrs, attribute)
		}
	}
	return attrs, nil
}


// getNestedResourcesWithAttribute returns a slice of resources with the given resource type and the attribute if it exists.
func getNestedResourcesWithAttribute(ctx *terraform.Evaluator, module *terraform.Module, resourceType, nestedBlockType, attributeName string) ([]*hclext.Block, hcl.Diagnostics) {
	content, diags := module.PartialContent(&hclext.BodySchema{
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
	}, ctx)
	if diags.HasErrors() {
		return nil, diags
	}
	filteredResources := make([]*hclext.Block, 0, len(content.Blocks))
	for _, resource := range content.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		filteredResources = append(filteredResources, resource)
	}
	return filteredResources, nil
}

// getNestedBlockAttrs returns a slice of attributes with the given attribute name from the nested blocks of the given resource type.
func getNestedBlockAttrs(ctx *terraform.Evaluator, module *terraform.Module, resourceType, nestedBlockType, attributeName string) ([]*hclext.Attribute, hcl.Diagnostics) {
	content, diags := module.PartialContent(&hclext.BodySchema{
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
	}, ctx)
	if diags.HasErrors() {
		return nil, diags
	}
	attrs := make([]*hclext.Attribute, 0, len(content.Blocks))
	for _, resource := range content.Blocks {
		if resource.Labels[0] != resourceType {
			continue
		}
		for _, block := range resource.Body.Blocks {
			if attr := getAttrFromBlock(block, attributeName); attr != nil {
				attrs = append(attrs, attr)
			}
		}
	}
	return attrs, nil
}

// getAttrFromBlock returns the attribute with the given attribute name from the block.
func getAttrFromBlock(block *hclext.Block, attributeName string) *hclext.Attribute {
	attribute, exists := block.Body.Attributes[attributeName]
	if !exists {
		return nil
	}
	return attribute
}

func fetchAttrsAndContext(r AttrValueRule, runner tflint.Runner) (*terraform.Evaluator, []*hclext.Attribute, hcl.Diagnostics) {
	// If we are using the tflint test runner then we need to create a new memory file system
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
		CallStack:      terraform.NewCallStack(),
		ModulePath:     addrs.RootModuleInstance,
	}

	if r.GetNestedBlockType() != nil {
		attrs, diags := getNestedBlockAttrs(ctx, config.Module, r.GetResourceType(), *r.GetNestedBlockType(), r.GetAttributeName())
		return ctx, attrs, diags
	}

	attrs, diags := getSimpleAttrs(config.Module, r.GetResourceType(), r.GetAttributeName(), ctx)

	return ctx, attrs, diags
}

func fetchResourcesAndContext(r AttrValueRule, runner tflint.Runner) (*terraform.Evaluator, []*hclext.Block, hcl.Diagnostics) {
	// If we are using the tflint test runner then we need to create a new memory file system
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
		CallStack:      terraform.NewCallStack(),
		ModulePath:     addrs.RootModuleInstance,
	}

	if r.GetNestedBlockType() != nil {
		resources, diags := getNestedResourcesWithAttribute(ctx, config.Module, r.GetResourceType(), *r.GetNestedBlockType(), r.GetAttributeName())
		return ctx, resources, diags
	}

	resources, diags := getSimpleResourcesWithAttribute(config.Module, r.GetResourceType(), r.GetAttributeName(), ctx)

	return ctx, resources, diags
}
