package wtype

// head
type LHHead struct {
	Name         string
	Manufacturer string
	ID           string
	Adaptor      *LHAdaptor
	Params       *LHChannelParameter
	TipLoading   TipLoadingBehaviour // defines the way in which tips are loaded into the head
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
	lastAddress := ZeroWellCoords()
	for channel, address := range addresses {
		//indicates the channel should not be used
		if address.IsZero() {
			if !lastAddress.IsZero() {
				// add to the well counter for the last address so that we leave
				// a space if we come back to it with a later channel
				counter.Increment(lastAddress)
			}
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
		lastAddress = address
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
