// This package exists to facilitate running workflows as tests, in
// particular where you wish to have several workflows within the same
// package that you wish to test and so they must share various
// resources such as directories and flags.
//
// This package should only be imported by other testing code because
// it necessarily declares various globals (e.g. flags) which will
// mess up your life if you're not careful.
package testing

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/antha-lang/antha/laboratory"
)

var (
	outDirPtr = flag.String("outdir", "", "Directory to write to (default: a temporary directory will be created)")
)

func NewTestLabBuilder(t *testing.T, inDir string, fh io.ReadCloser) *laboratory.LaboratoryBuilder {
	outDir := ""
	if outDirPtr != nil {
		outDir = *outDirPtr
	}
	if outDir != "" {
		if err := os.MkdirAll(outDir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	// outDir will be shared for all the tests, so we need to be
	// private within this. If outDir is "" then this is still safe.
	if d, err := ioutil.TempDir(outDir, "antha-test"); err != nil {
		t.Fatal(err)
	} else {
		outDir = d
	}

	labBuild := laboratory.EmptyLaboratoryBuilder(func(err error) { t.Fatal(err) })
	if err := labBuild.Setup(fh, inDir, outDir); err != nil {
		labBuild.Fatal(err)
	}

	return labBuild
}
