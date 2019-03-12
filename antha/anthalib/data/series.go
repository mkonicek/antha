package data

import (
	"reflect"

	"github.com/pkg/errors"
)

type advanceable interface {
	Next() bool // false = end iteration
}

// the generic value iterator.  Pays the cost of interface pointer on each value
type iterator interface {
	advanceable
	Value() interface{} // always must be implemented
}

// Series is a named sequence of values. for larger datasets this sequence may
// be loaded lazily (eg memory map) or may even be unbounded
type Series struct {
	col ColumnName
	// typically a scalar type
	typ  reflect.Type
	read func(*seriesIterCache) iterator
	meta seriesMeta
}

// seriesMeta captures differing series backend capabilities
type seriesMeta interface {
	// IsMaterialized = true if the Series is not lazy
	IsMaterialized() bool
}

// boundedMeta is implemented by bounded series metadata
type boundedMeta interface {
	seriesMeta
	// ExactSize can return -1 if size is not known
	ExactSize() int
	// MaxSize should always return >=0
	MaxSize() int
}

// NewSeriesFromSlice converts a slice of scalars to a new Series.
// notNull is a bit mask indicating which values are not null; notNull == nil implies that all values are not null.
func NewSeriesFromSlice(col ColumnName, values interface{}, notNull []bool) (*Series, error) {
	// for now, constructing native series (because it doesn't cause any copying - as opposed to arrow series)
	return newNativeSeriesFromSlice(col, values, notNull)
}

func (s *Series) assignableTo(typ reflect.Type) error {
	if !s.typ.AssignableTo(typ) {
		return errors.Errorf("column %s of type %v cannot be iterated as %v", s.col, s.typ, typ)
	}
	return nil
}

func (s *Series) convertibleTo(typ reflect.Type) error {
	if !s.typ.ConvertibleTo(typ) {
		return errors.Errorf("column %s of type %v cannot be converted to %v", s.col, s.typ, typ)
	}
	return nil
}
