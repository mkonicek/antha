package data

type readRow struct {
	schema        *Schema
	colReader     []iterator
	index         Index
	iteratorCache *seriesIterCache
}

func (rr *readRow) fill(series []*Series) {
	rr.schema = newSchema(series)
	for _, ser := range series {
		rr.colReader = append(rr.colReader, rr.iteratorCache.Ensure(ser))
	}
}

func (rr *readRow) Value() Row {
	return Row{
		index:  rr.index,
		schema: rr.schema,
		values: rr.rawValue(),
	}
}

// rawValue gets a row in the raw form (just values without metadata)
func (rr *readRow) rawValue() raw {
	r := make([]interface{}, len(rr.colReader))
	for c, sRead := range rr.colReader {
		r[c] = sRead.Value()
	}
	return r
}

// nothing here is threadsafe!
type seriesIterCache struct {

	// stores all iterators that need to be advanced together (ie form a node)
	// TODO if this were an ordered map, then it would be possible to make the assumption that dependencies
	// have been advanced in extension columns, allowing us to cache values.
	cache map[*Series]iterator

	// stores the initialized group states.
	groups map[*seriesGroup]*seriesGroupState
}

func newSeriesIterCache() *seriesIterCache {
	return &seriesIterCache{
		cache:  make(map[*Series]iterator),
		groups: make(map[*seriesGroup]*seriesGroupState),
	}
}

// EnsureGroup calls the group initializer exactly once
func (c *seriesIterCache) EnsureGroup(ig *seriesGroup) *seriesGroupState {
	if g, found := c.groups[ig]; found {
		return g
	}
	g := ig.createState()
	c.groups[ig] = g
	return g
}

func (c *seriesIterCache) Advance() bool {
	if len(c.cache) == 0 {
		return false
	}
	for _, sRead := range c.cache {
		if !sRead.Next() {
			return false
		}
	}
	return true
}

func (c *seriesIterCache) Ensure(ser *Series) iterator {
	if seriesRead, found := c.cache[ser]; found {
		return seriesRead
	}

	seriesRead := ser.read(c)
	c.cache[ser] = seriesRead
	return seriesRead
}

type tableIterator struct {
	readRow
}

func newTableIterator(series []*Series) *tableIterator {
	iter := &tableIterator{readRow{
		index:         -1,
		iteratorCache: newSeriesIterCache(),
	}}
	iter.fill(series)
	return iter
}

func (iter *tableIterator) Next() bool {
	// all series we depend on need to be advanced
	if !iter.iteratorCache.Advance() {
		return false
	}
	iter.index++
	return true
}

/*
 * generic tools for implementing operations which need series group common state
 */

// While implementing an operation which requires shared iterators state (e.g. Filter or Append) one should do the following:
// - define a shared state object: it should implement `seriesGroupStateImpl` - i.e. make it possible to iterate over a shared data source
// - create a `seriesGroup` object (eagerly - at the time when the operation is requested) and supply it with a shared state constructor
//   (which, by contrast, will be called lazily - i.e. only when iteration starts)
// - use `seriesGroup.read` as a read function for your output series.

// seriesGroup represents a group of series which should be iterated together, eg. those that have
// a shared dependency on another node.
type seriesGroup struct {
	createStateImpl func() seriesGroupStateImpl
}

func newSeriesGroup(createStateImpl func() seriesGroupStateImpl) *seriesGroup {
	return &seriesGroup{
		createStateImpl: createStateImpl,
	}
}

func (ig *seriesGroup) createState() *seriesGroupState {
	return &seriesGroupState{
		impl:    ig.createStateImpl(),
		wasNext: true,
		pos:     -1,
	}
}

func (ig *seriesGroup) read(colIndex int) func(cache *seriesIterCache) iterator {
	return func(cache *seriesIterCache) iterator {
		return &seriesGroupIter{
			commonState: cache.EnsureGroup(ig),
			colIndex:    colIndex,
			pos:         -1,
		}
	}
}

// seriesGroupState is a shared state for iteration over a series group.
// Uses the common state interface implementation provided by an operation + stores the result of the latest Next() call and calculates iteration position.
type seriesGroupState struct {
	impl    seriesGroupStateImpl
	wasNext bool
	pos     Index
}

func (st *seriesGroupState) next() {
	st.wasNext = st.impl.Next()
	if st.wasNext {
		st.pos++
	}
}

// Operation-specific series group shared state methods - should be implemented by each operation itself.
// TODO: this looks closely related to tableIterator; does it make any sense to create a common interface for them?
type seriesGroupStateImpl interface {
	// advances the common state
	Next() bool
	// gives access to columns values at the current iteration position
	Value(colIndex int) interface{}
}

// seriesGroupIter is an iterator over one of the series of a series group.
// All such iterators refer to the same seriesGroupState.
type seriesGroupIter struct {
	commonState *seriesGroupState
	colIndex    int
	pos         Index
}

// Next advances the common state provided that iter.pos is already up to date.
func (iter *seriesGroupIter) Next() bool {
	// see if we need to discard the current shared state
	retain := iter.pos != iter.commonState.pos
	if !retain {
		iter.commonState.next()
		iter.pos = iter.commonState.pos
	}
	return iter.commonState.wasNext
}

// Value reads the column value.
func (iter *seriesGroupIter) Value() interface{} {
	return iter.commonState.impl.Value(iter.colIndex)
}
