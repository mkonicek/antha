package wtype

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

func A1ArrayFromWells(wells []*LHWell) []string {
	return A1ArrayFromWellCoords(WCArrayFromWells(wells))
}

func WCArrayFromWells(wells []*LHWell) []WellCoords {
	ret := make([]WellCoords, 0, len(wells))

	for _, w := range wells {
		if w == nil {
			continue
		}

		ret = append(ret, w.Crds)
	}

	return ret
}

func WCArrayFromStrings(arr []string) []WellCoords {
	ret := make([]WellCoords, len(arr))

	for i, s := range arr {
		ret[i] = MakeWellCoords(s)
	}

	return ret
}

func A1ArrayFromWellCoords(arr []WellCoords) []string {
	ret := make([]string, len(arr))
	for i, v := range arr {
		ret[i] = v.FormatA1()
	}
	return ret
}

// make an array of these from an array of strings

func MakeWellCoordsArray(sa []string) []WellCoords {
	r := make([]WellCoords, len(sa))

	for i := 0; i < len(sa); i++ {
		r[i] = MakeWellCoords(sa[i])
	}

	return r
}

func WCArrayCols(wcA []WellCoords) []int {
	return squashedIntFromWCA(wcA, 0)
}

func WCArrayRows(wcA []WellCoords) []int {
	return squashedIntFromWCA(wcA, 1)
}

func containsInt(i int, ia []int) bool {
	for _, ii := range ia {
		if i == ii {
			return true
		}
	}
	return false
}

func squashedIntFromWCA(wcA []WellCoords, which int) []int {
	ret := make([]int, 0, len(wcA))
	for _, wc := range wcA {
		v := wc.X
		if which == 1 {
			v = wc.Y
		}

		// ignore nils

		if v == -1 {
			continue
		}

		if !containsInt(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

// convenience comparison operator

func CompareStringWellCoordsCol(sw1, sw2 string) int {
	w1 := MakeWellCoords(sw1)
	w2 := MakeWellCoords(sw2)
	return CompareWellCoordsCol(w1, w2)
}

func CompareWellCoordsCol(w1, w2 WellCoords) int {
	dx := w1.X - w2.X
	dy := w1.Y - w2.Y

	if dx < 0 {
		return -1
	} else if dx > 0 {
		return 1
	}

	if dy < 0 {
		return -1
	} else if dy > 0 {
		return 1
	} else {
		return 0
	}
}

func CompareStringWellCoordsRow(sw1, sw2 string) int {
	w1 := MakeWellCoords(sw1)
	w2 := MakeWellCoords(sw2)
	return CompareWellCoordsRow(w1, w2)
}

func CompareWellCoordsRow(w1, w2 WellCoords) int {
	dx := w1.X - w2.X
	dy := w1.Y - w2.Y

	if dy < 0 {
		return -1
	} else if dy > 0 {
		return 1
	}
	if dx < 0 {
		return -1
	} else if dx > 0 {
		return 1
	} else {
		return 0
	}
}

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

func ZeroWellCoords() WellCoords {
	return WellCoords{-1, -1}
}
func (wc WellCoords) IsZero() bool {
	return wc.Equals(ZeroWellCoords())
}

func MatchString(s1, s2 string) bool {
	m, _ := regexp.MatchString(s1, s2)
	return m
}

func (wc WellCoords) Equals(w2 WellCoords) bool {
	if wc.X == w2.X && wc.Y == w2.Y {
		return true
	}

	return false
}

func MakeWellCoords(wc string) WellCoords {
	// try each one in turn

	r := MakeWellCoordsA1(wc)

	zero := WellCoords{-1, -1}

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoords1A(wc)

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoordsXY(wc)

	return r
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	re := regexp.MustCompile(`^([A-Z]{1,})([0-9]{1,2})$`)
	matches := re.FindStringSubmatch(a1)

	if matches == nil {
		return WellCoords{-1, -1}
	}
	/*
		re, _ := regexp.Compile("[A-Z]{1,}")
		ix := re.FindIndex([]byte(a1))
		endC := ix[1]
	*/

	X := wutil.ParseInt(matches[2]) - 1
	Y := wutil.AlphaToNum(matches[1]) - 1

	return WellCoords{X, Y}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	re := regexp.MustCompile(`^([0-9]{1,2})([A-Z]{1,})$`)
	matches := re.FindStringSubmatch(a1)

	if matches == nil {
		return WellCoords{-1, -1}
	}

	Y := wutil.AlphaToNum(matches[2]) - 1
	X := wutil.ParseInt(matches[1]) - 1
	return WellCoords{X, Y}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	r := WellCoords{wutil.ParseInt(y[1:]) - 1, wutil.ParseInt(x[1:]) - 1}

	if r.X < 0 || r.Y < 0 {
		return WellCoords{-1, -1}
	}

	return r
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	if tx == nil || len(tx) != 2 || len(tx[0]) == 0 || len(tx[1]) == 0 {
		return WellCoords{-1, -1}
	}
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc WellCoords) FormatXY() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}

func (wc WellCoords) Format1A() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X+1) + wutil.NumToAlpha(wc.Y+1)
}

func (wc WellCoords) FormatA1() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return wutil.NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}

// WellNumber returns the index of the well coordinates on a platetype based on
// looking up the number of wells in the X, Y directions of the platetype.
// Setting byRow to true will count along each sequential row rather down each sequential column.
// e.g.
// if byRow == true: A1 = 0, A2 = 1, A12 = 11
// if byRow == false: A1 = 0, B1 = 1, E1 = 4
func (wc WellCoords) WellNumber(platetype *Plate, byRow bool) int {
	return wc.wellNumber(platetype.WlsX, platetype.WlsY, byRow)
}

func (wc WellCoords) wellNumber(xLength, yLength int, byRow bool) int {
	if wc.X < 0 || wc.Y < 0 {
		return -1
	}
	if byRow {
		return (yLength*(wc.Y) + wc.X)
	}
	return (xLength*(wc.X) + wc.Y)
}

func (wc WellCoords) ColNumString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X + 1)
}

func (wc WellCoords) RowLettString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return wutil.NumToAlpha(wc.Y + 1)
}

// comparison operators

func (wc WellCoords) RowLessThan(wc2 WellCoords) bool {
	if wc.Y == wc2.Y {
		return wc.X < wc2.X
	}
	return wc.Y < wc2.Y
}

func (wc WellCoords) ColLessThan(wc2 WellCoords) bool {
	if wc.X == wc2.X {
		return wc.Y < wc2.Y
	}
	return wc.X < wc2.X
}

// convenience structure to allow sorting

type WellCoordArrayCol []WellCoords
type WellCoordArrayRow []WellCoords

func (wca WellCoordArrayCol) Len() int           { return len(wca) }
func (wca WellCoordArrayCol) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayCol) Less(i, j int) bool { return wca[i].RowLessThan(wca[j]) }

func (wca WellCoordArrayRow) Len() int           { return len(wca) }
func (wca WellCoordArrayRow) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayRow) Less(i, j int) bool { return wca[i].ColLessThan(wca[j]) }

func canAdd(coords []WellCoords, c WellCoords) bool {
	if len(coords) == 0 {
		return true
	}
	if len(coords) == 1 {
		return (c.Y == coords[0].Y && c.X-coords[0].X == 1) || (c.X == coords[0].X && c.Y-coords[0].Y == 1)
	}
	//if vertical run
	if coords[0].X == coords[len(coords)-1].X {
		return c.X == coords[0].X && c.Y == coords[len(coords)-1].Y+1
	}
	//if horizontal run
	return c.Y == coords[0].Y && c.X == coords[len(coords)-1].X+1
}

func runToString(coords []WellCoords) string {
	if len(coords) == 0 {
		return ""
	}
	if len(coords) == 1 {
		return coords[0].FormatA1()
	}
	return coords[0].FormatA1() + "-" + coords[len(coords)-1].FormatA1()
}

//HumanizeWellCoords convenience function to make displaying a slice of WellCoords more human readable
func HumanizeWellCoords(coords []WellCoords) string {
	s := make([]string, 0, len(coords))
	run := make([]WellCoords, 0, len(coords))
	for _, coord := range coords {
		if coord.IsZero() {
			continue
		}
		if !canAdd(run, coord) {
			s = append(s, runToString(run))
			run = make([]WellCoords, 0, len(coords))
		}
		run = append(run, coord)
	}
	s = append(s, runToString(run))
	return strings.Join(s, ",")
}

type WellCoordSlice []WellCoords

//Trim remove nil well coords from begining and end of the slice
func (self WellCoordSlice) Trim() WellCoordSlice {
	var start, end int
	for i, wc := range self {
		if !wc.IsZero() {
			if start == end { //only true if we've not seen a nil
				start = i
			}
			end = i + 1
		}
	}
	return self[start:end]
}
