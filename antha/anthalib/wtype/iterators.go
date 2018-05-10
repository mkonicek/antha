package wtype

//iterators.go includes generalised replacements for plateiterator.go

//AddressIterator iterates through the addresses in an Addressable
type AddressIterator interface {
	Next() WellCoords
	Curr() WellCoords
	MoveTo(WellCoords)
	Reset()
	Valid() bool
}

//AddressSliceIterator iterates through slices of addresses
type AddressSliceIterator interface {
	Next() []WellCoords
	Curr() []WellCoords
	MoveTo(WellCoords)
	Reset()
	Valid() bool
}

type VerticalDirection int

const (
	BottomToTop VerticalDirection = -1
	TopToBottom                   = 1
)

type HorizontalDirection int

const (
	LeftToRight HorizontalDirection = 1
	RightToLeft                     = -1
)

type MajorOrder int

const (
	RowWise MajorOrder = iota
	ColumnWise
)

//GetAddressIterator which iterates through the addresses in addr in order order, moving in directions ver and hor
//when all addresses are returned, resets to the first address if repeat is true, otherwise Valid() returns false
func NewAddressIterator(addr Addressable, order MajorOrder, ver VerticalDirection, hor HorizontalDirection, repeat bool) AddressIterator {
	return newSimpleIterator(addr, order, ver, hor, repeat)
}

//NewColumnIterator get an iterator which iterates through the columns of the addressible, optionally repeating
func NewColumnIterator(addr Addressable, ver VerticalDirection, hor HorizontalDirection, repeat bool) AddressSliceIterator {
	it := newSimpleIterator(addr, ColumnWise, ver, hor, repeat)
	return newChunkedIterator(it, addr.NRows())
}

//NewRowIterator get an iterator which iterates through the columns of the addressible, optionally repeating
func NewRowIterator(addr Addressable, ver VerticalDirection, hor HorizontalDirection, repeat bool) AddressSliceIterator {
	it := newSimpleIterator(addr, RowWise, ver, hor, repeat)
	return newChunkedIterator(it, addr.NCols())
}

type updateFn func(WellCoords) WellCoords

type simpleIterator struct {
	curr   WellCoords
	first  WellCoords
	update updateFn
	reset  bool
	addr   Addressable
}

func newSimpleIterator(addr Addressable, order MajorOrder, ver VerticalDirection, hor HorizontalDirection, repeat bool) *simpleIterator {
	start := WellCoords{}
	if ver == BottomToTop {
		start.Y = addr.NRows() - 1
	}
	if hor == RightToLeft {
		start.X = addr.NCols() - 1
	}

	it := simpleIterator{
		curr:  start,
		first: start,
		reset: repeat,
		addr:  addr,
	}
	if order == RowWise {
		it.update = getRowWiseUpdate(hor, ver, addr)
	} else {
		it.update = getColWiseUpdate(hor, ver, addr)
	}
	return &it
}

func getRowWiseUpdate(hor HorizontalDirection, ver VerticalDirection, a Addressable) updateFn {
	dx := int(hor)
	dy := int(ver)
	nCols := a.NCols()

	if dx > 0 {
		ret := func(wc WellCoords) WellCoords {
			wc.X += dx
			if wc.X >= nCols {
				wc.X -= nCols
				wc.Y += dy
			}
			return wc
		}
		return ret
	}
	ret := func(wc WellCoords) WellCoords {
		wc.X += dx
		if wc.X < 0 {
			wc.X += nCols
			wc.Y += dy
		}
		return wc
	}
	return ret
}

func getColWiseUpdate(hor HorizontalDirection, ver VerticalDirection, a Addressable) updateFn {
	dx := int(hor)
	dy := int(ver)
	nRows := a.NRows()

	if dy > 0 {
		ret := func(wc WellCoords) WellCoords {
			wc.Y += dy
			if wc.Y >= nRows {
				wc.Y -= nRows
				wc.X += dx
			}
			return wc
		}
		return ret
	}
	ret := func(wc WellCoords) WellCoords {
		wc.Y += dy
		if wc.Y < 0 {
			wc.Y += nRows
			wc.X += dx
		}
		return wc
	}
	return ret
}

//Next get the next value in the iterator
func (self *simpleIterator) Next() WellCoords {
	self.curr = self.update(self.curr)
	if self.reset && !self.Valid() {
		self.curr = self.first
	}
	return self.curr
}

//Curr get the current value in the iterator
func (self *simpleIterator) Curr() WellCoords {
	return self.curr
}

//Valid addressable contain the current well coordinates
func (self *simpleIterator) Valid() bool {
	return self.addr.AddressExists(self.curr)
}

func (self *simpleIterator) MoveTo(wc WellCoords) {
	self.curr = wc
}

func (self *simpleIterator) Reset() {
	self.curr = self.first
}

type chunkedIterator struct {
	it       *simpleIterator
	chunkLen int
	curr     []WellCoords
}

func (self *chunkedIterator) Next() []WellCoords {
	i := 0
	for wc := self.it.Curr(); i < self.chunkLen; wc = self.it.Next() {
		self.curr[i] = wc
		i += 1
	}
	return self.curr
}

func (self *chunkedIterator) Curr() []WellCoords {
	return self.curr
}

func (self *chunkedIterator) MoveTo(wc WellCoords) {
	self.it.MoveTo(wc)
}

func (self *chunkedIterator) Reset() {
	self.it.Reset()
}

func (self *chunkedIterator) Valid() bool {
	for _, wc := range self.curr {
		if !self.it.addr.AddressExists(wc) {
			return false
		}
	}
	return true
}

func newChunkedIterator(it *simpleIterator, chunkLen int) *chunkedIterator {
	curr := make([]WellCoords, 0, chunkLen)
	for wc := it.Curr(); len(curr) < chunkLen; wc = it.Next() {
		curr = append(curr, wc)
	}

	return &chunkedIterator{
		it:       it,
		chunkLen: chunkLen,
		curr:     curr,
	}
}
