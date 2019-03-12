package data

type slicer struct {
	*seriesSlice
	iterator
	pos         Index
	commonState *seriesIterCache
}

func (iter *slicer) Next() bool {
	for iter.pos+1 < iter.start {
		if !iter.iterator.Next() {
			return false
		}
		iter.pos++
	}
	iter.pos++
	if iter.pos == iter.end {
		return false
	}
	// see if we exhausted the underlying series
	return iter.iterator.Next()
}

type seriesSlice struct {
	start, end Index
	group      *iterGroup
	// this is the wrapped, underlying series
	wrapped *Series
}

func (ss *seriesSlice) length() int {
	return int(ss.end - ss.start)
}

func (ss *seriesSlice) read(cache *seriesIterCache) iterator {
	sl := &slicer{
		seriesSlice: ss,
		pos:         -1,
		commonState: cache.EnsureGroup(ss.group).(*seriesIterCache),
	}
	// the wrapped iterator is placed in nested cache
	sl.iterator = sl.commonState.Ensure(ss.wrapped)
	return sl
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
