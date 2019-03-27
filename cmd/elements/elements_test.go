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

		// compileDir must be shared between CompileElements and GoTest, and must be distinct from bundleDir:
		compileDir := filepath.Join(outDir, "compile")
		bundleDir := filepath.Join(outDir, "bundle")
		t.Run("CompileAndTest", func(t *testing.T) {
			compileElements(t, l, inDir, compileDir, wf)
			// go test relies on the checkout of the elements so it makes
			// some sense for that to depend on the
			// checkout/transpilation/compilation of the elements.
			goTest(t, l, compileDir)
		})
		t.Run("Bundles", func(t *testing.T) { bundles(t, l, inDir, bundleDir, wf) })
	}
}

func compileElements(t *testing.T, l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) {
	if cb, err := composer.NewComposerBase(l, inDir, outDir); err != nil {
		t.Fatal(err)
	} else {
		defer cb.CloseLogs()
		if err := cb.ComposeMainAndRun(wf, true, false, false); err != nil {
			t.Fatal(err) //                 ^^ keep is true because we need the source for go test
		}
	}
}

func goTest(t *testing.T, l *logger.Logger, outDir string) {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = filepath.Join(outDir, "src")
	cmd.Env = composer.SetEnvGoPath(os.Environ(), outDir)

	if err := composer.RunAndLogCommand(cmd, l.With("cmd", "go test").Log); err != nil {
		t.Fatal(err)
	}
}

func bundles(t *testing.T, l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) {
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
					t.Run(f.Name, func(t *testing.T) {
						// calling parallel means that t.Run won't wait for this test to finish.
						t.Parallel()
						runBundle(t, l, wfTest, f.Name, outDir)
					})
					return nil
				}
			} else {
				return nil // not json
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	// once this function returns, testing will re-schedule all the
	// tests marked Parallel() and wait for them all.
}

func runBundle(t *testing.T, l *logger.Logger, wf *workflow.Workflow, bundleName, outDir string) {
	if outDir, err := ioutil.TempDir(outDir, ""); err != nil {
		t.Fatal(err)
	} else {
		l.Log("bundle", bundleName, "outdir", outDir)
		if cb, err := composer.NewComposerBase(l, filepath.Join(outDir, "src", filepath.Dir(bundleName)), outDir); err != nil {
			t.Fatal(err)
		} else if err := cb.ComposeMainAndRun(wf, true, true, true); err != nil {
			//                                     keep and run and linkedDrivers
			cb.CloseLogs()
			t.Fatal(err)
		} else {
			cb.CloseLogs()
			if err := os.RemoveAll(outDir); err != nil { // tidy up iff the test was successful to avoid exhausting disk space!
				t.Fatal(err)
			}
		}
	}
}
