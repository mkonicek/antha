package liquidhandling

// consts for generic liquid handling types
const (
	LLLiquidHandler string = "LLLiquidHandler" // requires detailed programming e.g. move, aspirate, move dispense etc.
	HLLiquidHandler string = "HLLiquidHandler" // can orchestrate liquid transfers itself
)

// consts for tip requirements of liquid handlers
const (
	DisposableTips              string = "Disposable" // disposable system	-- needs tip boxes & waste
	FixedTips                   string = "Fixed"      // fixed tip system	-- needs tip wash
	MixedDisposableAndFixedTips string = "Mixed"      // both disposable and mixed	-- needs all of the above
	NoTips                      string = "None"       // does not use tips
)
