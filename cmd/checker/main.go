// Package checker checks code
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	checks := GetPassesChecks()
	staticChecks := GetStaticChecks()
	checks = append(checks, staticChecks...)
	checks = append(checks, NoExitAnalyzer)
	multichecker.Main(checks...)
}
