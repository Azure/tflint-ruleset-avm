package waf_test

import (
	"os"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockFs(c string) afero.Afero {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "main.tf", []byte(c), os.ModePerm)
	return afero.Afero{Fs: fs}
}

func TestGetRules(t *testing.T) {
	rules := waf.GetRules()
	assert.Truef(t, len(rules) > 0, "rules should not be empty")
}
