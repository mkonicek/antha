// pcr_parser
package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func pcrReactionfromcsv(designFile string, sequenceFile string) (pcrReaction []pcr.Reaction, err error) {

	designedconstructs := readPCRDesign(designFile)
	sequences := readPCRDesign(sequenceFile)
	for _, c := range designedconstructs {
		var newpcrReaction pcr.Reaction
		newpcrReaction.ReactionName = c[0]
		newpcrReaction.Template.Nm = c[1]
		newpcrReaction.PrimerPair[0].Nm = c[2]
		newpcrReaction.PrimerPair[1].Nm = c[3]
		pcrReaction = append(pcrReaction, newpcrReaction)
	}

	var status string

	for b, _ := range pcrReaction {
		var x int
		for _, c := range sequences {
			var y int
			if b == 0 {
				for _, d := range sequences { //check for duplicate entries in list
					if d[0] == c[0] {
						y++
						if y > 1 {
							status = (status + "Part " + c[0] + " defined more than once in Sheet1 please specify a single entry. ")
						}
					}
				}
			}

			if c[0] == pcrReaction[b].Template.Nm {
				pcrReaction[b].Template.Seq = c[1]
				x++
			} else if c[0] == pcrReaction[b].PrimerPair[0].Nm {
				pcrReaction[b].PrimerPair[0].Seq = c[1]
				x++
			} else if c[0] == pcrReaction[b].PrimerPair[1].Nm {
				pcrReaction[b].PrimerPair[1].Seq = c[1]
				x++
			}

		}
		if x < 3 {
			status = (status + "Please specify all parts for reaction " + pcrReaction[b].ReactionName + " in Sheet1. ")
		}
	}
	if status != "" {
		err = fmt.Errorf(status)
	}
	return
}

func readPCRDesign(filename string) [][]string {

	var constructs [][]string

	csvfile, err := os.Open(filename)

	if err != nil {
		panic(err)
		return constructs
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {
		var parts []string

		if len(each[0]) > 0 {
			if string(strings.TrimSpace(each[0])[0]) != "#" {
				for _, p := range each {
					if p != "" {
						parts = append(parts, strings.TrimSpace(p))
					}
				}
				constructs = append(constructs, parts)
			}
		}
	}

	return constructs
}
