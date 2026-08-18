package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"

	"github.com/Azure/tflint-ruleset-avm/azurerm/project"
	"github.com/ahmetb/go-linq/v3"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/lonegunmanb/terraform-azurerm-schema/v4/generated"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(AzurermArgOrderRule)

// AzurermArgOrderRule checks whether the arguments in a block are sorted in azure doc order
type AzurermArgOrderRule struct {
	tflint.DefaultRule
}

func (r *AzurermArgOrderRule) Enabled() bool {
	return false
}

func (r *AzurermArgOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *AzurermArgOrderRule) Check(runner tflint.Runner) error {
	return Check(runner, r.CheckFile)
}

func (r *AzurermArgOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// NewAzurermArgOrderRule returns a new rule
func NewAzurermArgOrderRule() *AzurermArgOrderRule {
	return &AzurermArgOrderRule{}
}

// Name returns the rule name
func (r *AzurermArgOrderRule) Name() string {
	return "azurerm_arg_order"
}

// CheckFile checks whether the arguments in a block are sorted in codex order
func (r *AzurermArgOrderRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip azurerm_arg_order since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var err error
	for _, block := range blocks {
		var subErr error
		typeWanted := linq.From([]string{"provider", "resource", "data"}).Contains(block.Type)
		if !typeWanted {
			continue
		}
		isAzProviderBlock := block.Type == "provider" && block.Labels[0] == "azurerm"
		collection := generated.Resources
		if block.Type == "data" {
			collection = generated.DataSources
		}
		_, isAzBlock := collection[block.Labels[0]]
		if typeWanted && (isAzProviderBlock || isAzBlock) {
			subErr = r.visitAzBlock(runner, block)
		}
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *AzurermArgOrderRule) visitAzBlock(runner tflint.Runner, azBlock *hclsyntax.Block) error {
	emitter := func(block Block) error {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("Arguments are expected to be sorted in following order:\n%s", block.ToString()),
			block.DefRange(),
		)
	}
	file, _ := runner.GetFile(azBlock.Range().Filename)
	b := BuildResourceBlock(azBlock, file, emitter)
	return b.CheckBlock()
}
