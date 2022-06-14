package main

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

const packagePrefix = "SA"
const checkStringIndexToContains = "S1003"
const checkFunctionReturnErrorLast = "ST1008"
const checkRightCompareOrder = "ST1017"

func GetStaticChecks() []*analysis.Analyzer {
	var checks []*analysis.Analyzer
	additionalChecks := map[string]bool{
		checkStringIndexToContains:   true,
		checkFunctionReturnErrorLast: true,
		checkRightCompareOrder:       true,
	}

	for _, v := range staticcheck.Analyzers {
		_, ok := additionalChecks[v.Name]
		if strings.Contains(v.Name, packagePrefix) || ok {
			checks = append(checks, v)
		}
	}
	return checks
}
