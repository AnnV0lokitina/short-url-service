package main

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"testing"
)

func TestGetPassesChecks(t *testing.T) {
	checks := GetPassesChecks()
	assert.IsType(t, checks, []*analysis.Analyzer{})
}
