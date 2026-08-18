package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// generator collects rule metadata from the central rule register and writes
// a markdown table to RULES.md in the repo root. Intended to be invoked via
// `go generate` (see generate.go in the root package).
func main() {
	// Copy slice so we can sort without mutating original backing array (defensive)
	rs := append([]tflint.Rule(nil), rules.Rules...)

	// Sort by rule.Name() ascending (stable for deterministic output)
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Name() < rs[j].Name() })

	f, err := os.Create("RULES.md")
	if err != nil {
		panic(err)
	}
	defer f.Close() //nolint:errcheck // best-effort close of output file in a CLI generator

	w := bufio.NewWriter(f)
	defer w.Flush() //nolint:errcheck // best-effort flush; any error will surface on the deferred f.Close below

	fmt.Fprintln(w, "# Rules Reference")                                                                                                                                                                  //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	fmt.Fprintln(w)                                                                                                                                                                                       //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	fmt.Fprintln(w, "This document lists all rules currently registered in this ruleset. The Enabled column reflects the default state (some consolidated rules are wrapped to be disabled by default).") //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	fmt.Fprintln(w)                                                                                                                                                                                       //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	fmt.Fprintln(w, "| Name | Enabled | Severity | Link |")                                                                                                                                               //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	fmt.Fprintln(w, "| ---- | ------- | -------- | ---- |")                                                                                                                                               //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush

	for _, r := range rs {
		name := r.Name()
		enabled := "false"
		if r.Enabled() { // default state
			enabled = "true"
		}
		// Normalize severity textual representation.
		severity := severityString(r.Severity())
		link := linkOrDash(r)
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n", name, enabled, severity, link) //nolint:errcheck // writes to bufio.Writer; errors are caught at Flush
	}
}

func linkOrDash(r tflint.Rule) string {
	// Link() may return empty string. Some rules may not implement Link(); rely on interface existence.
	// Detect empty and return '-'. Escape pipes by replacing with %7C (very unlikely).
	l := r.Link()
	if l == "" {
		return "-"
	}
	// Prevent table breakage (simple sanitization)
	l = strings.ReplaceAll(l, "|", "%7C")
	return fmt.Sprintf("[%s](%s)", displayLinkText(l), l)
}

func displayLinkText(l string) string {
	// Show domain/path tail for brevity
	if len(l) > 60 {
		return l[:57] + "..."
	}
	return l
}

func severityString(s tflint.Severity) string {
	// Mirror tflint severities if possible; fallback to numeric.
	switch s {
	case tflint.ERROR:
		return "ERROR"
	case tflint.WARNING:
		return "WARNING"
	case tflint.NOTICE:
		return "NOTICE"
	default:
		return fmt.Sprintf("%d", s)
	}
}
