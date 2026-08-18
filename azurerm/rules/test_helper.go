package rules

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"reflect"
	"strings"
	"testing"
)

// AssertIssues is an assertion helper for comparing issues.
func AssertIssues(t *testing.T, expected helper.Issues, actual helper.Issues) {
	opts := []cmp.Option{
		// Range field will be ignored because it's not that important in tests such as positions
		cmpopts.IgnoreTypes(hcl.Range{}),
		messageComparer(),
	}
	if !cmp.Equal(expected, actual, opts...) {
		t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(expected, actual, opts...))
	}
}

func messageComparer() cmp.Option {
	return cmp.Comparer(func(x, y helper.Issue) bool {
		return reflect.TypeOf(x.Rule) == reflect.TypeOf(y.Rule) && pruneMessage(x.Message) == pruneMessage(y.Message)
	})
}

func pruneMessage(msg string) string {
	msg = strings.ReplaceAll(msg, " ", "")
	msg = strings.ReplaceAll(msg, "\t", "")
	return msg
}
