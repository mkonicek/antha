package data

import (
	"reflect"
)

type conversion struct {
	*Table
	newType reflect.Type
}

func (c *conversion) convert(s *Series) (*Series, error) {
	if err := s.convertibleTo(c.newType); err != nil {
		return nil, err
	}

	return &Series{
		col:  s.col,
		typ:  c.newType,
		meta: newExtendSeriesMeta([]*Series{s}, false),
		read: func(cache *seriesIterCache) iterator {
			return &convertSeries{newType: c.newType, wrapped: cache.Ensure(s)}
		},
	}, nil
}

type convertSeries struct {
	newType reflect.Type
	wrapped iterator
}

// rely on iterator cache side effects for why we are not advancing here
func (i *convertSeries) Next() bool { return true }

func (i *convertSeries) Value() interface{} {
	value := i.wrapped.Value()
	if value == nil {
		return nil
	}
	return reflect.ValueOf(value).Convert(i.newType).Interface()
}
