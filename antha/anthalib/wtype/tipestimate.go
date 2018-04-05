package wtype

// A Tip Estimate provides information on how many tips and tip boxes of a given type are expected to be used
type TipEstimate struct {
	TipType   string // identifier describing which tips are to be used
	NTips     int    // count of tips used total
	NTipBoxes int    // count of tip boxes to be used
}
