package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ComponentMatch struct {
	Matches []Match
}
type Match struct {
	IDs  []string
	WCs  []string
	Vols []wunit.Volume
	M    []int
	Sc   float64
}

type mt struct {
	Sc float64
	Vl float64
	Bk int
}

func align(want, got ComponentVector, independent bool) Match {
	mat := make([][]mt, len(want))

	mxmx := -999999.0
	mxi := -1
	mxj := -1

	for i := 0; i < len(want); i++ {
		mat[i] = make([]mt, len(got))

		if want[i] == nil {
			continue
		}

		for j := 0; j < len(got); j++ {
			if got[j] == nil {
				continue
			}

			if want[i].CName == got[j].CName && !want[i].Volume().IsZero() {
				v1 := want[i].Volume().Dup()
				v2 := got[j].Volume().Dup()

				if v1.GreaterThan(v2) {
					mat[i][j].Vl = v2.ConvertToString("ul")
					v2 = wunit.ZeroVolume()
				} else {
					mat[i][j].Vl = v1.ConvertToString("ul")
					v2.Subtract(v1)
				}
				mat[i][j].Sc = v2.ConvertToString("ul")

				mx := 0.0
				bk := 0
				if i > 0 && j > 0 {
					mx = mat[i-1][j-1].Sc
					bk = 2

					if independent {
						// gaps allowed
						if mat[i-1][j].Sc > mx {
							mx = mat[i-1][j].Sc
							bk = 1
						}
						if mat[i][j-1].Sc > mx {
							mx = mat[i][j-1].Sc
							bk = 3
						}
					}
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
	gVs := got.GetVols()

	// get the best

	i := mxi
	j := mxj

	for {
		IDs[i] = gIDs[j]
		WCs[i] = gWCs[j]
		Vols[i] = gVs[j]
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
					want[i].Vol = 0.0
				} else {
					want[i].Vol -= match.Vols[i].ConvertToString(want[i].Vunit)
				}
				got[match.M[i]].Vol -= want[i].Vol
				c += 1
			}
		}

		if c == len(match.WCs) {
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
