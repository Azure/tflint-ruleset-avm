// Package integration provides integration tests for tflint.
// Make sure to install tflint & the plugin.
//
// To install the plugin use `make install`.
package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name                  string
		Command               *exec.Cmd
		Dir                   string
		ExpectedIssueRuleName *string
	}{
		{
			Name:    "interface-private-endpoint",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "interface-private-endpoint",
		},
		{
			Name:                  "interface-private-endpoint-incorrect",
			Command:               exec.Command("tflint", "--format", "json", "--force"),
			Dir:                   "interface-private-endpoint-incorrect",
			ExpectedIssueRuleName: p("private_endpoints"),
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

			var stdout, stderr bytes.Buffer
			tc.Command.Stdout = &stdout
			tc.Command.Stderr = &stderr
			if err := tc.Command.Run(); err != nil {
				t.Fatalf("%s, stdout=%s stderr=%s", err, stdout.String(), stderr.String())
			}
		})
	}
}

func p[T any](v T) *T {
	return &v
}

func IsWindowsResultExist() bool {
	_, err := os.Stat("result_windows.json")
	return !os.IsNotExist(err)
}
