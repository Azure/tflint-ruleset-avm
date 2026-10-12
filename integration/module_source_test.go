package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModuleSourceIntegration(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		scope         string
		submodule     bool
		useChdir      bool
		severity      string
		wantSeverity  string
		expectedError string
	}{
		{name: "root direct submodule", source: "./modules/route-table-route"},
		{name: "root sibling rejected", source: "../route-table-route", wantSeverity: "error"},
		{name: "explicit root sibling rejected", source: "../route-table-route", scope: "root", wantSeverity: "error"},
		{name: "submodule defaults to root", source: "../route-table-route", submodule: true, wantSeverity: "error"},
		{name: "route table sibling", source: "../route-table-route", scope: "submodule", submodule: true},
		{name: "security group sibling with chdir", source: "../network-security-group-security-rule", scope: "submodule", submodule: true, useChdir: true},
		{name: "sibling trailing slash", source: "../route-table-route/", scope: "submodule", submodule: true},
		{name: "nested local root rejected", source: "./modules/nested", scope: "submodule", submodule: true, wantSeverity: "error"},
		{name: "nested sibling rejected", source: "../route-table-route/nested", scope: "submodule", submodule: true, wantSeverity: "error"},
		{name: "grandparent rejected", source: "../../peer", scope: "submodule", submodule: true, wantSeverity: "error"},
		{name: "severity override", source: "../route-table-route/nested", scope: "submodule", submodule: true, severity: "notice", wantSeverity: "info"},
		{name: "invalid scope", source: "./modules/route-table-route", scope: "invalid", expectedError: `module_scope must be one of "root" or "submodule", got "invalid"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixtureDir := t.TempDir()
			rootDir := filepath.Join(fixtureDir, "repo")
			targetDir := rootDir
			if tc.submodule {
				targetDir = filepath.Join(rootDir, "modules", "caller")
			}
			for _, dir := range []string{
				targetDir,
				filepath.Join(fixtureDir, "route-table-route"),
				filepath.Join(rootDir, "modules", "route-table-route", "nested"),
				filepath.Join(rootDir, "modules", "network-security-group-security-rule"),
				filepath.Join(rootDir, "modules", "caller", "modules", "nested"),
				filepath.Join(rootDir, "peer"),
			} {
				require.NoError(t, os.MkdirAll(dir, 0755))
			}
			require.NoError(t, os.WriteFile(
				filepath.Join(targetDir, "main.tf"),
				fmt.Appendf(nil, "module \"part\" {\n  source = %q\n}\n", tc.source),
				0600,
			))

			config := `plugin "terraform" {
  enabled = false
}
plugin "avm" {
  enabled = true
}
rule "avm_terraform_module_source_required" {
  enabled = true
`
			if tc.scope != "" {
				config += fmt.Sprintf("  module_scope = %q\n", tc.scope)
			}
			if tc.severity != "" {
				config += fmt.Sprintf("  severity = %q\n", tc.severity)
			}
			config += "}\n"
			require.NoError(t, os.WriteFile(filepath.Join(targetDir, ".tflint.hcl"), []byte(config), 0600))

			args := []string{"--format=json", "--only=avm_terraform_module_source_required", "--minimum-failure-severity=warning"}
			workingDir := targetDir
			if tc.useChdir {
				workingDir = rootDir
				args = append([]string{"--chdir=" + targetDir}, args...)
			}
			cmd := exec.Command("tflint", args...)
			cmd.Dir = workingDir
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			if tc.wantSeverity == "error" || tc.expectedError != "" {
				var exitError *exec.ExitError
				require.ErrorAs(t, err, &exitError, "stdout=%s stderr=%s", stdout.String(), stderr.String())
			} else {
				require.NoError(t, err, "stdout=%s stderr=%s", stdout.String(), stderr.String())
			}
			var output tflintOutput
			require.NoError(t, json.Unmarshal(stdout.Bytes(), &output), "stdout=%s stderr=%s", stdout.String(), stderr.String())
			if tc.expectedError != "" {
				require.Len(t, output.Errors, 1)
				assert.Contains(t, output.Errors[0].Message, tc.expectedError)
				assert.Empty(t, output.Issues)
				return
			}
			assert.Empty(t, output.Errors)
			if tc.wantSeverity == "" {
				assert.Empty(t, output.Issues)
				return
			}
			require.Len(t, output.Issues, 1)
			assert.Equal(t, "avm_terraform_module_source_required", output.Issues[0].Rule.Name)
			assert.Equal(t, tc.wantSeverity, output.Issues[0].Rule.Severity)
			localPrefix := "./modules/"
			if tc.scope == "submodule" {
				localPrefix = "../"
			}
			assert.Contains(t, output.Issues[0].Message, localPrefix+"<name>")
		})
	}
}
