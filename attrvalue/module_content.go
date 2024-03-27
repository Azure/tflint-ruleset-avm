package attrvalue

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint/terraform"
	"github.com/terraform-linters/tflint/terraform/addrs"
)

type AttrValueRule interface {
	GetResourceType() string
	GetNestedBlockType() *string
	GetAttributeName() string
}

// getSimpleAttrs returns a slice of attributes with the given attribute name from the resources of the given resource type.
func getSimpleAttrs(ctx *terraform.Evaluator, module *terraform.Module, resourceType string, attributeName string) ([]*hclext.Attribute, hcl.Diagnostics) {
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
		if attribute := getAttrFromBlock(resource, attributeName); attribute != nil {
			attrs = append(attrs, attribute)
		}
	}
	return attrs, nil
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
		if len(resource.Labels) != 2 && resource.Labels[0] != resourceType {
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

	if r.GetNestedBlockType() != nil {
		attrs, diags := getNestedBlockAttrs(ctx, config.Module, r.GetResourceType(), *r.GetNestedBlockType(), r.GetAttributeName())
		return ctx, attrs, diags
	}

	attrs, diags := getSimpleAttrs(ctx, config.Module, r.GetResourceType(), r.GetAttributeName())

	return ctx, attrs, diags
}
