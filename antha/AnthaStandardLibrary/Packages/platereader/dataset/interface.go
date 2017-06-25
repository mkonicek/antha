// interface
package dataset

import (
	"time"
)

// minimal interface to support existing antha elements which use plate reader data (AddPlateReder_Results and AddGFPODPlateReaderResults)
type PlateReaderData interface {
	BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error)
	ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}, readingtypekeyword string) (average float64, err error)
	FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error)
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

///////
type AbsorbanceData interface {
	BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error)
	AbsorbanceReading(wellname string, wavelength int) (average float64, err error)
	FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error)
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

type TimeCourseData interface {
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

//////////////////////////
