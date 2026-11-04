package rules

import (
	"fmt"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type moduleClassRuleConfig struct {
	Severity         string   `hclext:"severity,optional"`
	ModuleClass      string   `hclext:"module_class,optional"`
	IgnoredResources []string `hclext:"ignored_resources,optional"`
}

// decodeModuleClass reads module_class from the rule configuration and defaults to "resource".
func decodeModuleClass(runner tflint.Runner, ruleName string) (string, error) {
	config := moduleClassRuleConfig{ModuleClass: "resource"}
	if err := runner.DecodeRuleConfig(ruleName, &config); err != nil {
		return "", err
	}
	return config.ModuleClass, nil
}

// checkResourceModuleRoot runs check only for root modules configured as resource modules.
func checkResourceModuleRoot(runner tflint.Runner, ruleName, moduleClass string, check func(tflint.Runner) error) error {
	switch moduleClass {
	case "resource":
	case "pattern", "utility":
		return nil
	default:
		return fmt.Errorf(
			"invalid configuration for rule %q: module_class must be one of %q, %q, or %q, got %q",
			ruleName, "resource", "pattern", "utility", moduleClass,
		)
	}

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}
	return check(runner)
}

// moduleBlocks returns the top-level HCL blocks of every module file, in file name order.
func moduleBlocks(runner tflint.Runner) ([]*hclsyntax.Block, error) {
	var blocks []*hclsyntax.Block
	err := ForFiles(runner, func(_ tflint.Runner, file *hcl.File) error {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			return nil
		}
		blocks = append(blocks, body.Blocks...)
		return nil
	})
	return blocks, err
}

// isAzapiManagedResource reports whether block is a non-ignored AzAPI managed resource. Data sources and azapi_resource_list are excluded.
func isAzapiManagedResource(runner tflint.Runner, block *hclsyntax.Block) bool {
	return block.Type == "resource" &&
		len(block.Labels) == 2 &&
		strings.HasPrefix(block.Labels[0], "azapi_") &&
		block.Labels[0] != "azapi_resource_list" &&
		!common.IsIgnoredResourceBlock(runner, block.Type, block.Labels)
}

// azapiManagedResources returns the AzAPI managed resources in blocks.
func azapiManagedResources(runner tflint.Runner, blocks []*hclsyntax.Block) []*hclsyntax.Block {
	var resources []*hclsyntax.Block
	for _, block := range blocks {
		if isAzapiManagedResource(runner, block) {
			resources = append(resources, block)
		}
	}
	return resources
}
