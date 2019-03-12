package data

import (
	"github.com/pkg/errors"
)

// Row represents a materialized record.
type Row struct {
	Index  Index
	Values []Observation
}

// Rows are materialized table data, suitable for printing for example.
type Rows struct {
	Data   []Row
	Schema Schema
}

// Observation accesses a column value by name instead of column index.
func (r Row) Observation(c ColumnName) (Observation, error) {
	_, o, err := r.get(c)
	return o, err
}

func (r Row) get(c ColumnName) (int, Observation, error) {
	// TODO more efficiently access schema
	for i, o := range r.Values {
		if o.ColumnName() == c {
			return i, o, nil
		}
	}
	return 0, Observation{}, errors.New("no column " + string(c))
}

// raw returns the row in the raw form (just values without metadata)
func (r Row) raw() raw {
	raw := newRaw(len(r.Values))
	for i, obs := range r.Values {
		raw[i] = obs.Interface()
	}
	return raw
}

// Observation holds an arbitrary, nullable column value.
type Observation struct {
	col   *ColumnName
	value interface{}
}

// ColumnName returns the column name for the value.
func (o Observation) ColumnName() ColumnName {
	return *o.col
}

// IsNull returns true if the value is null.
func (o Observation) IsNull() bool {
	return o.value == nil
}

// Interface returns the underlying representation of the value.
func (o Observation) Interface() interface{} {
	return o.value
}

// raw is a compact internal representation of Row (just values without metadata)
type raw []interface{}

func newRaw(size int) raw {
	return make([]interface{}, size)
}

func (r raw) row(index Index, schema Schema) Row {
	obs := make([]Observation, len(r))
	for i, value := range r {
		obs[i] = Observation{&schema.Columns[i].Name, value}
	}
	return Row{Index: index, Values: obs}
}

// project projects a row according to p
func (r raw) project(p projection) raw {
	newR := newRaw(len(p.newToOld))
	for new, old := range p.newToOld {
		newR[new] = r[old]
	}
	return newR
}

// projection efficiently (without dealing with strings) defines an ordered subset of the set of columns
type projection struct {
	newToOld []int // new column index -> old column index
}

func newProjection(s Schema, cols ...ColumnName) (projection, error) {
	newToOld := make([]int, len(cols))
	for new, col := range cols {
		old, err := s.ColIndex(col)
		if err != nil {
			return projection{}, err
		}
		newToOld[new] = old
	}
	return projection{newToOld: newToOld}, nil
}

func mustNewProjection(s Schema, cols ...ColumnName) projection {
	p, err := newProjection(s, cols...)
	handle(err)
	return p
}
