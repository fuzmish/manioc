package manioctypechecker_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/fuzmish/manioc/linter/manioctypechecker"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis/analysistest"
)

func normalizeModFile(path string) (io.Reader, error) {
	// load
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// parse
	file, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	// find all relative path and convert to absolute path
	for _, r := range append([]*modfile.Replace{}, file.Replace...) {
		if modfile.IsDirectoryPath(r.New.Path) {
			vOld := r.Old
			vNew := r.New
			ret, err := filepath.Abs(filepath.Join(dir, vNew.Path))
			if err != nil {
				return nil, err
			}
			if err := file.DropReplace(vOld.Path, vOld.Version); err != nil {
				return nil, err
			}
			if err := file.AddReplace(vOld.Path, vOld.Version, ret, vNew.Version); err != nil {
				return nil, err
			}
		}
	}
	// serialize
	out, err := file.Format()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdataPath := analysistest.TestData()
	testPackage := "a"
	// Since `testutil.WithModules` will copy the test sources to a temporary location,
	// it will break the module replacement specification by relative local paths.
	// To prevent this, here we attempt to generate the content of go.mod file at runtime,
	// replacing the relative paths with absolute paths.
	modfile, err := normalizeModFile(filepath.Join(testdataPath, "src", testPackage, "go.mod"))
	if err != nil {
		t.Fatal(err)
		return
	}
	// run tests
	testdata := testutil.WithModules(t, testdataPath, modfile)
	analysistest.Run(t, testdata, manioctypechecker.Analyzer, testPackage)
}
