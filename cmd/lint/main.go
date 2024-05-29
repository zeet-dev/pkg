package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/zeet-dev/pkg/linters"
)

// See Makefile for usage
func main() {
	var analyzers = []*analysis.Analyzer{
		linters.ErrorsAsAnalyzer,
		linters.ErrorsStackAnalyzer,
	}

	multichecker.Main(analyzers...)
}
