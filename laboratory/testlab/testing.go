// This package exists to facilitate running workflows as tests, in
// particular where you wish to have several workflows within the same
// package that you wish to test and so they must share various
// resources such as directories and flags.
//
// This package should only be imported by other testing code because
// it necessarily declares various globals (e.g. flags) which will
// mess up your life if you're not careful.
//
// Also in here are various helper methods, particularly WithTestLab
// that allows you to write "normal go tests" that need to work with
// the laboratory, for example for tests involving files and so forth.
package testlab

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/workflow"
)

var (
	outDirPtr = flag.String("outdir", "", "Directory to write to (default: a temporary directory will be created)")

	sharedInventoryGuard sync.Once
	sharedInventory      *inventory.Inventory
)

func NewTestLabBuilder(t *testing.T, inDir string, fh io.ReadCloser) *laboratory.LaboratoryBuilder {
	inv := EnsureSharedInventory()
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
	if err := labBuild.Setup(fh, inDir, outDir, inv); err != nil {
		labBuild.Fatal(err)
	}

	return labBuild
}

func WithTestLab(t *testing.T, inDir string, callbacks *TestElementCallbacks) {
	wf := workflow.EmptyWorkflow()
	wf.JobId = workflow.JobId("testing")
	wfBuf := new(bytes.Buffer)
	if err := wf.ToWriter(wfBuf, false); err != nil {
		t.Fatal(err)
	}

	inv := EnsureSharedInventory()

	labBuild := laboratory.EmptyLaboratoryBuilder(func(err error) { t.Fatal(err) })
	if err := labBuild.Setup(ioutil.NopCloser(wfBuf), inDir, "", inv); err != nil {
		labBuild.Fatal(err)
	}
	defer labBuild.Decommission()
	NewTestElement(t, labBuild, callbacks)
	if err := labBuild.RunElements(); err != nil {
		labBuild.Fatal(err)
	}
}

func EnsureSharedInventory() *inventory.Inventory {
	sharedInventoryGuard.Do(func() {
		id := id.NewIDGenerator("testing")
		sharedInventory = inventory.NewInventory(id)
		sharedInventory.LoadForWorkflow(nil)
	})
	return sharedInventory
}

type TestElement struct {
	t   *testing.T
	cbs *TestElementCallbacks
}

type TestElementCallbacks struct {
	Name       string
	Setup      func(*laboratory.Laboratory) error
	Steps      func(*laboratory.Laboratory) error
	Analysis   func(*laboratory.Laboratory) error
	Validation func(*laboratory.Laboratory) error
}

func NewTestElement(t *testing.T, installer laboratory.ElementInstaller, cbs *TestElementCallbacks) *TestElement {
	elem := &TestElement{
		t:   t,
		cbs: cbs,
	}
	installer.InstallElement(elem)
	return elem
}

func (te *TestElement) Name() workflow.ElementInstanceName {
	name := te.cbs.Name
	if name == "" {
		name = te.t.Name()
	}
	return workflow.ElementInstanceName(name)
}

func (te *TestElement) TypeName() workflow.ElementTypeName {
	return workflow.ElementTypeName("TestElement")
}

func (te *TestElement) Setup(lab *laboratory.Laboratory) error {
	if te.cbs.Setup != nil {
		return te.cbs.Setup(lab)
	} else {
		return nil
	}
}

func (te *TestElement) Steps(lab *laboratory.Laboratory) error {
	if te.cbs.Steps != nil {
		return te.cbs.Steps(lab)
	} else {
		return nil
	}
}
func (te *TestElement) Analysis(lab *laboratory.Laboratory) error {
	if te.cbs.Analysis != nil {
		return te.cbs.Analysis(lab)
	} else {
		return nil
	}
}
func (te *TestElement) Validation(lab *laboratory.Laboratory) error {
	if te.cbs.Validation != nil {
		return te.cbs.Validation(lab)
	} else {
		return nil
	}
}
