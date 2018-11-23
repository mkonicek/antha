package liquidhandling

// consts for generic liquid handling types
type LiquidHandlerLevel int

const (
	LLLiquidHandler LiquidHandlerLevel = iota // requires detailed programming e.g. move, aspirate, move dispense etc.
	HLLiquidHandler                           // can orchestrate liquid transfers itself
)

var lhLevels = map[LiquidHandlerLevel]string{
	LLLiquidHandler: "low level",
	HLLiquidHandler: "high level",
}

func (lhl LiquidHandlerLevel) String() string {
	if r, ok := lhLevels[lhl]; ok {
		return r
	}
	panic("unknown liquid handler level")
}

func (lhl LiquidHandlerLevel) IsValid() bool {
	_, ok := lhLevels[lhl]
	return ok
}

// TipType types of tips used by liquid handlers
type TipType int

const (
	DisposableTips              TipType = iota // disposable system	-- needs tip boxes & waste
	FixedTips                                  // fixed tip system	-- needs tip wash
	MixedDisposableAndFixedTips                // both disposable and mixed	-- needs all of the above
	NoTips                                     // does not use tips
)

var tipNames = map[TipType]string{
	DisposableTips:              "Disposable",
	FixedTips:                   "Fixed",
	MixedDisposableAndFixedTips: "Mixed",
	NoTips:                      "None",
}

func (tt TipType) String() string {
	if r, ok := tipNames[tt]; ok {
		return r
	}
	panic("unknown tip type")
}

func (tt TipType) IsValid() bool {
	_, ok := tipNames[tt]
	return ok
}
