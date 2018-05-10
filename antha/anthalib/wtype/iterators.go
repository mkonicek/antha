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

type UpdateFn func(WellCoords) WellCoords

type simpleIterator struct {
	curr   WellCoords
	first  WellCoords
	update UpdateFn
	reset  bool
	addr   Addressable
}

func getRowWiseUpdate(hor HorizontalDirection, ver VerticalDirection, a Addressable) UpdateFn {
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

func getColWiseUpdate(hor HorizontalDirection, ver VerticalDirection, a Addressable) UpdateFn {
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

//GetAddressIterator which iterates through the addresses in addr in order order, moving in directions ver and hor
//when all addresses are returned, resets to the first address if reset is true, otherwise Valid() returns false
func GetAddressIterator(addr Addressable, order MajorOrder, ver VerticalDirection, hor HorizontalDirection, reset bool) AddressIterator {
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
		reset: reset,
		addr:  addr,
	}
	if order == RowWise {
		it.update = getRowWiseUpdate(hor, ver, addr)
	} else {
		it.update = getColWiseUpdate(hor, ver, addr)
	}
	return &it
}
