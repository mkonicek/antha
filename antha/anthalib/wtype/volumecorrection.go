package wtype

import "github.com/antha-lang/antha/antha/anthalib/wunit"

type VolumeCorrection struct {
	Type     string
	Volume   wunit.Volume
	Location string // refactor point: location strings to types
}
