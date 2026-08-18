package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"

	"github.com/Azure/tflint-ruleset-avm/azurerm/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/lonegunmanb/terraform-azurerm-schema/v4/generated"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(AzurermResourceTagRule)

// AzurermResourceTagRule checks whether the tags arg is specified if supported
type AzurermResourceTagRule struct {
	tflint.DefaultRule
}

func (r *AzurermResourceTagRule) Name() string {
	return "azurerm_resource_tag"
}

func (r *AzurermResourceTagRule) Enabled() bool {
	return false
}

func (r *AzurermResourceTagRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *AzurermResourceTagRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *AzurermResourceTagRule) Check(runner tflint.Runner) error {
	return Check(runner, r.CheckFile)
}

// NewAzurermResourceTagRule returns a new rule
func NewAzurermResourceTagRule() *AzurermResourceTagRule {
	return &AzurermResourceTagRule{}
}

// CheckFile checks whether the tags arg is specified if supported
func (r *AzurermResourceTagRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip azurerm_resource_tag since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var err error
	for _, block := range blocks {
		var subErr error
		switch block.Type {
		case "resource":
			subErr = r.visitAzResource(runner, block)
		}
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *AzurermResourceTagRule) visitAzResource(runner tflint.Runner, azBlock *hclsyntax.Block) error {
	resourceSchema, isAzureResource := generated.Resources[azBlock.Labels[0]]
	if !isAzureResource {
		return nil
	}
	_, isTagSupported := resourceSchema.Block.Attributes["tags"]
	_, isTagSet := azBlock.Body.Attributes["tags"]
	if isTagSupported && !isTagSet {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("`tags` argument is not set but supported in resource `%s`", azBlock.Labels[0]),
			azBlock.DefRange(),
		)
	}
	return nil
}
