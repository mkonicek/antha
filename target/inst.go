package target

import (
	"github.com/antha-lang/antha/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type Inst interface {
	Device() Device
	DependsOn() []Inst
	SetDependsOn([]Inst)
	GetTimeEstimate() float64
}

type CmpError struct {
	Dev     Device
	Depends []Inst
	Error   error
}

func (a *CmpError) Device() Device {
	return a.Dev
}

func (a *CmpError) DependsOn() []Inst {
	return a.Depends
}

func (a *CmpError) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *CmpError) GetTimeEstimate() float64 {
	return 0.0
}

type Incubate struct {
	Dev     Device
	Depends []Inst
	Files   Files
	Time    float64
}

func (a *Incubate) Device() Device {
	return a.Dev
}

func (a *Incubate) DependsOn() []Inst {
	return a.Depends
}

func (a *Incubate) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *Incubate) GetTimeEstimate() float64 {
	return a.Time
}

// TODO: merge with microArch/report?
type Mix struct {
	Dev             Device
	Depends         []Inst
	Request         *lh.LHRequest
	Properties      *liquidhandling.LHProperties
	FinalProperties *liquidhandling.LHProperties
	Final           map[string]string // Map from ids in Properties to FinalProperties
	Files           Files
}

func (a *Mix) Device() Device {
	return a.Dev
}

func (a *Mix) DependsOn() []Inst {
	return a.Depends
}

func (a *Mix) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *Mix) GetTimeEstimate() float64 {
	est := 0.0

	if a.Request != nil {
		est = a.Request.TimeEstimate
	}

	return est
}

type Manual struct {
	Dev     Device
	Label   string
	Details string
	Depends []Inst
	Time    float64
}

func (a *Manual) DependsOn() []Inst {
	return a.Depends
}

func (a *Manual) Device() Device {
	return a.Dev
}

func (a *Manual) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *Manual) GetTimeEstimate() float64 {
	return a.Time
}

// Run calls on device
type Run struct {
	Dev     Device
	Label   string
	Details string
	Depends []Inst
	Calls   []driver.Call
}

func (a *Run) DependsOn() []Inst {
	return a.Depends
}

func (a *Run) Device() Device {
	return a.Dev
}

func (a *Run) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *Run) GetTimeEstimate() float64 {
	return 0.0
}

// Virtual instruction to hang dependencies on
type Wait struct {
	Depends []Inst
}

func (a *Wait) Device() Device {
	return nil
}

func (a *Wait) DependsOn() []Inst {
	return a.Depends
}

func (a *Wait) SetDependsOn(x []Inst) {
	a.Depends = x
}

func (a *Wait) GetTimeEstimate() float64 {
	return 0.0
}
