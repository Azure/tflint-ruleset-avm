package rules_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const rootModuleSourceMessage = "The `source` property must reference an AVM registry module (`Azure/` containing `avm-`) or a direct local module (`./modules/<name>`)"

const siblingModuleSourceMessage = "The `source` property must reference an AVM registry module (`Azure/` containing `avm-`) or a direct local module (`../<name>`)"

func TestModuleSource(t *testing.T) {
	cases := []struct {
		desc   string
		config string
		issues helper.Issues
	}{
		{
			desc: "source exists, ok",
			config: `module "other-module" {
  source  = "Azure/avm-res-keyvault-vault/azurerm"
  version = "0.5.3"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "no version, ok",
			config: `module "other-module" {
  source  = "Azure/avm-res-keyvault-vault/azurerm"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "no source, not ok",
			config: `module "other-module" {
}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewModuleSourceRule(),
					Message: "The `source` field should be declared in the `module` block",
				},
			},
		},
		{
			desc: "git reference, not ok",
			config: `module "other-module" {
  source  = "git::https://Azure/terraform-azurerm-avm-res-keyvault-vault.git"
  version = "0.5.3"
}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewModuleSourceRule(),
					Message: rootModuleSourceMessage,
				},
			},
		},
		{
			desc: "github reference, not ok",
			config: `module "other-module" {
  source  = "github.com/Azure/terraform-azurerm-avm-res-keyvault-vault"
  version = "0.5.3"
}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewModuleSourceRule(),
					Message: rootModuleSourceMessage,
				},
			},
		},
		{
			desc: "source sub module, ok",
			config: `module "other-module" {
  source  = "./modules/my-sub-module"
  version = "0.5.3"
}`,
			issues: helper.Issues{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := rules.NewModuleSourceRule()
			filename := "terraform.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

func TestModuleSourceScopes(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		rootAllowed      bool
		submoduleAllowed bool
	}{
		{name: "resource registry", source: "Azure/avm-res-keyvault-vault/azurerm", rootAllowed: true, submoduleAllowed: true},
		{name: "pattern registry", source: "Azure/avm-ptn-alz/azure", rootAllowed: true, submoduleAllowed: true},
		{name: "utility registry", source: "Azure/avm-utl-regions/azurerm", rootAllowed: true, submoduleAllowed: true},
		{name: "external registry submodule", source: "Azure/avm-res-network-virtualnetwork/azurerm//modules/subnet", rootAllowed: true, submoduleAllowed: true},
		{name: "direct local submodule", source: "./modules/route-table", rootAllowed: true},
		{name: "direct local submodule trailing slash", source: "./modules/route-table/", rootAllowed: true},
		{name: "route table sibling", source: "../route-table-route", submoduleAllowed: true},
		{name: "network security group sibling", source: "../network-security-group-security-rule", submoduleAllowed: true},
		{name: "sibling trailing slash", source: "../part/", submoduleAllowed: true},
		{name: "sibling with dot in name", source: "../part.v2", submoduleAllowed: true},
		{name: "sibling with underscore", source: "../sub_part", submoduleAllowed: true},
		{name: "empty", source: ""},
		{name: "current root", source: "."},
		{name: "current root slash", source: "./"},
		{name: "parent root", source: ".."},
		{name: "parent root slash", source: "../"},
		{name: "modules directory", source: "./modules"},
		{name: "modules directory slash", source: "./modules/"},
		{name: "other local directory", source: "./part"},
		{name: "example directory", source: "./examples/part"},
		{name: "grandparent", source: "../../x"},
		{name: "nested sibling", source: "../a/b"},
		{name: "sibling dot", source: "../."},
		{name: "sibling parent", source: "../.."},
		{name: "sibling dot slash", source: ".././"},
		{name: "sibling parent slash", source: "../../"},
		{name: "sibling normalized traversal", source: "../a/../b"},
		{name: "sibling doubled separator", source: "..//part"},
		{name: "sibling repeated trailing slash", source: "../part//"},
		{name: "nested local root", source: "./modules/a/modules/b"},
		{name: "nested local directory", source: "./modules/a/b"},
		{name: "local dot", source: "./modules/."},
		{name: "local parent", source: "./modules/.."},
		{name: "local traversal", source: "./modules/../part"},
		{name: "local normalized traversal", source: "./modules/a/../b"},
		{name: "local doubled separator", source: "./modules//part"},
		{name: "local repeated trailing slash", source: "./modules/part//"},
		{name: "Windows sibling", source: `..\part`},
		{name: "Windows grandparent", source: `..\..\part`},
		{name: "mixed sibling separators", source: `../part\child`},
		{name: "mixed sibling traversal", source: `../..\part`},
		{name: "Windows local module", source: `.\modules\part`},
		{name: "mixed local separators", source: `./modules/part\child`},
		{name: "mixed local traversal", source: `./modules/..\part`},
		{name: "drive relative sibling", source: "../C:part"},
		{name: "drive relative local", source: "./modules/C:part"},
		{name: "Windows sibling dot normalization", source: "../..."},
		{name: "Windows sibling parent normalization", source: "../.. "},
		{name: "Windows sibling dot space normalization", source: "../. "},
		{name: "Windows local parent normalization", source: "./modules/.. "},
		{name: "sibling trailing dot", source: "../part."},
		{name: "sibling trailing space", source: "../part "},
		{name: "local trailing dot", source: "./modules/part."},
		{name: "sibling null byte", source: "../part\x00"},
		{name: "local null byte", source: "./modules/part\x00"},
		{name: "absolute Unix path", source: "/modules/part"},
		{name: "absolute Windows path", source: `C:\modules\part`},
		{name: "git URL", source: "git::https://github.com/Azure/terraform-azure-avm-res-example-widget.git//modules/part"},
		{name: "GitHub URL", source: "github.com/Azure/terraform-azure-avm-res-example-widget"},
		{name: "non AVM registry module", source: "Azure/compute/azurerm"},
		{name: "non Azure registry module", source: "other/avm-res-example-widget/azure"},
	}

	for _, scope := range []string{"", "root", "submodule"} {
		scopeName := scope
		if scopeName == "" {
			scopeName = "default"
		}
		for _, tc := range tests {
			t.Run(scopeName+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				source, err := json.Marshal(tc.source)
				require.NoError(t, err)
				files := map[string]string{
					"main.tf": fmt.Sprintf("module \"part\" {\n  source = %s\n}\n", source),
				}
				if scope != "" {
					files[".tflint.hcl"] = fmt.Sprintf(`rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = %q
}`, scope)
				}
				runner := helper.TestRunner(t, files)
				rule := common.NewConfigurableRule(rules.NewModuleSourceRule())

				require.NoError(t, rule.Check(runner))

				allowed := tc.rootAllowed
				message := rootModuleSourceMessage
				if scope == "submodule" {
					allowed = tc.submoduleAllowed
					message = siblingModuleSourceMessage
				}
				if allowed {
					assert.Empty(t, runner.Issues)
					return
				}
				require.Len(t, runner.Issues, 1)
				assert.Equal(t, rule.Name(), runner.Issues[0].Rule.Name())
				assert.Equal(t, tflint.ERROR, runner.Issues[0].Rule.Severity())
				assert.Equal(t, message, runner.Issues[0].Message)
				assert.Equal(t, hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 18},
					End:      hcl.Pos{Line: 2, Column: 9, Byte: 24},
				}, runner.Issues[0].Range)
			})
		}
	}
}

type scopedSourceRunner struct {
	tflint.Runner
	modulePath addrs.Module
}

func (r scopedSourceRunner) GetModulePath() (addrs.Module, error) {
	return r.modulePath, nil
}

func TestModuleSourceSkipsChildModules(t *testing.T) {
	for _, scope := range []string{"root", "submodule"} {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{
				"main.tf": `module "part" {
  source = "../../part"
}`,
				".tflint.hcl": fmt.Sprintf(`rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = %q
}`, scope),
			})
			rule := common.NewConfigurableRule(rules.NewModuleSourceRule())

			require.NoError(t, rule.Check(scopedSourceRunner{Runner: runner, modulePath: addrs.Module{"child"}}))
			assert.Empty(t, runner.Issues)
		})
	}
}

func TestModuleSourceRejectsInvalidModuleScope(t *testing.T) {
	for _, expression := range []string{`""`, `"Root"`, `"Submodule"`, `"submodule "`, `"example"`, `"other"`, "null", "true", "1", "[]", "{}"} {
		for _, modulePath := range []addrs.Module{nil, {"child"}} {
			t.Run(expression+"/"+modulePath.String(), func(t *testing.T) {
				t.Parallel()
				runner := helper.TestRunner(t, map[string]string{
					".tflint.hcl": fmt.Sprintf(`rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = %s
}`, expression),
				})
				rule := common.NewConfigurableRule(rules.NewModuleSourceRule())

				err := rule.Check(scopedSourceRunner{Runner: runner, modulePath: modulePath})

				require.Error(t, err)
				if expression[0] == '"' {
					assert.ErrorContains(t, err, `module_scope must be one of "root" or "submodule"`)
				}
				assert.Empty(t, runner.Issues)
			})
		}
	}
}

func TestModuleSourceSeverity(t *testing.T) {
	for _, severity := range []string{"warning", "notice", "critical"} {
		t.Run(severity, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{
				"main.tf": `module "part" {
  source = "../a/b"
}`,
				".tflint.hcl": fmt.Sprintf(`rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = "submodule"
  severity     = %q
}`, severity),
			})
			rule := common.NewConfigurableRule(rules.NewModuleSourceRule())

			err := rule.Check(runner)

			if severity == "critical" {
				require.ErrorContains(t, err, `severity must be one of "error", "warning", or "notice"`)
				assert.Empty(t, runner.Issues)
				return
			}
			require.NoError(t, err)
			require.Len(t, runner.Issues, 1)
			wantSeverity := tflint.WARNING
			if severity == "notice" {
				wantSeverity = tflint.NOTICE
			}
			assert.Equal(t, wantSeverity, runner.Issues[0].Rule.Severity())
			assert.Equal(t, siblingModuleSourceMessage, runner.Issues[0].Message)
		})
	}
}

func TestModuleSourceCheckWithoutWrapper(t *testing.T) {
	for _, scope := range []string{"root", "submodule", "invalid"} {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{
				"main.tf": `module "part" {
  source = "../part"
}`,
				".tflint.hcl": fmt.Sprintf(`rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = %q
}`, scope),
			})

			err := rules.NewModuleSourceRule().Check(runner)

			if scope == "invalid" {
				require.ErrorContains(t, err, `module_scope must be one of "root" or "submodule"`)
				assert.Empty(t, runner.Issues)
				return
			}
			require.NoError(t, err)
			if scope == "submodule" {
				assert.Empty(t, runner.Issues)
				return
			}
			require.Len(t, runner.Issues, 1)
			assert.Equal(t, rootModuleSourceMessage, runner.Issues[0].Message)
		})
	}
}
