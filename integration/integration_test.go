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
	"slices"
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
		Range   struct {
			Filename string `json:"filename"`
			Start    struct {
				Line int `json:"line"`
			} `json:"start"`
		} `json:"range"`
	} `json:"issues"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func TestResourceExclusionsIntegration(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name          string
		dir           string
		rule          string
		wantCount     int
		expectFailure bool
		filename      string
		lines         []int
	}{
		{
			name:          "wiring applies only to business resources",
			dir:           "resource-exclusions-wiring",
			rule:          "avm_interface_resource_types",
			wantCount:     2,
			expectFailure: true,
			filename:      "main.tf",
			lines:         []int{6, 10},
		},
		{
			name:          "mixed resources retain a fileless definition diagnostic",
			dir:           "resource-exclusions-fileless",
			rule:          "avm_interface_resource_types",
			wantCount:     1,
			expectFailure: true,
		},
		{
			name:      "telemetry only is inapplicable",
			dir:       "resource-exclusions-telemetry-only",
			rule:      "avm_interface_resource_types",
			wantCount: 0,
		},
		{
			name:      "raw heredocs honor managed and data selectors",
			dir:       "resource-exclusions-raw",
			rule:      "avm_terraform_literal_heredoc_disallowed",
			wantCount: 2,
			filename:  "main.tf",
			lines:     []int{14, 20},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := exec.Command(
				"tflint",
				"--chdir="+filepath.Join(dir, test.dir),
				"--format",
				"json",
				"--only="+test.rule,
				"--minimum-failure-severity=warning",
			)
			var stdout, stderr bytes.Buffer
			command.Stdout = &stdout
			command.Stderr = &stderr
			if err := command.Run(); (err != nil) != test.expectFailure {
				t.Fatalf("unexpected tflint exit state: %v, stdout=%s stderr=%s", err, stdout.String(), stderr.String())
			}

			var output tflintOutput
			if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
				t.Fatalf("failed to parse tflint JSON output: %s\nstdout=%s", err, stdout.String())
			}
			if len(output.Issues) != test.wantCount {
				t.Fatalf("expected %d issues, got %s", test.wantCount, stdout.String())
			}
			for _, issue := range output.Issues {
				if issue.Rule.Name != test.rule {
					t.Fatalf("unexpected issue %#v", issue)
				}
				if test.rule == "avm_terraform_literal_heredoc_disallowed" {
					continue
				}
				if !bytes.Contains([]byte(issue.Message), []byte("variable `resource_types`")) &&
					!bytes.Contains([]byte(issue.Message), []byte("using variable `resource_types`")) {
					t.Fatalf("unexpected issue message %q", issue.Message)
				}
			}
			if test.filename != "" {
				actualLines := make([]int, 0, len(output.Issues))
				for _, issue := range output.Issues {
					if filepath.Base(issue.Range.Filename) != test.filename {
						t.Fatalf("expected diagnostic in %s, got %#v", test.filename, issue.Range)
					}
					actualLines = append(actualLines, issue.Range.Start.Line)
				}
				if !slices.Equal(actualLines, test.lines) {
					t.Fatalf("expected diagnostic lines %v, got %v", test.lines, actualLines)
				}
			} else if test.wantCount > 0 && (output.Issues[0].Range.Filename != "" || output.Issues[0].Range.Start.Line != 0) {
				t.Fatalf("expected the missing variables.tf diagnostic to be explicitly fileless, got %#v", output.Issues[0].Range)
			}
		})
	}
}

type expectedIssue struct {
	Name     string
	Severity string
}

const invalidSeverityError = `severity must be one of "error", "warning", or "notice", got "critical"`

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name           string
		Dir            string
		OnlyRule       string
		UseChdir       bool
		UseRecursive   bool
		ExpectedIssues []expectedIssue
		ExpectFailure  bool
		ExpectedError  string
	}{
		{
			Name: "interface-private-endpoint",
			Dir:  "interface-private-endpoint",
		},
		{
			Name: "interface-private-endpoint-deprecated",
			Dir:  "interface-private-endpoint-deprecated",
			ExpectedIssues: []expectedIssue{
				{Name: "avm_interface_private_endpoints_deprecated", Severity: "info"},
			},
		},
		{
			Name: "interface-private-endpoint-incorrect",
			Dir:  "interface-private-endpoint-incorrect",
			ExpectedIssues: []expectedIssue{
				{Name: "avm_interface_private_endpoints", Severity: "error"},
			},
			ExpectFailure: true,
		},
		{
			Name: "azapi-required-interfaces",
			Dir:  "azapi-required-interfaces",
		},
		{
			Name: "azapi-required-interfaces-incorrect",
			Dir:  "azapi-required-interfaces-incorrect",
			ExpectedIssues: []expectedIssue{
				{Name: "avm_interface_ignore_body_changes", Severity: "error"},
				{Name: "avm_interface_resource_types", Severity: "info"},
				{Name: "avm_interface_retry", Severity: "error"},
				{Name: "avm_interface_timeouts", Severity: "error"},
			},
			ExpectFailure: true,
		},
		{
			Name: "azapi-required-interfaces-child-parent",
			Dir:  "azapi-required-interfaces-child",
		},
		{
			Name: "azapi-required-interfaces-child-module",
			Dir:  filepath.Join("azapi-required-interfaces-child", "child"),
		},
		{
			Name:          "invalid-severity",
			Dir:           "invalid-severity",
			ExpectFailure: true,
			ExpectedError: invalidSeverityError,
		},
		{
			Name:     "output-resource-id-pattern-chdir",
			Dir:      "output-resource-id-pattern",
			OnlyRule: "avm_output_resource_id_required",
			UseChdir: true,
		},
		{
			Name:         "output-resource-id-pattern-recursive",
			Dir:          "output-resource-id-pattern",
			OnlyRule:     "avm_output_resource_id_required",
			UseRecursive: true,
		},
		{
			Name:     "output-resource-id-resource-chdir",
			Dir:      "output-resource-id-resource",
			OnlyRule: "avm_output_resource_id_required",
			UseChdir: true,
			ExpectedIssues: []expectedIssue{
				{Name: "avm_output_resource_id_required", Severity: "error"},
			},
			ExpectFailure: true,
		},
		{
			Name:         "output-resource-id-resource-recursive",
			Dir:          "output-resource-id-resource",
			OnlyRule:     "avm_output_resource_id_required",
			UseRecursive: true,
			ExpectedIssues: []expectedIssue{
				{Name: "avm_output_resource_id_required", Severity: "error"},
			},
			ExpectFailure: true,
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

			if !tc.UseChdir {
				if err := os.Chdir(testDir); err != nil {
					t.Fatal(err)
				}
			}

			args := []string{"--format", "json", "--minimum-failure-severity=warning"}
			if tc.OnlyRule != "" {
				args = append(args, "--only="+tc.OnlyRule)
			}
			if tc.UseChdir {
				args = append([]string{"--chdir=" + testDir}, args...)
			} else if tc.UseRecursive {
				args = append([]string{"--recursive"}, args...)
			}
			cmd := exec.Command("tflint", args...)
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
			if tc.ExpectedError != "" {
				found := false
				for _, outputError := range out.Errors {
					if bytes.Contains([]byte(outputError.Message), []byte(tc.ExpectedError)) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf(
						"expected error %q was not reported, stdout=%s stderr=%s",
						tc.ExpectedError,
						stdout.String(),
						stderr.String(),
					)
				}
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
