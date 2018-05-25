package wtype

type LHGlobals struct {
	MIN_REASONABLE_VOLUME_UL float64
	VOL_RESOLUTION_DIGITS    int
}

var (
	Globals = LHGlobals{
		MIN_REASONABLE_VOLUME_UL: 0.01,
		VOL_RESOLUTION_DIGITS:    2,
	}

	MAGICBARRIERPROMPTSTRING = "___BARRIER___"
)
