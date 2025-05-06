package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(RequiredProvidersRule)

type RequiredProvidersRule struct {
	tflint.DefaultRule
}

func NewRequiredProvidersRule() *RequiredProvidersRule {
	return &RequiredProvidersRule{}
}

func (t *RequiredProvidersRule) Name() string {
	return "required_providers_tfnfr26"
}

func (t *RequiredProvidersRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr26---category-code-style---providers-in-required_providers"
}

func (t *RequiredProvidersRule) Enabled() bool {
	return true
}

func (t *RequiredProvidersRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *RequiredProvidersRule) Check(r tflint.Runner) error {
	tFile, err := r.GetFile("terraform.tf")
	if err != nil {
		return err
	}

	body, ok := tFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	var errList error
	for _, block := range body.Blocks {
		if block.Type != "terraform" {
			continue
		}

		if subErr := t.checkBlock(r, block); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *RequiredProvidersRule) checkBlock(r tflint.Runner, block *hclsyntax.Block) error {
	isRequiredProvidersDeclared := false
	var errList error

	for _, nestedBlock := range block.Body.Blocks {
		if nestedBlock.Type != "required_providers" {
			continue
		}

		isRequiredProvidersDeclared = true
		errList = multierror.Append(errList, t.checkRequiredProvidersArgOrder(r, nestedBlock))
		errList = multierror.Append(errList, t.checkRequiredProvidersVersion(r, nestedBlock))
	}

	if isRequiredProvidersDeclared {
		return nil
	}

	return r.EmitIssue(
		t,
		"The `required_providers` field should be declared in `terraform` block",
		block.DefRange(),
	)
}

func (t *RequiredProvidersRule) checkRequiredProvidersArgOrder(r tflint.Runner, providerBlock *hclsyntax.Block) error {
	file, _ := r.GetFile(providerBlock.Range().Filename)
	var providerNames []string
	providerParamTxts := make(map[string]string)
	providerParamIssues := helper.Issues{}
	providers := providerBlock.Body.Attributes

	for _, config := range attributesByLines(providers) {
		sortedMap, sorted := printSortedAttrTxt(file.Bytes, config)
		name := config.Name
		providerParamTxts[name] = sortedMap
		providerNames = append(providerNames, name)

		if !sorted {
			providerParamIssues = append(providerParamIssues, &helper.Issue{
				Rule:    t,
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

		return r.EmitIssue(
			t,
			fmt.Sprintf("The arguments of `required_providers` are expected to be sorted as follows:\n%s", sortedRequiredProviderTxt),
			providerBlock.DefRange(),
		)
	}

	var errList error
	for _, issue := range providerParamIssues {
		if subErr := r.EmitIssue(issue.Rule, issue.Message, issue.Range); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func attributesByLines(attributes hclsyntax.Attributes) []*hclsyntax.Attribute {
	var attrs []*hclsyntax.Attribute
	for _, attr := range attributes {
		attrs = append(attrs, attr)
	}

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Range().Start.Line < attrs[j].Range().Start.Line
	})

	return attrs
}

func RemoveSpaceAndLine(str string) string {
	newStr := strings.ReplaceAll(str, " ", "")
	newStr = strings.ReplaceAll(newStr, "\t", "")
	newStr = strings.ReplaceAll(newStr, "\n", "")

	return newStr
}

func printSortedAttrTxt(src []byte, attr *hclsyntax.Attribute) (string, bool) {
	isSorted := true
	exp, isMap := attr.Expr.(*hclsyntax.ObjectConsExpr)
	if !isMap {
		return string(attr.Range().SliceBytes(src)), isSorted
	}

	var keys []string
	object := make(map[string]string)
	for _, item := range exp.Items {
		key := string(item.KeyExpr.Range().SliceBytes(src))
		value := fmt.Sprintf("%s = %s", key, string(item.ValueExpr.Range().SliceBytes(src)))
		keys = append(keys, key)
		object[key] = value
	}

	isSorted = sort.StringsAreSorted(keys)
	if !isSorted {
		sort.Strings(keys)
	}

	var objectAttrs []string
	for _, key := range keys {
		objectAttrs = append(objectAttrs, object[key])
	}

	sortedExpTxt := strings.Join(objectAttrs, "\n")
	var sortedAttrTxt string
	if RemoveSpaceAndLine(sortedExpTxt) == "" {
		sortedAttrTxt = fmt.Sprintf("%s = {}", attr.Name)
	} else {
		sortedAttrTxt = fmt.Sprintf("%s = {\n%s\n}", attr.Name, sortedExpTxt)
	}
	formattedTxt := string(hclwrite.Format([]byte(sortedAttrTxt)))

	return formattedTxt, isSorted
}

func (t *RequiredProvidersRule) checkRequiredProvidersVersion(r tflint.Runner, providerBlock *hclsyntax.Block) error {
	var errList error
	file, _ := r.GetFile(providerBlock.Range().Filename)

	for _, v := range providerBlock.Body.Attributes {
		if provider, ok := v.Expr.(*hclsyntax.ObjectConsExpr); ok {
			for _, item := range provider.Items {
				attrType := string(item.KeyExpr.Range().SliceBytes(file.Bytes))
				if attrType != "version" {
					continue
				}

				attrVal := string(item.ValueExpr.Range().SliceBytes(file.Bytes))
				if !strings.Contains(attrVal, "~>") && !(strings.Contains(attrVal, ">") && strings.Contains(attrVal, "<")) {
					errList = multierror.Append(errList, r.EmitIssue(
						t,
						"The `version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
						provider.Range(),
					))
				}
			}
		} else if provider, ok := v.Expr.(*hclsyntax.TemplateExpr); ok {
			versionVal, diags := provider.Value(nil)
			if diags.HasErrors() {
				errList = multierror.Append(errList, r.EmitIssue(
					t,
					diags.Error(),
					provider.Range(),
				))
			}

			version := versionVal.AsString()
			if !strings.Contains(version, "~>") && !(strings.Contains(version, ">") && strings.Contains(version, "<")) {
				errList = multierror.Append(errList, r.EmitIssue(
					t,
					"The provider version constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
					provider.Range(),
				))
			}
		} else {
			errList = multierror.Append(errList, r.EmitIssue(
				t,
				"The provider only supports string type and block type",
				provider.Range(),
			))
		}
	}

	return errList
}
