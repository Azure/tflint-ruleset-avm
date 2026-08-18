package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"sort"
	"strings"
)

var headMetaArgPriority = map[string]int{"for_each": 0, "count": 0, "provider": 1}
var tailMetaArgPriority = map[string]int{"lifecycle": 0, "depends_on": 1}

// IsHeadMeta checks whether a name represents a type of head Meta arg
func IsHeadMeta(argName string) bool {
	_, isHeadMeta := headMetaArgPriority[argName]
	return isHeadMeta
}

// IsTailMeta checks whether a name represents a type of tail Meta arg
func IsTailMeta(argName string) bool {
	_, isTailMeta := tailMetaArgPriority[argName]
	return isTailMeta
}

func ref(hr hcl.Range) *hcl.Range {
	return &hr
}

func ForFiles(runner tflint.Runner, action func(tflint.Runner, *hcl.File) error) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := action(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

// PrintSortedAttrTxt print the sorted hcl text of an attribute
func PrintSortedAttrTxt(src []byte, attr *hclsyntax.Attribute) (string, bool) {
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

// RemoveSpaceAndLine remove space, "\t" and "\n" from the given string
func RemoveSpaceAndLine(str string) string {
	newStr := strings.ReplaceAll(str, " ", "")
	newStr = strings.ReplaceAll(newStr, "\t", "")
	newStr = strings.ReplaceAll(newStr, "\n", "")
	return newStr
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
