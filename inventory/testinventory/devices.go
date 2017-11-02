package testinventory

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
)

// heights in mm
const (
	offset                                      = 0.25
	gilsonoffsetpcrplate                        = 2.0 // 2.136
	gilsonoffsetgreiner                         = 2.0
	riserheightinmm                             = 40.0 - offset
	shallowriserheightinmm                      = 20.25 - offset
	shallowriser18heightinmm                    = 18.75 - offset
	coolerheight                                = 16.0
	isofreezecoolerheight                       = 10.0
	pcrtuberack496                              = 28.0
	valueformaxheadtonotintoDSWplatewithp20tips = 4.5
	bioshake96welladaptorheight                 = 4.5
	bioshakestandardadaptorheight               = 5.0
)

const (
	incubatoroffset = -1.58
)

var (
	incubatorheightinmm = devices.Shaker["3000 T-elm"]["Height"]*1000 + incubatoroffset
	inhecoincubatorinmm = devices.Shaker["InhecoStaticOnDeck"]["Height"] * 1000
)

// defaultDevices are default devices upon which an sbs format plate may be placed
var defaultDevices = map[string]device{
	"riser40": riser{
		Name:         "riser40",
		Manufacturer: "Cybio",
		Heightinmm:   riserheightinmm,
		Synonyms:     []string{"riser40", "riser"},
	},

	"riser20": riser{
		Name:         "riser20",
		Manufacturer: "Gilson",
		Heightinmm:   shallowriserheightinmm,
		Synonyms:     []string{"riser20", "shallowriser"},
	},

	"riser18": riser{
		Name:         "riser18",
		Manufacturer: "Gilson",
		Heightinmm:   shallowriser18heightinmm,
		Synonyms:     []string{"riser18", "shallowriser18"},
	},

	"bioshake": incubator{
		Riser: riser{
			Name:         "bioshake",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm,
			Synonyms:     []string{"bioshake"},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": []string{"position_1"},
		},
	},

	"bioshake_96well_adaptor": incubator{
		Riser: riser{
			Name:         "bioshake_96well_adaptor",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm + bioshake96welladaptorheight,
			Synonyms:     []string{"bioshake_96well_adaptor"},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": []string{"position_1"},
		},
	},

	"bioshake_standard_adaptor": incubator{
		Riser: riser{Name: "bioshake_standard_adaptor",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm + bioshakestandardadaptorheight,
			Synonyms:     []string{"bioshake_standard_adaptor"},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": []string{"position_1"},
		},
	},
}

type device interface {
	GetConstraints() constraints
	GetSynonyms() []string
	GetHeightInmm() float64
	GetRiser() riser
	GetName() string
}

// Constraints map device type to allowed positions for a device
type constraints map[string][]string

// A riser is an SBS format object upon which a plate can be placed.
type riser struct {
	Name         string
	Manufacturer string
	Heightinmm   float64
	Synonyms     []string
}

func (r riser) GetRiser() riser {
	return r
}

func (r riser) GetConstraints() constraints {
	return nil
}

func (r riser) GetSynonyms() []string {
	return r.Synonyms
}

func (r riser) GetHeightInmm() float64 {
	return r.Heightinmm
}

func (r riser) GetName() string {
	return r.Name
}

// An incubator is an SBS format device upon which a plate can be placed with
// constraints
type incubator struct {
	Riser               riser
	Properties          map[string]float64
	PositionConstraints constraints // map device to positions where the device is restricted; if empty no restrictions are expected
}

func (i incubator) GetRiser() riser {
	return i.Riser
}

func (i incubator) GetConstraints() constraints {
	return i.PositionConstraints
}

func (i incubator) GetSynonyms() []string {
	return i.Riser.Synonyms
}

func (i incubator) GetHeightInmm() float64 {
	return i.Riser.Heightinmm
}

func (i incubator) GetName() string {
	return i.Riser.Name
}
