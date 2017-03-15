package wtype

import (
	"fmt"
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

type mt struct {
	Sc float64
	Vl float64
	Bk int
}

func printMat(mat [][]mt) {
	fmt.Println("*****")
	for _, v := range mat {
		for _, x := range v {
			fmt.Printf("%-5.1f:%-1d ", x.Sc, x.Bk)
		}

		fmt.Println()
	}
	fmt.Println("-----")
}

func align(want, got ComponentVector, independent bool) Match {
	// ensure things are ok

	for i, v := range want {
		if v == nil {
			want[i] = NewLHComponent()
		}

		g := want[i]

		if g == nil {
			got[i] = NewLHComponent()
		}
	}

	mat := make([][]mt, len(want))

	mxmx := -999999.0
	mxi := -1
	mxj := -1

	for i := 0; i < len(want); i++ {
		mat[i] = make([]mt, len(got))

		// if we must be contiguous then we skip all cells aligned
		// to rows with zeroes in 'want'

		for j := 0; j < len(got); j++ {
			// only allow gaps if independent is set
			if (got[j].CName == "" || want[i].CName != got[j].CName) && !independent {
				continue
			}

			if want[i].CName != got[j].CName {
				// we set this to no volume
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
					mat[i][j].Sc = v2.ConvertToString("ul")
				}
			}

			mx := 0.0
			bk := 0
			if i > 0 && j > 0 {
				mx = mat[i-1][j-1].Sc
				bk = 2

				/*
					if independent {
						if want[i-1] == nil || want[i-1].CName == "" || want[i-1].Vol == 0.0 {
							if mat[i-1][j].Sc > mx {
								mx = mat[i-1][j].Sc
								bk = 1
							}
						}
						if mat[i][j-1].Sc > mx {
							mx = mat[i][j-1].Sc
							bk = 3
						}
					}
				*/
			}
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
	}

	m := Match{IDs: IDs, WCs: WCs, Vols: Vols, M: Ms, Sc: mxmx}

	if mxi < 0 || mxj < 0 || mxmx <= 0.0 {
		return m
	}

	// get stuff out

	gIDs := got.GetPlateIds()
	gWCs := got.GetWellCoords()
	//gVs := got.GetVols()

	// get the best

	i := mxi
	j := mxj

	for {
		if want[i].Vol == 0 && mat[i][j].Bk == 0 {
			break
		}
		IDs[i] = gIDs[j]
		WCs[i] = gWCs[j]
		Vols[i] = wunit.NewVolume(mat[i][j].Vl, "ul")
		Ms[i] = j

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

	return m
}

func matchComponents(want, got ComponentVector, independent bool) (ComponentMatch, error) {
	// not sure of the algorithm here:
	// we want to match as many as possible in one go
	// then clean up the others

	m := ComponentMatch{Matches: make([]Match, 0, 1)}

	for {
		match := align(want, got, independent)
		if match.Sc <= 0.0 {
			break
		}
		m.Matches = append(m.Matches, match)

		// deplete
		c := 0
		for i := 0; i < len(match.WCs); i++ {
			if match.WCs[i] != "" {
				if got[match.M[i]].Vol >= want[i].Vol {
					got[match.M[i]].Vol -= want[i].Vol
					want[i].Vol = 0.0
				} else {
					got[match.M[i]].Vol -= want[i].Vol
					want[i].Vol -= match.Vols[i].ConvertToString(want[i].Vunit)
				}
				c += 1
			}
		}

		if c == len(match.WCs) || c == 0 {
			break
		}
	}

	return m, nil
}

func scoreMatch(m ComponentMatch, independent bool) float64 {
	s := 0.0

	for _, mtch := range m.Matches {
		for i := 0; i < len(mtch.Vols); i++ {
			s += mtch.Vols[i].RawValue()
		}
	}

	s /= float64(len(m.Matches))

	return s
}
