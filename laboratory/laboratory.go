package laboratory

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/target/qpcrdevice"
	"github.com/antha-lang/antha/target/shakerincubator"
	"github.com/antha-lang/antha/target/woplatereader"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
)

type Element interface {
	Name() workflow.ElementInstanceName
	TypeName() workflow.ElementTypeName

	Setup(*Laboratory) error
	Steps(*Laboratory) error
	Analysis(*Laboratory) error
	Validation(*Laboratory) error
}

type LaboratoryBuilder struct {
	inDir    string
	outDir   string
	workflow *workflow.Workflow

	elemLock      sync.Mutex
	elemsUnrun    int64
	elements      map[Element]*ElementBase
	nextElementId uint64

	// This lock is here to serialize access to errors (append and
	// Pack()), and to avoid races around closing Errored chan
	errLock sync.Mutex
	errors  utils.ErrorSlice
	Errored chan struct{}

	Completed chan struct{}

	Fatal func(error) // Fatal is a field and not a method so that we can dynamically change it

	*lineMapManager
	Logger *logger.Logger
	logFH  *os.File

	effects *effects.LaboratoryEffects

	instrs effects.Insts
}

func EmptyLaboratoryBuilder(fatalFunc func(error)) *LaboratoryBuilder {
	labBuild := &LaboratoryBuilder{
		elements:  make(map[Element]*ElementBase),
		Errored:   make(chan struct{}),
		Completed: make(chan struct{}),
		Fatal:     fatalFunc,

		lineMapManager: NewLineMapManager(),
		Logger:         logger.NewLogger(),
	}
	if fatalFunc == nil {
		// we wrap in a func here because we may change the value of
		// Logger (see SetupWorkflow) and so don't want to capture an
		// old value here.
		labBuild.Fatal = func(err error) { labBuild.Logger.Fatal(err) }
	}
	return labBuild
}

func parseFlags() (inDir, outDir string) {
	flag.StringVar(&inDir, "indir", "", "Path to directory from which to read input files")
	flag.StringVar(&outDir, "outdir", "", "Path to directory in which to write output files")
	flag.Parse()
	return
}

func NewLaboratoryBuilder(fh io.ReadCloser) *LaboratoryBuilder {
	labBuild := EmptyLaboratoryBuilder(nil)
	inDir, outDir := parseFlags()
	if err := labBuild.Setup(fh, inDir, outDir); err != nil {
		labBuild.Fatal(err)
	}
	return labBuild
}

func (labBuild *LaboratoryBuilder) Setup(fh io.ReadCloser, inDir, outDir string) error {
	return utils.ErrorFuncs{
		func() error { return labBuild.SetupWorkflow(fh) },
		func() error { return labBuild.SetupPaths(inDir, outDir) },
		func() error { return labBuild.SetupEffects() },
	}.Run()
}

func (labBuild *LaboratoryBuilder) SetupWorkflow(fh io.ReadCloser) error {
	// Got to load in the workflow first so we gain access to the JobId.
	if wf, err := workflow.WorkflowFromReaders(fh); err != nil {
		return err
	} else if err := wf.Validate(); err != nil {
		return err
	} else {
		labBuild.workflow = wf
		labBuild.Logger = labBuild.Logger.With("jobId", wf.JobId)
		return nil
	}
}

func (labBuild *LaboratoryBuilder) SetupPaths(inDir, outDir string) error {
	labBuild.inDir, labBuild.outDir = inDir, outDir

	// Make sure we have a valid working outDir:
	if labBuild.outDir == "" {
		if d, err := ioutil.TempDir("", "antha-run-outputs"); err != nil {
			return err
		} else {
			labBuild.outDir = d
		}
	}
	labBuild.Logger.Log("outdir", labBuild.outDir)

	// Create subdirs within it:
	for _, leaf := range []string{"elements", "data", "devices"} {
		if err := os.MkdirAll(filepath.Join(labBuild.outDir, leaf), 0700); err != nil {
			return err
		}
	}

	// Switch the logger over to write to disk too:
	if logFH, err := os.OpenFile(filepath.Join(labBuild.outDir, "logs.txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		return err
	} else {
		labBuild.logFH = logFH
		labBuild.Logger.SwapWriters(logFH, os.Stderr)
	}

	// Sort out inDir:
	if labBuild.inDir == "" {
		// We do this to make certain that we have a root path to join
		// onto so we can't permit reading arbitrary parts of the
		// filesystem.
		if d, err := ioutil.TempDir("", "antha-run-inputs"); err != nil {
			return err
		} else {
			labBuild.inDir = d
		}
	}
	labBuild.Logger.Log("indir", labBuild.inDir)
	return nil
}

func (labBuild *LaboratoryBuilder) SetupEffects() error {
	if fm, err := effects.NewFileManager(filepath.Join(labBuild.inDir, "data"), filepath.Join(labBuild.outDir, "data")); err != nil {
		return err
	} else {
		labBuild.effects = effects.NewLaboratoryEffects(string(labBuild.workflow.JobId), fm)
		labBuild.effects.Inventory.LoadForWorkflow(labBuild.workflow)
		return nil
	}
}

func (labBuild *LaboratoryBuilder) SaveErrors() error {
	if err := labBuild.Errors(); err != nil {
		return ioutil.WriteFile(filepath.Join(labBuild.outDir, "errors.txt"), []byte(err.Error()), 0400)
	} else {
		return nil
	}
}

func (labBuild *LaboratoryBuilder) Decommission() {
	if labBuild.logFH != nil {
		labBuild.Logger.SwapWriters(os.Stderr)
		if err := labBuild.logFH.Sync(); err != nil {
			labBuild.Logger.Log("msg", "Error when syncing log file handle", "error", err)
		}
		if err := labBuild.logFH.Close(); err != nil {
			labBuild.Logger.Log("msg", "Error when closing log file handle", "error", err)
		}
		labBuild.logFH = nil
	}
}

func (labBuild *LaboratoryBuilder) Compile() {
	if devices, err := labBuild.connectDevices(); err != nil {
		labBuild.Fatal(err)

	} else {
		defer devices.Close()

		// We have to do this this late because we need the connections
		// to the plugins established to eg figure out if the device
		// supports prompting.
		human.New(labBuild.effects.IDGenerator).DetermineRole(devices)

		devDir := filepath.Join(labBuild.outDir, "devices")

		if nodes, err := labBuild.effects.Maker.MakeNodes(labBuild.effects.Trace.Instructions()); err != nil {
			labBuild.Fatal(err)

		} else if instrs, err := codegen.Compile(labBuild.effects, devDir, devices, nodes); err != nil {
			labBuild.Fatal(err)

		} else {
			labBuild.instrs = instrs
		}
	}
}

func (labBuild *LaboratoryBuilder) connectDevices() (*target.Target, error) {
	tgt := target.New()
	if global, err := mixer.NewGlobalMixerConfig(labBuild.effects.Inventory, &labBuild.workflow.Config.GlobalMixer); err != nil {
		return nil, err
	} else {
		err := utils.ErrorSlice{
			mixer.NewGilsonPipetMaxInstances(labBuild.Logger, tgt, labBuild.effects.Inventory, global, labBuild.workflow.Config.GilsonPipetMax),
			mixer.NewTecanInstances(labBuild.Logger, tgt, labBuild.effects.Inventory, global, labBuild.workflow.Config.Tecan),
			mixer.NewCyBioInstances(labBuild.Logger, tgt, labBuild.effects.Inventory, global, labBuild.workflow.Config.CyBio),
			mixer.NewLabcyteInstances(labBuild.Logger, tgt, labBuild.effects.Inventory, global, labBuild.workflow.Config.Labcyte),
			mixer.NewHamiltonInstances(labBuild.Logger, tgt, labBuild.effects.Inventory, global, labBuild.workflow.Config.Hamilton),
			qpcrdevice.NewQPCRInstances(tgt, labBuild.workflow.Config.QPCR),
			shakerincubator.NewShakerIncubatorsInstances(tgt, labBuild.workflow.Config.ShakerIncubator),
			woplatereader.NewWOPlateReaderInstances(tgt, labBuild.workflow.Config.PlateReader),
		}.Pack()
		if err != nil {
			return nil, err
		} else if err := tgt.Connect(labBuild.workflow); err != nil {
			tgt.Close()
			return nil, err
		} else {
			return tgt, nil
		}
	}
}

func (labBuild *LaboratoryBuilder) Export() {
	if err := export(labBuild.effects.IDGenerator, labBuild.outDir, labBuild.instrs); err != nil {
		labBuild.Fatal(err)
	}
}

// This interface exists just to allow both the lab builder and
// laboratory itself to contain the InstallElement method. We need the
// method on both due to the dynamic inter-element calls.
type ElementInstaller interface {
	InstallElement(Element)
}

func (labBuild *LaboratoryBuilder) InstallElement(e Element) {
	eb := labBuild.NewElementBase(e)
	labBuild.elemLock.Lock()
	defer labBuild.elemLock.Unlock()
	labBuild.elements[e] = eb
	labBuild.elemsUnrun++
}

func (labBuild *LaboratoryBuilder) AddConnection(src, dst Element, fun func()) error {
	if ebSrc, found := labBuild.elements[src]; !found {
		return fmt.Errorf("Unknown src element: %v", src)
	} else if ebDst, found := labBuild.elements[dst]; !found {
		return fmt.Errorf("Unknown dst element: %v", dst)
	} else {
		ebDst.AddBlockedInput()
		ebSrc.AddOnExit(func() {
			fun()
			ebDst.InputReady()
		})
		return nil
	}
}

// Run all the installed elements.
func (labBuild *LaboratoryBuilder) RunElements() error {
	labBuild.elemLock.Lock()
	if labBuild.elemsUnrun == 0 {
		labBuild.elemLock.Unlock()
		close(labBuild.Completed)
		return nil

	} else {
		for _, eb := range labBuild.elements {
			eb.AddOnExit(labBuild.elementCompleted)
		}
		for _, eb := range labBuild.elements {
			go eb.Run(labBuild.makeLab(eb, labBuild.Logger))
		}
		labBuild.elemLock.Unlock()
		<-labBuild.Completed

		return labBuild.Errors()
	}
}

func (labBuild *LaboratoryBuilder) elementCompleted() {
	labBuild.elemLock.Lock()
	defer labBuild.elemLock.Unlock()
	labBuild.elemsUnrun--
	if labBuild.elemsUnrun == 0 {
		close(labBuild.Completed)
	}
}

func (labBuild *LaboratoryBuilder) recordError(err error) {
	labBuild.errLock.Lock()
	defer labBuild.errLock.Unlock()
	labBuild.errors = append(labBuild.errors, err)
	select { // we keep the lock here to avoid a race to close
	case <-labBuild.Errored:
	default:
		close(labBuild.Errored)
	}
}

func (labBuild *LaboratoryBuilder) Errors() error {
	select {
	case <-labBuild.Errored:
		labBuild.errLock.Lock()
		defer labBuild.errLock.Unlock()
		return labBuild.errors.Pack()
	default:
		return nil
	}
}

type Laboratory struct {
	labBuild *LaboratoryBuilder
	element  Element

	Logger *logger.Logger
	*effects.LaboratoryEffects
}

func (labBuild *LaboratoryBuilder) makeLab(eb *ElementBase, logger *logger.Logger) *Laboratory {
	e := eb.element
	return &Laboratory{
		labBuild:          labBuild,
		element:           e,
		Logger:            logger.With("id", eb.id, "name", e.Name(), "type", e.TypeName()),
		LaboratoryEffects: labBuild.effects,
	}
}

func (lab *Laboratory) InstallElement(e Element) {
	lab.labBuild.InstallElement(e)
}

// Only for use when you're in an element and want to call another element.
func (lab *Laboratory) CallSteps(e Element) error {
	// it should already be in the map because the element constructor
	// will have called through to InstallElement which would have
	// added it.
	eb, found := lab.labBuild.elements[e]
	if !found {
		return fmt.Errorf("CallSteps called on unknown element '%s'", e.Name())
	}

	finished := make(chan struct{})
	eb.AddOnExit(func() { close(finished) })
	eb.AddOnExit(lab.labBuild.elementCompleted)

	// take the root logger (from labBuild) and build up from there.
	logger := lab.labBuild.Logger.With("parentName", lab.element.Name(), "parentType", lab.element.TypeName())
	go eb.Run(lab.labBuild.makeLab(eb, logger), eb.element.Steps)
	<-finished

	select {
	case <-lab.labBuild.Errored:
		return lab.labBuild.errors
	default:
		return nil
	}
}

func (lab *Laboratory) error(err error) {
	lab.labBuild.recordError(err)
	lab.Logger.Log("error", err.Error())
}

func (lab *Laboratory) errorf(fmtStr string, args ...interface{}) {
	lab.error(fmt.Errorf(fmtStr, args...))
}

// ElementBase
type ElementBase struct {
	// every element has a unique id to ensure we don't collide on names.
	id uint64
	// count of inputs that are not yet ready (plus 1)
	pendingCount int64
	// this gets closed when all inputs become ready
	InputsReady chan struct{}
	// funcs to run when this element is completed
	onExit []func()
	// the actual element
	element Element
}

func (labBuild *LaboratoryBuilder) NewElementBase(e Element) *ElementBase {
	return &ElementBase{
		id:           atomic.AddUint64(&labBuild.nextElementId, 1),
		pendingCount: 1,
		InputsReady:  make(chan struct{}),
		element:      e,
	}
}

func (eb *ElementBase) Run(lab *Laboratory, funs ...func(*Laboratory) error) {
	defer eb.Completed(lab)
	eb.InputReady()

	if len(funs) == 0 {
		funs = []func(*Laboratory) error{
			eb.element.Setup,
			eb.element.Steps,
			eb.element.Analysis,
			eb.element.Validation,
		}
	}

	defer func() {
		if res := recover(); res != nil {
			lab.errorf("panic: %v", res)
			// Use println because of the embedded \n in the Stack Trace
			fmt.Println(lab.labBuild.lineMapManager.ElementStackTrace())
		}
	}()

	select {
	case <-eb.InputsReady:
		lab.Logger.Log("progress", "starting")
		for _, fun := range funs {
			select {
			case <-lab.labBuild.Errored:
				return
			default:
				if err := fun(lab); err != nil {
					lab.error(err)
					return
				}
			}
		}
	case <-lab.labBuild.Errored:
		return
	}
}

func (eb *ElementBase) Completed(lab *Laboratory) {
	if err := eb.Save(lab); err != nil {
		lab.error(err)
	}
	lab.Logger.Log("progress", "completed")
	funs := eb.onExit
	eb.onExit = nil
	for _, fun := range funs {
		fun()
	}
}

func (eb *ElementBase) Save(lab *Laboratory) error {
	p := filepath.Join(lab.labBuild.outDir, "elements", fmt.Sprintf("%d_%s.json", eb.id, eb.element.Name()))
	if fh, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		return err
	} else {
		defer fh.Close()
		return json.NewEncoder(fh).Encode(eb.element)
	}
}

func (eb *ElementBase) InputReady() {
	if atomic.AddInt64(&eb.pendingCount, -1) == 0 {
		// we've done the transition from 1 -> 0. By definition, we're
		// the only routine that can do that, so we don't need to be
		// careful.
		close(eb.InputsReady)
	}
}

func (eb *ElementBase) AddBlockedInput() {
	atomic.AddInt64(&eb.pendingCount, 1)
}

func (eb *ElementBase) AddOnExit(fun func()) {
	eb.onExit = append(eb.onExit, fun)
}
