// Package target provides the construction of a target machine from a
// collection of devices
package target

import (
	"errors"

	"github.com/antha-lang/antha/laboratory/effects"
)

var (
	errNoTarget = errors.New("no target configuration found")
)

const (
	// DriverSelectorV1Name is the basic selector name for device plugins
	// (drivers)
	DriverSelectorV1Name = "antha.driver.v1.TypeReply.type"
)

// Well known device plugins (drivers) selectors
var (
	DriverSelectorV1Human = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.human.v1.Human",
	}
	DriverSelectorV1ShakerIncubator = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.shakerincubator.v1.ShakerIncubator",
	}
	DriverSelectorV1Mixer = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.mixer.v1.Mixer",
	}
	DriverSelectorV1Prompter = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.prompter.v1.Prompter",
	}
	DriverSelectorV1DataSource = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.datasource.v1.DataSource",
	}
	DriverSelectorV1WriteOnlyPlateReader = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.platereader.v1.PlateReader",
	}
	DriverSelectorV1QPCRDevice = effects.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.quantstudio.v1.QuantStudioService",
	}
)

type targetKey int

const theTargetKey targetKey = 0

// Target for execution (collection of devices).
type Target struct {
	Devices []effects.Device
}

// New creates a new target
func New() *Target {
	return &Target{}
}

func (a *Target) canCompile(d effects.Device, reqs ...effects.Request) bool {
	for _, req := range reqs {
		if !d.CanCompile(req) {
			return false
		}
	}
	return true
}

// CanCompile returns the devices that can compile the given set of requests
func (a *Target) CanCompile(reqs ...effects.Request) (r []effects.Device) {
	for _, d := range a.Devices {
		if a.canCompile(d, reqs...) {
			r = append(r, d)
		}
	}
	return
}

// AddDevice adds a device to the target configuration
func (a *Target) AddDevice(d effects.Device) {
	a.Devices = append(a.Devices, d)
}
