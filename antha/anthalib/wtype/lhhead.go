package wtype

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
	var params *LHChannelParameter
	var adaptor *LHAdaptor
	if keepIDs {
		params = head.Params.DupKeepIDs()
		adaptor = head.Adaptor.DupKeepIDs()
	} else {
		params = head.Params.Dup()
		adaptor = head.Adaptor.Dup()
	}
	h := NewLHHead(head.Name, head.Manufacturer, params)
	h.Adaptor = adaptor
	h.TipLoading = head.TipLoading
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

//GetWellTargets get the offset from the center of the well for each channel that
//can access the well simultaneously
//returns nil if the well cannot fit multiple channels
func (head *LHHead) GetWellTargets(well *LHWell) []Coordinates2D {
	//no well targets if there's no adaptor loaded
	if head.Adaptor == nil {
		return nil
	}
	channelPositions := head.Adaptor.GetMostCompactChannelPositions()
	channelRadius := 3.0 //should come from head

	//total size of channels, including radius
	channelSize := channelPositions[len(channelPositions)-1].Add(Coordinates2D{X: 2 * channelRadius, Y: 2 * channelRadius})

	if wellSize := well.GetSize(); wellSize.X < channelSize.X || wellSize.Y < channelSize.Y {
		return nil
	}

	center := Coordinates2D{}
	for _, pos := range channelPositions {
		center = center.Add(pos)
	}
	center = center.Divide(float64(len(channelPositions)))

	for i := range channelPositions {
		channelPositions[i] = channelPositions[i].Subtract(center)
	}

	return channelPositions
}

//CanReach return true if the head can service the given addresses in the given object
//simultaneously.
//addresses is a slice of well addresses which should be serviced by successive channels of
//the head, eg. ["A1", "B1", "", "D1",] means channels 0, 1, and 3 should address wells
//A1, B1, and D1 respectively and channels 2 and 4-7 should not be used.
//Repeated addresses (e.g. ["A1", "A1", "A1"]) imply multiple tips per well, with
//exact positioning for each tip calculated with LHHead.GetWellTargets().
//Addresses are not reordered, and so ["A1", "B1"] != ["B1", "A1"].
func (head *LHHead) CanReach(plate *LHPlate, addresses []WellCoords) bool {
	if head.Adaptor == nil {
		return false
	}
	if len(addresses) > head.GetParams().Multi {
		return false
	}
	if len(addresses) == 0 {
		return true
	}

	wellTargets := head.GetWellTargets(plate.Welltype)

	//get the real world position of the addresses
	coords := make([]*Coordinates2D, head.GetParams().Multi)
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
			crd2D = crd2D.Add(wellTargets[channel%len(wellTargets)])
		}

		coords[channel] = &crd2D
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
	supportedRange := NewRectangle(head.Adaptor.GetSmallestChannelOffset(), head.Adaptor.GetLargestChannelOffset())
	positionAccuracy := 0.1 //should be a property of the head
	for _, offset := range offsets {
		if offset != nil && supportedRange.distanceOutside(*offset) > positionAccuracy {
			return false
		}
	}

	//hjk: TODO check that offsets are equal if the head suppots only uniform offsets

	return true
}
