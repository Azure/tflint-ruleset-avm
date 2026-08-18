package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"sort"
	"strings"
)

var _ tflint.Rule = &TerraformRequiredProvidersDeclarationRule{}

// TerraformRequiredProvidersDeclarationRule checks whether the required_providers block is declared in terraform block and whether the args of it are sorted in alphabetic order
type TerraformRequiredProvidersDeclarationRule struct {
	tflint.DefaultRule
}

func (r *TerraformRequiredProvidersDeclarationRule) Enabled() bool {
	return false
}

func (r *TerraformRequiredProvidersDeclarationRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformRequiredProvidersDeclarationRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// NewTerraformRequiredProvidersDeclarationRule returns a new rule
func NewTerraformRequiredProvidersDeclarationRule() *TerraformRequiredProvidersDeclarationRule {
	return &TerraformRequiredProvidersDeclarationRule{}
}

// Name returns the rule name
func (r *TerraformRequiredProvidersDeclarationRule) Name() string {
	return "terraform_required_providers_declaration"
}

func (r *TerraformRequiredProvidersDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	var err error
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_required_providers_declaration check since it's not hcl file")
		return nil
	}
	filename := body.Range().Filename
	if isOverrideTfFile(filename) {
		logger.Debug("skip terraform_required_version_declaration check since it's override file")
		return nil
	}
	blocks := body.Blocks
	for _, block := range blocks {
		switch block.Type {
		case "terraform":
			if subErr := r.checkBlock(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func (r *TerraformRequiredProvidersDeclarationRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	isRequiredProvidersDeclared := false
	var err error
	for _, nestedBlock := range block.Body.Blocks {
		switch nestedBlock.Type {
		case "required_providers":
			isRequiredProvidersDeclared = true
			err = multierror.Append(err, r.checkRequiredProvidersArgOrder(runner, nestedBlock))
		}
	}
	if isRequiredProvidersDeclared {
		return nil
	}
	return runner.EmitIssue(
		r,
		"The `required_providers` field should be declared in `terraform` block",
		block.DefRange(),
	)
}

func (r *TerraformRequiredProvidersDeclarationRule) checkRequiredProvidersArgOrder(runner tflint.Runner, providerBlock *hclsyntax.Block) error {
	file, _ := runner.GetFile(providerBlock.Range().Filename)
	var providerNames []string
	providerParamTxts := make(map[string]string)
	providerParamIssues := helper.Issues{}
	providers := providerBlock.Body.Attributes
	for _, config := range attributesByLines(providers) {
		sortedMap, sorted := PrintSortedAttrTxt(file.Bytes, config)
		name := config.Name
		providerParamTxts[name] = sortedMap
		providerNames = append(providerNames, name)
		if !sorted {
			providerParamIssues = append(providerParamIssues, &helper.Issue{
				Rule:    r,
				Message: fmt.Sprintf("Parameters of provider `%s` are expected to be sorted as follows:\n%s", name, sortedMap),
				Range:   config.NameRange,
			})
		}
	}
	sort.Slice(providerNames, func(x, y int) bool {
		providerX := providers[providerNames[x]]
		providerY := providers[providerNames[y]]
		if providerX.Range().Start.Line == providerY.Range().Start.Line {
			return providerX.Range().Start.Column < providerY.Range().Start.Column
		}
		return providerX.Range().Start.Line < providerY.Range().Start.Line
	})
	if !sort.StringsAreSorted(providerNames) {
		sort.Strings(providerNames)
		var sortedProviderParamTxts []string
		for _, providerName := range providerNames {
			sortedProviderParamTxts = append(sortedProviderParamTxts, providerParamTxts[providerName])
		}
		sortedProviderParamTxt := strings.Join(sortedProviderParamTxts, "\n")
		var sortedRequiredProviderTxt string
		if RemoveSpaceAndLine(sortedProviderParamTxt) == "" {
			sortedRequiredProviderTxt = fmt.Sprintf("%s {}", providerBlock.Type)
		} else {
			sortedRequiredProviderTxt = fmt.Sprintf("%s {\n%s\n}", providerBlock.Type, sortedProviderParamTxt)
		}
		sortedRequiredProviderTxt = string(hclwrite.Format([]byte(sortedRequiredProviderTxt)))
		return runner.EmitIssue(
			r,
			fmt.Sprintf("The arguments of `required_providers` are expected to be sorted as follows:\n%s", sortedRequiredProviderTxt),
			providerBlock.DefRange(),
		)
	}
	var err error
	for _, issue := range providerParamIssues {
		if subErr := runner.EmitIssue(issue.Rule, issue.Message, issue.Range); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
