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
	p    *LHPlate
	rule func(WellCoords, *LHPlate) WellCoords
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
	if it.cur.X >= it.p.WellsX() || it.cur.X < 0 {
		return false
	}

	if it.cur.Y >= it.p.WellsY() || it.cur.Y < 0 {
		return false
	}

	return true
}

func (it *BasicPlateIterator) Next() WellCoords {
	it.cur = it.rule(it.cur, it.p)
	return it.cur
}
func (it *BasicPlateIterator) SetStartTo(wc WellCoords) {
	it.fst = wc
}

func (it *BasicPlateIterator) SetCurTo(wc WellCoords) {
	it.cur = wc
}

func DownOneColumn(wc WellCoords, p *LHPlate) WellCoords {
	wc.Y += 1
	return wc
}

func AlongOneRow(wc WellCoords, p *LHPlate) WellCoords {
	wc.X += 1
	return wc
}

func NextInRowOnce(wc WellCoords, p *LHPlate) WellCoords {
	wc.X += 1
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y += 1
	}
	return wc
}
func NextInRow(wc WellCoords, p *LHPlate) WellCoords {
	wc.X += 1
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y += 1
	}
	if wc.Y >= p.WellsY() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}

func NextInColumn(wc WellCoords, p *LHPlate) WellCoords {
	wc.Y += 1
	if wc.Y >= p.WellsY() {
		wc.Y = 0
		wc.X += 1
	}
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}
func NextInColumnOnce(wc WellCoords, p *LHPlate) WellCoords {
	//fmt.Println(wc.FormatA1(), " ", "X: ", wc.X, " Y: ", wc.Y, "WX: ", p.WellsX(), " WY: ", p.WellsY())
	wc.Y += 1
	if wc.Y >= p.WellsY() {
		wc.Y = 0
		wc.X += 1
	}
	return wc
}

func NewColumnWiseIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInColumn
	bi.p = p
	return &bi
}
func NewOneTimeColumnWiseIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInColumnOnce
	bi.p = p
	return &bi
}

func NewRowWiseIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextInRow
	bi.p = p
	return &bi
}
func NewOneTimeRowWiseIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.p = p
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

func NewColVectorIterator(p *LHPlate, multi int) VectorPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.p = p
	bi.rule = NextInColumnOnce
	rule := NewMultiIteratorRule(multi)
	mi := MultiPlateIterator{bi, multi, rule, LHVChannel}
	return &mi
}

func NewRowVectorIterator(p *LHPlate, multi int) VectorPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.p = p
	bi.rule = NextInRowOnce
	rule := NewMultiIteratorRule(multi)
	mi := MultiPlateIterator{bi, multi, rule, LHHChannel}
	return &mi
}
