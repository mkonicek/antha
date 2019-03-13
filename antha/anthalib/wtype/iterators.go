package wtype

//iterators.go includes generalised replacements for plateiterator.go

//AddressIterator iterates through the addresses in an Addressable
type AddressIterator interface {
	Next() WellCoords
	Curr() WellCoords
	MoveTo(WellCoords)
	Reset()
	GetAddressable() Addressable
	Valid() bool
}

//AddressSliceIterator iterates through slices of addresses
type AddressSliceIterator interface {
	Next() []WellCoords
	Curr() []WellCoords
	MoveTo(WellCoords)
	Reset()
	GetAddressable() Addressable
	Valid() bool
}

type VerticalDirection int

const (
	BottomToTop VerticalDirection = -1
	TopToBottom VerticalDirection = 1
)

func (s VerticalDirection) String() string {
	switch s {
	case BottomToTop:
		return "bottom to top"
	case TopToBottom:
		return "top to bottom"
	}
	return ""
}

type HorizontalDirection int

const (
	LeftToRight HorizontalDirection = 1
	RightToLeft HorizontalDirection = -1
)

func (s HorizontalDirection) String() string {
	switch s {
	case LeftToRight:
		return "left to right"
	case RightToLeft:
		return "right to left"
	}
	return ""
}

type MajorOrder int

const (
	RowWise MajorOrder = iota
	ColumnWise
)

var majorOrderNames = map[MajorOrder]string{
	RowWise:    "row wise",
	ColumnWise: "column wise",
}

func (s MajorOrder) String() string {
	return majorOrderNames[s]
}

//NewAddressIterator iterates through the addresses in addr in order order, moving in directions ver and hor
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

//GetTickingIterator iterates through the addresses in addr in order order, moving in directions ver and hor, returning the output in slices of length chunkSize.
//In the direction specified by order, steps are of size stepSize (which is assumed to a factor of the number of adresses in that direction), and repeat such that
//every address is returned. Each address is repeated repeatAddresses times.
//For example, if addr is of size 4x4, then the iterator returned by
//  NewTickingIterator(addr, RowOrder, TopToBottom, LeftToRight, false, 4, 2, 2)
//returns
//  [[A1,A1,C1,C2],[B1,B1,D1,D1],[A2,A2, ...
func NewTickingIterator(addr Addressable, order MajorOrder, ver VerticalDirection, hor HorizontalDirection, repeat bool, chunkSize, stepSize, repeatAddresses int) AddressSliceIterator {
	it := newSteppingIterator(addr, order, ver, hor, repeat, stepSize)
	tick := newRepeatingIterator(it, repeatAddresses)
	chunk := newChunkedIterator(tick, chunkSize)
	return chunk
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
		it.update = getRowWiseUpdate(int(hor), int(ver), addr)
	} else {
		it.update = getColWiseUpdate(int(hor), int(ver), addr)
	}
	return &it
}

func newSteppingIterator(addr Addressable, order MajorOrder, ver VerticalDirection, hor HorizontalDirection, repeat bool, stepSize int) *simpleIterator {
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
		it.update = getRowWiseUpdate(int(hor)*stepSize, int(ver), addr)
	} else {
		it.update = getColWiseUpdate(int(hor), int(ver)*stepSize, addr)
	}
	return &it
}

func getRowWiseUpdate(dx, dy int, a Addressable) updateFn {
	nCols := a.NCols()

	if dx > 0 {
		ret := func(wc WellCoords) WellCoords {
			if wc.X == nCols-1 {
				wc.X = 0
				wc.Y += dy
			} else {
				wc.X += dx
				if wc.X >= nCols {
					wc.X = wc.X%nCols + 1
				}
			}
			return wc
		}
		return ret
	}
	ret := func(wc WellCoords) WellCoords {
		if wc.X == 0 {
			wc.X = nCols - 1
			wc.Y += dy
		} else {
			wc.X += dx
			if wc.X < 0 {
				wc.X += nCols - 1
			}
		}
		return wc
	}
	return ret
}

func getColWiseUpdate(dx, dy int, a Addressable) updateFn {
	nRows := a.NRows()

	if dy > 0 {
		ret := func(wc WellCoords) WellCoords {
			if wc.Y == nRows-1 {
				wc.Y = 0
				wc.X += dx
			} else {
				wc.Y += dy
				if wc.Y >= nRows {
					wc.Y = wc.Y%nRows + 1
				}
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

func (self *simpleIterator) GetAddressable() Addressable {
	return self.addr
}

func (self *simpleIterator) Reset() {
	self.curr = self.first
}

type chunkedIterator struct {
	it       AddressIterator
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

func (self *chunkedIterator) GetAddressable() Addressable {
	return self.it.GetAddressable()
}

func (self *chunkedIterator) Valid() bool {
	for _, wc := range self.curr {
		if !self.it.GetAddressable().AddressExists(wc) {
			return false
		}
	}
	return true
}

func newChunkedIterator(it AddressIterator, chunkLen int) *chunkedIterator {
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

type repeatingIterator struct {
	it         AddressIterator
	repeats    int
	currRepeat int
}

func newRepeatingIterator(it AddressIterator, repeats int) *repeatingIterator {
	return &repeatingIterator{
		it:      it,
		repeats: repeats,
	}
}

func (self *repeatingIterator) Next() WellCoords {
	self.currRepeat += 1
	if self.currRepeat < self.repeats {
		return self.it.Curr()
	}
	self.currRepeat = 0
	return self.it.Next()
}

func (self *repeatingIterator) Curr() WellCoords {
	return self.it.Curr()
}

func (self *repeatingIterator) MoveTo(wc WellCoords) {
	self.it.MoveTo(wc)
}

func (self *repeatingIterator) Reset() {
	self.it.Reset()
}

func (self *repeatingIterator) GetAddressable() Addressable {
	return self.it.GetAddressable()
}

func (self *repeatingIterator) Valid() bool {
	return self.it.Valid()
}
