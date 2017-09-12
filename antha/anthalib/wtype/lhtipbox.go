// liquidhandling/lhtypes.Go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package wtype

import "fmt"

/* tip box */

type LHTipbox struct {
	ID         string
	Boxname    string
	Type       string
	Mnfr       string
	Nrows      int
	Ncols      int
	Height     float64
	Tiptype    *LHTip
	AsWell     *LHWell
	NTips      int
	Tips       [][]*LHTip
	TipXOffset float64
	TipYOffset float64
	TipXStart  float64
	TipYStart  float64
	TipZStart  float64
}

func NewLHTipbox(nrows, ncols int, height float64, manufacturer, boxtype string, tiptype *LHTip, well *LHWell, tipxoffset, tipyoffset, tipxstart, tipystart, tipzstart float64) *LHTipbox {
	var tipbox LHTipbox
	//tipbox.ID = "tipbox-" + GetUUID()
	tipbox.ID = GetUUID()
	tipbox.Type = boxtype
	tipbox.Boxname = fmt.Sprintf("%s_%s", boxtype, tipbox.ID[1:len(tipbox.ID)-2])
	tipbox.Mnfr = manufacturer
	tipbox.Nrows = nrows
	tipbox.Ncols = ncols
	tipbox.Tips = make([][]*LHTip, ncols)
	tipbox.NTips = tipbox.Nrows * tipbox.Ncols
	tipbox.Height = height
	tipbox.Tiptype = tiptype
	tipbox.AsWell = well
	for i := 0; i < ncols; i++ {
		tipbox.Tips[i] = make([]*LHTip, nrows)
	}
	tipbox.TipXOffset = tipxoffset
	tipbox.TipYOffset = tipyoffset
	tipbox.TipXStart = tipxstart
	tipbox.TipYStart = tipystart
	tipbox.TipZStart = tipzstart
	return initialize_tips(&tipbox, tiptype)
}

func (tb LHTipbox) GetID() string {
	return tb.ID
}

func (tb LHTipbox) Output() string {
	s := ""
	for j := 0; j < tb.NumRows(); j++ {
		for i := 0; i < tb.NumCols(); i++ {
			if tb.Tips[i][j] == nil {
				s += "."
			} else if tb.Tips[i][j].Dirty {
				s += "*"
			} else {
				s += "o"
			}
		}
		s += "\n"
	}

	return s
}

func (tb LHTipbox) String() string {
	return fmt.Sprintf(
		`LHTipbox {
ID        : %s,
Boxname   : %s,
Type      : %s,
Mnfr      : %s,
Nrows     : %d,
Ncols     : %d,
Height    : %f,
Tiptype   : %p,
AsWell    : %v,
NTips     : %d,
Tips      : %p,
TipXOffset: %f,
TipYOffset: %f,
TipXStart : %f,
TipYStart : %f,
TipZStart : %f,
}`,
		tb.ID,
		tb.Boxname,
		tb.Type,
		tb.Mnfr,
		tb.Nrows,
		tb.Ncols,
		tb.Height,
		tb.Tiptype,
		tb.AsWell,
		tb.NTips,
		tb.Tips,
		tb.TipXOffset,
		tb.TipYOffset,
		tb.TipXStart,
		tb.TipYStart,
		tb.TipZStart,
	)
}

//lazy sunuva
func (tb *LHTipbox) Dup() *LHTipbox {
	tb2 := NewLHTipbox(tb.Nrows, tb.Ncols, tb.Height, tb.Mnfr, tb.Type, tb.Tiptype, tb.AsWell, tb.TipXOffset, tb.TipYOffset, tb.TipXStart, tb.TipYStart, tb.TipZStart)

	for i := 0; i < len(tb.Tips); i++ {
		for j := 0; j < len(tb.Tips[i]); j++ {
			t := tb.Tips[i][j]
			if t == nil {
				tb2.Tips[i][j] = nil
			} else {
				tb2.Tips[i][j] = t.Dup()
			}
		}
	}

	return tb2
}

// @implement named

func (tb *LHTipbox) GetName() string {
	return tb.Boxname
}

func (tb *LHTipbox) N_clean_tips() int {
	c := 0
	for j := 0; j < tb.Nrows; j++ {
		for i := 0; i < tb.Ncols; i++ {
			if tb.Tips[i][j] != nil && !tb.Tips[i][j].Dirty {
				c += 1
			}
		}
	}
	return c
}

func trim(ba []bool) []bool {
	r := make([]bool, 0, len(ba))
	s := -1
	e := -1

	for i := 0; i < len(ba); i++ {
		if ba[i] {
			if s == -1 {
				s = i
			}

			e = i
		}
	}

	for i := s; i <= e; i++ {
		r = append(r, ba[i])
	}

	return r
}

// returns wells with tips in
// -- this needs to align the tips to get with the channels
// since there's no strict requirement to fix here
func (tb *LHTipbox) GetTipsMasked(mask []bool, ori int) []string {
	if ori == LHVChannel {
		for i := 0; i < tb.NumCols(); i++ {
			r := tb.searchCleanTips(i, trim(mask), ori)
			if r != nil && len(r) != 0 {
				tb.Remove(r)
				return r
			}
		}
	} else if ori == LHHChannel {
		for i := 0; i < tb.NumRows(); i++ {
			r := tb.searchCleanTips(i, trim(mask), ori)
			if r != nil && len(r) != 0 {
				tb.Remove(r)
				return r
			}
		}
	}

	// not found or unknown orientation
	return []string{}
}

func (tb *LHTipbox) Remove(sa []string) bool {
	ar := WCArrayFromStrings(sa)

	for _, wc := range ar {
		if wc.X < 0 {
			continue
		}
		if wc.X >= len(tb.Tips) || wc.Y >= len(tb.Tips[wc.X]) || tb.Tips[wc.X][wc.Y] == nil {
			return false
		}

		tb.Tips[wc.X][wc.Y] = nil
	}

	return true
}

func inflateMask(mask []bool, offset, size int) []bool {
	r := make([]bool, size)

	for i := 0; i < len(mask); i++ {
		r[i+offset] = mask[i]
	}

	return r
}

func maskToWellCoords(mask []bool, offset, ori int) []string {
	wc := make([]WellCoords, len(mask))

	for i := 0; i < len(mask); i++ {
		wc[i] = WellCoords{X: -1, Y: -1}

		curWC := WellCoords{X: -1, Y: -1}

		if ori == LHVChannel {
			curWC = WellCoords{X: offset, Y: i}
		} else if ori == LHHChannel {
			curWC = WellCoords{X: i, Y: offset}
		}

		if mask[i] {
			wc[i] = curWC
		}
	}

	r := make([]string, len(wc))

	for i := 0; i < len(wc); i++ {
		if wc[i].X != -1 {
			r[i] = wc[i].FormatA1()
		}
	}

	return r
}

func (tb *LHTipbox) searchCleanTips(offset int, mask []bool, ori int) []string {
	r := make([]string, 0, 1)

	if ori == LHVChannel {
		df := tb.NumRows() - len(mask) + 1
		for i := 0; i < df; i++ {
			m := inflateMask(mask, i, tb.NumRows())
			if tb.hasCleanTips(offset, m, ori) {
				return maskToWellCoords(m, offset, ori)
			}
		}
	} else if ori == LHHChannel {
		df := tb.NumCols() - len(mask) + 1
		for i := 0; i < df; i++ {
			m := inflateMask(mask, i, tb.NumCols())
			if tb.hasCleanTips(offset, m, ori) {
				return maskToWellCoords(m, offset, ori)
			}
		}

	}

	return r
}

// fails iff for true mask[i] there is no corresponding clean tip
func (tb *LHTipbox) hasCleanTips(offset int, mask []bool, ori int) bool {
	if ori == LHVChannel {
		for i := 0; i < len(mask); i++ {
			if mask[i] && (tb.Tips[offset][i] == nil || tb.Tips[offset][i].Dirty) {
				return false
			}
		}

		return true
	} else if ori == LHHChannel {
		for i := 0; i < len(mask); i++ {
			if mask[i] && (tb.Tips[i][offset] == nil || !tb.Tips[i][offset].Dirty) {
				return false
			}
		}

		return true
	}

	return false
}

// deprecated shortly
func (tb *LHTipbox) GetTips(mirror bool, multi, orient int) []string {
	// this removes the tips as well
	var ret []string = nil
	if orient == LHHChannel {
		for j := 0; j < tb.Nrows; j++ {
			c := 0
			s := -1
			for i := 0; i < tb.Ncols; i++ {
				if tb.Tips[i][j] != nil && !tb.Tips[i][j].Dirty {
					c += 1
					if s == -1 {
						s = i
					}
				}
			}

			if c >= multi {
				ret = make([]string, multi)
				for i := 0; i < multi; i++ {
					tb.Tips[i+s][j] = nil
					wc := WellCoords{i + s, j}
					ret[i] = wc.FormatA1()
				}
				break
			}
		}

	} else if orient == LHVChannel {
		// find the first column with a contiguous set of at least multi
		for i := 0; i < tb.Ncols; i++ {
			c := 0
			s := -1
			// if we're picking up < the maxium number of tips we need to be careful
			// that there are no tips beneath the ones we're picking up

			for j := tb.Nrows - 1; j >= 0; j-- {
				if tb.Tips[i][j] != nil { // && !tb.Tips[i][j].Dirty
					c += 1
					if s == -1 {
						s = j
					}
				} else {
					if s != -1 {
						break // we've reached a gap
					}
				}
			}

			if c >= multi {
				ret = make([]string, 0, multi)
				n := 0
				for j := s; j >= 0; j-- {
					tb.Tips[i][j] = nil
					wc := WellCoords{i, j}
					//fmt.Println(j, "Getting TIP from ", wc.FormatA1())
					ret = append(ret, wc.FormatA1())
					n += 1
					if n >= multi {
						break
					}
				}

				//fmt.Println("RET: ", ret)
				break
			}
		}

	}

	tb.NTips -= len(ret)
	return reverse(ret)
}

func reverse(ar []string) []string {
	ret := make([]string, 0, len(ar))
	for k := len(ar) - 1; k >= 0; k-- {
		ret = append(ret, ar[k])
	}
	return ret
}
func (tb *LHTipbox) Refresh() {
	initialize_tips(tb, tb.Tiptype)
}

func (tb *LHTipbox) Refill() {
	tb.Refresh()
}

func initialize_tips(tipbox *LHTipbox, tiptype *LHTip) *LHTipbox {
	nr := tipbox.Nrows
	nc := tipbox.Ncols
	for i := 0; i < nc; i++ {
		for j := 0; j < nr; j++ {
			tipbox.Tips[i][j] = CopyTip(*tiptype)
		}
	}
	tipbox.NTips = tipbox.Nrows * tipbox.Ncols
	return tipbox
}

// @implement SBSLabware

func (tipbox *LHTipbox) NumRows() int {
	return tipbox.Nrows
}
func (tipbox *LHTipbox) NumCols() int {
	return tipbox.Ncols
}

func (tipbox *LHTipbox) PlateHeight() float64 {
	return tipbox.Height
}
