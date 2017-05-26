// pcr_parser
package parser

import (
	"fmt"
	"path/filepath"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/pcr"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//This function takes in a pcr design file (.xlsx or .csv) and converts this into an array of type PCRReactions.
func ParsePCRExcel(designfile wtype.File) ([]pcr.Reaction, error) {

	data, err := designfile.ReadAll()
	var pcrreaction []pcr.Reaction

	if err != nil {
		err = fmt.Errorf(err.Error())
	} else {

		switch {
		case filepath.Ext(designfile.Name) == ".xlsx":
			csvfile1, err := xlsxparserBinary(data, 0, "Sheet2")
			csvfile2, err := xlsxparserBinary(data, 1, "Sheet2")

			pcrreaction, err = pcrReactionfromcsv(csvfile2.Name(), csvfile1.Name())
			return pcrreaction, err

		//case filepath.Ext(designfile.Name) == ".csv":
		//	pcrreaction = pcrReactionfromcsv(designfile.Name)
		//	return pcrreaction, err

		case filepath.Ext(designfile.Name) != ".xlsx": //".csv" && filepath.Ext(designfile.Name) != ".xlsx":
			err = fmt.Errorf("File format not supported please use .xlsx file. ")
			return nil, err
		}
		return pcrreaction, err
	}
	return pcrreaction, err
}
