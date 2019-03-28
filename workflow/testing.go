package workflow

import (
	"encoding/json"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Testing contains data and configuration for testing.
type Testing struct {
	MixTaskChecks []MixTaskCheck
}

// MixStepCheck contains check data for a single mix step.
type MixTaskCheck struct {
	Instructions json.RawMessage
	Outputs      map[string]*wtype.Plate
	TimeEstimate time.Duration
}
