package main

import (
	"github.com/fuzmish/manioc/linter/manioctypechecker"
	"golang.org/x/tools/go/analysis"
)

//nolint:deadcode,gochecknoglobals
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		manioctypechecker.Analyzer,
	}
}
