package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ tflint.Rule = new(TerraformRequiredVersionRule)

const recommendedRequiredVersionFormat = "use `~> #.#` or `>= #.#.#, < #.#.#`"

// noUpperBoundProbe is a Terraform version far enough in the future that any upper bound excludes it.
var noUpperBoundProbe = goversion.Must(goversion.NewVersion("999999.0.0"))

// TerraformRequiredVersionRule checks that `required_version` in `terraform.tf` bounds both the minimum version and the maximum major version.
type TerraformRequiredVersionRule struct {
	tflint.DefaultRule
}

// NewTerraformRequiredVersionRule returns a new rule.
func NewTerraformRequiredVersionRule() *TerraformRequiredVersionRule {
	return new(TerraformRequiredVersionRule)
}

func (r *TerraformRequiredVersionRule) Name() string {
	return "avm_terraform_required_version_constraint"
}

func (r *TerraformRequiredVersionRule) Enabled() bool {
	return true
}

func (r *TerraformRequiredVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformRequiredVersionRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR25/"
}

func (r *TerraformRequiredVersionRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformRequiredVersionRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_required_version_constraint since it is not an HCL file")
		return nil
	}
	if fileBaseName(body) != "terraform.tf" {
		return nil
	}
	// The single `terraform` block requirement is enforced by avm_terraform_configuration_file_required.
	var first *hclsyntax.Block
	for _, block := range body.Blocks {
		if block.Type != "terraform" {
			continue
		}
		if first == nil {
			first = block
		}
		if attr, ok := block.Body.Attributes["required_version"]; ok {
			return r.checkConstraint(runner, attr)
		}
	}
	if first == nil {
		return nil
	}
	return runner.EmitIssue(r, "The `terraform` block in `terraform.tf` must define `required_version`", first.DefRange())
}

func (r *TerraformRequiredVersionRule) checkConstraint(runner tflint.Runner, attr *hclsyntax.Attribute) error {
	v, ok := staticValue(attr, cty.String)
	if !ok {
		return nil
	}
	raw := v.AsString()
	constraints, parseErr := goversion.NewConstraint(raw)
	if parseErr != nil {
		return runner.EmitIssue(r, fmt.Sprintf("`required_version` must be a valid version constraint, got `%s`", raw), attr.SrcRange)
	}

	var err error
	emit := func(msg string) {
		if subErr := runner.EmitIssue(r, msg, attr.SrcRange); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}

	lower := lowerBound(constraints)
	if lower == nil {
		emit(fmt.Sprintf("`required_version` must include a minimum Terraform version, got `%s`; %s", raw, recommendedRequiredVersionFormat))
		if constraints.Check(noUpperBoundProbe) {
			emit(fmt.Sprintf("`required_version` must include a maximum major Terraform version, got `%s`; %s", raw, recommendedRequiredVersionFormat))
		}
		return err
	}

	nextMajor := goversion.Must(goversion.NewVersion(fmt.Sprintf("%d.0.0", lower.Segments()[0]+1)))
	if constraints.Check(nextMajor) {
		emit(fmt.Sprintf("`required_version` must exclude Terraform %s and later, got `%s`; %s", nextMajor, raw, recommendedRequiredVersionFormat))
	}
	return err
}

// lowerBound returns the highest version used by a lower-bound operator (`>=`, `>`, `~>`, `=`, or none), or nil if there is none.
func lowerBound(constraints goversion.Constraints) *goversion.Version {
	var lower *goversion.Version
	for _, c := range constraints {
		s := strings.TrimSpace(c.String())
		var rest string
		switch {
		case strings.HasPrefix(s, "~>"):
			rest = s[2:]
		case strings.HasPrefix(s, ">="):
			rest = s[2:]
		case strings.HasPrefix(s, "!="), strings.HasPrefix(s, "<"):
			continue
		case strings.HasPrefix(s, ">"), strings.HasPrefix(s, "="):
			rest = s[1:]
		default:
			rest = s
		}
		ver, verErr := goversion.NewVersion(strings.TrimSpace(rest))
		if verErr != nil {
			continue
		}
		if lower == nil || ver.GreaterThan(lower) {
			lower = ver
		}
	}
	return lower
}
