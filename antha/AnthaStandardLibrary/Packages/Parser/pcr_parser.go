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

// ParsePCRExc3el takes in a pcr design file (.xlsx or .csv) and converts this
// into an array of type PCRReactions.
func ParsePCRExcel(designfile wtype.File) ([]pcr.Reaction, error) {

	data, err := designfile.ReadAll()
	var pcrreaction []pcr.Reaction

	if err != nil {
		err = fmt.Errorf(err.Error())
	} else {

		switch {
		case filepath.Ext(designfile.Name) == ".xlsx":
			csvfile1, err := xlsxparserBinary(data, 0, "Sheet2")
			if err != nil {
				return nil, err
			}
			csvfile2, err := xlsxparserBinary(data, 1, "Sheet2")
			if err != nil {
				return nil, err
			}

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

func pcrReactionfromcsv(designFile string, sequenceFile string) (pcrReactions []pcr.Reaction, err error) {
	var errs []string
	pcrDesigns := readPCRDesign(designFile)
	partSequences := readPCRDesign(sequenceFile)
	for _, pcrDesignFields := range pcrDesigns {
		var newpcrReaction pcr.Reaction
		if len(pcrDesignFields) > 0 {
			newpcrReaction.ReactionName = pcrDesignFields[0]
			if len(pcrDesignFields) > 1 {
				newpcrReaction.Template.Nm = pcrDesignFields[1]
			} else {
				errs = append(errs, fmt.Sprintf("no template specified for %s", newpcrReaction.ReactionName))
			}
			if len(pcrDesignFields) > 2 {
				newpcrReaction.PrimerPair[0].Nm = pcrDesignFields[2]
			} else {
				errs = append(errs, fmt.Sprintf("no primer 1 specified for %s", newpcrReaction.ReactionName))
			}
			if len(pcrDesignFields) > 3 {
				newpcrReaction.PrimerPair[1].Nm = pcrDesignFields[3]
			} else {
				errs = append(errs, fmt.Sprintf("no primer 2 specified for %s", newpcrReaction.ReactionName))
			}
			pcrReactions = append(pcrReactions, newpcrReaction)
		}
	}

	for i, pcrReaction := range pcrReactions {
		var reactionPartCounter int
		for _, sequenceField := range partSequences {
			var y int
			if i == 0 {
				for _, seqField := range partSequences { //check for duplicate entries in list
					if seqField[0] == sequenceField[0] {
						y++
						if y > 1 {
							errs = append(errs, "Part "+sequenceField[0]+" defined more than once in Sheet1 please specify a single entry. ")
						}
					}
				}
			}

			if sequenceField[0] == pcrReaction.Template.Nm {
				pcrReaction.Template.Seq = sequenceField[1]
				reactionPartCounter++
			} else if sequenceField[0] == pcrReaction.PrimerPair[0].Nm {
				if len(sequenceField) > 1 {
					pcrReaction.PrimerPair[0].Seq = sequenceField[1]
				}
				reactionPartCounter++
			} else if sequenceField[0] == pcrReaction.PrimerPair[1].Nm {
				if len(sequenceField) > 1 {
					pcrReaction.PrimerPair[0].Seq = sequenceField[1]
				}
				reactionPartCounter++
			}

		}
		if reactionPartCounter < 3 {
			errs = append(errs, "Please specify all parts for reaction "+pcrReactions[i].ReactionName+" in Sheet1. ")
		}
	}
	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "\n"))
	}
	return
}

func readPCRDesign(filename string) [][]string {

	var constructs [][]string

	csvfile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer csvfile.Close() //nolint

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// sanity check, display to standard output
	for _, record := range rawCSVdata {
		var parts []string

		if len(record[0]) > 0 {
			if string(strings.TrimSpace(record[0])[0]) != "#" {
				for _, partName := range record {
					trimmed := strings.TrimSpace(partName)
					if trimmed != "" {
						parts = append(parts, trimmed)
					}
				}
				constructs = append(constructs, parts)
			}
		}
	}

	return constructs
}
