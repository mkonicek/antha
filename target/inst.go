package target

import (
	"time"

	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

// An Inst is a instruction
type Inst interface {
	// Device that this instruction was generated for
	Device() Device
	// DependsOn returns instructions that this instruction depends on
	DependsOn() []Inst
	// SetDependsOn updates DependsOn
	SetDependsOn([]Inst)
}

// An Initializer is an instruction with initialization instructions
type Initializer interface {
	GetInitializers() []Inst
}

// A Finalizer is an instruction with finalization instructions
type Finalizer interface {
	GetFinalizers() []Inst
}

// A TimeEstimator is an instruction that can estimate its own execution time
type TimeEstimator interface {
	// GetTimeEstimate returns a time estimate for this instruction in seconds
	GetTimeEstimate() float64
}

type dependsMixin struct {
	Depends []Inst
}

// DependsOn implements an Inst
func (a *dependsMixin) DependsOn() []Inst {
	return a.Depends
}

// SetDependsOn implements an Inst
func (a *dependsMixin) SetDependsOn(x []Inst) {
	a.Depends = x
}

type noDeviceMixin struct{}

// Device implements an Inst
func (a noDeviceMixin) Device() Device {
	return nil
}

// An Order is a task to order physical components
type Order struct {
	Manual
	Mixes []*Mix
}

// A PlatePrep is a task to setup plates
type PlatePrep struct {
	Manual
	Mixes []*Mix
}

// A SetupMixer is a task to setup a mixer
type SetupMixer struct {
	Manual
	Mixes []*Mix
}

// A SetupIncubator is a task to setup an incubator
type SetupIncubator struct {
	Manual
}

var (
	_ TimeEstimator = (*Mix)(nil)
	_ Initializer   = (*Mix)(nil)
)

// A Mix is a task that runs a mixer
type Mix struct {
	dependsMixin

	Dev             Device
	Request         *lh.LHRequest
	Properties      *liquidhandling.LHProperties
	FinalProperties *liquidhandling.LHProperties
	Final           map[string]string // Map from ids in Properties to FinalProperties
	Files           Files
	Initializers    []Inst
}

// Device implements an Inst
func (a *Mix) Device() Device {
	return a.Dev
}

// GetTimeEstimate implements a TimeEstimator
func (a *Mix) GetTimeEstimate() float64 {
	est := 0.0

	if a.Request != nil {
		est = a.Request.TimeEstimate
	}

	return est
}

// GetInitializers implements an Initializer
func (a *Mix) GetInitializers() []Inst {
	return a.Initializers
}

// A Manual is human-aided interaction
type Manual struct {
	dependsMixin

	Dev     Device
	Label   string
	Details string
}

// Device implements an Inst
func (a *Manual) Device() Device {
	return a.Dev
}

var (
	_ Finalizer   = (*Run)(nil)
	_ Initializer = (*Run)(nil)
)

// Run calls on device
type Run struct {
	dependsMixin

	Dev     Device
	Label   string
	Details string
	Calls   []driver.Call
	// Additional instructions to add to beginning of instruction stream.
	// Instructions are assumed to depend in FIFO order.
	Initializers []Inst
	// Additional instructions to add to end of instruction stream.
	// Instructions are assumed to depend in LIFO order.
	Finalizers []Inst
}

// Device implements an Inst
func (a *Run) Device() Device {
	return a.Dev
}

// GetInitializers implements an Initializer instruction
func (a *Run) GetInitializers() []Inst {
	return a.Initializers
}

// GetFinalizers implements a Finalizer instruction
func (a *Run) GetFinalizers() []Inst {
	return a.Finalizers
}

// Prompt is manual prompt instruction
type Prompt struct {
	dependsMixin
	noDeviceMixin

	Message string
}

// Wait is a virtual instruction to hang dependencies on. A better name might
// been no-op.
type Wait struct {
	dependsMixin
	noDeviceMixin
}

// TimedWait is a wait for a period of time.
type TimedWait struct {
	dependsMixin
	noDeviceMixin

	Duration time.Duration
}

// SequentialOrder takes a set of instructions with out any dependencies and
// modifies them to follow sequential order
func SequentialOrder(insts ...Inst) []Inst {
	for idx, inst := range insts {
		if idx == 0 {
			continue
		}
		inst.SetDependsOn([]Inst{insts[idx-1]})
	}

	return insts
}
