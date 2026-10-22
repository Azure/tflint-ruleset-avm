package common

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
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

	reference, diagnostics := addrs.ParseRef(traversal)
	if diagnostics.HasErrors() || len(reference.Remaining) != 0 {
		return resourceAddress{}, fmt.Errorf(
			"ignored_resources entry %q must use resource_type.resource_name, resource.resource_type.resource_name, or data.resource_type.resource_name",
			value,
		)
	}

	var resource addrs.Resource
	switch subject := reference.Subject.(type) {
	case addrs.Resource:
		resource = subject
	case addrs.ResourceInstance:
		if subject.Key != addrs.NoKey {
			return resourceAddress{}, fmt.Errorf(
				"ignored_resources entry %q must use resource_type.resource_name, resource.resource_type.resource_name, or data.resource_type.resource_name",
				value,
			)
		}
		resource = subject.Resource
	default:
		return resourceAddress{}, fmt.Errorf(
			"ignored_resources entry %q must use resource_type.resource_name, resource.resource_type.resource_name, or data.resource_type.resource_name",
			value,
		)
	}

	blockType := "resource"
	if resource.Mode == addrs.DataResourceMode {
		blockType = "data"
	}
	return resourceAddress{blockType: blockType, typeName: resource.Type, name: resource.Name}, nil
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

func (r *ignoredResourceRunner) IsIgnoredResourceBlock(blockType string, labels []string) bool {
	return isIgnoredAddress(blockType, labels, r.ignored)
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

// IsIgnoredResourceBlock reports whether a static resource or data block is configured as ignored.
func IsIgnoredResourceBlock(runner tflint.Runner, blockType string, labels []string) bool {
	filteringRunner, ok := runner.(interface {
		IsIgnoredResourceBlock(string, []string) bool
	})
	return ok && filteringRunner.IsIgnoredResourceBlock(blockType, labels)
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
	return isIgnoredAddress(block.Type, block.Labels, ignored)
}

func isIgnoredAddress(blockType string, labels []string, ignored ignoredResources) bool {
	if len(labels) != 2 || (blockType != "resource" && blockType != "data") {
		return false
	}
	_, ok := ignored[resourceAddress{
		blockType: blockType,
		typeName:  labels[0],
		name:      labels[1],
	}]
	return ok
}
