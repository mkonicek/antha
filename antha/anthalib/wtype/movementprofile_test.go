package wtype

import (
	"fmt"
	"math"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type LASetVelocityTest struct {
	Velocity    wunit.Velocity
	ShouldError bool
}

func (self *LASetVelocityTest) expectingError(err error) bool {
	return (err != nil) == self.ShouldError
}

func (self *LASetVelocityTest) Run(t *testing.T, la *LinearAcceleration) {
	if err := la.SetVelocity(self.Velocity); !self.expectingError(err) {
		t.Errorf("SetVelocity(%v): expecting error = %t, got error %v", self.Velocity, self.ShouldError, err)
	}
}

type LASetVelocityTests []*LASetVelocityTest

func (self LASetVelocityTests) Run(t *testing.T, la *LinearAcceleration) {
	for _, test := range self {
		test.Run(t, la)
	}
}

type LASetAccelerationTest struct {
	Acceleration wunit.Acceleration
	ShouldError  bool
}

func (self *LASetAccelerationTest) expectingError(err error) bool {
	return (err != nil) == self.ShouldError
}

func (self *LASetAccelerationTest) Run(t *testing.T, la *LinearAcceleration) {
	if err := la.SetAcceleration(self.Acceleration); !self.expectingError(err) {
		t.Errorf("SetAcceleration(%v): expecting error = %t, got error %v", self.Acceleration, self.ShouldError, err)
	}
}

type LASetAccelerationTests []*LASetAccelerationTest

func (self LASetAccelerationTests) Run(t *testing.T, la *LinearAcceleration) {
	for _, test := range self {
		test.Run(t, la)
	}
}

type LAGetTimeToTravelTest struct {
	Distance            wunit.Length
	ExpectedTimeSeconds float64
	Tolerance           float64
}

func (self *LAGetTimeToTravelTest) Run(t *testing.T, la *LinearAcceleration) {
	timeInS := la.GetTimeToTravel(self.Distance).MustInStringUnit("s").RawValue()
	if math.Abs(self.ExpectedTimeSeconds-timeInS) > self.Tolerance {
		t.Errorf("GetTimeToTravel(%v): got %f s expected %f s", self.Distance, timeInS, self.ExpectedTimeSeconds)
	}
}

type LAGetTimeToTravelTests []*LAGetTimeToTravelTest

func (self LAGetTimeToTravelTests) Run(t *testing.T, la *LinearAcceleration) {
	for _, test := range self {
		test.Run(t, la)
	}
}

type LinearAccelerationTest struct {
	Input           *LinearAcceleration
	ShouldError     bool //should initialisation result in an error
	SetVelocity     LASetVelocityTests
	SetAcceleration LASetAccelerationTests
	GetTimeToTravel LAGetTimeToTravelTests
}

func (self *LinearAccelerationTest) expectingError(err error) bool {
	return (err != nil) == self.ShouldError
}

func (self *LinearAccelerationTest) Run(t *testing.T) {
	t.Run(fmt.Sprintf("V=[%v-%v],A=[%v-%v]", self.Input.MinSpeed, self.Input.MaxSpeed, self.Input.MinAcceleration, self.Input.MaxAcceleration), func(t *testing.T) {
		if la, err := NewLinearAcceleration(self.Input.MinSpeed, self.Input.Speed, self.Input.MaxSpeed, self.Input.MinAcceleration, self.Input.Acceleration, self.Input.MaxAcceleration); !self.expectingError(err) {
			t.Errorf("while validating: expecting error %t, got error %v", self.ShouldError, err)
		} else if !self.ShouldError {
			self.SetVelocity.Run(t, la)
			self.SetAcceleration.Run(t, la)
			self.GetTimeToTravel.Run(t, la)
		}
	})
}

type LinearAccelerationTests []*LinearAccelerationTest

func (self LinearAccelerationTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

func TestLinearAcceleration(t *testing.T) {
	LinearAccelerationTests{
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(1.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(1.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(1.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			SetVelocity: LASetVelocityTests{
				{
					Velocity:    wunit.NewVelocity(0.5, "mm/s"),
					ShouldError: true,
				},
				{
					Velocity:    wunit.NewVelocity(10.5, "mm/s"),
					ShouldError: true,
				},
				{
					Velocity: wunit.NewVelocity(5.0, "mm/s"),
				},
			},
			SetAcceleration: LASetAccelerationTests{
				{
					Acceleration: wunit.NewAcceleration(0.5, "mm/s^2"),
					ShouldError:  true,
				},
				{
					Acceleration: wunit.NewAcceleration(10.5, "mm/s^2"),
					ShouldError:  true,
				},
				{
					Acceleration: wunit.NewAcceleration(5.0, "mm/s^2"),
				},
			},
			GetTimeToTravel: LAGetTimeToTravelTests{
				{ //constantly accelerating or decelerating to full speed
					Distance:            wunit.NewLength(5, "mm"),
					ExpectedTimeSeconds: 2.0,
					Tolerance:           1.0e-5,
				},
				{ //constantly accelerating or decelerating to half speed
					Distance:            wunit.NewLength(2.5, "mm"),
					ExpectedTimeSeconds: math.Sqrt(2.0),
					Tolerance:           1.0e-5,
				},
				{ //1 second at constant velocity
					Distance:            wunit.NewLength(10, "mm"),
					ExpectedTimeSeconds: 3.0,
					Tolerance:           1.0e-5,
				},
			},
		},
	}.Run(t)
}
