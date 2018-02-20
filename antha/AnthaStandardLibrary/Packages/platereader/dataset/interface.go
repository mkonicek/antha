// interface
package dataset

import (
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type ReaderOption interface {
	AddCondition(condition interface{}) error
}

// minimal interface to support existing antha elements which use plate reader data (AddPlateReder_Results)
type PlateReaderData interface {
	ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}) (average float64, err error)
}

///////
type AbsorbanceData interface {
	Absorbance(wellname string, wavelength int, options ...ReaderOption) (average wtype.Absorbance, err error)
}

type AbsorbanceTimeCourseData interface {
	AbsorbanceData
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

type FluorescenceData interface {
}

// minimal interface to support existing fluoresence based antha elements which use plate reader data (AddGFPODPlateReaderResults)
type FluorescenceTimeCourseData interface {
	FluorescenceData
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

//////////////////////////
