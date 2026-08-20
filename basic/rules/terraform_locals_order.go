package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformLocalsOrderRule checks whether comments use the preferred syntax
type TerraformLocalsOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformLocalsOrderRule returns a new rule
func NewTerraformLocalsOrderRule() *TerraformLocalsOrderRule {
	return &TerraformLocalsOrderRule{}
}

// Name returns the rule name
func (r *TerraformLocalsOrderRule) Name() string {
	return "terraform_locals_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformLocalsOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformLocalsOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformLocalsOrderRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR32/"
}

// Check checks whether single line comments is used
func (r *TerraformLocalsOrderRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformLocalsOrderRule) checkFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_locals_order check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var err error
	for _, block := range blocks {
		if block.Type != "locals" {
			continue
		}
		if subErr := r.checkLocalsOrder(runner, block); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformLocalsOrderRule) checkLocalsOrder(runner tflint.Runner, block *hclsyntax.Block) error {
	attributes := r.attributesInLineOrder(block)
	if r.sorted(attributes) {
		return nil
	}
	return r.suggestedOrder(runner, block, attributes)
}

func (r *TerraformLocalsOrderRule) suggestedOrder(runner tflint.Runner, block *hclsyntax.Block, attributes []*hclsyntax.Attribute) error {
	sort.Slice(attributes, func(x, y int) bool {
		return attributes[x].Name < attributes[y].Name
	})
	file, err := runner.GetFile(block.Range().Filename)
	if err != nil {
		return err
	}
	var localsHclTxts []string
	for _, a := range attributes {
		localsHclTxts = append(localsHclTxts, string(a.SrcRange.SliceBytes(file.Bytes)))
	}
	localsHclTxt := strings.Join(localsHclTxts, "\n")
	localsHclTxt = fmt.Sprintf("%s {\n%s\n}", block.Type, localsHclTxt)
	formattedTxt := string(hclwrite.Format([]byte(localsHclTxt)))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended locals order:\n%s", formattedTxt),
		block.DefRange(),
	)
}

func (r *TerraformLocalsOrderRule) sorted(attributes []*hclsyntax.Attribute) bool {
	var attrNames []string
	for _, a := range attributes {
		attrNames = append(attrNames, a.Name)
	}
	return sort.StringsAreSorted(attrNames)
}

func (r *TerraformLocalsOrderRule) attributesInLineOrder(block *hclsyntax.Block) []*hclsyntax.Attribute {
	var attributes []*hclsyntax.Attribute
	for _, attribute := range block.Body.Attributes {
		attributes = append(attributes, attribute)
	}
	sort.Slice(attributes, func(x, y int) bool {
		posX := attributes[x].SrcRange.Start
		posY := attributes[y].SrcRange.Start
		if posX.Line == posY.Line {
			return posX.Column < posY.Column
		}
		return posX.Line < posY.Line
	})
	return attributes
}
