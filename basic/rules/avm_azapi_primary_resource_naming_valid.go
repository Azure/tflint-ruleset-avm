package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(AzapiPrimaryResourceNamingRule)

// AzapiPrimaryResourceNamingRule checks that exactly one AzAPI managed resource, the primary resource, is named `this`.
type AzapiPrimaryResourceNamingRule struct {
	tflint.DefaultRule
}

// NewAzapiPrimaryResourceNamingRule returns a new rule.
func NewAzapiPrimaryResourceNamingRule() *AzapiPrimaryResourceNamingRule {
	return new(AzapiPrimaryResourceNamingRule)
}

func (r *AzapiPrimaryResourceNamingRule) Name() string {
	return "avm_azapi_primary_resource_naming_valid"
}

func (r *AzapiPrimaryResourceNamingRule) Enabled() bool {
	return true
}

func (r *AzapiPrimaryResourceNamingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AzapiPrimaryResourceNamingRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFRMNFR2/"
}

// Check decodes the module class when the rule is used without the ruleset wrapper.
func (r *AzapiPrimaryResourceNamingRule) Check(runner tflint.Runner) error {
	moduleClass, err := decodeModuleClass(runner, r.Name())
	if err != nil {
		return err
	}
	return r.CheckWithModuleClass(runner, moduleClass)
}

// CheckWithModuleClass checks resource module roots and skips pattern and utility modules.
func (r *AzapiPrimaryResourceNamingRule) CheckWithModuleClass(runner tflint.Runner, moduleClass string) error {
	return checkResourceModuleRoot(runner, r.Name(), moduleClass, r.checkModule)
}

func (r *AzapiPrimaryResourceNamingRule) checkModule(runner tflint.Runner) error {
	blocks, err := moduleBlocks(runner)
	if err != nil {
		return err
	}
	resources := azapiManagedResources(runner, blocks)
	if len(resources) == 0 {
		return nil
	}

	var named []*hclsyntax.Block
	for _, resource := range resources {
		if resource.Labels[1] == "this" {
			named = append(named, resource)
		}
	}

	if len(named) == 0 {
		return runner.EmitIssue(
			r,
			"The primary AzAPI resource must be named `this`, for example `azapi_resource.this`",
			resources[0].DefRange(),
		)
	}

	primary := named[0]
	for _, resource := range named {
		if resource.Labels[0] == "azapi_resource" {
			primary = resource
			break
		}
	}

	var errList error
	for _, resource := range named {
		if resource == primary {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("`%s.this` must not be named `this` because `%s.this` is the primary resource; name it after what it represents", resource.Labels[0], primary.Labels[0]),
			resource.DefRange(),
		)
		if subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}
	return errList
}
