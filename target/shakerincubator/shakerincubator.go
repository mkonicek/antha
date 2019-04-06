package shakerincubator

import (
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/driver"
	shakerincubator "github.com/antha-lang/antha/driver/antha_shakerincubator_v1"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/handler"
	"github.com/antha-lang/antha/workflow"
)

var (
	_ effects.Device = (*ShakerIncubator)(nil)
)

func NewShakerIncubatorsInstances(tgt *target.Target, config workflow.ShakerIncubatorConfig) error {
	for id := range config.Devices {
		if err := tgt.AddDevice(New(id)); err != nil {
			return err
		}
	}
	return nil
}

// A ShakerIncubator is a device that can shake and incubate things
type ShakerIncubator struct {
	id workflow.DeviceInstanceID
	handler.GenericHandler
}

// New returns a new shaker incubator
func New(id workflow.DeviceInstanceID) *ShakerIncubator {
	si := &ShakerIncubator{
		id: id,
	}
	si.GenericHandler = handler.GenericHandler{
		Labels: []effects.NameValue{
			target.DriverSelectorV1ShakerIncubator,
		},
		GenFunc: si.generate,
	}
	return si
}

func (a *ShakerIncubator) Id() workflow.DeviceInstanceID {
	return a.id
}

func (a *ShakerIncubator) carrierOpen() driver.Call {
	return driver.Call{
		Method: "/antha.shakerincubator.v1.ShakerIncubator/CarrierOpen",
		Args:   &shakerincubator.Blank{},
		Reply:  &shakerincubator.BoolReply{},
	}
}

func (a *ShakerIncubator) carrierClose() driver.Call {
	return driver.Call{
		Method: "/antha.shakerincubator.v1.ShakerIncubator/CarrierClose",
		Args:   &shakerincubator.Blank{},
		Reply:  &shakerincubator.BoolReply{},
	}
}

func (a *ShakerIncubator) reset() []driver.Call {
	return []driver.Call{
		{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/ShakeStop",
			Args:   &shakerincubator.Blank{},
			Reply:  &shakerincubator.BoolReply{},
		},
		{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/TemperatureReset",
			Args:   &shakerincubator.Blank{},
			Reply:  &shakerincubator.BoolReply{},
		},
		{
			Method: "/antha.shakerincubator.v1.ShakerIncubator/CarrierOpen",
			Args:   &shakerincubator.Blank{},
			Reply:  &shakerincubator.BoolReply{},
		},
	}
}

func (a *ShakerIncubator) temperatureSet(temp wunit.Temperature) driver.Call {
	return driver.Call{
		Method: "/antha.shakerincubator.v1.ShakerIncubator/TemperatureSet",
		Args: &shakerincubator.TemperatureSettings{
			Temperature: temp.RawValue(), // in C
		},
		Reply: &shakerincubator.BoolReply{},
	}
}

func (a *ShakerIncubator) shakeStart(rate wunit.Rate, length wunit.Length) driver.Call {
	if length.IsNil() {
		length = wunit.NewLength(3.0/1000.0, "m")
	}
	return driver.Call{
		Method: "/antha.shakerincubator.v1.ShakerIncubator/ShakeStart",
		Args: &shakerincubator.ShakerSettings{
			Frequency: rate.SIValue(),
			Radius:    length.SIValue(),
		},
		Reply: &shakerincubator.BoolReply{},
	}
}

func (a *ShakerIncubator) generate(cmd interface{}) (effects.Insts, error) {
	inc, ok := cmd.(*effects.IncubateInst)
	if !ok {
		return nil, fmt.Errorf("expecting %T found %T instead", inc, cmd)
	}

	initializers := effects.Insts{
		&target.Run{
			DeviceMixin: effects.DeviceMixin{Dev: a},
			Label:       "open incubator carrier",
			Calls: []driver.Call{
				a.carrierOpen(),
			},
		},

		&target.Prompt{
			Message: "close incubator carrier?",
		},

		&target.Run{
			DeviceMixin: effects.DeviceMixin{Dev: a},
			Label:       "close incubator carrier",
			Calls: []driver.Call{
				a.carrierClose(),
			},
		},
	}

	finalizers := effects.Insts{
		&target.Run{
			DeviceMixin: effects.DeviceMixin{Dev: a},
			Label:       "turn off incubator",
			Calls:       a.reset(),
		},
	}

	var insts effects.Insts
	if !inc.PreTime.IsNil() {
		var calls []driver.Call
		if !inc.PreTemp.IsNil() {
			calls = append(calls, a.temperatureSet(inc.PreTemp))
		}
		if !inc.PreShakeRate.IsNil() {
			calls = append(calls, a.shakeStart(inc.PreShakeRate, inc.PreShakeRadius))
		}

		insts = append(insts,
			&target.Run{
				DeviceMixin: effects.DeviceMixin{Dev: a},
				Label:       "pre incubate",
				Calls:       calls,
			},
			&target.TimedWait{
				Duration: time.Duration(inc.PreTime.Seconds() * float64(time.Second)),
			},
		)
	}

	var calls []driver.Call
	if !inc.Temp.IsNil() {
		calls = append(calls, a.temperatureSet(inc.Temp))
	}
	if !inc.ShakeRate.IsNil() {
		calls = append(calls, a.shakeStart(inc.ShakeRate, inc.ShakeRadius))
	}

	insts = append(insts, &target.Run{
		DeviceMixin:  effects.DeviceMixin{Dev: a},
		Label:        "start incubator",
		Calls:        calls,
		Initializers: initializers,
		Finalizers:   finalizers,
	})

	if !inc.Time.IsNil() {
		insts = append(insts, &target.TimedWait{
			Duration: time.Duration(inc.Time.Seconds() * float64(time.Second)),
		})
	}

	insts.SequentialOrder()
	return insts, nil
}
