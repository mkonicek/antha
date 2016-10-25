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
