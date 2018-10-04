package testinventory

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"math"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/inventory"
)

type platetest struct {
	TestPlateName  string
	ExpectedHeight float64
	ExpectedZStart float64
}

var tests = []platetest{
	{TestPlateName: "reservoir", ExpectedZStart: 0.0, ExpectedHeight: 40.0},
	{TestPlateName: "pcrplate_skirted", ExpectedZStart: MinimumZHeightPermissableForLVPipetMax, ExpectedHeight: 15.5},
	{TestPlateName: "pcrplate", ExpectedZStart: MinimumZHeightPermissableForLVPipetMax, ExpectedHeight: 15.5},
	{TestPlateName: "greiner384", ExpectedZStart: 2.5, ExpectedHeight: 14.0},
	{TestPlateName: "Nuncon12well", ExpectedZStart: 4.0, ExpectedHeight: 19.0},
	{TestPlateName: "Nuncon12wellAgar", ExpectedZStart: 9.0, ExpectedHeight: 19.0},
	{TestPlateName: "strip_tubes_0.2ml", ExpectedZStart: 0.0, ExpectedHeight: 15.5},
}

var testsofPlateWithRiser = []platetest{
	{TestPlateName: "pcrplate_with_cooler", ExpectedZStart: coolerheight + MinimumZHeightPermissableForLVPipetMax, ExpectedHeight: 15.5 + coolerheight},
	{TestPlateName: "pcrplate_with_isofreeze_cooler", ExpectedZStart: isofreezecoolerheight, ExpectedHeight: 15.5 + isofreezecoolerheight - MinimumZHeightPermissableForLVPipetMax},
	{TestPlateName: "pcrplate_skirted_with_isofreeze_cooler", ExpectedZStart: isofreezecoolerheight + 2.0, ExpectedHeight: 15.5 - gilsonoffsetpcrplate + isofreezecoolerheight + 3.4 - 0.036},
	{TestPlateName: "pcrplate_with_496rack", ExpectedZStart: pcrtuberack496HeightInmm, ExpectedHeight: 15.5 + pcrtuberack496HeightInmm - MinimumZHeightPermissableForLVPipetMax},
	{TestPlateName: "pcrplate_semi_skirted_with_496rack", ExpectedZStart: pcrtuberack496HeightInmm + 1.0, ExpectedHeight: 15.5 + pcrtuberack496HeightInmm},
	{TestPlateName: "strip_tubes_0.2ml_with_496rack", ExpectedZStart: pcrtuberack496HeightInmm - 2.5, ExpectedHeight: 15.5 + pcrtuberack496HeightInmm - 2.5},
	{TestPlateName: "FluidX700ulTubes_with_FluidX_high_profile_rack", ExpectedZStart: 2, ExpectedHeight: 26.736 + fluidXhighProfileRackHeight},
}

func TestAddRiser(t *testing.T) {
	t.Skip()
	ctx := NewContext(context.Background())

	for _, test := range tests {
		for _, device := range defaultDevices {
			testPlate, err := inventory.NewPlate(ctx, test.TestPlateName)
			if err != nil {
				t.Error(err)
				continue
			}

			testname := test.TestPlateName + "_" + device.GetName()

			newPlates := addRiser(testPlate, device)
			if e, f := 0, len(newPlates); e == f {
				if !doNotAddThisRiserToThisPlate(testPlate, device) {
					t.Errorf("expected some plates resulting from adding riser %s to plate %s but none found", device.GetName(), testPlate.Type)
				}
				continue
			}

			newPlate := newPlates[0]
			if e, f := testname, newPlate.Type; e != f {
				t.Errorf("expected %s but found %s", e, f)
			}

			offset := platespecificoffset[test.TestPlateName]

			// check that the height is as expected using default inventory
			if testPlate.Height() != test.ExpectedHeight {
				t.Error(
					"for", test.TestPlateName, "\n",
					"Expected plate height:", test.ExpectedHeight, "\n",
					"got:", testPlate.Height(), "\n",
				)
			}

			// check that the height is as expected with riser added
			if f, e := newPlate.Height(), test.ExpectedHeight; e != f {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate height:", e, "\n",
					"got:", f, "\n",
				)
			}

			// now test z offsets
			if testPlate.WellZStart != test.ExpectedZStart {
				t.Error(
					"for", test.TestPlateName, "\n",
					"Expected plate ZStart:", test.ExpectedZStart, "\n",
					"got:", testPlate.WellZStart, "\n",
				)
			}

			if f, e := newPlate.WellZStart, test.ExpectedZStart+device.GetHeightInmm()-offset+plateRiserSpecificOffset(testPlate, device); e != f {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate ZStart:", test.ExpectedZStart, "+",
					"device:", device.GetHeightInmm(), "=", e, "\n",
					"got:", f, "\n",
				)
			}

			if f, e := testPlate.WellZStart, test.ExpectedZStart; e != f {
				t.Error(
					"for", "no device", "\n",
					"testname", test.TestPlateName, "\n",
					"Expected plate height:", e, "\n",
					"got:", f, "\n",
				)
			}
		}
	}
}

type testdevice struct {
	name                string
	constraintdevice    string
	constraintposition1 string
	height              float64
}

var testdevices = []testdevice{
	{name: "bioshake", constraintdevice: "Pipetmax", constraintposition1: "position_1", height: 55.92},
}

type deviceExceptions map[string][]string // key is device name, exceptions are the plates which will give a result which differs from norm

var exceptions deviceExceptions = map[string][]string{
	"bioshake":                  {"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
	"bioshake_96well_adaptor":   {"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
	"bioshake_standard_adaptor": {"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
}

func TestDeviceMethods(t *testing.T) {
	t.Skip()
	for _, device := range testdevices {

		_, ok := defaultDevices[device.name]

		if !ok {
			t.Error(
				"for", device.name, "\n",
				"not found in devices", defaultDevices, "\n",
			)
		} else {
			c := defaultDevices[device.name].GetConstraints()
			h := defaultDevices[device.name].GetHeightInmm()
			//r := Devices[device].GetRiser()

			if constraints, ok := c[device.constraintdevice]; !ok || constraints[0] != device.constraintposition1 {
				t.Error(
					"for", device.name, "\n",
					"Constraints", c, "\n",
					"expected key", device.constraintdevice, "\n",
					"expected 1st position", device.constraintposition1, "\n",
				)
			}

			if h != device.height {
				t.Error(
					"for", device.name, "\n",
					"expectd height", device.height, "\n",
					"got", h, "\n",
				)
			}
		}

	}

}

func TestSetConstraints(t *testing.T) {
	t.Skip()
	ctx := NewContext(context.Background())

	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}

	for _, testplate := range GetPlates(ctx) {
		for _, device := range defaultDevices {

			if device.GetConstraints() == nil {
				continue
			}

			if search.InStrings(exceptions[device.GetName()], testplate.Type) {
				continue
			}

			newPlates := addRiser(testplate, device)

			if strings.Contains(testplate.Type, device.GetName()) {
				if e, f := 0, len(newPlates); e != f {
					t.Errorf("expecting %d plates found %d", e, f)
					continue
				}
			} else if !containsRiser(testplate) {
				if e, f := 1, len(newPlates); e != f {
					if !doNotAddThisRiserToThisPlate(testplate, device) {
						t.Errorf("expecting %d plates found %d", e, f)
					}
					continue
				} else if e, f := testplate.Type+"_"+device.GetName(), newPlates[0].Type; e != f {
					t.Errorf("expecting type %s found %s", e, f)
					continue
				}
			} else {
				continue
			}

			for _, testplate := range newPlates {
				//positionsinterface, found := testplate.Welltype.Extra[platform]
				positions, found := testplate.IsConstrainedOn(platform)

				if doNotAddThisRiserToThisPlate(testplate, device) {
					// skip this case
				} else if !found || len(positions) == 0 {
					t.Error(
						"for", device, "\n",
						"testname", testplate.Type, "\n",
						"Constraints found?", found, "\n",
						"Positions: ", positions, "\n",
						"expected positions: ", expectedpositions, "\n",
						"Constraints expected :", device.GetConstraints()[platform], "\n",
						"Constraints got :", testplate.Welltype.Extra[platform], "\n",
					)
				} else if len(positions) != len(expectedpositions) {
					t.Error(
						"for", device, "\n",
						"testname", testplate.Type, "\n",
						"Positions got: ", positions, "\n",
						"Positions expected: ", expectedpositions, "\n",
						"Constraints expected :", device.GetConstraints()[platform], "\n",
						"Constraints got :", testplate.Welltype.Extra[platform], "\n",
					)
				}
			}
		}
	}
}

func TestGetConstraints(t *testing.T) {
	t.Skip()
	ctx := NewContext(context.Background())

	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}
	for _, testplate := range GetPlates(ctx) {
		for _, device := range defaultDevices {

			if device.GetConstraints() == nil {
				continue
			}

			if search.InStrings(exceptions[device.GetName()], testplate.Type) {
				continue
			}
			var testname string
			if strings.Contains(testplate.Type, device.GetName()) {
				testname = testplate.Type
			} else if !containsRiser(testplate) {
				testname = testplate.Type + "_" + device.GetName()
			} else {
				continue
			}

			testplate, err := inventory.NewPlate(ctx, testname)
			if err != nil {
				if !doNotAddThisRiserToThisPlate(testplate, device) {
					t.Error(err)
				}
				continue
			}

			//positionsinterface, found := testplate.Welltype.Extra[platform]
			//positions, ok := positionsinterface.([]string)

			positions, found := testplate.IsConstrainedOn(platform)
			if !found || positions == nil || len(positions) != len(expectedpositions) || positions[0] != expectedpositions[0] {
				if doNotAddThisRiserToThisPlate(testplate, device) && len(device.GetConstraints()[platform]) > 0 {
					t.Error(
						"for", device, "\n",
						"testname", testname, "\n",
						"Constraints found", found, "\n",
						"Positions: ", positions, "\n",
						"expected positions: ", expectedpositions, "\n",
						"Constraints expected :", device.GetConstraints()[platform], "\n",
						"Constraints got :", testplate.Welltype.Extra[platform], "\n",
					)
				}
			}
		}
	}
}

func TestPlateZs(t *testing.T) {
	t.Skip()
	ctx := NewContext(context.Background())

	var allTests []platetest

	allTests = append(allTests, testsofPlateWithRiser...)
	allTests = append(allTests, tests...)

	for _, test := range allTests {

		testplate, err := inventory.NewPlate(ctx, test.TestPlateName)
		if err != nil {
			t.Error(err)
			continue
		}

		if testplate.WellZStart != test.ExpectedZStart {
			t.Error(
				"for", test.TestPlateName, "\n",
				"expected ZStart: ", test.ExpectedZStart, "\n",
				"got ZStart:", testplate.WellZStart, "\n",
			)
		}

		// check that the height is as expected using default inventory
		if math.Abs(testplate.Height()-test.ExpectedHeight) > 0.001 {
			t.Error(
				"for", test.TestPlateName, "\n",
				"Expected plate height:", test.ExpectedHeight, "\n",
				"got:", testplate.Height(), "\n",
			)
		}
	}
}

func addRiser(plate *wtype.Plate, riser device) (plates []*wtype.Plate) {
	if containsRiser(plate) || doNotAddThisRiserToThisPlate(plate, riser) {
		return
	}

	for _, risername := range riser.GetSynonyms() {
		var dontaddrisertothisplate bool

		newplate := plate.Dup()
		riserheight := riser.GetHeightInmm()
		if offset, found := platespecificoffset[plate.Type]; found {
			riserheight = riserheight - offset
		}

		riserheight = riserheight + plateRiserSpecificOffset(plate, riser)

		newplate.WellZStart = plate.WellZStart + riserheight
		newname := plate.Type + "_" + risername
		newplate.Type = newname
		if riser.GetConstraints() != nil {
			// duplicate well before adding constraint to prevent applying
			// constraint to all common &Welltype on other plates

			for device, allowedpositions := range riser.GetConstraints() {
				newwell := newplate.Welltype.Dup()
				newplate.Welltype = newwell
				_, ok := newwell.Extra[device]
				if !ok {
					newplate.SetConstrained(device, allowedpositions)
				} else {
					dontaddrisertothisplate = true
				}
			}
		}

		if !dontaddrisertothisplate {
			plates = append(plates, newplate)
		}
	}

	return
}

// was devices.go

// heights in mm
const (
	offset                   = 0.25
	gilsonoffsetpcrplate     = 2.0 // 2.136
	gilsonoffsetgreiner      = 2.0
	riserheightinmm          = 40.0 - offset
	shallowriserheightinmm   = 20.25 - offset
	shallowriser18heightinmm = 18.75 - offset
	coolerheight             = 16.0
	isofreezecoolerheight    = 10.0
	pcrtuberack496HeightInmm = 28.0
	//valueformaxheadtonotintoDSWplatewithp20tips = 4.5
	bioshake96welladaptorheight        = 4.5
	bioshakestandardadaptorheight      = 5.0
	appliedbiosystemsmagbeadbaseheight = 12.0 //height of just plate base, upon which most skirted plates can rest
	fluidXhighProfileRackHeight        = 2.0 - MinimumZHeightPermissableForLVPipetMax
)

const (
	incubatoroffset = -1.58
)

var (
	incubatorheightinmm = devices.Shaker["3000 T-elm"]["Height"]*1000 + incubatoroffset
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

func doNotAddThisRiserToThisPlate(plate *wtype.Plate, riser device) bool {

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

func plateRiserSpecificOffset(plate *wtype.Plate, riser device) float64 {

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

// was plateAccessories.go

// The height below which an error will be generated
// when attempting to perform transfers with low volume head and tips (0.5 - 20ul) on the Gilson PipetMax.
const MinimumZHeightPermissableForLVPipetMax = 0.636

//var commonwelltypes

var platespecificoffset = map[string]float64{
	"pcrplate_skirted": gilsonoffsetpcrplate,
	"greiner384":       gilsonoffsetgreiner,
	"costar48well":     3.0,
	"Nuncon12well":     11.0, // this must be wrong!! check z start without riser properly
	"Nuncon12wellAgar": 11.0, // this must be wrong!! check z start without riser properly
	"VWR12well":        3.0,
}

// function to check if a platename already contains a riser
func containsRiser(plate *wtype.Plate) bool {
	for _, dev := range defaultDevices {
		for _, synonym := range dev.GetSynonyms() {
			if strings.Contains(plate.Type, "_"+synonym) {
				return true
			}
		}
	}

	return false
}
