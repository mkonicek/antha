package wtype

type PlateIterator interface {
	Rewind() WellCoords
	Next() WellCoords
	Curr() WellCoords
	Valid() bool
	SetStartTo(WellCoords)
	SetCurTo(WellCoords)
}

type VectorPlateIterator interface {
	Rewind() []WellCoords
	Next() []WellCoords
	Curr() []WellCoords
	Valid() bool
	SetStartTo(WellCoords)
	SetCurTo(WellCoords)
}

type BasicPlateIterator struct {
	fst  WellCoords
	cur  WellCoords
	a    Addressable
	rule func(WellCoords, Addressable) WellCoords
}

type MultiPlateIterator struct {
	BasicPlateIterator
	multi int
	rule  func(PlateIterator) []WellCoords
	ori   int
}

func (mpi *MultiPlateIterator) Curr() []WellCoords {
	wc := mpi.BasicPlateIterator.Curr()

	/*
		wa := make([]WellCoords, mpi.multi)
		for i := 0; i < mpi.multi; i++ {
			wa[i] = mpi.BasicPlateIterator.Curr()
			mpi.BasicPlateIterator.Next()
		}
	*/

	wa := mpi.rule(&mpi.BasicPlateIterator)

	mpi.SetCurTo(wc)

	return wa
}

func (mpi *MultiPlateIterator) Rewind() []WellCoords {
	mpi.BasicPlateIterator.Rewind()
	return mpi.Curr()
}

func (mpi *MultiPlateIterator) Next() []WellCoords {
	_ = mpi.rule(&mpi.BasicPlateIterator)
	return mpi.Curr()
}

func (mpi *MultiPlateIterator) Valid() bool {
	if !mpi.BasicPlateIterator.Valid() {
		return false
	}
	wc := mpi.BasicPlateIterator.Curr()

	valid := true

	wa := mpi.Curr()

	for i := 0; i < len(wa); i++ {
		if (mpi.ori == LHVChannel && wa[i].X != wc.X) || (mpi.ori == LHHChannel && wa[i].Y != wc.Y) {
			valid = false
			break
		}
	}

	// reset
	mpi.BasicPlateIterator.cur = wc

	return valid
}

func (it *BasicPlateIterator) Rewind() WellCoords {
	it.cur = it.fst
	return it.cur
}
func (it *BasicPlateIterator) Curr() WellCoords {
	return it.cur
}

func (it *BasicPlateIterator) Valid() bool {
	return it.a.AddressExists(it.cur)
}

func (it *BasicPlateIterator) Next() WellCoords {
	it.cur = it.rule(it.cur, it.a)
	return it.cur
}
func (it *BasicPlateIterator) SetStartTo(wc WellCoords) {
	it.fst = wc
}

func (it *BasicPlateIterator) SetCurTo(wc WellCoords) {
	it.cur = wc
}

func DownOneColumn(wc WellCoords, a Addressable) WellCoords {
	wc.Y += 1
	return wc
}

func AlongOneRow(wc WellCoords, a Addressable) WellCoords {
	wc.X += 1
	return wc
}

func NextInRowOnce(wc WellCoords, a Addressable) WellCoords {
	wc.X += 1
	if wc.X >= a.NCols() {
		wc.X = 0
		wc.Y += 1
	}
	return wc
}
func NextInRow(wc WellCoords, a Addressable) WellCoords {
	wc.X += 1
	if wc.X >= a.NCols() {
		wc.X = 0
		wc.Y += 1
	}
	if wc.Y >= a.NRows() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}

func NextInColumn(wc WellCoords, a Addressable) WellCoords {
	wc.Y += 1
	if wc.Y >= a.NRows() {
		wc.Y = 0
		wc.X += 1
	}
	if wc.X >= a.NCols() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}
func NextInColumnOnce(wc WellCoords, a Addressable) WellCoords {
	//fmt.Println(wc.FormatA1(), " ", "X: ", wc.X, " Y: ", wc.Y, "WX: ", a.NCols(), " WY: ", a.NRows())
	wc.Y += 1
	if wc.Y >= a.NRows() {
		wc.Y = 0
		wc.X += 1
	}
	return wc
}

func NewColumnWiseIterator(a Addressable) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInColumn
	bi.a = a
	return &bi
}
func NewOneTimeColumnWiseIterator(a Addressable) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInColumnOnce
	bi.a = a
	return &bi
}

func NewRowWiseIterator(a Addressable) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInRow
	bi.a = a
	return &bi
}
func NewOneTimeRowWiseIterator(a Addressable) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.a = a
	bi.rule = NextInRowOnce
	return &bi
}

func NewMultiIteratorRule(multi int) func(PlateIterator) []WellCoords {
	return func(pi PlateIterator) []WellCoords {
		wc := make([]WellCoords, multi)
		for i := 0; i < multi; i++ {
			wc[i] = pi.Curr()
			pi.Next()
		}

		return wc
	}
}

// still an issue with this generating out-of-bounds wells in single mode?!
type TickingColVectorIterator struct {
	Plate  *LHPlate
	Multi  int
	Ticker *Ticker
	start  int
}

func (tcvi *TickingColVectorIterator) Rewind() []WellCoords {
	tcvi.Ticker = &Ticker{Val: tcvi.start, TickEvery: tcvi.Ticker.TickEvery, TickBy: tcvi.Ticker.TickBy}
	return tcvi.Curr()
}

func (tcvi *TickingColVectorIterator) Next() []WellCoords {
	tcvi.Ticker = &Ticker{Val: tcvi.Ticker.Val + 1, TickEvery: tcvi.Ticker.TickEvery, TickBy: tcvi.Ticker.TickBy}

	tickRaw := tcvi.Ticker.Dup()

	end := 0

	for i := 0; i < tcvi.Multi-1; i++ {
		tickRaw.Tick()
		end = tickRaw.Val
		if end/tcvi.Plate.WellsY() > tcvi.Ticker.Val/tcvi.Plate.WellsY() {
			tcvi.Ticker = &Ticker{Val: end, TickEvery: tcvi.Ticker.TickEvery, TickBy: tcvi.Ticker.TickBy}
			break
		}
	}

	if end >= tcvi.Plate.WellsX()*tcvi.Plate.WellsY() {
		return []WellCoords{}
	}

	return tcvi.Curr()
}

func (tcvi *TickingColVectorIterator) Curr() []WellCoords {
	offsets := make([]int, tcvi.Multi)
	save := tcvi.Ticker.Dup()
	for i := 0; i < tcvi.Multi; i++ {
		v := tcvi.Ticker.Val
		if v >= tcvi.Plate.Nwells {
			return []WellCoords{}
		}
		offsets[i] = tcvi.Ticker.Val
		tcvi.Ticker.Tick()
	}
	tcvi.Ticker = save
	return tcvi.Plate.GetWellCoordsFromOrdering(offsets, BYCOLUMN)
}

func (tcvi *TickingColVectorIterator) Valid() bool {
	mx := tcvi.Plate.WellsX()*tcvi.Plate.WellsY() - 1

	tck := tcvi.Ticker.Dup()

	for i := 0; i < tcvi.Multi-1; i++ {
		if tck.Val > mx {
			return false
		}
		tck.Tick()
	}

	wcs := tcvi.Curr()

	if len(wcs) == 0 {
		return false
	}
	col := -1

	//fmt.Println(A1ArrayFromWellCoords(wcs))
	for _, wc := range wcs {
		if wc.X < 0 || wc.Y < 0 {
			return false
		}

		// are all rows and columns in bounds?
		if wc.X >= tcvi.Plate.WellsX() || wc.Y >= tcvi.Plate.WellsY() {
			return false
		}

		// are they all in the same column?
		if col == -1 {
			col = wc.X
		} else if col != wc.X {
			return false
		}
	}

	return true
}

func (tcvi *TickingColVectorIterator) SetStartTo(wc WellCoords) {
	tcvi.start = tcvi.Plate.GetOrderingFromA1WellCoords([]string{wc.FormatA1()}, BYCOLUMN)[0]
}
func (tcvi *TickingColVectorIterator) SetCurTo(wc WellCoords) {
	v := tcvi.Plate.GetOrderingFromA1WellCoords([]string{wc.FormatA1()}, BYCOLUMN)[0]
	tcvi.Ticker = &Ticker{Val: v, TickEvery: tcvi.Ticker.TickEvery, TickBy: tcvi.Ticker.TickBy}
}

func NewTickingColVectorIterator(p *LHPlate, multi, tpw, wpt int) VectorPlateIterator {
	ticker := &Ticker{TickEvery: tpw, TickBy: wpt}
	return &TickingColVectorIterator{Plate: p, Multi: multi, Ticker: ticker}
}

func NewColVectorIterator(a Addressable, multi int) VectorPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.a = a
	bi.rule = NextInColumnOnce
	rule := NewMultiIteratorRule(multi)
	mi := MultiPlateIterator{bi, multi, rule, LHVChannel}
	return &mi
}

func NewRowVectorIterator(a Addressable, multi int) VectorPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.a = a
	bi.rule = NextInRowOnce
	rule := NewMultiIteratorRule(multi)
	mi := MultiPlateIterator{bi, multi, rule, LHHChannel}
	return &mi
}
