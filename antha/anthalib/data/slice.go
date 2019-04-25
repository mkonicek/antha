package data

// Table.Slice implementation.
func sliceTable(t *Table, start, end Index) *Table {
	newTable := newFromTable(t, t.sortKey...)
	// the common state generator func
	group := newSeriesGroup(func() seriesGroupStateImpl {
		return &sliceState{
			start:     start,
			end:       end,
			source:    newTableIterator(t.series),
			sourcePos: -1,
		}
	})
	for i, ser := range t.series {
		// metadata
		m := &seriesSlice{
			start:   start,
			end:     end,
			wrapped: ser,
		}
		// slice series
		newTable.series[i] = &Series{
			typ:  ser.typ,
			col:  ser.col,
			read: group.read(i),
			meta: m,
		}
	}
	return newTable
}

// Common state of several slice series. They share the source table iterator (including series iterators cache).
type sliceState struct {
	start, end Index
	source     *tableIterator
	sourcePos  Index
}

func (st *sliceState) Next() bool {
	for st.sourcePos+1 < st.start {
		if !st.source.Next() {
			return false
		}
		st.sourcePos++
	}
	if st.sourcePos+1 == st.end || !st.source.Next() {
		return false
	}
	st.sourcePos++
	return true
}

func (st *sliceState) Value(colIndex int) interface{} {
	return st.source.colReader[colIndex].Value()
}

// Slice series metadata.
type seriesSlice struct {
	start, end Index
	// this is the wrapped, underlying series
	wrapped *Series
}

func (ss *seriesSlice) length() int {
	return int(ss.end - ss.start)
}

func (ss *seriesSlice) IsMaterialized() bool { return false }

func (ss *seriesSlice) ExactSize() int {
	length := ss.length()

	if length == 0 {
		return 0
	}
	if b, ok := ss.wrapped.meta.(boundedMeta); ok {
		w := b.ExactSize()
		if w == -1 {
			return -1
		}
		wrappedLen := w - int(ss.start)
		if wrappedLen < length {
			return wrappedLen
		}
		return length
	}
	return -1
}

func (ss *seriesSlice) MaxSize() int {
	length := ss.length()
	if b, ok := ss.wrapped.meta.(boundedMeta); ok {
		w := b.MaxSize() - int(ss.start)
		if w < length {
			return w
		}
	}
	return length
}

var _ boundedMeta = (*seriesSlice)(nil)
