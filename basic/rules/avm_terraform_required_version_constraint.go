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

	lower, upper := constraintBounds(constraints)
	if lower == nil {
		emit(fmt.Sprintf("`required_version` must include a minimum Terraform version, got `%s`; %s", raw, recommendedRequiredVersionFormat))
		if upper == nil {
			emit(fmt.Sprintf("`required_version` must include a maximum major Terraform version, got `%s`; %s", raw, recommendedRequiredVersionFormat))
		}
		return err
	}

	nextMajor := goversion.Must(goversion.NewVersion(fmt.Sprintf("%d.0.0", lower.Segments()[0]+1)))
	if !upper.below(nextMajor) {
		emit(fmt.Sprintf("`required_version` must exclude Terraform %s and later, got `%s`; %s", nextMajor, raw, recommendedRequiredVersionFormat))
	}
	return err
}

// upperBound is the tightest maximum version a constraint allows.
type upperBound struct {
	version   *goversion.Version
	inclusive bool
}

// below reports whether every version allowed by the bound is lower than limit. A nil bound has no maximum.
func (u *upperBound) below(limit *goversion.Version) bool {
	if u == nil {
		return false
	}
	if u.inclusive {
		return u.version.LessThan(limit)
	}
	return u.version.LessThanOrEqual(limit)
}

func (u *upperBound) tighten(version *goversion.Version, inclusive bool) *upperBound {
	if u == nil || version.LessThan(u.version) || (version.Equal(u.version) && !inclusive) {
		return &upperBound{version: version, inclusive: inclusive}
	}
	return u
}

// constraintBounds returns the highest lower bound and the tightest upper bound across all constraints. `!=` never bounds a range, so it is ignored.
func constraintBounds(constraints goversion.Constraints) (lower *goversion.Version, upper *upperBound) {
	raiseLower := func(v *goversion.Version) {
		if lower == nil || v.GreaterThan(lower) {
			lower = v
		}
	}
	for _, c := range constraints {
		op, rest := splitConstraintOperator(strings.TrimSpace(c.String()))
		ver, verErr := goversion.NewVersion(rest)
		if verErr != nil {
			continue
		}
		switch op {
		case "!=":
		case "<":
			upper = upper.tighten(ver, false)
		case "<=":
			upper = upper.tighten(ver, true)
		case ">", ">=":
			raiseLower(ver)
		case "~>":
			raiseLower(ver)
			// `~> X` with a single segment has no maximum. With more segments, the maximum is at most the next major.
			if versionSegmentCount(rest) > 1 {
				nextMajor := goversion.Must(goversion.NewVersion(fmt.Sprintf("%d.0.0", ver.Segments()[0]+1)))
				upper = upper.tighten(nextMajor, false)
			}
		default:
			raiseLower(ver)
			upper = upper.tighten(ver, true)
		}
	}
	return lower, upper
}

func splitConstraintOperator(s string) (op, rest string) {
	for _, candidate := range []string{"~>", ">=", "<=", "!=", ">", "<", "="} {
		if strings.HasPrefix(s, candidate) {
			return candidate, strings.TrimSpace(s[len(candidate):])
		}
	}
	return "", s
}

// versionSegmentCount counts the numeric segments written in a version, ignoring pre-release and metadata suffixes.
func versionSegmentCount(v string) int {
	core, _, _ := strings.Cut(v, "-")
	core, _, _ = strings.Cut(core, "+")
	return len(strings.Split(core, "."))
}
