package main

import (
	"github.com/fuzmish/manioc/linter/manioctypechecker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(manioctypechecker.Analyzer)
}
