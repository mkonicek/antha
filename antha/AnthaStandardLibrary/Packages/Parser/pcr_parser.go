// pcr_parser
package parser

import (
	"fmt"
	"path/filepath"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/PCR"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//This function takes in a pcr design file (.xlsx or .csv) and converts this into an array of type PCRReactions.
func ParsePCRExcel(designfile wtype.File) ([]PCR.PCRReaction, error) {

	data, err := designfile.ReadAll()
	var pcrreaction []PCR.PCRReaction
	if err != nil {
		err = fmt.Errorf(err.Error())
	} else {

		switch {
		case filepath.Ext(designfile.Name) == ".xlsx":

			csvfile, err := xlsxparserBinary(data, 0, "Sheet1") //this returns a file os.file

			pcrreaction = pcrReactionfromcsv(csvfile.Name())
			return pcrreaction, err

		case filepath.Ext(designfile.Name) == ".csv":
			pcrreaction = pcrReactionfromcsv(designfile.Name)
			return pcrreaction, err

		case filepath.Ext(designfile.Name) != ".csv" && filepath.Ext(designfile.Name) != ".xlsx":
			err = fmt.Errorf("File format not supported please use .csv or .xlsx file. ")
			return nil, err
		}
		return pcrreaction, err
	}
	return pcrreaction, err
}
