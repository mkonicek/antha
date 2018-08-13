package wtype

import "testing"

func assertIDsEqual(t *testing.T, a, b *LHAdaptor) {
	if a.ID != b.ID {
		t.Error("Adaptor.ID changed")
	}
	if a.Params.ID != b.Params.ID {
		t.Error("Adaptor.Params.ID changed")
	}

	for i := range a.Tips {
		if aT, bT := a.Tips[i], b.Tips[i]; aT != nil && bT != nil {
			if aT.ID != bT.ID {
				t.Errorf("Adaptor.Tips[%d] id changed", i)
			}
		}
	}
}

func assertIDsChanged(t *testing.T, a, b *LHAdaptor) {
	if a.ID == b.ID {
		t.Error("Adaptor.ID not changed")
	}
	if a.Params.ID == b.Params.ID {
		t.Error("Adaptor.Params.ID not changed")
	}

	for i := range a.Tips {
		if aT, bT := a.Tips[i], b.Tips[i]; aT != nil && bT != nil {
			if aT.ID == bT.ID {
				t.Errorf("Adaptor.Tips[%d] id not changed", i)
			}
		}
	}
}

func assertAdaptorTipLengths(t *testing.T, a, b *LHAdaptor, msg string) {
	if aT, bT := len(a.Tips), len(b.Tips); aT != bT {
		t.Errorf("Adaptor.Tips %s not the same length, %d != %d", msg, aT, bT)
		return
	}
	for i := range a.Tips {
		if aNil, bNil := a.Tips[i] != nil, b.Tips[i] != nil; aNil != bNil {
			t.Errorf("Adaptor.Tips[%d] (%s) tip presence didn't match, %t, %t", i, msg, aNil, bNil)
		}
	}
}

func TestLHAdaptorDup(t *testing.T) {
	params := &LHChannelParameter{
		Multi: 8,
	}

	adaptor := NewLHAdaptor("testName", "testMfr", params)

	//add some tips
	for i := 0; i < 8; i = i + 2 {
		adaptor.AddTip(i, makeTipForTest())
	}

	newIDs := adaptor.Dup()
	oldIDs := adaptor.DupKeepIDs()

	assertAdaptorTipLengths(t, adaptor, newIDs, "newIDs")
	assertAdaptorTipLengths(t, adaptor, oldIDs, "oldIDs")

	assertIDsChanged(t, adaptor, newIDs)
	assertIDsEqual(t, adaptor, oldIDs)

	//parameter changes shouldn't affect copies
	params.Multi = 5
	if newIDs.Params.Multi != 8 {
		t.Error("params not duplicated for newIDs")
	}
	if oldIDs.Params.Multi != 8 {
		t.Error("params not duplicated for keepIds")
	}

	//changing tip info shouldn't affect copies
	for _, tip := range adaptor.Tips {
		if tip != nil {
			tip.ID = "CHANGEDID"
		}
	}
	for i := range newIDs.Tips {
		if tip := newIDs.Tips[i]; tip != nil {
			if tip.ID == "CHANGEDID" {
				t.Error("Tips not duplicated for newIDs")
				break
			}
		}
	}
	for i := range newIDs.Tips {
		if tip := oldIDs.Tips[i]; tip != nil {
			if tip.ID == "CHANGEDID" {
				t.Error("Tips not duplicated for oldIDs")
				break
			}
		}
	}

	//removing tips shouldn't affect copies
	adaptor.RemoveTips()
	if g, e := newIDs.NumTipsLoaded(), 4; g != e {
		t.Errorf("expected %d tips, got %d for newIds", e, g)
	}
	if g, e := oldIDs.NumTipsLoaded(), 4; g != e {
		t.Errorf("expected %d tips, got %d for oldIDs", e, g)
	}

}
