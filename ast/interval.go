package ast

type Interval struct {
	min, max float64
}

func (a *Interval) Meet(x *Interval) *Interval {
	if x == nil && a == nil {
		return nil
	}
	if x == nil {
		return a
	}
	if a == nil {
		return x
	}

	max := a.max
	if max < x.max {
		max = x.max
	}
	min := a.min
	if min > x.min {
		min = x.min
	}
	return &Interval{
		min: min,
		max: max,
	}
}

func (a *Interval) Contains(x *Interval) bool {
	if x == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a.min <= x.min && x.max <= a.max
}

// Create the interval [min, max]
func NewInterval(min, max float64) *Interval {
	if min > max {
		min, max = max, min
	}
	return &Interval{
		min: min,
		max: max,
	}
}

// Create the interval [a, a]
func NewPoint(a float64) *Interval {
	return &Interval{
		min: a,
		max: a,
	}
}
