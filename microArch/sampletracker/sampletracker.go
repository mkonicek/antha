package sampletracker

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var st *SampleTracker

type SampleTracker struct {
	records  map[string]string
	forwards map[string]string
	plates   map[string]*wtype.LHPlate
}

func newSampleTracker() *SampleTracker {
	r := make(map[string]string)
	f := make(map[string]string)
	p := make(map[string]*wtype.LHPlate)
	st := SampleTracker{r, f, p}
	return &st
}

func GetSampleTracker() *SampleTracker {
	if st == nil {
		st = newSampleTracker()
	}

	return st
}

func (st *SampleTracker) SetInputPlate(p *wtype.LHPlate) {
	st.plates[p.ID] = p

	for _, w := range p.HWells {
		if !w.Empty() {
			st.SetLocationOf(w.WContents.ID, w.WContents.Loc)
			w.SetUserAllocated()
		}
	}
}

// this is destructive, i.e. once asked for that's it
// that's one way to make it thread-safe...
func (st *SampleTracker) GetInputPlates() []*wtype.LHPlate {
	var ret []*wtype.LHPlate
	if len(st.plates) == 0 {
		return ret
	}
	ret = make([]*wtype.LHPlate, 0, len(st.plates))

	for _, p := range st.plates {
		ret = append(ret, p)
	}

	st.plates = make(map[string]*wtype.LHPlate)

	return ret
}

func (st *SampleTracker) SetLocationOf(ID string, loc string) {
	//	fmt.Println("LOCATION OF ", ID, " SET TO ", loc)
	st.records[ID] = loc
}

func (st *SampleTracker) GetLocationOf(ID string) (string, bool) {
	//fmt.Println("GET LOCATION OF :", ID)
	if ID == "" {
		return "", false
	}

	s, ok := st.records[ID]

	// look to see if there's a forwarding address
	// can this lead to an out of date location???

	if !ok {
		return st.GetLocationOf(st.forwards[ID])
	}

	return s, ok
}

func (st *SampleTracker) UpdateIDOf(ID string, newID string) {
	_, ok := st.records[ID]
	if ok {
		st.records[newID] = st.records[ID]
	} else {
		// set up a forward
		// actually a backward...
		st.forwards[newID] = ID
	}
}
