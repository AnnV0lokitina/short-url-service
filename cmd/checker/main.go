package main

import (
	"strings"

	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	checks := GetPassesChecks()
	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Name, "SA") {
			checks = append(checks, v)
		}
	}
	multichecker.Main(checks...)
}
