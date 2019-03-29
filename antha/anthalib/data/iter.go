package data

type readRow struct {
	cols          []*ColumnName
	colReader     []iterator
	index         Index
	iteratorCache *seriesIterCache
}

func (rr *readRow) fill(series []*Series) {
	for _, ser := range series {
		rr.cols = append(rr.cols, &ser.col)
		rr.colReader = append(rr.colReader, rr.iteratorCache.Ensure(ser))
	}
}

func (rr *readRow) Value() Row {
	r := Row{Index: rr.index}
	for c, sRead := range rr.colReader {
		r.Values = append(r.Values, Observation{col: rr.cols[c], value: sRead.Value()})
	}
	return r
}

// rawValue gets a row in the raw form (just values without metadata)
func (rr *readRow) rawValue() raw {
	r := make([]interface{}, len(rr.colReader))
	for c, sRead := range rr.colReader {
		r[c] = sRead.Value()
	}
	return r
}

// iterGroup is called to initialize the shared state of multiple related iterators, eg. those that have
// a shared dependency on another node.
type iterGroup struct {
	// TODO interface seems related to sync.Once
	init func() interface{}
}

// nothing here is threadsafe!
type seriesIterCache struct {

	// stores all iterators that need to be advanced together (ie form a node)
	// TODO if this were an ordered map, then it would be possible to make the assumption that dependencies
	// have been advanced in extension columns, allowing us to cache values.
	cache map[*Series]iterator

	// stores the initialized group states.
	groups map[*iterGroup]interface{}
}

func newSeriesIterCache() *seriesIterCache {
	return &seriesIterCache{
		cache:  make(map[*Series]iterator),
		groups: make(map[*iterGroup]interface{}),
	}
}

// EnsureGroup calls the group initializer exactly once
func (c *seriesIterCache) EnsureGroup(ig *iterGroup) interface{} {
	if g, found := c.groups[ig]; found {
		return g
	}
	g := ig.init()
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
