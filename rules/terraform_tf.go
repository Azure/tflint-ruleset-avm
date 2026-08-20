package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

var _ tflint.Rule = new(TerraformTfFileRule)

// TerraformTfFileRule checks that terraform.tf contains exactly one Terraform block.
type TerraformTfFileRule struct {
	tflint.DefaultRule
}

// NewTerraformTfFileRule returns a Terraform file-layout rule.
func NewTerraformTfFileRule() *TerraformTfFileRule {
	return new(TerraformTfFileRule)
}

func (r *TerraformTfFileRule) Name() string {
	return "terraform_tf_file"
}

func (r *TerraformTfFileRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR39/"
}

func (r *TerraformTfFileRule) Enabled() bool {
	return true
}

func (r *TerraformTfFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformTfFileRule) Check(runner tflint.Runner) error {
	tFile, err := runner.GetFile("terraform.tf")
	if err != nil {
		if strings.Contains(err.Error(), "file not found") {
			return runner.EmitIssue(r, "AVM Terraform modules must contain a `terraform.tf` file", hcl.Range{})
		}
		return err
	}
	if tFile == nil {
		return runner.EmitIssue(r, "AVM Terraform modules must contain a `terraform.tf` file", hcl.Range{})
	}
	body, ok := tFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	if len(body.Attributes) != 0 || len(body.Blocks) != 1 || body.Blocks[0].Type != "terraform" {
		return runner.EmitIssue(
			r,
			"`terraform.tf` must contain exactly one `terraform` block and no other top-level content",
			body.Range(),
		)
	}
	return nil
}
