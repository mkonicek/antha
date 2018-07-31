package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"reflect"
)

func CopyComponentArray(arin []*wtype.Liquid) []*wtype.Liquid {
	r := make([]*wtype.Liquid, len(arin))

	for i, v := range arin {
		r[i] = v.Dup()
	}

	return r
}

//GetChannelOffset get the smallest possible distance between successive channels
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

//GetMostCompactChannelPositions get the relative channel positions for the head
//in the most tightly bunched layout supported
func GetMostCompactChannelPositions(head *wtype.LHHead) []wtype.Coordinates {
	ret := make([]wtype.Coordinates, head.GetParams().Multi)
	offset := GetSmallestChannelOffset(head)
	last := wtype.Coordinates{}

	for i := range ret {
		last = last.Add(offset)
		ret[i] = last
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

	channelOffset := GetSmallestChannelOffset(head)
	machineTol := 0.1 //should also come from the head

	//check that the positioning requested for each channel is allowable
	var lastPos *wtype.Coordinates
	var expectedOffset wtype.Coordinates
	for _, pos := range coords {
		if lastPos == nil {
			lastPos = pos
			continue
		}
		expectedOffset = expectedOffset.Add(channelOffset)
		if pos != nil {
			actualOffset := pos.Subtract(*lastPos)
			//if the difference in the XY direction only is greater than machine error
			if actualOffset.Subtract(expectedOffset).AbsXY() > machineTol {
				return false
			}
			lastPos = pos
			expectedOffset = wtype.Coordinates{}
		}
	}

	return true
}

func TipsWellsAligned(head *wtype.LHHead, plt *wtype.Plate, wellsfrom []string) bool {

	// heads which can do independent multichanneling are dealt with separately
	if head.Adaptor.Params.Independent {
		return disContiguousTipsWellsAligned(head, plt, wellsfrom)
	} else {
		return contiguousTipsWellsAligned(head, plt, wellsfrom)
	}
}

func disContiguousTipsWellsAligned(head *wtype.LHHead, plt *wtype.Plate, wellsfrom []string) bool {
	prm := head.Adaptor.Params
	// inflate wellsfrom to full multichannel size
	fullWellsFrom, ok := expandWellsFrom(prm.Orientation, *plt, wellsfrom)

	// not in right orientation, not in single column etc. etc
	if !ok {
		return false
	}

	// get the wells pointed at by the channels

	channelWells := ChannelWells(prm, plt, fullWellsFrom)

	// now does this work?

	return reflect.DeepEqual(channelWells, fullWellsFrom)
}

func isInArr(s string, sa []string) bool {
	for _, ss := range sa {
		if ss == s {
			return true
		}
	}

	return false
}

func expandWellsFrom(orientation wtype.ChannelOrientation, plt wtype.Plate, wellsfrom []string) ([]string, bool) {
	wcArr := wtype.WCArrayFromStrings(wellsfrom)

	var wells []*wtype.LHWell

	// get column or row

	switch orientation {
	case wtype.LHHChannel:
		wc := wtype.WCArrayRows(wcArr)
		if len(wc) != 1 {
			return []string{}, false
		}

		rowIndex := wc[0]

		if rowIndex < 0 || rowIndex >= len(plt.Rows) {
			return []string{}, false
		}

		wells = plt.Rows[rowIndex]

	case wtype.LHVChannel:
		wc := wtype.WCArrayCols(wcArr)
		if len(wc) != 1 {
			return []string{}, false
		}

		colIndex := wc[0]

		if colIndex < 0 || colIndex >= len(plt.Cols) {
			return []string{}, false
		}

		wells = plt.Cols[colIndex]

	default:
		return []string{}, false
	}

	WCs := wtype.A1ArrayFromWells(wells)

	ret := make([]string, len(WCs))

	for i := 0; i < len(WCs); i++ {
		if isInArr(WCs[i], wellsfrom) {
			ret[i] = WCs[i]
		}
	}

	return ret, true
}

func assertWFContiguousNonEmpty(sa []string) bool {
	// this needs to specifically disallow arrays containing "X", ""+, "X"
	found := false
	gap := false

	for _, s := range sa {
		if s == "" {
			if found {
				gap = true
			}
		} else {
			if gap {
				return false
			}
			found = true
		}
	}

	return found
}

func contiguousTipsWellsAligned(head *wtype.LHHead, plt *wtype.Plate, wellsfrom []string) bool {

	//can't multi channel with 2-7 wells per column
	if plt.NRows() != 1 && plt.NRows()%8 != 0 {
		return false
	}
	// guarantee wellsfrom is contiguous and has at least one non-""
	// trailing "" are OK
	if !assertWFContiguousNonEmpty(wellsfrom) {
		return false
	}

	// if we've only got one row, check that there are well targets for us to move into
	if plt.NRows() == 1 {
		return plt.AreWellTargetsEnabled(head.GetParams().Multi, 9.0)
	}

	// if this is something like a standard sbs-format plate, i.e. wells in a single, continuous space
	// 1) find well coords for channels given parameters
	// 2) compare to wells requested

	prm := head.Adaptor.Params
	channelwells := ChannelWells(prm, plt, wellsfrom)

	// only works if all are filled

	return wutil.StringArrayEqual(channelwells, wellsfrom)
}

func ChannelsUsed(wf []string) []bool {
	ret := make([]bool, len(wf))

	for i := 0; i < len(wf); i++ {
		if wf[i] != "" {
			ret[i] = true
		}
	}

	return ret
}

func ChannelWells(prm *wtype.LHChannelParameter, plt *wtype.Plate, wellsfrom []string) []string {
	channelsused := ChannelsUsed(wellsfrom)

	firstwell := ""

	for i := 0; i < len(wellsfrom); i++ {
		if channelsused[i] {
			firstwell = wellsfrom[i]
			break
		}
	}

	if firstwell == "" {
		panic("Empty channel array passed to transferinstruction")
	}

	tipsperwell, wellskip := TipsPerWell(*prm, *plt)

	tipwells := make([]string, len(wellsfrom))

	wc := wtype.MakeWellCoords(firstwell)

	fwc := wc.Y

	if prm.Orientation == wtype.LHHChannel {
		fwc = wc.X
	}

	ticker := wtype.Ticker{TickEvery: tipsperwell, TickBy: wellskip + 1, Val: fwc}

	for i := 0; i < len(wellsfrom); i++ {
		if channelsused[i] {
			tipwells[i] = wc.FormatA1()
		}

		ticker.Tick()

		if prm.Orientation == wtype.LHVChannel {
			wc.Y = ticker.Val
		} else {
			wc.X = ticker.Val
		}
	}

	return tipwells
}

func TipsPerWell(prm wtype.LHChannelParameter, p wtype.Plate) (int, int) {
	// assumptions:

	// 1) sbs format plate
	// 2) head pitch matches usual 96-well format

	if prm.Multi == 1 {
		return 1, 0
	}

	nwells := 1
	ntips := prm.Multi

	if prm.Orientation == wtype.LHVChannel {
		if prm.Multi != 8 {
			panic("Unsupported V head format (must be 8)")
		}

		nwells = p.WellsY()
	} else if prm.Orientation == wtype.LHHChannel {
		if prm.Multi != 12 {
			panic("Unsupported H head format (must be 12)")
		}
		nwells = p.WellsX()
	} else {
		panic("unknown orientation")
	}

	// how many  tips fit into one well
	// how many wells are skipped between each tip

	// examples:
	// 1	8	: {1,0}  (single channel, 96-well plate)
	// 8	8	: {1,0}  (8 channels, 96-well)
	// 8	16	: {1,1}  (8 channels, 384-well)
	// 8	32	: {1,3}  (8 channels, 1536 plate)
	// 8    4	: {2,0}  (8 channels, 24 plate)
	// 8 	1	: {8,0}  (8 channels, 12 well strip)
	// 8	2	: {3,0}  (8 channels, 6 or 8 well plate)

	tpw := 1

	if ntips > nwells {
		tpw = ntips / nwells
	}

	wellskip := 0

	if nwells > ntips {
		wellskip = (nwells / ntips) - 1
	}

	return tpw, wellskip
}

func FirstIndexInStrArray(s string, a []string) int {
	for i, v := range a {
		if v == s {
			return i
		}
	}

	return -1
}

func physicalTipCheck(head *wtype.LHHead, plt *wtype.Plate, wellsFrom []string) bool { //nolint used by tests
	// assumptions
	// 1 - first tip is aligned with the middle of first well
	// 2 - each subsequent tip is a constant (9mm) distance from the first in whichever orientation
	// 3 - tips are always the same size (deliberate overestimate) at 8mm diameter
	// the aim is to be conservative: the consequences of being wrong will be slightly less multichanneling

	// parameters for assumptions above
	// XXX TODO these need to come from the robot and inventory
	coneDist := 9.0
	tipRadius := 4.0

	// trim wellsFrom

	trim := func(sa []string) []string {
		r := make([]string, 0, len(sa))

		for _, v := range sa {
			if v != "" {
				r = append(r, v)
			}
		}

		return r
	}

	trimmedWellsFrom := trim(wellsFrom)

	// get well coords

	wellCoords := make(wtype.PointSet, 0, len(trimmedWellsFrom))

	for _, wfs := range trimmedWellsFrom {
		wc := wtype.MakeWellCoords(wfs)
		wca, _ := plt.WellCoordsToCoords(wc, wtype.BottomReference)
		wellCoords = append(wellCoords, wca)
	}

	// normalize to first well

	wellCoords = wellCoords.CentreTo(wellCoords[0])

	// now iterate through the coordinates for the tips and ensure they are *inside* their corresponding wells

	tipCoords := wtype.Coordinates{}

	tipIncr := wtype.Coordinates{Y: coneDist}

	if head.Adaptor.Params.Orientation == wtype.LHHChannel {
		tipIncr = wtype.Coordinates{X: coneDist}
	}

	for i := range trimmedWellsFrom {
		if tipClash(tipCoords, tipRadius, wellCoords[i], plt.Welltype) {
			return false
		}

		tipCoords = tipCoords.Add(tipIncr)
	}

	return true
}

func tipClash(tipCoords wtype.Coordinates, tipRadius float64, wellCoords wtype.Coordinates, wellType *wtype.LHWell) bool {
	if wellType.Shape().H == wellType.Shape().W {
		// square or round wells
		dist := tipCoords.Subtract(wellCoords).Abs()
		dist += tipRadius
		wellRadius := wellType.Shape().H / 2.0

		return wellRadius < dist
	} else {
		// not presently implemented
		return true
	}
}
