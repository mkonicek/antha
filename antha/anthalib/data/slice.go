package data

// Table.Slice implementation.
func sliceTable(t *Table, start, end Index) *Table {
	newTable := newFromTable(t, t.sortKey...)
	// the common state generator func
	group := &iterGroup{func() interface{} {
		return &sliceState{
			start:     start,
			end:       end,
			source:    newTableIterator(t.series),
			sourcePos: -1,
		}
	}}
	for i, ser := range t.series {
		// metadata
		m := &seriesSlice{
			start:    start,
			end:      end,
			colIndex: i,
			wrapped:  ser,
			group:    group,
		}
		// slice series
		newTable.series[i] = &Series{
			typ:  ser.typ,
			col:  ser.col,
			read: m.read,
			meta: m,
		}
	}
	return newTable
}

// Common state of several slice series. They share the source table iterator (including series iterators cache).
type sliceState struct {
	start, end Index
	source     *tableIterator
	sourceNext bool
	sourcePos  Index
}

func (st *sliceState) advance() {
	for st.sourcePos+1 < st.start {
		if !st.source.Next() {
			st.sourceNext = false
			return
		}
		st.sourcePos++
	}
	if st.sourcePos+1 == st.end || !st.source.Next() {
		st.sourceNext = false
		return
	}
	st.sourcePos++
	st.sourceNext = true
}

func (st *sliceState) pos() Index {
	if st.sourcePos == -1 {
		return -1
	}
	return st.sourcePos - st.start
}

// Slice series iterator.
type sliceIter struct {
	commonState *sliceState
	pos         Index
	colIndex    int
}

// Next advances the common state provided that iter.pos is already up to date.
func (iter *sliceIter) Next() bool {
	// see if we need to discard the current shared state
	retain := iter.pos != iter.commonState.pos()
	if !retain {
		iter.commonState.advance()
		iter.pos = iter.commonState.pos()
	}
	return iter.commonState.sourceNext
}

// Value reads the cached column value.
func (iter *sliceIter) Value() interface{} {
	return iter.commonState.source.colReader[iter.colIndex].Value()
}

// Slice series metadata.
type seriesSlice struct {
	start, end Index
	colIndex   int
	group      *iterGroup
	// this is the wrapped, underlying series
	wrapped *Series
}

func (ss *seriesSlice) length() int {
	return int(ss.end - ss.start)
}

func (ss *seriesSlice) read(cache *seriesIterCache) iterator {
	return &sliceIter{
		commonState: cache.EnsureGroup(ss.group).(*sliceState),
		pos:         -1,
		colIndex:    ss.colIndex,
	}
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
