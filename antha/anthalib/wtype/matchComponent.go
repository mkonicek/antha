package wtype

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// TODO --> deal with, e.g., 384 well plates

type ComponentMatch struct {
	Matches []Match
}
type Match struct {
	IDs  []string       // PlateIDs in 'got' array
	WCs  []string       // Wellcoords in 'got' array
	Vols []wunit.Volume // vols (before suck) in 'got'
	M    []int          // offsets in 'got' array
	Sc   float64        // total score for this match
}

func (m Match) Equals(m2 Match) bool {
	eqV := func(va, v2 []wunit.Volume) bool {
		if len(va) != len(v2) {
			return false
		}
		for i := 0; i < len(va); i++ {
			if !va[i].EqualTo(v2[i]) {
				return false
			}
		}
		return true
	}

	return reflect.DeepEqual(m.IDs, m2.IDs) && reflect.DeepEqual(m.WCs, m2.WCs) && reflect.DeepEqual(m.M, m2.M) && eqV(m.Vols, m2.Vols) && m.Sc == m2.Sc
}

type mt struct {
	Sc float64
	Vl float64
	Bk int
}

func printMat(mat [][]mt) {
	fmt.Println("*****")
	for i, v := range mat {
		for j := range v {
			fmt.Printf("(%d,%d):%-5.1f:%-1d:%-5.1f ", i, j, mat[i][j].Sc, mat[i][j].Bk, mat[i][j].Vl)
		}

		fmt.Println()
	}
	fmt.Println("-----")
}

func align(want, got ComponentVector, independent, debug bool) Match {
	for i, v := range want {
		if v == nil {
			want[i] = NewLHComponent()
		}

		if i >= len(got) {
			continue
		}

		g := got[i]

		if g == nil {
			got[i] = NewLHComponent()
		}
	}

	mat := make([][]mt, len(want))

	mxmx := -999999.0
	mxi := -1
	mxj := -1

	for i := 0; i < len(want); i++ {
		if want[i] == nil {
			continue
		}
		mat[i] = make([]mt, len(got))

		// if we must be contiguous then we skip all cells aligned
		// to rows with zeroes in 'want'

		for j := 0; j < len(got); j++ {
			// only allow gaps if independent is set
			if got[j] == nil {
				continue
			}
			// CName might not always match
			//if (got[j].CName == "" || want[i].CName != got[j].CName) && !independent {
			if (got[j].CName == "" || !want[i].Matches(got[j])) && !independent {
				continue
			}

			//if want[i].CName != got[j].CName {
			if !want[i].Matches(got[j]) {
				mat[i][j].Vl = 0.0
				mat[i][j].Sc = 0.0
			} else {
				v1 := want[i].Volume().Dup()
				v2 := got[j].Volume().Dup()

				if v1.GreaterThan(v2) {
					mat[i][j].Vl = v2.ConvertToString("ul")
					v2 = wunit.ZeroVolume()
				} else {
					mat[i][j].Vl = v1.ConvertToString("ul")
					v2.Subtract(v1)
				}

				if !want[i].Volume().IsZero() {
					//	mat[i][j].Sc = v2.ConvertToString("ul")
					mat[i][j].Sc = mat[i][j].Vl
				}
			}

			mx := 0.0
			bk := 0
			if i > 0 && j > 0 && mat[i-1][j-1].Sc > mx {
				mx = mat[i-1][j-1].Sc
				bk = 2
			}

			/*
				// get several things from the same place
				// if this is a trough it's fine to do parallel
				// otherwise the code in transferblock forces single channeling
				if i > 0 && mat[i-1][j].Sc > mx {
					mx = mat[i-1][j].Sc
					bk = 1
				}
			*/

			mat[i][j].Sc += mx
			mat[i][j].Bk = bk

			if mat[i][j].Sc > mxmx {
				mxmx = mat[i][j].Sc
				mxi = i
				mxj = j
			}

		}
	}

	IDs := make([]string, len(want))
	WCs := make([]string, len(want))
	Vols := make([]wunit.Volume, len(want))
	Ms := make([]int, len(want))

	for i := 0; i < len(want); i++ {
		Ms[i] = -1
		Vols[i] = wunit.ZeroVolume()
	}

	m := Match{IDs: IDs, WCs: WCs, Vols: Vols, M: Ms, Sc: mxmx}

	if mxi < 0 || mxj < 0 || mxmx <= 0.0 {
		return m
	}

	// get stuff out

	gIDs := got.GetPlateIds()
	gWCs := got.GetWellCoords()

	// get the best

	i := mxi
	j := mxj

	for {
		if want[i].Vol == 0 && mat[i][j].Bk == 0 || mat[i][j].Vl == 0 && !independent {
			break
		}

		if mat[i][j].Vl != 0 {
			IDs[i] = gIDs[j]
			WCs[i] = gWCs[j]
			Ms[i] = j
		}
		Vols[i] = wunit.NewVolume(mat[i][j].Vl, "ul")

		bk := mat[i][j].Bk

		if bk == 0 {
			break
		} else if bk == 1 {
			i = i - 1
		} else if bk == 2 {
			i = i - 1
			j = j - 1
		} else if bk == 3 {
			j = j - 1
		}
	}

	if debug {
		printMat(mat)
	}

	return m
}

var NotFoundError = errors.New("not found")

// IsNotFound returns true if the underlying error is NotFoundError
func IsNotFound(err error) bool {
	return errors.Cause(err) == NotFoundError
}

// matchComponents takes one bite each time... the best it can find
// needs to be run repeatedly to pick everything up
// TODO: needs to supply more options
func MatchComponents(want, got ComponentVector, independent, debug bool) (Match, error) {
	match := align(want, got, independent, debug)

	if match.Sc <= 0.0 {
		return Match{}, NotFoundError
	}

	return match, nil
}
