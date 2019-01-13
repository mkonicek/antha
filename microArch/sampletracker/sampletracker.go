package sampletracker

import (
	"encoding/json"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

// SampleTracker record the location of components generated during element execution
// as well as any explicitly set input plates
type SampleTracker struct {
	lock     sync.Mutex
	records  map[string]string
	forwards map[string]string
	plates   []*wtype.Plate
}

type sampleTrackerSerializable struct {
	Records  map[string]string
	Forwards map[string]string
	Plates   []*wtype.Plate
}

func NewSampleTracker() *SampleTracker {
	return &SampleTracker{
		records:  make(map[string]string),
		forwards: make(map[string]string),
	}
}

func (st *SampleTracker) MarshalJSON() ([]byte, error) {
	st.lock.Lock()
	defer st.lock.Unlock()

	sts := &sampleTrackerSerializable{
		Records:  st.records,
		Forwards: st.forwards,
		Plates:   st.plates,
	}

	return json.Marshal(sts)
}

func (st *SampleTracker) UnmarshalJSON(bs []byte) error {
	if string(bs) == "null" {
		return nil

	} else {
		st.lock.Lock()
		defer st.lock.Unlock()
		sts := &sampleTrackerSerializable{}
		if err := json.Unmarshal(bs, sts); err != nil {
			return err
		} else {
			st.records = sts.Records
			st.forwards = sts.Forwards
			st.plates = sts.Plates
			return nil
		}
	}
}

// SetInputPlate declare the given plate as an input to the experiment
// recording the id and location of every sample in it
func (st *SampleTracker) SetInputPlate(idGen *id.IDGenerator, p *wtype.Plate) {
	st.lock.Lock()
	defer st.lock.Unlock()

	st.plates = append(st.plates, p)

	for _, w := range p.HWells {
		if !w.IsEmpty(idGen) {
			st.setLocationOf(w.WContents.ID, w.WContents.Loc)
			w.SetUserAllocated()
		}
	}
}

// GetInputPlates return a list of all input plates explicitly set during element
// execution
func (st *SampleTracker) GetInputPlates() []*wtype.Plate {
	st.lock.Lock()
	defer st.lock.Unlock()

	return st.plates
}

func (st *SampleTracker) setLocationOf(ID string, loc string) {
	st.records[ID] = loc
}

// SetLocationOf set the string encoded location of the component with the given ID
func (st *SampleTracker) SetLocationOf(ID string, loc string) {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.setLocationOf(ID, loc)
}

func (st *SampleTracker) getLocationOf(ID string) (string, bool) {
	if ID == "" {
		return "", false
	}

	s, ok := st.records[ID]

	// look to see if there's a forwarding address
	// can this lead to an out of date location???

	if !ok {
		return st.getLocationOf(st.forwards[ID])
	}

	return s, ok
}

// GetLocationOf return the string location of the component with the given ID.
// If no such component is known, the returned location will be the empty string
func (st *SampleTracker) GetLocationOf(ID string) (string, bool) {
	st.lock.Lock()
	defer st.lock.Unlock()

	return st.getLocationOf(ID)
}

// UpdateIDOf add newID as an alias for ID, such that both refer to the same location
func (st *SampleTracker) UpdateIDOf(ID string, newID string) {
	st.lock.Lock()
	defer st.lock.Unlock()
	_, ok := st.records[ID]
	if ok {
		st.records[newID] = st.records[ID]
	} else {
		// set up a forward
		// actually a backward...
		st.forwards[newID] = ID
	}
}
