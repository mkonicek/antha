package wtype

import "fmt"

// A Tip Estimate provides information on how many tips and tip boxes of a given type are expected to be used
type TipEstimate struct {
	TipType   string // identifier describing which tips are to be used
	NTips     int    // count of tips used total
	NTipBoxes int    // count of tip boxes to be used
}

func (self TipEstimate) String() string {
	return fmt.Sprintf("%s: %d tips in %d boxes", self.TipType, self.NTips, self.NTipBoxes)
}
