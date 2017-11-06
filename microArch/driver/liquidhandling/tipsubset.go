package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type TipSubset struct {
	Mask    []bool
	TipType string
	Channel *wtype.LHChannelParameter
}
