package liquidhandling

// consts for generic liquid handling types
type LiquidHandlerLevel uint8

const (
	LLLiquidHandler LiquidHandlerLevel = iota // requires detailed programming e.g. move, aspirate, move dispense etc.
	HLLiquidHandler                           // can orchestrate liquid transfers itself
	maxLiquidHandler
)

var lhLevels = []string{
	LLLiquidHandler: "low level",
	HLLiquidHandler: "high level",
}

func (lhl LiquidHandlerLevel) String() string {
	if lhl < maxLiquidHandler {
		return lhLevels[lhl]
	}
	panic("unknown liquid handler level")
}

func (lhl LiquidHandlerLevel) IsValid() bool {
	return lhl < maxLiquidHandler
}

// TipType types of tips used by liquid handlers
type TipType uint8

const (
	DisposableTips              TipType = iota // disposable system	-- needs tip boxes & waste
	FixedTips                                  // fixed tip system	-- needs tip wash
	MixedDisposableAndFixedTips                // both disposable and mixed	-- needs all of the above
	NoTips                                     // does not use tips
	maxTipTypes
)

var tipNames = []string{
	DisposableTips:              "Disposable",
	FixedTips:                   "Fixed",
	MixedDisposableAndFixedTips: "Mixed",
	NoTips:                      "None",
}

func (tt TipType) String() string {
	if tt < maxTipTypes {
		return tipNames[tt]
	}
	panic("unknown tip type")
}

func (tt TipType) IsValid() bool {
	return tt < maxTipTypes
}
