package wtype

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func TestMovementbehaviourJSONSerialise(t *testing.T) {

	wait1, err := NewGenericAction(wunit.NewTime(5.0, "min"))
	if err != nil {
		t.Fatal(err)
	}
	wait2, err := NewGenericAction(wunit.NewTime(2.0, "ms"))
	if err != nil {
		t.Fatal(err)
	}

	safety, err := NewMoveToSafetyHeightAction(wunit.NewLength(104, "mm"))
	if err != nil {
		t.Fatal(err)
	}

	x, err := NewLinearAcceleration(wunit.NewVelocity(0.1, "mm/s"), wunit.NewVelocity(15, "mm/s"), wunit.NewVelocity(100, "mm/s"), wunit.NewAcceleration(0.0, "mm/s^2"), wunit.NewAcceleration(50, "mm/s^2"), wunit.NewAcceleration(500, "mm/s^2"))
	if err != nil {
		t.Fatal(err)
	}

	y, err := NewLinearAcceleration(wunit.NewVelocity(0.1, "mm/s"), wunit.NewVelocity(20, "mm/s"), wunit.NewVelocity(100, "mm/s"), wunit.NewAcceleration(0.0, "mm/s^2"), wunit.NewAcceleration(60, "mm/s^2"), wunit.NewAcceleration(500, "mm/s^2"))
	if err != nil {
		t.Fatal(err)
	}

	z, err := NewLinearAcceleration(wunit.NewVelocity(0.1, "mm/s"), wunit.NewVelocity(3, "mm/s"), wunit.NewVelocity(10, "mm/s"), wunit.NewAcceleration(0.0, "mm/s^2"), wunit.NewAcceleration(20, "mm/s^2"), wunit.NewAcceleration(50, "mm/s^2"))
	if err != nil {
		t.Fatal(err)
	}

	var b MovementBehaviour

	if a, err := NewMovementBehaviour(x, y, z, [][]Dimension{{XDim, YDim}, {ZDim}}, []MovementAction{wait2, safety}, []MovementAction{wait1}); err != nil {
		t.Fatal(err)
	} else if bytes, err := json.Marshal(a); err != nil {
		t.Fatal(err)
	} else if err := json.Unmarshal(bytes, &b); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(a, &b) {
		t.Errorf("deserialised MovementBehaviour didn't match:\n e: %v\n g:%v", a, b)
	}
}

type DurationToMoveBetweenTest struct {
	Initial                 Coordinates
	Final                   Coordinates
	ExpectedDurationSeconds float64
	Tolerance               float64
}

func (self *DurationToMoveBetweenTest) Run(t *testing.T, mb *MovementBehaviour) {
	duration := mb.DurationToMoveBetween(self.Initial, self.Final)
	if got := duration.MustInStringUnit("s").RawValue(); math.Abs(got-self.ExpectedDurationSeconds) > self.Tolerance {
		t.Errorf("DurationToMoveBetween(%v, %v): got %v, expected %f s", self.Initial, self.Final, duration, self.ExpectedDurationSeconds)
	}
}

type DurationToMoveBetweenTests []DurationToMoveBetweenTest

func (self DurationToMoveBetweenTests) Run(t *testing.T, mb *MovementBehaviour) {
	for _, test := range self {
		test.Run(t, mb)
	}
}

type MovementBehaviourTest struct {
	Name        string //since string of MovementBehaviour would be too long
	Input       *MovementBehaviour
	ShouldError bool //true if initialisation should raise an error
	Tests       DurationToMoveBetweenTests
}

func (self *MovementBehaviourTest) expectingError(err error) bool {
	return (err != nil) == self.ShouldError
}

func (self *MovementBehaviourTest) Run(t *testing.T) {
	t.Run(self.Name, func(t *testing.T) {
		if mb, err := NewMovementBehaviour(self.Input.Profiles[XDim], self.Input.Profiles[YDim], self.Input.Profiles[ZDim], self.Input.Order, self.Input.BeforeMove, self.Input.AfterMove); !self.expectingError(err) {
			t.Errorf("in constructor: expecting error %t, got error %v", self.ShouldError, err)
			return
		} else if !self.ShouldError {
			self.Tests.Run(t, mb)
		}
	})
}

type MovementBehaviourTests []MovementBehaviourTest

func (self MovementBehaviourTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

func TestMovementBehaviour(t *testing.T) {
	profile, err := NewLinearAcceleration(
		wunit.NewVelocity(1.0, "mm/s"),
		wunit.NewVelocity(5.0, "mm/s"),
		wunit.NewVelocity(10.0, "mm/s"),
		wunit.NewAcceleration(1.0, "mm/s^2"),
		wunit.NewAcceleration(5.0, "mm/s^2"),
		wunit.NewAcceleration(10.0, "mm/s^2"),
	)
	if err != nil {
		t.Fatal(err)
	}

	waitOneSec, err := NewGenericAction(wunit.NewTime(1.0, "s"))
	if err != nil {
		t.Fatal(err)
	}

	waitTwoSec, err := NewGenericAction(wunit.NewTime(2.0, "s"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewGenericAction(wunit.NewTime(-2.0, "s")); err == nil {
		t.Error("no error produced for GenericAction with negative time")
	}

	moveToSafeHeight, err := NewMoveToSafetyHeightAction(wunit.NewLength(5, "mm"))
	if err != nil {
		t.Fatal(err)
	}

	MovementBehaviourTests{
		{
			Name: "missing Y profile",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, nil, profile},
				Order:    [][]Dimension{{XDim}, {YDim}, {ZDim}},
			},
			ShouldError: true,
		},
		{
			Name: "missing XDim",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, profile, profile},
				Order:    [][]Dimension{{YDim}, {ZDim}},
			},
			ShouldError: true,
		},
		{
			Name: "repeated ZDim",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, profile, profile},
				Order:    [][]Dimension{{XDim, ZDim}, {YDim}, {ZDim}},
			},
			ShouldError: true,
		},
		{
			Name: "several ordering errors",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, profile, profile},
				Order:    [][]Dimension{{XDim, ZDim}, {ZDim}},
			},
			ShouldError: true,
		},
		{
			Name: "OK simple",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, profile, profile},
				Order:    [][]Dimension{{XDim}, {YDim}, {ZDim}},
			},
			Tests: DurationToMoveBetweenTests{
				{
					Initial: Coordinates{},
					Final:   Coordinates{5, 5, 5},
					ExpectedDurationSeconds: 6, //movements carried out separately so 2s each
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{},
					Final:   Coordinates{5, 2.5, 5},
					ExpectedDurationSeconds: 4 + math.Sqrt(2),
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{},
					Final:   Coordinates{},
					ExpectedDurationSeconds: 0,
					Tolerance:               1e-5,
				},
			},
		},
		{
			Name: "OK XY Simultaneous",
			Input: &MovementBehaviour{
				Profiles: MovementProfile{profile, profile, profile},
				Order:    [][]Dimension{{XDim, YDim}, {ZDim}},
			},
			Tests: DurationToMoveBetweenTests{
				{
					Initial: Coordinates{},
					Final:   Coordinates{5, 5, 5},
					ExpectedDurationSeconds: 4,
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{},
					Final:   Coordinates{5, 2.5, 5},
					ExpectedDurationSeconds: 4, //shorter distance has no effect as waiting for X
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{},
					Final:   Coordinates{},
					ExpectedDurationSeconds: 0,
					Tolerance:               1e-5,
				},
			},
		},
		{
			Name: "OK with Generic",
			Input: &MovementBehaviour{
				Profiles:   MovementProfile{profile, profile, profile},
				Order:      [][]Dimension{{XDim, YDim}, {ZDim}},
				BeforeMove: []MovementAction{waitOneSec},
				AfterMove:  []MovementAction{waitTwoSec},
			},
			Tests: DurationToMoveBetweenTests{
				{
					Initial: Coordinates{},
					Final:   Coordinates{X: 5, Y: 5, Z: 5},
					ExpectedDurationSeconds: 7,
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{X: 10, Y: 10, Z: 10},
					Final:   Coordinates{X: 10, Y: 10, Z: 10},
					ExpectedDurationSeconds: 3, //assume that no machines interperet "move to current pos" as noop and not run actions...
					Tolerance:               1e-5,
				},
			},
		},
		{
			Name: "OK with SafetyHeight",
			Input: &MovementBehaviour{
				Profiles:   MovementProfile{profile, profile, profile},
				Order:      [][]Dimension{{XDim, YDim}, {ZDim}},
				BeforeMove: []MovementAction{waitOneSec, moveToSafeHeight},
				AfterMove:  []MovementAction{waitTwoSec},
			},
			Tests: DurationToMoveBetweenTests{
				{
					Initial: Coordinates{},
					Final:   Coordinates{X: 5, Y: 5},
					ExpectedDurationSeconds: 9,
					Tolerance:               1e-5,
				},
				{
					Initial: Coordinates{},
					Final:   Coordinates{},
					ExpectedDurationSeconds: 7,
					Tolerance:               1e-5,
				},
			},
		},
	}.Run(t)
}

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

type LATimeToTravelBetweenTest struct {
	Start               wunit.Length
	End                 wunit.Length
	ExpectedTimeSeconds float64
	Tolerance           float64
}

func (self *LATimeToTravelBetweenTest) Run(t *testing.T, la *LinearAcceleration) {
	timeInS := la.TimeToTravelBetween(self.Start, self.End).MustInStringUnit("s").RawValue()
	if math.Abs(self.ExpectedTimeSeconds-timeInS) > self.Tolerance {
		t.Errorf("TimeToTravelBetween(%v, %v): got %f s expected %f s", self.Start, self.End, timeInS, self.ExpectedTimeSeconds)
	}
}

type LATimeToTravelBetweenTests []*LATimeToTravelBetweenTest

func (self LATimeToTravelBetweenTests) Run(t *testing.T, la *LinearAcceleration) {
	for _, test := range self {
		test.Run(t, la)
	}
}

type LinearAccelerationTest struct {
	Input               *LinearAcceleration
	ShouldError         bool //should initialisation result in an error
	SetVelocity         LASetVelocityTests
	SetAcceleration     LASetAccelerationTests
	TimeToTravelBetween LATimeToTravelBetweenTests
}

func (self *LinearAccelerationTest) expectingError(err error) bool {
	return (err != nil) == self.ShouldError
}

func (self *LinearAccelerationTest) Run(t *testing.T) {
	t.Run(self.Input.String(), func(t *testing.T) {
		if la, err := NewLinearAcceleration(self.Input.MinSpeed, self.Input.Speed, self.Input.MaxSpeed, self.Input.MinAcceleration, self.Input.Acceleration, self.Input.MaxAcceleration); !self.expectingError(err) {
			t.Errorf("while validating: expecting error %t, got error %v", self.ShouldError, err)
		} else if !self.ShouldError {
			self.SetVelocity.Run(t, la)
			self.SetAcceleration.Run(t, la)
			self.TimeToTravelBetween.Run(t, la)
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
				Speed:           wunit.NewVelocity(0.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(1.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(1.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(1.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(1.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(3.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(1.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(3.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(1.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(1.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(-1.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(3.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(1.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(1.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(11.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(-2.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(0.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(0.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			ShouldError: true,
		},
		{
			Input: &LinearAcceleration{
				MinSpeed:        wunit.NewVelocity(0.0, "mm/s"),
				Speed:           wunit.NewVelocity(4.0, "mm/s"),
				MaxSpeed:        wunit.NewVelocity(10.0, "mm/s"),
				MinAcceleration: wunit.NewAcceleration(0.0, "mm/s^2"),
				Acceleration:    wunit.NewAcceleration(1.0, "mm/s^2"),
				MaxAcceleration: wunit.NewAcceleration(10.0, "mm/s^2"),
			},
			SetVelocity: LASetVelocityTests{
				{
					Velocity:    wunit.NewVelocity(-0.5, "mm/s"),
					ShouldError: true,
				},
				{
					Velocity:    wunit.NewVelocity(10.5, "mm/s"),
					ShouldError: true,
				},
				{
					Velocity:    wunit.NewVelocity(0, "mm/s"),
					ShouldError: true,
				},
				{
					Velocity: wunit.NewVelocity(5.0, "mm/s"),
				},
			},
			SetAcceleration: LASetAccelerationTests{
				{
					Acceleration: wunit.NewAcceleration(0, "mm/s^2"),
					ShouldError:  true,
				},
				{
					Acceleration: wunit.NewAcceleration(-0.5, "mm/s^2"),
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
			TimeToTravelBetween: LATimeToTravelBetweenTests{
				{ //constantly accelerating or decelerating to full speed
					Start:               wunit.NewLength(5, "mm"),
					End:                 wunit.NewLength(10, "mm"),
					ExpectedTimeSeconds: 2.0,
					Tolerance:           1.0e-5,
				},
				{ //constantly accelerating or decelerating to half speed
					Start:               wunit.NewLength(2.5, "mm"),
					End:                 wunit.NewLength(5, "mm"),
					ExpectedTimeSeconds: math.Sqrt(2.0),
					Tolerance:           1.0e-5,
				},
				{ //1 second at constant velocity
					Start:               wunit.NewLength(10, "mm"),
					End:                 wunit.NewLength(20, "mm"),
					ExpectedTimeSeconds: 3.0,
					Tolerance:           1.0e-5,
				},
			},
		},
	}.Run(t)
}
