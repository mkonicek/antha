package main

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
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
	pcrtuberack496HeightInmm                    = 28.0
	valueformaxheadtonotintoDSWplatewithp20tips = 4.5
	bioshake96welladaptorheight                 = 4.5
	bioshakestandardadaptorheight               = 5.0
	appliedbiosystemsmagbeadbaseheight          = 12.0 //height of just plate base, upon which most skirted plates can rest
	//appliedbiosystemsmagbeadtotalheight         = 17.0 //height of base and well, in which other plates can rest
	fluidXhighProfileRackHeight = 2.0 - MinimumZHeightPermissableForLVPipetMax
)

const (
	incubatoroffset = -1.58
)

var (
	incubatorheightinmm = devices.Shaker["3000 T-elm"]["Height"]*1000 + incubatoroffset
	//inhecoincubatorinmm = devices.Shaker["InhecoStaticOnDeck"]["Height"] * 1000
)

// defaultDevices are default devices upon which an sbs format plate may be placed
var defaultDevices = map[string]device{
	"riser40": riser{
		Name:         "riser40",
		Manufacturer: "Cybio",
		Heightinmm:   riserheightinmm,
		Synonyms:     []string{"riser40", "riser"},
		PlateConstraints: plateConstraints{
			NotThesePlates: []plateWithConstraint{
				{
					Name:          "FluidX700ulTubes",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate_semi_skirted",
					SpecialOffset: 0.0,
				},
				{
					Name:          "strip_tubes_0.2ml",
					SpecialOffset: 0.0,
				},
			},
		},
	},

	"riser20": riser{
		Name:         "riser20",
		Manufacturer: "Gilson",
		Heightinmm:   shallowriserheightinmm,
		Synonyms:     []string{"riser20", "shallowriser"},
		PlateConstraints: plateConstraints{
			NotThesePlates: []plateWithConstraint{
				{
					Name:          "FluidX700ulTubes",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate_semi_skirted",
					SpecialOffset: 0.0,
				},
				{
					Name:          "strip_tubes_0.2ml",
					SpecialOffset: 0.0,
				},
			},
		},
	},

	"riser18": riser{
		Name:         "riser18",
		Manufacturer: "Gilson",
		Heightinmm:   shallowriser18heightinmm,
		Synonyms:     []string{"riser18", "shallowriser18"},
		PlateConstraints: plateConstraints{
			NotThesePlates: []plateWithConstraint{
				{
					Name:          "FluidX700ulTubes",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate",
					SpecialOffset: 0.0,
				},
				{
					Name:          "pcrplate_semi_skirted",
					SpecialOffset: 0.0,
				},
				{
					Name:          "strip_tubes_0.2ml",
					SpecialOffset: 0.0,
				},
			},
		},
	},

	"with_496rack": riser{
		Name:         "with_496rack",
		Manufacturer: "Gilson",
		Heightinmm:   pcrtuberack496HeightInmm,
		Synonyms:     []string{"with_496rack"},
		PlateConstraints: plateConstraints{
			OnlyThesePlates: []plateWithConstraint{
				{
					Name:          "pcrplate",
					SpecialOffset: -MinimumZHeightPermissableForLVPipetMax,
				},
				{
					Name:          "pcrplate_semi_skirted",
					SpecialOffset: 0.0,
				},
				{
					Name:          "strip_tubes_0.2ml",
					SpecialOffset: -2.5,
				},
			},
		},
	},

	"with_AB_magnetic_ring_stand": riser{
		Name:         "with_AB_magnetic_ring_stand",
		Manufacturer: "Applied Biosystems",
		Heightinmm:   appliedbiosystemsmagbeadbaseheight,
		Synonyms:     []string{"with_AB_magnetic_ring_stand"},
		PlateConstraints: plateConstraints{
			OnlyThesePlates: []plateWithConstraint{
				{
					Name:          "TwistDNAPlate",
					SpecialOffset: 0.75,
				},
				{
					Name:          "GreinerSWVBottom",
					SpecialOffset: 0.25,
				},
				{
					Name:          "Nunc_96_deepwell_1ml",
					SpecialOffset: 3.30,
				},
			},
		},
	},

	"with_FluidX_high_profile_rack": riser{
		Name:         "with_FluidX_high_profile_rack",
		Manufacturer: "FluidX",
		Heightinmm:   fluidXhighProfileRackHeight,
		Synonyms:     []string{"with_FluidX_high_profile_rack"},
		PlateConstraints: plateConstraints{
			OnlyThesePlates: []plateWithConstraint{
				{
					Name:          "FluidX700ulTubes",
					SpecialOffset: 0.0,
				},
			},
		},
	},

	"bioshake": incubator{
		Riser: riser{
			Name:         "bioshake",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm,
			Synonyms:     []string{"bioshake"},
			PlateConstraints: plateConstraints{
				NotThesePlates: []plateWithConstraint{
					{
						Name:          "FluidX700ulTubes",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate_semi_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "strip_tubes_0.2ml",
						SpecialOffset: 0.0,
					},
				},
			},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": {"position_1"},
		},
	},

	"bioshake_96well_adaptor": incubator{
		Riser: riser{
			Name:         "bioshake_96well_adaptor",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm + bioshake96welladaptorheight,
			Synonyms:     []string{"bioshake_96well_adaptor"},
			PlateConstraints: plateConstraints{
				OnlyThesePlates: []plateWithConstraint{
					{
						Name:          "pcrplate",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate_semi_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "strip_tubes_0.2ml",
						SpecialOffset: 0.0,
					},
				},
			},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": {"position_1"},
		},
	},

	"bioshake_standard_adaptor": incubator{
		Riser: riser{Name: "bioshake_standard_adaptor",
			Manufacturer: "QInstruments",
			Heightinmm:   incubatorheightinmm + bioshakestandardadaptorheight,
			Synonyms:     []string{"bioshake_standard_adaptor"},
			PlateConstraints: plateConstraints{
				NotThesePlates: []plateWithConstraint{
					{
						Name:          "FluidX700ulTubes",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate_semi_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "strip_tubes_0.2ml",
						SpecialOffset: 0.0,
					},
				},
			},
		},
		Properties: devices.Shaker["3000 T-elm"],
		PositionConstraints: map[string][]string{
			"Pipetmax": {"position_1"},
		},
	},

	"with_cooler": incubator{
		Riser: riser{
			Name:         "with_cooler",
			Manufacturer: "Eppendorf",
			Heightinmm:   coolerheight,
			Synonyms:     []string{"with_cooler"},
			PlateConstraints: plateConstraints{
				OnlyThesePlates: []plateWithConstraint{
					{
						Name:          "pcrplate",
						SpecialOffset: 0.0,
					},
					{
						Name:          "pcrplate_skirted",
						SpecialOffset: 3.4,
					},
					{
						Name:          "pcrplate_semi_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "strip_tubes_0.2ml",
						SpecialOffset: 0.0,
					},
				},
			},
		},
		Properties: map[string]float64{
			"Height": 0.0,
		},
		PositionConstraints: map[string][]string{},
	},

	"with_isofreeze_cooler": incubator{
		Riser: riser{
			Name:         "with_isofreeze_cooler",
			Manufacturer: "Isofreeze",
			Heightinmm:   isofreezecoolerheight,
			Synonyms:     []string{"with_isofreeze_cooler"},
			PlateConstraints: plateConstraints{
				OnlyThesePlates: []plateWithConstraint{
					{
						Name:          "pcrplate",
						SpecialOffset: -MinimumZHeightPermissableForLVPipetMax,
					},
					{
						Name:          "pcrplate_skirted",
						SpecialOffset: 3.4 - 0.036,
					},
					{
						Name:          "pcrplate_semi_skirted",
						SpecialOffset: 0.0,
					},
					{
						Name:          "strip_tubes_0.2ml",
						SpecialOffset: 0.0,
					},
				},
			},
		},
		Properties: map[string]float64{
			"Height": 0.0,
		},
		PositionConstraints: map[string][]string{},
	},
}

func doNotAddThisRiserToThisPlate(plate *wtype.LHPlate, riser device) bool {

	if plate == nil {
		return true
	}

	platedeviceConstraints := riser.GetPlateConstraints()

	if len(platedeviceConstraints.OnlyThesePlates) > 0 {
		for _, plateWithConstraints := range platedeviceConstraints.OnlyThesePlates {
			if plate.Type == plateWithConstraints.Name {
				return false
			}
		}
		return true
	}

	if len(platedeviceConstraints.NotThesePlates) > 0 {
		for _, plateWithConstraints := range platedeviceConstraints.NotThesePlates {
			if plate.Type == plateWithConstraints.Name {
				return true
			}
		}
	}
	return false
}

func plateRiserSpecificOffset(plate *wtype.LHPlate, riser device) float64 {

	if plate == nil {
		return 0.0
	}

	platedeviceConstraints := riser.GetPlateConstraints()

	if len(platedeviceConstraints.OnlyThesePlates) > 0 {
		for _, plateWithConstraints := range platedeviceConstraints.OnlyThesePlates {
			if plate.Type == plateWithConstraints.Name {
				return plateWithConstraints.SpecialOffset
			}
		}
		return 0.0
	}

	return 0.0
}

type device interface {
	GetConstraints() constraints
	GetSynonyms() []string
	GetHeightInmm() float64
	GetRiser() riser
	GetName() string
	GetPlateConstraints() plateConstraints
}

// Constraints map device type to allowed positions for a device
type constraints map[string][]string

// plateConstraints specifies constraints around which plates are compatible with
// a riser.
type plateConstraints struct {
	// OnlyThesePlates lists a subset of plates for which the riser is only compatible with.
	// If this list is not empty only those plates will be valid options with the riser.
	OnlyThesePlates []plateWithConstraint
	// NotThesePlates lists a subset of plates for which the riser is not compatible with.
	// This will only be evaluated if the OnlyThesePlates field is empty.
	// If the NotThesePlates list is not empty these plates will not be added with the riser.
	NotThesePlates []plateWithConstraint
}

type plateWithConstraint struct {
	// Name of the plate which has a riser constraint
	Name string
	// Any plate specific offset, in mm, which should be added to the riser height.
	// For example, if the riser is a tube rack and the specified plate
	// has very narrow tubes which sit low in the riser then a special offset can be added here to adjust for this.
	// In the example case a negative number would be used to reduce the effectie riser height.
	SpecialOffset float64
}

// A riser is an SBS format object upon which a plate can be placed.
type riser struct {
	Name             string
	Manufacturer     string
	Heightinmm       float64
	Synonyms         []string
	PlateConstraints plateConstraints
}

func (r riser) GetRiser() riser {
	return r
}

func (r riser) GetConstraints() constraints {
	return nil
}

func (r riser) GetPlateConstraints() plateConstraints {
	return r.PlateConstraints
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

func (i incubator) GetPlateConstraints() plateConstraints {
	return i.Riser.PlateConstraints
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
