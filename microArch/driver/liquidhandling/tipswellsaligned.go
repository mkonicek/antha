package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func CopyComponentArray(arin []*wtype.Liquid) []*wtype.Liquid {
	r := make([]*wtype.Liquid, len(arin))

	for i, v := range arin {
		r[i] = v.Dup()
	}

	return r
}

//GetSmallestChannelOffset get the smallest possible distance between successive channels
func GetSmallestChannelOffset(head *wtype.LHHead) wtype.Coordinates {

	//hjk: currently assume a constant fixed offset between channels
	//     this will need updating when we support better reporting of head
	//     capabilities wrt. independent multi-channelling
	channelStep := 9.0 //9mm between channels
	channelOffset := wtype.Coordinates{}
	if head.GetParams().Orientation == wtype.LHVChannel {
		channelOffset.Y = channelStep
	} else {
		channelOffset.X = channelStep
	}

	return channelOffset
}

//GetLargestChannelOffset get the largest possible distance between successive channels
func GetLargestChannelOffset(head *wtype.LHHead) wtype.Coordinates {
	//equal to smallest if independent
	if !head.GetParams().Independent {
		return GetSmallestChannelOffset(head)
	}

	//completely arbitrary for now since we don't report this
	return GetSmallestChannelOffset(head).Multiply(2.0)
}

//GetMostCompactChannelPositions get the relative channel positions for the head
//in the most tightly bunched layout supported
func GetMostCompactChannelPositions(head *wtype.LHHead) []wtype.Coordinates {
	ret := make([]wtype.Coordinates, head.GetParams().Multi)
	offset := GetSmallestChannelOffset(head)
	last := wtype.Coordinates{}

	for i := range ret {
		ret[i] = last
		last = last.Add(offset)
	}

	return ret
}

//GetWellTargets get the offset from the center of the well for each channel that
//can access the well simultaneously
//returns nil if the well cannot fit multiple channels
func GetWellTargets(head *wtype.LHHead, well *wtype.LHWell) []wtype.Coordinates {
	channelPositions := GetMostCompactChannelPositions(head)
	channelRadius := 3.0 //should come from head

	//total size of channels, including radius
	channelSize := channelPositions[len(channelPositions)-1].Add(wtype.Coordinates{X: 2 * channelRadius, Y: 2 * channelRadius})

	if wellSize := well.GetSize(); wellSize.X < channelSize.X || wellSize.Y < channelSize.Y {
		return nil
	}

	center := wtype.Coordinates{}
	for _, pos := range channelPositions {
		center = center.Add(pos)
	}
	center = center.Divide(float64(len(channelPositions)))

	for i := range channelPositions {
		channelPositions[i] = channelPositions[i].Subtract(center)
	}

	return channelPositions
}

//CanHeadReach return true if the head can service the given addresses in the given object
//simultaneously.
//addresses is a slice of well addresses which should be serviced by successive channels of
//the head, eg. ["A1", "B1", "", "D1",] means channels 0, 1, and 3 should address wells
//A1, B1, and D1 respectively and channels 2 and 4-7 should not be used
func CanHeadReach(head *wtype.LHHead, plate *wtype.LHPlate, addresses []wtype.WellCoords) bool {
	if len(addresses) > head.GetParams().Multi {
		return false
	}
	if len(addresses) == 0 {
		return false
	}

	wellTargets := GetWellTargets(head, plate.Welltype)

	//get the real world position of the addresses
	coords := make([]*wtype.Coordinates, head.GetParams().Multi)
	for channel, address := range addresses {
		//indicates the channel should not be used
		if address.IsZero() {
			continue
		}

		//we're not particularly interested in the z dimension, so just use topreference
		crd, ok := plate.WellCoordsToCoords(address, wtype.TopReference)
		if !ok {
			//can't address a well which doesn't exist
			return false
		}

		//add well target offset if the well supports multiple tips at once
		if wellTargets != nil {
			crd = crd.Add(wellTargets[channel%len(wellTargets)])
		}

		coords[channel] = &crd
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
	offsets := make([]*wtype.Coordinates, len(coords)-1)
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
	smallestOffset := GetSmallestChannelOffset(head)
	largestOffset := GetLargestChannelOffset(head)
	positionAccuracy := 0.1 //should be a property of the head
	for _, offset := range offsets {
		if offset != nil && distanceOutsideSquare(smallestOffset, largestOffset, *offset) > positionAccuracy {
			return false
		}
	}

	//hjk: TODO check that offsets are equal if the head suppots only uniform offsets

	return true
}

//distanceOutsideSquare return a lower bound on how far position is outside the
//square defined by the corners lowerLeft,topRight
//where a is the bottom left and b is the top right corner
//value will be negative if position is inside the square
//z-value for all arguments is ignored
func distanceOutsideSquare(lowerLeft, topRight, pos wtype.Coordinates) float64 {
	ret := lowerLeft.X - pos.X
	if r := pos.X - topRight.X; r > ret {
		ret = r
	}
	if r := lowerLeft.Y - pos.Y; r > ret {
		ret = r
	}
	if r := pos.Y - topRight.Y; r > ret {
		ret = r
	}
	return ret
}
