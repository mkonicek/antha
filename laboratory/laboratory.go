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

	elemLock   sync.Mutex
	elemsUnrun int64
	elements   map[Element]*ElementBase

	// This lock is here to serialize access to errors (append and
	// Pack()), and to avoid races around closing Errored chan
	errLock sync.Mutex
	errors  utils.ErrorSlice
	Errored chan struct{}

	Completed chan struct{}

	*lineMapManager
	Logger      *logger.Logger
	logFH       *os.File
	FileManager *FileManager

	*effects.LaboratoryEffects

	instrs effects.Insts
}

func NewLaboratoryBuilder(fh io.ReadCloser) *LaboratoryBuilder {
	labBuild := &LaboratoryBuilder{
		elements:  make(map[Element]*ElementBase),
		Errored:   make(chan struct{}),
		Completed: make(chan struct{}),

		lineMapManager: NewLineMapManager(),
		Logger:         logger.NewLogger(),
	}

	// Got to load in the workflow first so we gain access to the JobId.
	if wf, err := workflow.WorkflowFromReaders(fh); err != nil {
		labBuild.Fatal(err)
	} else if err := wf.Validate(); err != nil {
		labBuild.Fatal(err)
	} else {
		labBuild.workflow = wf
		labBuild.Logger = labBuild.Logger.With("jobId", wf.JobId)
	}

	flag.StringVar(&labBuild.inDir, "indir", "", "Path to directory from which to read input files")
	flag.StringVar(&labBuild.outDir, "outdir", "", "Path to directory in which to write output files")
	flag.Parse()

	if labBuild.outDir == "" {
		if d, err := ioutil.TempDir("", fmt.Sprintf("antha-run-out-%s", labBuild.workflow.JobId)); err != nil {
			labBuild.Fatal(err)
		} else {
			labBuild.outDir = d
		}
	}
	labBuild.Logger.Log("outdir", labBuild.outDir)
	for _, leaf := range []string{"elements", "data", "devices"} {
		if err := os.MkdirAll(filepath.Join(labBuild.outDir, leaf), 0700); err != nil {
			labBuild.Fatal(err)
		}
	}

	if logFH, err := os.OpenFile(filepath.Join(labBuild.outDir, "logs.txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		labBuild.Fatal(err)
	} else {
		labBuild.logFH = logFH
		labBuild.Logger.SwapWriters(logFH, os.Stderr)
	}

	if labBuild.inDir == "" {
		// We do this to make certain that we have a root path to join
		// onto so we can't permit reading arbitrary parts of the
		// filesystem.
		if d, err := ioutil.TempDir("", fmt.Sprintf("antha-run-in-%s", labBuild.workflow.JobId)); err != nil {
			labBuild.Fatal(err)
		} else {
			labBuild.inDir = d
		}
	}
	labBuild.Logger.Log("indir", labBuild.inDir)

	if fm, err := NewFileManager(filepath.Join(labBuild.inDir, "data"), filepath.Join(labBuild.outDir, "data")); err != nil {
		labBuild.Fatal(err)
	} else {
		labBuild.FileManager = fm
	}
	labBuild.LaboratoryEffects = effects.NewLaboratoryEffects(string(labBuild.workflow.JobId))

	// TODO: discuss this: not sure if we want to do this based off
	// zero plate types defined, or if we want an explicit flag or
	// something?
	if len(labBuild.workflow.Inventory.PlateTypes) == 0 {
		labBuild.Logger.Log("PlateTypeInventory", "loading built in plate types")
		labBuild.Inventory.PlateTypes.LoadLibrary()
	} else {
		labBuild.Inventory.PlateTypes.SetPlateTypes(labBuild.workflow.Inventory.PlateTypes)
	}

	labBuild.Inventory.TipBoxes.LoadLibrary()

	return labBuild
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
		human.New(labBuild.IDGenerator).DetermineRole(devices)

		devDir := filepath.Join(labBuild.outDir, "devices")

		if nodes, err := labBuild.Maker.MakeNodes(labBuild.Trace.Instructions()); err != nil {
			labBuild.Fatal(err)

		} else if instrs, err := codegen.Compile(labBuild.LaboratoryEffects, devDir, devices, nodes); err != nil {
			labBuild.Fatal(err)

		} else {
			labBuild.instrs = instrs
		}
	}
}

func (labBuild *LaboratoryBuilder) connectDevices() (*target.Target, error) {
	tgt := target.New()
	if global, err := mixer.NewGlobalMixerConfig(labBuild.Inventory, &labBuild.workflow.Config.GlobalMixer); err != nil {
		return nil, err
	} else {
		err := utils.ErrorSlice{
			mixer.NewGilsonPipetMaxInstances(labBuild.Logger, tgt, labBuild.Inventory, global, labBuild.workflow.Config.GilsonPipetMax),
			mixer.NewTecanInstances(labBuild.Logger, tgt, labBuild.Inventory, global, labBuild.workflow.Config.Tecan),
			mixer.NewCyBioInstances(labBuild.Logger, tgt, labBuild.Inventory, global, labBuild.workflow.Config.CyBio),
			mixer.NewLabcyteInstances(labBuild.Logger, tgt, labBuild.Inventory, global, labBuild.workflow.Config.Labcyte),
			mixer.NewHamiltonInstances(labBuild.Logger, tgt, labBuild.Inventory, global, labBuild.workflow.Config.Hamilton),
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
	labBuild.FileManager.SummarizeWritten(labBuild.Logger)
	if err := export(labBuild.IDGenerator, labBuild.outDir, labBuild.instrs); err != nil {
		labBuild.Fatal(err)
	}
}

// This interface exists just to allow both the lab builder and
// laboratory itself to trivially contain the InstallElement
// method. It's used in transpiled elements.
type ElementInstaller interface {
	InstallElement(Element)
}

func (labBuild *LaboratoryBuilder) InstallElement(e Element) {
	eb := NewElementBase(e)
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
		for e, eb := range labBuild.elements {
			go eb.Run(labBuild.makeLab(e, labBuild.Logger))
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

func (labBuild *LaboratoryBuilder) Fatal(err error) {
	labBuild.Logger.Fatal(err)
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
	*LaboratoryBuilder
	element Element
	Logger  *logger.Logger
}

func (labBuild *LaboratoryBuilder) makeLab(e Element, logger *logger.Logger) *Laboratory {
	return &Laboratory{
		LaboratoryBuilder: labBuild,
		element:           e,
		Logger:            logger.With("name", e.Name(), "type", e.TypeName()),
	}
}

// Only for use when you're in an element and want to call another element.
func (lab *Laboratory) CallSteps(e Element) error {
	eb := NewElementBase(e)

	finished := make(chan struct{})
	eb.AddOnExit(func() { close(finished) })
	eb.AddOnExit(lab.elementCompleted)

	// take the root logger (from labBuild) and build up from there.
	logger := lab.LaboratoryBuilder.Logger.With("parentName", lab.element.Name(), "parentType", lab.element.TypeName())
	go eb.Run(lab.makeLab(e, logger), eb.element.Steps)
	<-finished

	select {
	case <-lab.Errored:
		return lab.errors
	default:
		return nil
	}
}

func (lab *Laboratory) error(err error) {
	lab.recordError(err)
	lab.Logger.Log("error", err.Error())
}

func (lab *Laboratory) errorf(fmtStr string, args ...interface{}) {
	lab.error(fmt.Errorf(fmtStr, args...))
}

// ElementBase
type ElementBase struct {
	// count of inputs that are not yet ready (plus 1)
	pendingCount int64
	// this gets closed when all inputs become ready
	InputsReady chan struct{}
	// funcs to run when this element is completed
	onExit []func()
	// the actual element
	element Element
}

func NewElementBase(e Element) *ElementBase {
	return &ElementBase{
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
			fmt.Println(lab.lineMapManager.ElementStackTrace())
		}
	}()

	select {
	case <-eb.InputsReady:
		lab.Logger.Log("progress", "starting")
		for _, fun := range funs {
			select {
			case <-lab.Errored:
				return
			default:
				if err := fun(lab); err != nil {
					lab.error(err)
					return
				}
			}
		}
	case <-lab.Errored:
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
	p := filepath.Join(lab.outDir, "elements", fmt.Sprintf("%s.json", eb.element.Name()))
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
