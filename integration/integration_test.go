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
			Name     string `json:"name"`
			Severity string `json:"severity"`
		} `json:"rule"`
		Message string `json:"message"`
	} `json:"issues"`
	Errors []any `json:"errors"`
}

type expectedIssue struct {
	Name     string
	Severity string
}

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name           string
		Dir            string
		ExpectedIssues []expectedIssue
		ExpectFailure  bool
	}{
		{
			Name: "interface-private-endpoint",
			Dir:  "interface-private-endpoint",
		},
		{
			Name: "interface-private-endpoint-deprecated",
			Dir:  "interface-private-endpoint-deprecated",
			ExpectedIssues: []expectedIssue{
				{Name: "deprecated_private_endpoints_interface", Severity: "info"},
			},
		},
		{
			Name: "interface-private-endpoint-incorrect",
			Dir:  "interface-private-endpoint-incorrect",
			ExpectedIssues: []expectedIssue{
				{Name: "private_endpoints", Severity: "error"},
			},
			ExpectFailure: true,
		},
		{
			Name: "basic-ext",
			Dir:  "basic-ext",
			ExpectedIssues: []expectedIssue{
				{Name: "terraform_variable_separate", Severity: "info"},
			},
		},
		{
			Name: "azurerm-ext",
			Dir:  "azurerm-ext",
			ExpectedIssues: []expectedIssue{
				{Name: "azurerm_arg_order", Severity: "info"},
			},
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

			cmd := exec.Command("tflint", "--format", "json", "--minimum-failure-severity=warning")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if tc.ExpectFailure && err == nil {
				t.Fatalf("tflint succeeded but a failure was expected, stdout=%s stderr=%s", stdout.String(), stderr.String())
			}
			if !tc.ExpectFailure && err != nil {
				t.Fatalf("tflint failed to execute: %s, stdout=%s stderr=%s",
					err, stdout.String(), stderr.String())
			}

			var out tflintOutput
			if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
				t.Fatalf("failed to parse tflint JSON output: %s\nstdout=%s",
					err, stdout.String())
			}

			got := make([]expectedIssue, 0, len(out.Issues))
			for _, iss := range out.Issues {
				got = append(got, expectedIssue{Name: iss.Rule.Name, Severity: iss.Rule.Severity})
			}
			sort.Slice(got, func(i, j int) bool {
				if got[i].Name == got[j].Name {
					return got[i].Severity < got[j].Severity
				}
				return got[i].Name < got[j].Name
			})

			want := append([]expectedIssue(nil), tc.ExpectedIssues...)
			sort.Slice(want, func(i, j int) bool {
				if want[i].Name == want[j].Name {
					return want[i].Severity < want[j].Severity
				}
				return want[i].Name < want[j].Name
			})

			if !equalIssues(got, want) {
				t.Fatalf(
					"unexpected issues reported\n  want: %v\n  got : %v\n  raw stdout: %s",
					want, got, stdout.String(),
				)
			}
		})
	}
}

func equalIssues(a, b []expectedIssue) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Severity != b[i].Severity {
			return false
		}
	}
	return true
}
