package main

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"testing"
)

func TestGetStaticChecks(t *testing.T) {
	checks := GetStaticChecks()
	assert.IsType(t, checks, []*analysis.Analyzer{})
}
