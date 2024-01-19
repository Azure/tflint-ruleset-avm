package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformDotTfRule)

type TerraformDotTfRule struct {
	tflint.DefaultRule
}

func NewTerraformDotTfRule() *TerraformDotTfRule {
	return new(TerraformDotTfRule)
}

func (t *TerraformDotTfRule) Name() string {
	return "tfnfr26"
}

func (t *TerraformDotTfRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr26---category-code-style---providers-must-be-declared-in-the-required_providers-block-in-terraformtf-and-must-have-a-constraint-on-minimum-and-maximum-major-version"
}

func (t *TerraformDotTfRule) Enabled() bool {
	return false
}

func (t *TerraformDotTfRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *TerraformDotTfRule) Check(r tflint.Runner) error {
	tFile, err := r.GetFile("terraform.tf")
	if err != nil {
		return err
	}
	if tFile == nil {
		return r.EmitIssue(t, "All avm Terraform modules must contain `terraform.tf` file", hcl.Range{})
	}
	body := tFile.Body.(*hclsyntax.Body)
	terraformBlockFound := false
	for _, b := range body.Blocks {
		if b.Type == "terraform" {
			terraformBlockFound = true
		} else {
			return r.EmitIssue(t, "`terraform.tf` file must contain `terraform` block only", body.Range())
		}
	}
	if !terraformBlockFound {
		err := r.EmitIssue(t, "`terraform.tf` file must contain `terraform` block only", body.Range())
		if err != nil {
			return err
		}
	}
	return nil
}
