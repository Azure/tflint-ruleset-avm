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
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name    string
		Command *exec.Cmd
		Dir     string
	}{
		{
			Name:    "optional-defaults",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "optional-defaults",
		},
		{
			Name:    "optional-defaults-incorrect",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "optional-defaults-incorrect",
		},
		{
			Name:    "simplevaluerule-null-value",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "simplevaluerule-null-value",
		},
		{
			Name:    "unknownrule-null",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "unknownrule-null",
		},
		{
			Name:    "unknownrule-null-incorrect",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "unknownrule-null-incorrect",
		},
	}

	dir, _ := os.Getwd()
	defer os.Chdir(dir)

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

			var b []byte
			var err error
			if runtime.GOOS == "windows" && IsWindowsResultExist() {
				b, err = os.ReadFile(filepath.Join(testDir, "result_windows.json"))
			} else {
				b, err = os.ReadFile(filepath.Join(testDir, "result.json"))
			}
			if err != nil {
				t.Fatal(err)
			}

			var expected interface{}
			if err := json.Unmarshal(b, &expected); err != nil {
				t.Fatal(err)
			}

			var got interface{}
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, expected); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func IsWindowsResultExist() bool {
	_, err := os.Stat("result_windows.json")
	return !os.IsNotExist(err)
}
