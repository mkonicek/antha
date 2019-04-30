package target

import (
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/utils"
)

// An Initializer is an instruction with initialization instructions
type Initializer interface {
	GetInitializers() []ast.Inst
}

// A Finalizer is an instruction with finalization instructions
type Finalizer interface {
	GetFinalizers() []ast.Inst
}

// A TimeEstimator is an instruction that can estimate its own execution time
type TimeEstimator interface {
	// GetTimeEstimate returns a time estimate for this instruction in seconds
	GetTimeEstimate() float64
}

// A TipEstimator is an instruction that uses tips and provides information on how many
type TipEstimator interface {
	// GetTipEstimates returns an estimate of how many tips this instruction will use
	GetTipEstimates() []wtype.TipEstimate
}

type dependsMixin struct {
	Depends []ast.Inst
}

// DependsOn implements an Inst
func (a *dependsMixin) DependsOn() []ast.Inst {
	return a.Depends
}

// SetDependsOn implements an Inst
func (a *dependsMixin) SetDependsOn(x ...ast.Inst) {
	a.Depends = x
}

// AppendDependsOn implements an Inst
func (a *dependsMixin) AppendDependsOn(x ...ast.Inst) {
	a.Depends = append(a.Depends, x...)
}

type noDeviceMixin struct{}

// Device implements an Inst
func (a noDeviceMixin) Device() ast.Device {
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
	// Corresponding mix
	Mix              *Mix
	IncubationPlates []*wtype.Plate
}

var (
	_ TimeEstimator = (*Mix)(nil)
	_ Initializer   = (*Mix)(nil)
)

// A Mix is a task that runs a mixer
type Mix struct {
	dependsMixin

	Dev             ast.Device
	Request         *lh.LHRequest
	Properties      *liquidhandling.LHProperties
	FinalProperties *liquidhandling.LHProperties
	Final           map[string]string // Map from ids in Properties to FinalProperties
	Files           Files
	Initializers    []ast.Inst
	Summary         *MixSummary
}

// Device implements an Inst
func (a *Mix) Device() ast.Device {
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

// GetTipEstimates implements a TipEstimator
func (a *Mix) GetTipEstimates() []wtype.TipEstimate {
	ret := []wtype.TipEstimate{}

	if a.Request != nil {
		ret = make([]wtype.TipEstimate, len(a.Request.TipsUsed))
		copy(ret, a.Request.TipsUsed)
	}

	return ret
}

// GetInitializers implements an Initializer
func (a *Mix) GetInitializers() []ast.Inst {
	return a.Initializers
}

// MixSummary contains all the summary information which describes the mix instruction
type MixSummary struct {
	// Layout is a validated JSON summary of the deck layouts before and after
	// the mix operation takes place, suitable for consumption by the front end.
	// The JSON schema is available at microArch/driver/liquidhandling/schemas/layout.schema.json
	Layout []byte
	// Actions is a validated JSON summary of the steps taken during the liquidhandling
	// operation, suitable for consumption by the front end.
	// The JSON schema is available at microArch/driver/liquidhandling/schemas/actions.schema.json
	Actions []byte
}

// NewMixSummary construct a new MixSummary object from the instructions and initial and final robot states
// an error is returned if the parameters are invalid of if either summary object fails JSON-schema validation
func NewMixSummary(itree *liquidhandling.ITree, initial *liquidhandling.LHProperties, final *liquidhandling.LHProperties, idMap map[string]string) (*MixSummary, error) {
	layout, layoutErr := lh.SummarizeLayout(initial, final, idMap)
	actions, actionsErr := lh.SummarizeActions(initial, itree)
	return &MixSummary{
		Layout:  layout,
		Actions: actions,
	}, utils.ErrorSlice{layoutErr, actionsErr}.Pack()
}

// A Manual is human-aided interaction
type Manual struct {
	dependsMixin

	Dev     ast.Device
	Label   string
	Details string
}

// Device implements an Inst
func (a *Manual) Device() ast.Device {
	return a.Dev
}

var (
	_ Finalizer   = (*Run)(nil)
	_ Initializer = (*Run)(nil)
)

// Run calls on device
type Run struct {
	dependsMixin

	Dev     ast.Device
	Label   string
	Details string
	Calls   []driver.Call
	// Additional instructions to add to beginning of instruction stream.
	// Instructions are assumed to depend in FIFO order.
	Initializers []ast.Inst
	// Additional instructions to add to end of instruction stream.
	// Instructions are assumed to depend in LIFO order.
	Finalizers []ast.Inst
}

// Device implements an Inst
func (a *Run) Device() ast.Device {
	return a.Dev
}

// GetInitializers implements an Initializer instruction
func (a *Run) GetInitializers() []ast.Inst {
	return a.Initializers
}

// GetFinalizers implements a Finalizer instruction
func (a *Run) GetFinalizers() []ast.Inst {
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
	noDeviceMixin
	dependsMixin
}

// TimedWait is a wait for a period of time.
type TimedWait struct {
	dependsMixin
	noDeviceMixin

	Duration time.Duration
}
