package data

// Code generated by gen.py. DO NOT EDIT.

import "github.com/pkg/errors"

// MatchFloat64 implements a filter on float64 columns.
type MatchFloat64 func(...float64) bool

// Float64 matches the named column values as float64 arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) Float64(fn MatchFloat64, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeFloat64, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]float64, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(float64)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// Float64 matches the named column values as float64 arguments.
func (o *MustFilterOn) Float64(m MatchFloat64, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.Float64(m, assertions...)
	handle(err)
	return t
}

// MatchInt64 implements a filter on int64 columns.
type MatchInt64 func(...int64) bool

// Int64 matches the named column values as int64 arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) Int64(fn MatchInt64, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeInt64, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]int64, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(int64)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// Int64 matches the named column values as int64 arguments.
func (o *MustFilterOn) Int64(m MatchInt64, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.Int64(m, assertions...)
	handle(err)
	return t
}

// MatchInt implements a filter on int columns.
type MatchInt func(...int) bool

// Int matches the named column values as int arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) Int(fn MatchInt, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeInt, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]int, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(int)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// Int matches the named column values as int arguments.
func (o *MustFilterOn) Int(m MatchInt, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.Int(m, assertions...)
	handle(err)
	return t
}

// MatchString implements a filter on string columns.
type MatchString func(...string) bool

// String matches the named column values as string arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) String(fn MatchString, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeString, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]string, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(string)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// String matches the named column values as string arguments.
func (o *MustFilterOn) String(m MatchString, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.String(m, assertions...)
	handle(err)
	return t
}

// MatchBool implements a filter on bool columns.
type MatchBool func(...bool) bool

// Bool matches the named column values as bool arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) Bool(fn MatchBool, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeBool, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]bool, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(bool)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// Bool matches the named column values as bool arguments.
func (o *MustFilterOn) Bool(m MatchBool, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.Bool(m, assertions...)
	handle(err)
	return t
}

// MatchTimestampMillis implements a filter on TimestampMillis columns.
type MatchTimestampMillis func(...TimestampMillis) bool

// TimestampMillis matches the named column values as TimestampMillis arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) TimestampMillis(fn MatchTimestampMillis, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeTimestampMillis, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]TimestampMillis, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(TimestampMillis)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// TimestampMillis matches the named column values as TimestampMillis arguments.
func (o *MustFilterOn) TimestampMillis(m MatchTimestampMillis, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.TimestampMillis(m, assertions...)
	handle(err)
	return t
}

// MatchTimestampMicros implements a filter on TimestampMicros columns.
type MatchTimestampMicros func(...TimestampMicros) bool

// TimestampMicros matches the named column values as TimestampMicros arguments.
// If any column is nil the filter is automatically false.
// If given any SchemaAssertions, they are called now and may have side effects.
func (o *FilterOn) TimestampMicros(fn MatchTimestampMicros, assertions ...SchemaAssertion) (*Table, error) {
	if err := o.checkSchema(typeTimestampMicros, assertions...); err != nil {
		return nil, errors.Wrapf(err, "can't filter %+v with %+v", o.t, fn)
	}

	projection := mustNewProjection(o.t.schema, o.cols...)

	matchGen := func() rawMatch {
		return func(r raw) bool {
			matchVals := make([]TimestampMicros, len(o.cols))
			for new, old := range projection.newToOld {
				val := r[old]
				if val == nil {
					return false
				}
				matchVals[new] = val.(TimestampMicros)
			}
			return fn(matchVals...)
		}
	}

	return filterTable(matchGen, o.t), nil
}

// TimestampMicros matches the named column values as TimestampMicros arguments.
func (o *MustFilterOn) TimestampMicros(m MatchTimestampMicros, assertions ...SchemaAssertion) *Table {
	t, err := o.FilterOn.TimestampMicros(m, assertions...)
	handle(err)
	return t
}
