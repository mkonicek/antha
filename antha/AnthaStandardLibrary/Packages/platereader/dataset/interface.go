// interface
package dataset

import (
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// minimal interface to support existing antha elements which use plate reader data (AddPlateReder_Results)
type PlateReaderData interface {
	ReadingsAsAverage(wellname string, emexortime platereader.FilterOption, fieldvalue interface{}) (average float64, err error)
}

///////
type AbsorbanceData interface {
	Absorbance(wellname string, wavelength int, options ...interface{}) (average wtype.Absorbance, err error)
	AllAbsorbanceData() (map[string][]wtype.Absorbance, error)
}

type FluorescenceData interface {
	Fluorescence(wellname string, excitationWavelength, emissionWavelength int, options ...interface{}) (average float64, err error)
}

type TimeCourseData interface {
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

type AbsorbanceTimeCourseData interface {
	AbsorbanceData
	TimeCourseData
}

// minimal interface to support existing fluoresence based antha elements which use plate reader data (AddGFPODPlateReaderResults)
type FluorescenceTimeCourseData interface {
	FluorescenceData
	TimeCourseData
}

//////////////////////////
