package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"

	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformResourceDataArgLayoutRule struct {
	tflint.DefaultRule
}

// NewTerraformResourceDataArgLayoutRule returns a new rule
func NewTerraformResourceDataArgLayoutRule() *TerraformResourceDataArgLayoutRule {
	return &TerraformResourceDataArgLayoutRule{}
}

// Name returns the rule name
func (r *TerraformResourceDataArgLayoutRule) Name() string {
	return "terraform_resource_data_arg_layout"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformResourceDataArgLayoutRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformResourceDataArgLayoutRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformResourceDataArgLayoutRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the arguments/attributes in a block are sorted in azure doc Layout
func (r *TerraformResourceDataArgLayoutRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		subErr := r.visitFile(runner, file)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformResourceDataArgLayoutRule) visitFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_resource_data_arg_layout check since it's not hcl file")
		return nil
	}
	if body == nil {
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		switch block.Type {
		case "resource", "data":
			emitter := func(block Block) error {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("Arguments are expected to be arranged in following Layout:\n%s", block.ToString()),
					block.DefRange(),
				)
			}
			b := BuildResourceBlock(block, file, emitter)
			if subErr := b.CheckBlock(); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

//func (r *TerraformResourceDataArgLayoutRule) visitResourceDataBlock(runner tflint.Runner, block *hclsyntax.Block) error {
//	issue := new(helper.Issue)
//	r.visitBlock(runner, block, issue)
//	if !r.isIssueEmpty(issue) {
//		return runner.EmitIssue((*issue).Rule, (*issue).Message, (*issue).Range)
//	}
//	return nil
//}
//
//func (r *TerraformResourceDataArgLayoutRule) visitBlock(runner tflint.Runner, block *hclsyntax.Block, issue *helper.Issue) string {
//	file, _ := runner.GetFile(block.Range().Filename)
//	argGrps, isCorrectLayout := r.getSortedArgGrps(block)
//	var sortedArgHclTxts []string
//	isGapNeeded := false
//	for _, args := range argGrps {
//		if len(args) == 0 {
//			continue
//		}
//		if isGapNeeded {
//			sortedArgHclTxts = append(sortedArgHclTxts, "")
//		}
//		for _, arg := range args {
//			var argHclTxt string
//			if arg.Block == nil {
//				argHclTxt = string(arg.Range.SliceBytes(file.Bytes))
//			} else {
//				argHclTxt = r.visitBlock(runner, arg.Block, issue)
//			}
//			sortedArgHclTxts = append(sortedArgHclTxts, argHclTxt)
//		}
//		isGapNeeded = true
//	}
//	sortedBlockHclTxt := strings.Join(sortedArgHclTxts, "\n")
//	if strings.TrimSpace(sortedBlockHclTxt) == "" {
//		sortedBlockHclTxt = fmt.Sprintf("%s {}", r.getBlockHead(block))
//	} else {
//		sortedBlockHclTxt = fmt.Sprintf("%s {\n%s\n}", r.getBlockHead(block), sortedBlockHclTxt)
//	}
//	sortedBlockHclTxt = string(hclwrite.Format([]byte(sortedBlockHclTxt)))
//
//	if !isCorrectLayout {
//		issue.Rule = r
//		issue.Message = fmt.Sprintf("Arguments are expected to be arranged in following Layout:\n%s", sortedBlockHclTxt)
//		issue.Range = block.DefRange()
//	}
//	return sortedBlockHclTxt
//}
//
//func (r *TerraformResourceDataArgLayoutRule) getSortedArgGrps(block *hclsyntax.Block) ([][]Arg, bool) {
//	var headMetaArgs, tailMetaArgs, attrs, nestedBlocks []Arg
//	for attrName, attr := range block.Body.Attributes {
//		arg := Arg{
//			Name:  attrName,
//			Range: attr.Range(),
//		}
//		if IsHeadMeta(attrName) {
//			headMetaArgs = append(headMetaArgs, arg)
//		} else if IsTailMeta(attrName) {
//			tailMetaArgs = append(tailMetaArgs, arg)
//		} else {
//			attrs = append(attrs, arg)
//		}
//	}
//	for _, nestedBlock := range block.Body.Blocks {
//		blockName := r.getBlockHead(nestedBlock)
//		arg := Arg{
//			Name:  blockName,
//			Range: nestedBlock.Range(),
//			Block: nestedBlock,
//		}
//		if IsHeadMeta(blockName) {
//			headMetaArgs = append(headMetaArgs, arg)
//		} else if IsTailMeta(blockName) {
//			tailMetaArgs = append(tailMetaArgs, arg)
//		} else {
//			nestedBlocks = append(nestedBlocks, arg)
//		}
//	}
//	sort.Slice(headMetaArgs, func(i, j int) bool {
//		return GetHeadMetaPriority(headMetaArgs[i].Name) > GetHeadMetaPriority(headMetaArgs[j].Name)
//	})
//	sort.Slice(tailMetaArgs, func(i, j int) bool {
//		return GetTailMetaPriority(tailMetaArgs[i].Name) > GetTailMetaPriority(tailMetaArgs[j].Name)
//	})
//	attrs = GetArgsWithOriginalOrder(attrs)
//	nestedBlocks = GetArgsWithOriginalOrder(nestedBlocks)
//
//	argGrps := [][]Arg{headMetaArgs, attrs, nestedBlocks, tailMetaArgs}
//	var lastArgGrp, sortedArgs []Arg
//	isCorrectLayout := true
//	for _, argGrp := range argGrps {
//		if isCorrectLayout && len(lastArgGrp) > 0 && len(argGrp) > 0 {
//			if argGrp[0].Range.Start.Line-lastArgGrp[len(lastArgGrp)-1].Range.End.Line < 2 {
//				isCorrectLayout = false
//			}
//		}
//		if len(argGrp) > 0 {
//			lastArgGrp = argGrp
//		}
//		sortedArgs = append(sortedArgs, argGrp...)
//	}
//	isCorrectLayout = isCorrectLayout && reflect.DeepEqual(sortedArgs, GetArgsWithOriginalOrder(sortedArgs))
//	return argGrps, isCorrectLayout
//}
//
//func (r *TerraformResourceDataArgLayoutRule) getBlockHead(block *hclsyntax.Block) string {
//	heads := []string{block.Type}
//	for _, label := range block.Labels {
//		heads = append(heads, fmt.Sprintf("\"%s\"", label))
//	}
//	return strings.Join(heads, " ")
//}
//
//func (r *TerraformResourceDataArgLayoutRule) isIssueEmpty(issue *helper.Issue) bool {
//	return *issue == helper.Issue{}
//}
