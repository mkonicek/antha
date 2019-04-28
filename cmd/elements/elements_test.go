package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

var (
	keepPtr   = flag.Bool("keep", false, "Keep the test environment even if testing is successful")
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

	if inDir == "" && len(flag.Args()) == 0 {
		t.Skip("No workflow provided.")
	}

	if wf, err := makeWorkflow(l, nil, inDir, ""); err != nil {
		t.Fatal(err)
	} else if err := wf.Validate(); err != nil {
		t.Fatal(err)
	} else {

		t.Run("CompileAndTest", func(t *testing.T) {
			buildDir := filepath.Join(outDir, "compileAllElements")
			pkgs, err := compileElements(l, inDir, buildDir, wf)
			if err != nil {
				t.Fatal(err)
			}
			// go test relies on the checkout of the elements so it makes
			// some sense for that to depend on the
			// checkout/transpilation/compilation of the elements.
			if err := goTest(l, buildDir, wf, pkgs); err != nil {
				t.Fatal(err)
			}
			if keepPtr == nil || !*keepPtr {
				os.RemoveAll(buildDir)
			}
		})

		t.Run("Workflows", func(t *testing.T) {
			testingDir := filepath.Join(outDir, "testWorkflows")
			if err := workflows(l, inDir, testingDir, wf); err != nil {
				t.Fatal(err)
			}
			if keepPtr == nil || !*keepPtr {
				os.RemoveAll(testingDir)
			}
		})
	}
}

func compileElements(l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) (string, error) {
	if cb, err := composer.NewComposerBase(l, inDir, outDir); err != nil {
		return "", err
	} else {
		defer cb.CloseLogs()
		//                             /-- keep is true because we need the source for go test
		if err := cb.ComposeMainAndRun(true, false, false, wf); err != nil {
			return "", err
		}
		return cb.GoList()
	}
}

func goTest(l *logger.Logger, outDir string, wf *workflow.Workflow, coverPkgs string) error {
	for repoName := range wf.Repositories {
		cmd := exec.Command("go", "test", "-v")
		cmd.Dir = filepath.Join(outDir, "workflow")

		if coverPkgs != "" {
			cmd.Args = append(cmd.Args,
				"-covermode=atomic",
				fmt.Sprintf("-coverprofile=%s", filepath.Join(cmd.Dir, "cover.out")),
				fmt.Sprintf("-coverpkg=%s", coverPkgs),
			)
		}

		cmd.Args = append(cmd.Args, path.Join(string(repoName), "..."))

		if err := composer.RunAndLogCommand(cmd, composer.RawLogger); err != nil {
			return err
		}
	}
	return nil
}

func workflows(l *logger.Logger, inDir, outDir string, wf *workflow.Workflow) error {
	if cb, err := composer.NewComposerBase(l, inDir, outDir); err != nil {
		return err
	} else {
		defer cb.CloseLogs()
		tc := cb.NewTestsComposer(true, true, true) // keep, run, linked

		wfPaths, err := workflow.GatherPaths(nil, inDir)
		if err != nil {
			return err
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
						l.Log("repository", repoName, "added", f.Name)
						inDir := filepath.Join(outDir, "src", string(repoName), filepath.Dir(filepath.Dir(f.Name)))
						return tc.AddWorkflow(wfTest, inDir)
					}
				} else {
					return nil // not json
				}
			})
			if err != nil {
				return err
			}
		}
		return tc.ComposeTestsAndRun()
	}
}
