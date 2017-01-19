package wtype

func CopyComponentArray(arin []*LHComponent) []*LHComponent {
	r := make([]*LHComponent, len(arin))

	for i, v := range arin {
		r[i] = v.Dup()
	}

	return r
}

func canGet(want, got ComponentVector) bool {
	for i := 0; i < len(want); i++ {
		// is there, like, stuff where we need it?

		if want[i] == nil && got[i] == nil {
			continue
		} else if (want[i] == nil && got[i] != nil) || (want[i] != nil && got[i] == nil) {
			return false
		}

		// check the component type and junk

		if want[i].CName != got[i].CName {
			return false
		}

		// finally is there enough?

		if got[i].Volume().LessThan(want[i].Volume()) {
			return false
		}
	}

	// like, whatever
	return true
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
	// 1	8	: {1,0} (single channel, 96-well plate)
	// 8	8	: {1,0} (8-channels, 96-well)
	// 8	16	: {1,1} (8 channels, 384-well)
	// 8	32	: {1,3} (8 channels, 1536 plate)
	// 8    4	: {2,0} (8 channels, 24 plate)
	// 8 	1	: {8,0} (8 channels, 12 well strip)
	// 8	2	: {3,0} (8 channels, 6 or 8 well plate)

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
