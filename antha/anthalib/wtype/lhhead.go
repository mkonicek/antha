package wtype

import (
	"fmt"
	"github.com/pkg/errors"
)

// head
type LHHead struct {
	Name         string
	Manufacturer string
	ID           string
	Adaptor      *LHAdaptor
	Params       *LHChannelParameter
	//TipLoading defined the behaviour of the head when loading tips
	TipLoading TipLoadingBehaviour
}

//NewLHHead constructor for liquid handling heads
func NewLHHead(name, mf string, params *LHChannelParameter) *LHHead {
	return &LHHead{
		Name:         name,
		Manufacturer: mf,
		Params:       params,
	}
}

//Dup duplicate the head and adaptor, changing the IDs
func (head *LHHead) Dup() *LHHead {
	return head.dup(false)
}

//DupKeepIDs duplicate the head and adaptor, keeping the IDs the same
func (head *LHHead) DupKeepIDs() *LHHead {
	return head.dup(true)
}

func (head *LHHead) dup(keepIDs bool) *LHHead {
	h := NewLHHead(head.Name, head.Manufacturer, nil)
	h.TipLoading = head.TipLoading
	if keepIDs {
		h.ID = head.ID
		h.Params = head.Params.DupKeepIDs()
		h.Adaptor = head.Adaptor.DupKeepIDs()
	} else {
		h.Params = head.Params.Dup()
		h.Adaptor = head.Adaptor.Dup()
	}
	return h
}

//GetParams get the channel parameters of the head or the adaptor if one is loaded
func (lhh *LHHead) GetParams() *LHChannelParameter {
	if lhh.Adaptor == nil {
		return lhh.Params
	} else {
		return lhh.Adaptor.GetParams()
	}
}

//CanReach return true if the head can service the given addresses in the given object
//simultaneously.
//addresses is a slice of well addresses which should be serviced by successive channels of
//the head, eg. ["A1", "B1", "", "D1",] means channels 0, 1, and 3 should address wells
//A1, B1, and D1 respectively and channels 2 and 4-7 should not be used.
//Repeated addresses (e.g. ["A1", "A1", "A1"]) imply multiple tips per well, with
//exact positioning for each tip calculated with LHHead.GetWellTargets().
//Addresses are not reordered, and so ["A1", "B1"] != ["B1", "A1"].
func (head *LHHead) CanReach(plate *LHPlate, addresses WellCoordSlice) bool {
	if head.Adaptor == nil {
		return false
	}
	if !head.Adaptor.Params.Independent {
		addresses = addresses.Trim()
	}
	if len(addresses) > head.GetParams().Multi {
		return false
	}
	if len(addresses) == 0 {
		return true
	}

	wellTargets := head.Adaptor.GetWellTargets(plate.Welltype)

	//get the real world position of the addresses
	coords := make([]*Coordinates2D, head.GetParams().Multi)
	counter := make(addressCounter, len(addresses))
	for channel, address := range addresses {
		//indicates the channel should not be used
		if address.IsZero() {
			continue
		}

		//we're not particularly interested in the z dimension, so just use topreference
		crd, ok := plate.WellCoordsToCoords(address, TopReference)
		if !ok {
			//can't address a well which doesn't exist
			return false
		}

		crd2D := crd.To2D()

		//add well target offset if the well supports multiple tips at once
		if wellTargets != nil {
			crd2D = crd2D.Add(wellTargets[counter.GetCount(address)%len(wellTargets)])
		}

		coords[channel] = &crd2D

		counter.Increment(address)
	}

	if !head.GetParams().Independent {
		//non-independent heads can only use contiguous channels
		seenNil := false
		for _, crd := range coords {
			if seenNil && crd != nil {
				return false
			}
			seenNil = seenNil || crd == nil
		}
	}

	//offsets[i] is the distance between the i-1th and the ith coordinate
	offsets := make([]*Coordinates2D, len(coords)-1)
	lastI := -1
	for i, coord := range coords {
		if coord != nil {
			if lastI >= 0 {
				//divide the total offset since lastI evenly evenly
				offset := coord.Subtract(*coords[lastI]).Divide(float64(i - lastI))
				offsets[i-1] = &offset
			}
			lastI = i
		}
	}

	//check that each offset is within the supported range of the machine
	supportedRange := NewRectangle(head.Adaptor.GetSmallestChannelSpacing(), head.Adaptor.GetLargestChannelSpacing())
	positionAccuracy := 0.1 //should be a property of the head
	supportedRange = supportedRange.Expand(positionAccuracy)

	for _, offset := range offsets {
		if offset != nil && !supportedRange.Contains(*offset) {
			return false
		}
	}

	//hjk: TODO check that offsets are equal if the head suppots only uniform offsets

	return true
}

//count occrances of addresses
type addressCounter map[WellCoords]int

//GetCount how many times have we seen this address
func (self addressCounter) GetCount(wc WellCoords) int {
	if count, ok := self[wc]; ok {
		return count
	}
	return 0
}

//Increment the counter for this address
func (self addressCounter) Increment(wc WellCoords) {
	if _, ok := self[wc]; !ok {
		self[wc] = 0
	}
	self[wc] += 1
}

type SequentialTipLoadingBehaviour int

const (
	//NoSequentialTipLoading tips are loaded all at once, an error is raised if not possible
	NoSequentialTipLoading SequentialTipLoadingBehaviour = iota
	//ForwardSequentialTipLoading chunks of contiguous tips are loaded sequentially in the order encountered
	ForwardSequentialTipLoading
	//ReverseSequentialTipLoading chunks of contiguous tips are loaded sequentially in reverse order
	ReverseSequentialTipLoading
)

var sequentialTipLoadingBehaviourNames = map[SequentialTipLoadingBehaviour]string{
	NoSequentialTipLoading:      "no sequential tip loading",
	ForwardSequentialTipLoading: "forward sequential tip loading",
	ReverseSequentialTipLoading: "reverse sequential tip loading",
}

func (s SequentialTipLoadingBehaviour) String() string {
	return sequentialTipLoadingBehaviourNames[s]
}

//TipLoadingBehaviour describe the way in which tips are loaded
type TipLoadingBehaviour struct {
	//OverrideLoadTipsCommand true it the liquid handler will override which tips are loaded
	OverrideLoadTipsCommand bool
	//AutoRefillTipboxes are tipboxes automaticall refilled
	AutoRefillTipboxes bool
	//LoadingOrder are tips loaded ColumnWise or RowWise
	LoadingOrder MajorOrder
	//VerticalLoadingDirection the direction along which columns are loaded
	VerticalLoadingDirection VerticalDirection
	//HorizontalLoadingDirection the direction along which rows are loaded
	HorizontalLoadingDirection HorizontalDirection
	//ChunkingBehaviour how to load tips when the requested number aren't available contiguously
	ChunkingBehaviour SequentialTipLoadingBehaviour
}

//String get a string description for debuggin
func (s TipLoadingBehaviour) String() string {

	autoRefill := ""
	if !s.AutoRefillTipboxes {
		autoRefill = "no "
	}

	if !s.OverrideLoadTipsCommand {
		return fmt.Sprintf("%sauto-refilling, no loading override", autoRefill)
	}

	return fmt.Sprintf("%sauto-refilling, loading order: %v, %v, %v, %v", autoRefill, s.LoadingOrder, s.VerticalLoadingDirection, s.HorizontalLoadingDirection, s.ChunkingBehaviour)
}

// GetBehaviour get the coordinates of the tips which would be loaded by the tiploading behaviour
// when asked to load num tips from the given tipbox.
// Returns slices of wellcoords, e.g. [[G1, H1], [A2, B2]], where in this case [G1, H1] would be loaded
// simultaneously followed by [A2, B2].
// if num is greater than the tips remaining, the tipbox will be refilled if
// AutoRefillTipboxes is true, otherwise returns an error if there aren't enough tips
func (self *TipLoadingBehaviour) GetBehaviour(tb *LHTipbox, num int) ([][]WellCoords, error) {
	var ret [][]WellCoords

	if !tb.HasEnoughTips(num) {
		if self.AutoRefillTipboxes {
			tb.Refill()
		} else {
			return nil, errors.Errorf("not enough tips in tipbox: require %d, remaining %d", num, tb.N_clean_tips())
		}
	}

	it := NewAddressIterator(tb,
		self.LoadingOrder,
		self.VerticalLoadingDirection,
		self.HorizontalLoadingDirection,
		false)

	tipsRemaining := num
	var lastTipCoord WellCoords
	currChunk := make([]WellCoords, 0, num)
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		//start a new chunk if this chunk has something in it AND (we found an empty position OR we changed row and column)
		if len(currChunk) > 0 && (!tb.HasTipAt(wc) || (lastTipCoord.X != wc.X && lastTipCoord.Y != wc.Y)) {
			//keep the chunk if either this chunk provides all the tips we need or we can load it sequentially
			if !(self.ChunkingBehaviour == NoSequentialTipLoading && len(currChunk) < tipsRemaining) {
				ret = append(ret, currChunk)
				tipsRemaining -= len(currChunk)
			}
			currChunk = make([]WellCoords, 0, tipsRemaining)
		}
		//if we have all the chunks we need
		if len(currChunk) >= tipsRemaining {
			break
		}
		//add the next tip
		if tb.HasTipAt(wc) {
			currChunk = append(currChunk, wc)
			lastTipCoord = wc
		}
	}
	if len(currChunk) > 0 {
		ret = append(ret, currChunk)
		tipsRemaining -= len(currChunk)
	}

	if self.ChunkingBehaviour == ReverseSequentialTipLoading {
		//apparently this is actually the recommended way to reverse a list in place
		for i := len(ret)/2 - 1; i >= 0; i-- {
			opp := len(ret) - 1 - i
			ret[i], ret[opp] = ret[opp], ret[i]
		}

		for _, chunk := range ret {
			for i := len(chunk)/2 - 1; i >= 0; i-- {
				opp := len(chunk) - 1 - i
				chunk[i], chunk[opp] = chunk[opp], chunk[i]
			}
		}
	}

	if tipsRemaining > 0 {
		return ret, errors.New("not enough tips in tipbox")
	}

	return ret, nil
}
