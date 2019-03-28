package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

var (
	inDirPtr  = flag.String("indir", "", "Directory from which to read files (optional)")
	outDirPtr = flag.String("outdir", "", "Directory to write to (default: a temporary directory will be created)")
)

func TestElements(t *testing.T) {
	l := logger.NewLogger()

	outDir := ""
	if outDirPtr != nil {
		outDir = *outDirPtr
	}
	if outDir == "" { // outDir is shared between the tests, so we have to fix it here.
		if d, err := ioutil.TempDir("", "antha-test"); err != nil {
			t.Fatal(err)
		} else {
			outDir = d
		}
	}

	inDir := ""
	if inDirPtr != nil {
		inDir = *inDirPtr
	}

	if wf, err := makeWorkflow(l, nil, inDir, ""); err != nil {
		t.Fatal(err)
	} else if err := wf.Validate(); err != nil {
		t.Fatal(err)
	} else {

		t.Run("CompileAndTest", func(t *testing.T) {
			compileDir := filepath.Join(outDir, "compile")
			compileElements(t, l, inDir, compileDir, wf)
			// go test relies on the checkout of the elements so it makes
			// some sense for that to depend on the
			// checkout/transpilation/compilation of the elements.
			goTest(t, l, compileDir)
			os.RemoveAll(compileDir)
		})

		t.Run("Bundles", func(t *testing.T) {
			bundleDir := filepath.Join(outDir, "bundle")
			bundles(t, l, inDir, bundleDir, wf)
			os.RemoveAll(bundleDir)
		})
	}
}

func compileElements(t *testing.T, l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) {
	if cb, err := composer.NewComposerBase(l, inDir, outDir); err != nil {
		t.Fatal(err)
	} else {
		defer cb.CloseLogs()
		//                             /-- keep is true because we need the source for go test
		if err := cb.ComposeMainAndRun(true, false, false, wf); err != nil {
			t.Fatal(err)
		}
	}
}

func goTest(t *testing.T, l *logger.Logger, outDir string) {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = filepath.Join(outDir, "src")
	cmd.Env = composer.SetEnvGoPath(os.Environ(), outDir)

	if err := composer.RunAndLogCommand(cmd, l.With("cmd", "test").Log); err != nil {
		t.Fatal(err)
	}
}

func bundles(t *testing.T, l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) {
	if cb, err := composer.NewComposerBase(l, inDir, outDir); err != nil {
		t.Fatal(err)
	} else {
		defer cb.CloseLogs()
		tc := cb.NewTestsComposer(true, true, true)

		wfPaths, err := workflow.GatherPaths(nil, inDir)
		if err != nil {
			t.Fatal(err)
		}
		for repoName, repo := range wf.Repositories {
			err := repo.Walk(func(f *workflow.File) error {
				if filepath.Ext(f.Name) == ".json" {
					// attempt to parse it as a workflow, but don't worry too
					// much if we fail.
					if rc, err := f.Contents(); err != nil {
						return err // this is an error we need to report on!
					} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
						return err
					} else if wfTest, err := workflow.WorkflowFromReaders(append(rs, rc)...); err != nil {
						l.Log("repository", repoName, "skipping", f.Name, "error", err)
						return nil
					} else if err := wfTest.Validate(); err != nil {
						l.Log("repository", repoName, "skipping", f.Name, "error", err)
						return nil
					} else {
						return tc.AddWorkflow(wfTest)
					}
				} else {
					return nil // not json
				}
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		if err := tc.ComposeTestsAndRun(); err != nil {
			t.Fatal(err)
		}
	}
}
