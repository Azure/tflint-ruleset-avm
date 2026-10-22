package common

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type resourceAddress struct {
	blockType string
	typeName  string
	name      string
}

type ignoredResources map[resourceAddress]struct{}

type ignoredResourceRunner struct {
	tflint.Runner
	ignored ignoredResources
}

func parseIgnoredResources(values []string) (ignoredResources, error) {
	ignored := make(ignoredResources, len(values))
	for _, value := range values {
		address, err := parseResourceAddress(value)
		if err != nil {
			return nil, err
		}
		ignored[address] = struct{}{}
	}
	return ignored, nil
}

func parseResourceAddress(value string) (resourceAddress, error) {
	traversal, diagnostics := hclsyntax.ParseTraversalAbs(
		[]byte(value),
		"ignored_resources",
		hcl.Pos{Line: 1, Column: 1},
	)
	if diagnostics.HasErrors() {
		return resourceAddress{}, fmt.Errorf("ignored_resources entry %q must be a Terraform resource or data address", value)
	}

	parts := traversalParts(traversal)
	switch {
	case len(parts) == 2:
		return resourceAddress{blockType: "resource", typeName: parts[0], name: parts[1]}, nil
	case len(parts) == 3 && parts[0] == "resource":
		return resourceAddress{blockType: "resource", typeName: parts[1], name: parts[2]}, nil
	case len(parts) == 3 && parts[0] == "data":
		return resourceAddress{blockType: "data", typeName: parts[1], name: parts[2]}, nil
	default:
		return resourceAddress{}, fmt.Errorf(
			"ignored_resources entry %q must use resource_type.resource_name, resource.resource_type.resource_name, or data.resource_type.resource_name",
			value,
		)
	}
}

func traversalParts(traversal hcl.Traversal) []string {
	parts := make([]string, 0, len(traversal))
	for _, step := range traversal {
		switch step := step.(type) {
		case hcl.TraverseRoot:
			parts = append(parts, step.Name)
		case hcl.TraverseAttr:
			parts = append(parts, step.Name)
		default:
			return nil
		}
	}
	return parts
}

func (r *ignoredResourceRunner) GetModuleContent(
	schema *hclext.BodySchema,
	option *tflint.GetModuleContentOption,
) (*hclext.BodyContent, error) {
	content, err := r.Runner.GetModuleContent(schema, option)
	if err != nil {
		return nil, err
	}
	return filterBodyContent(content, r.ignored), nil
}

func (r *ignoredResourceRunner) GetResourceContent(
	resourceName string,
	schema *hclext.BodySchema,
	option *tflint.GetModuleContentOption,
) (*hclext.BodyContent, error) {
	content, err := r.Runner.GetResourceContent(resourceName, schema, option)
	if err != nil {
		return nil, err
	}
	return filterBodyContent(content, r.ignored), nil
}

func (r *ignoredResourceRunner) FilterIgnoredBlocks(blocks []*hclext.Block) []*hclext.Block {
	return filterBlocks(blocks, r.ignored)
}

// FilterIgnoredBlocks removes configured resource and data blocks when the runner supports exclusions.
func FilterIgnoredBlocks(runner tflint.Runner, blocks []*hclext.Block) []*hclext.Block {
	filteringRunner, ok := runner.(interface {
		FilterIgnoredBlocks([]*hclext.Block) []*hclext.Block
	})
	if !ok {
		return blocks
	}
	return filteringRunner.FilterIgnoredBlocks(blocks)
}

func filterBodyContent(content *hclext.BodyContent, ignored ignoredResources) *hclext.BodyContent {
	if content == nil {
		return nil
	}
	filtered := content.Copy()
	filtered.Blocks = filterBlocks(filtered.Blocks, ignored)
	return filtered
}

func filterBlocks(blocks []*hclext.Block, ignored ignoredResources) []*hclext.Block {
	filtered := make([]*hclext.Block, 0, len(blocks))
	for _, block := range blocks {
		if isIgnoredBlock(block, ignored) {
			continue
		}
		filtered = append(filtered, block)
	}
	return filtered
}

func isIgnoredBlock(block *hclext.Block, ignored ignoredResources) bool {
	if len(block.Labels) != 2 || (block.Type != "resource" && block.Type != "data") {
		return false
	}
	_, ok := ignored[resourceAddress{
		blockType: block.Type,
		typeName:  block.Labels[0],
		name:      block.Labels[1],
	}]
	return ok
}
