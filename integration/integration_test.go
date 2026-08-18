// Package integration provides integration tests for tflint.
// Make sure to install tflint & the plugin.
//
// To install the plugin use `make install`.
package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

// tflintOutput is a minimal structure matching the fields emitted by
// `tflint --format json` that we care about for integration assertions.
type tflintOutput struct {
	Issues []struct {
		Rule struct {
			Name string `json:"name"`
		} `json:"rule"`
		Message string `json:"message"`
	} `json:"issues"`
	Errors []any `json:"errors"`
}

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name string
		Dir  string
		// ExpectedRuleNames is the sorted set of rule names that the tflint run
		// should report. An empty slice means no issues are expected.
		ExpectedRuleNames []string
	}{
		{
			Name:              "interface-private-endpoint",
			Dir:               "interface-private-endpoint",
			ExpectedRuleNames: []string{},
		},
		{
			Name:              "interface-private-endpoint-incorrect",
			Dir:               "interface-private-endpoint-incorrect",
			ExpectedRuleNames: []string{"private_endpoints"},
		},
		{
			Name:              "basic-ext",
			Dir:               "basic-ext",
			ExpectedRuleNames: []string{"terraform_variable_separate"},
		},
		{
			Name:              "azurerm-ext",
			Dir:               "azurerm-ext",
			ExpectedRuleNames: []string{"azurerm_arg_order"},
		},
	}

	dir, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(dir)
	}()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			testDir := filepath.Join(dir, tc.Dir)

			t.Cleanup(func() {
				if err := os.Chdir(dir); err != nil {
					t.Fatal(err)
				}
			})

			if err := os.Chdir(testDir); err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command("tflint", "--format", "json", "--force")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("tflint failed to execute: %s, stdout=%s stderr=%s",
					err, stdout.String(), stderr.String())
			}

			var out tflintOutput
			if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
				t.Fatalf("failed to parse tflint JSON output: %s\nstdout=%s",
					err, stdout.String())
			}

			got := make([]string, 0, len(out.Issues))
			for _, iss := range out.Issues {
				got = append(got, iss.Rule.Name)
			}
			sort.Strings(got)

			want := append([]string(nil), tc.ExpectedRuleNames...)
			sort.Strings(want)

			if !equalStringSlices(got, want) {
				t.Fatalf(
					"unexpected issues reported\n  want rule names: %v\n  got rule names : %v\n  raw stdout: %s",
					want, got, stdout.String(),
				)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
