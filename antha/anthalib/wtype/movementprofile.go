package wtype

import (
	"github.com/pkg/errors"
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// MovementBehaviour describe movement behaviour in three dimensions
type MovementBehaviour struct {
	Profiles   map[Dimension]LinearMovementProfile
	Order      [][]Dimension //Order in which movement is carried out - dimensions in the same list are carried out simultaneously
	BeforeMove []MovementAction
	AfterMove  []MovementAction
}

// Duration calculate the time taken for the entire move
func (self *MovementBehaviour) DurationToMoveBetween(startPosition, endPosition Coordinates) wunit.Time {
	ret := wunit.NewTime(0.0, "s")

	for _, action := range self.BeforeMove {
		time, pos := action.Duration(startPosition, self)
		ret.IncrBy(time) //nolint
		startPosition = pos
	}

	for _, movementGroup := range self.Order {
		longest := wunit.NewTime(0.0, "s")
		for _, dim := range movementGroup {
			dist := wunit.NewLength(math.Abs(endPosition.Dim(dim)-startPosition.Dim(dim)), "mm")
			time := self.Profiles[dim].GetTimeToTravel(dist)
			if time.GreaterThan(longest) {
				longest = time
			}
		}
		ret.IncrBy(longest) // nolint
	}

	for _, action := range self.AfterMove {
		time, pos := action.Duration(startPosition, self)
		ret.IncrBy(time) // nolint
		startPosition = pos
	}

	return ret
}

// MovementAction an action that is carried out before or after a move operation
type MovementAction interface {
	// Duration calculate the time taken by the action and the final position of the head assembly,
	// given the current location of the head assembly and movement behaviour
	Duration(location Coordinates, behaviour *MovementBehaviour) (wunit.Time, Coordinates)
}

// GenericAction an action which happens in constant time at the begining or end of a move
// for example locking or unlocking a head
type GenericAction struct {
	Time wunit.Time
}

// Duration return the time taken by the action and the final position of the head assembly
func (self *GenericAction) Duration(location Coordinates, behaviour *MovementBehaviour) (wunit.Time, Coordinates) {
	return wunit.NewTime(self.Time.RawValue(), self.Time.Unit().PrefixedSymbol()), location
}

// MoveToSafetyHeightAction move in the Z-Direction to a machine specific safety height to avoid in-transit collisions
type MoveToSafetyHeightAction struct {
	SafetyHeight wunit.Length
}

// Duration return the time taken by the action and the final position of the head assembly
func (self *MoveToSafetyHeightAction) Duration(location Coordinates, behaviour *MovementBehaviour) (wunit.Time, Coordinates) {
	safetyHeightMm := self.SafetyHeight.MustInStringUnit("mm").RawValue()
	location.Z = safetyHeightMm
	return behaviour.Profiles[ZDim].GetTimeToTravel(wunit.NewLength(safetyHeightMm-location.Z, "mm")), location
}

// LinearMovementProfile movement behaviour is one direction only
type LinearMovementProfile interface {
	// GetTimeToTravel calculate the time required to travel the given distance
	// starting and ending at a complete stop
	GetTimeToTravel(wunit.Length) wunit.Time

	// SetVelocity set the maximum velocity for the movement
	SetVelocity(wunit.Velocity) error

	// SetAcceleration set the maximum acceleration for the movement
	SetAcceleration(wunit.Acceleration) error
}

// LinearAcceleration accelerates at constant acceleration at maximum acceleration
// until reaching maximum velocity, continues at maximum velocity, then decelerates
// continuously at maximum acceleration to a stop
type LinearAcceleration struct {
	MinSpeed        wunit.Velocity
	MaxSpeed        wunit.Velocity
	Speed           wunit.Velocity
	MinAcceleration wunit.Acceleration
	MaxAcceleration wunit.Acceleration
	Acceleration    wunit.Acceleration
}

// SetVelocity set the velocity
func (self *LinearAcceleration) SetVelocity(v wunit.Velocity) error {
	if v.LessThan(self.MinSpeed) || v.GreaterThan(self.MaxSpeed) {
		return errors.Errorf("cannot set velocity: %v is outside allowable rage [%v - %v]", v, self.MinSpeed, self.MaxSpeed)
	}
	self.Speed = wunit.NewVelocity(v.RawValue(), v.Unit().PrefixedSymbol())
	return nil
}

// SetAcceleration set the acceleration
func (self *LinearAcceleration) SetAcceleration(v wunit.Acceleration) error {
	if v.LessThan(self.MinAcceleration) || v.GreaterThan(self.MaxAcceleration) {
		return errors.Errorf("cannot set velocity: %v is outside allowable rage [%v - %v]", v, self.MinAcceleration, self.MaxAcceleration)
	}
	self.Acceleration = wunit.NewAcceleration(v.RawValue(), v.Unit().PrefixedSymbol())
	return nil
}

// GetTimeToTravel calculate the time required to travel the given distance
func (self *LinearAcceleration) GetTimeToTravel(distance wunit.Length) wunit.Time {
	//convert everything into SI units
	vMax := self.Speed.MustInStringUnit("m/s").RawValue()
	aMax := self.Acceleration.MustInStringUnit("m/s^2").RawValue()
	distanceM := distance.MustInStringUnit("m").RawValue()

	// for constant acceleration:
	//   \ddot{x} = aMax            (1)
	//    \dot{x} = aMax * t        (2)
	//         x  = aMax * t / 2    (3)
	// let t_1 = time at which \dot{x} = vMax and x = x_1; from (2)
	//   t_1 = vMax / aMax
	// sub t_1 into (3)
	//   x_1 = vMax^2 / (2 * aMax)
	distanceToMaxVelocity := vMax * vMax / (2.0 * aMax)

	//distance is long enough that there's a period of constant velocity in the middle
	if distanceAtConstantVelocity := (distanceM - 2.0*distanceToMaxVelocity); distanceAtConstantVelocity > 0.0 {
		timeForMaxVelocity := vMax / aMax
		timeAtConstantVelocity := distanceAtConstantVelocity / vMax
		return wunit.NewTime(2.0*timeForMaxVelocity+timeAtConstantVelocity, "s")
	} else {
		// from (3) and by symmetry
		return wunit.NewTime(math.Sqrt(2.0*distanceM/aMax), "s")
	}
}
