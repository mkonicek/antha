package wtype

import (
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

func CopyComponentArray(arin []*LHComponent) []*LHComponent {
	r := make([]*LHComponent, len(arin))

	for i, v := range arin {
		r[i] = v.Dup()
	}

	return r
}

// are tips going to align to wells?
func TipsWellsAligned(prm LHChannelParameter, plt LHPlate, wellsfrom []string) bool {
	// heads which can do independent multichanneling are dealt with separately

	if prm.Independent {
		return disContiguousTipsWellsAligned(prm, plt, wellsfrom)
	} else {
		return contiguousTipsWellsAligned(prm, plt, wellsfrom)
	}
}

func disContiguousTipsWellsAligned(prm LHChannelParameter, plt LHPlate, wellsfrom []string) bool {
	// inflate wellsfrom to full multichannel size
	fullWellsFrom, ok := expandWellsFrom(prm.Orientation, plt, wellsfrom)

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

func expandWellsFrom(orientation int, plt LHPlate, wellsfrom []string) ([]string, bool) {
	wcArr := WCArrayFromStrings(wellsfrom)

	var wells []*LHWell

	// get column or row

	switch orientation {
	case LHHChannel:
		wc := WCArrayRows(wcArr)
		if len(wc) != 1 {
			return []string{}, false
		}

		rowIndex := wc[0]

		if rowIndex < 0 || rowIndex >= len(plt.Rows) {
			return []string{}, false
		}

		wells = plt.Rows[rowIndex]

	case LHVChannel:
		wc := WCArrayCols(wcArr)
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

	WCs := A1ArrayFromWells(wells)

	ret := make([]string, len(WCs))

	for i := 0; i < len(WCs); i++ {
		if isInArr(WCs[i], wellsfrom) {
			ret[i] = WCs[i]
		}
	}

	return ret, true
}

func contiguousTipsWellsAligned(prm LHChannelParameter, plt LHPlate, wellsfrom []string) bool {
	// 1) find well coords for channels given parameters
	// 2) compare to wells requested

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

func ChannelWells(prm LHChannelParameter, plt LHPlate, wellsfrom []string) []string {
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

	tipsperwell, wellskip := TipsPerWell(prm, plt)

	tipwells := make([]string, len(wellsfrom))

	wc := MakeWellCoords(firstwell)

	fwc := wc.Y

	if prm.Orientation == LHHChannel {
		fwc = wc.X
	}

	ticker := Ticker{TickEvery: tipsperwell, TickBy: wellskip + 1, Val: fwc}

	for i := 0; i < len(wellsfrom); i++ {
		if channelsused[i] {
			tipwells[i] = wc.FormatA1()
		}

		ticker.Tick()

		if prm.Orientation == LHVChannel {
			wc.Y = ticker.Val
		} else {
			wc.X = ticker.Val
		}
	}

	return tipwells
}

func TipsPerWell(prm LHChannelParameter, p LHPlate) (int, int) {
	// assumptions:

	// 1) sbs format plate
	// 2) head pitch matches usual 96-well format

	if prm.Multi == 1 {
		return 1, 0
	}

	nwells := 1
	ntips := prm.Multi

	if prm.Orientation == LHVChannel {
		if prm.Multi != 8 {
			panic("Unsupported V head format (must be 8)")
		}

		nwells = p.WellsY()
	} else if prm.Orientation == LHHChannel {
		if prm.Multi != 12 {
			panic("Unsupported H head format (must be 12)")
		}
		nwells = p.WellsX()
	} else {
		// empty
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
