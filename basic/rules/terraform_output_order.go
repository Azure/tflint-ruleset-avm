package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"sort"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformOutputOrderRule checks whether the outputs are sorted in expected order
type TerraformOutputOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformOutputOrderRule returns a new rule
func NewTerraformOutputOrderRule() *TerraformOutputOrderRule {
	return &TerraformOutputOrderRule{}
}

// Name returns the rule name
func (r *TerraformOutputOrderRule) Name() string {
	return "terraform_output_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformOutputOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformOutputOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformOutputOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the outputs are sorted in expected order
func (r *TerraformOutputOrderRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkOutputOrder(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformOutputOrderRule) checkOutputOrder(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_output_order check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	firstOutputBlockRange := r.firstOutputRange(blocks)
	if firstOutputBlockRange == nil {
		return nil
	}
	if r.sorted(blocks) {
		return nil
	}
	return r.suggestedOrder(runner, file, blocks, firstOutputBlockRange)
}

func (r *TerraformOutputOrderRule) suggestedOrder(runner tflint.Runner, file *hcl.File, blocks hclsyntax.Blocks, firstOutputBlockRange *hcl.Range) error {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Labels[0] < blocks[j].Labels[0]
	})
	var sortedOutputHclTxts []string
	for _, b := range blocks {
		sortedOutputHclTxts = append(sortedOutputHclTxts, string(b.Range().SliceBytes(file.Bytes)))
	}
	sortedOutputHclBytes := hclwrite.Format([]byte(strings.Join(sortedOutputHclTxts, "\n\n")))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended output order:\n%s", sortedOutputHclBytes),
		*firstOutputBlockRange,
	)
}

func (r *TerraformOutputOrderRule) sorted(blocks hclsyntax.Blocks) bool {
	var outputNames []string
	for _, block := range blocks {
		switch block.Type {
		case "output":
			outputNames = append(outputNames, block.Labels[0])
		}
	}

	return sort.StringsAreSorted(outputNames)
}

func (r *TerraformOutputOrderRule) firstOutputRange(blocks hclsyntax.Blocks) *hcl.Range {
	for _, b := range blocks {
		switch b.Type {
		case "output":
			return ref(b.DefRange())
		}
	}
	return nil
}
