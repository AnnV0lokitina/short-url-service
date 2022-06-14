package main

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

const packagePrefix = "SA"
const includeCode = []string{"S1003"}

func GetPassesChecks() []*analysis.Analyzer {
	var checks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Name, packagePrefix) {
			checks = append(checks, v)
		}
	}

	return checks
}
