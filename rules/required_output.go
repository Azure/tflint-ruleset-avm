package rules

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Fs = &fsAdapter{
	Fs: afero.NewOsFs(),
}
var _ tflint.Rule = &RequiredOutputRule{}

type RequiredOutputRule struct {
	tflint.DefaultRule
	ruleName           string
	requiredOutputName string
	link               string
}

func NewRequiredOutputResourceRule() *RequiredOutputRule {
	return NewRequiredOutputRule("tffr2", "resource", "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr2---category-outputs---additional-terraform-outputs")
}

func NewRequiredOutputResourceIdRule() *RequiredOutputRule {
	return NewRequiredOutputRule("rmfr7", "resource_id", "https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs")
}

func NewRequiredOutputRule(ruleName, requiredOutputName, link string) *RequiredOutputRule {
	return &RequiredOutputRule{
		ruleName:           ruleName,
		requiredOutputName: requiredOutputName,
		link:               link,
	}
}

func (m *RequiredOutputRule) Name() string {
	return m.ruleName
}

func (m *RequiredOutputRule) Link() string {
	return m.link
}

func (m *RequiredOutputRule) Enabled() bool {
	return false
}

func (m *RequiredOutputRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (m *RequiredOutputRule) Check(runner tflint.Runner) error {
	wd, _ := runner.GetOriginalwd()
	module, diag := tfconfig.LoadModuleFromFilesystem(Fs, wd)
	if diag.HasErrors() {
		return diag
	}
	_, ok := module.Outputs[m.requiredOutputName]
	if !ok {
		return runner.EmitIssue(m, fmt.Sprintf("module owners MUST output the `%s` in their modules", m.requiredOutputName), hcl.Range{
			Filename: "outputs.tf",
		})
	}
	return nil
}

var _ tfconfig.FS = &fsAdapter{}

type fsAdapter struct {
	afero.Fs
}

func (b *fsAdapter) Open(name string) (tfconfig.File, error) {
	f, err := b.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (b *fsAdapter) ReadFile(name string) ([]byte, error) {
	return afero.ReadFile(b.Fs, name)
}

func (b *fsAdapter) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(b.Fs, dirname)
}
