// antha/AnthaStandardLibrary/Packages/Parser/RebaseParser.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// package genbank converts DNA sequence files in genbank format into a set of DNA sequences.
package genbank

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//Parses file of type .gb to DNASequence. Features are not annotated
func GenbankToFeaturelessDNASequence(sequenceFile wtype.File) (wtype.DNASequence, error) {
	data, err := sequenceFile.ReadAll()
	if err != nil {
		return wtype.DNASequence{}, err
	}
	var line string
	genbanklines := make([]string, 0)
	buffer := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	return handleGenbank(genbanklines)
}

//Parses a feature from a genbank file into a DNASequence.
func GenbankFeatureToDNASequence(file wtype.File, featurename string) (wtype.DNASequence, error) {
	var line string
	genbanklines := make([]string, 0)

	data, err := file.ReadAll()

	if err != nil {
		return wtype.DNASequence{}, err
	}

	buffer := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buffer)

	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	annotated, err := handleGenbank(genbanklines)
	if err != nil {
		return wtype.DNASequence{}, err
	}

	var standardseq wtype.DNASequence
	for _, feature := range annotated.Features {
		if strings.Contains(feature.Name, featurename) {
			standardseq.Nm = feature.Name
			standardseq.Seq = feature.DNASeq
			return standardseq, nil
		}
	}
	errstr := fmt.Sprint("Feature: ", featurename, "not found. ", "found these features: ", annotated.FeatureNames())
	return standardseq, fmt.Errorf(errstr)
}

// parses a genbank file into a DNASEquence making features from annotations
func GenbankToAnnotatedSeq(file wtype.File) (annotated wtype.DNASequence, err error) {
	data, err := file.ReadAll()
	if err != nil {
		return
	}
	annotated, err = GenbankContentsToAnnotatedSeq(data)
	return
}

// parses contents of a genbank file into a DNASEquence making features from annotations
func GenbankContentsToAnnotatedSeq(contentsinbytes []byte) (annotated wtype.DNASequence, err error) {
	var line string
	genbanklines := make([]string, 0)

	data := bytes.NewBuffer(contentsinbytes)

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	annotated, err = handleGenbank(genbanklines)

	return
}

func handleGenbank(lines []string) (annotatedseq wtype.DNASequence, err error) {
	if lines[0][0:5] == `LOCUS` {
		name, _, _, circular, _, err := locusLine(lines[0])

		if err != nil {
			return annotatedseq, err
		}

		seq := handleSequence(lines)

		features, err := handleFeatures(lines, seq, "DNA")
		if err != nil {
			return annotatedseq, err
		}
		annotatedseq, err = sequences.MakeAnnotatedSeq(name, seq, circular, features)
		if err != nil {
			return annotatedseq, err
		}
	} else {
		err = fmt.Errorf("no LOCUS found on first line")
	}
	return
}
func locusLine(line string) (name string, seqlength int, seqtype string, circular bool, date string, err error) {

	fields := strings.SplitN(line, " ", 2)
	restofline := fields[1]
	fields = strings.Split(restofline, " ")

	newarray := make([]string, 0)
	for _, s := range fields {
		if s != "" && s != " " {
			newarray = append(newarray, s)
		}
	}
	fields = newarray
	if len(fields) > 1 {
		if len(fields) < 5 {
			err = errors.New("the locusline does not contain enough elements or is not formatted correctly. Please check file.")
			return
		}
		name = fields[0]
		i, newerr := strconv.Atoi(fields[1])
		if newerr != nil {
			err = newerr
		}
		seqlength = i
		seqtype = fields[3]

		if fields[4] == "circular" {
			circular = true
		} else {
			circular = false
		}

		if len(fields) > 5 {
			date = fields[5]
		} else {
			date = "No date supplied"
		}
		return
	} else {
		err = fmt.Errorf("invalid genbank locus line: \"%s\"", line)
	}

	return
}
func cleanup(line string) (cleanarray []string) {
	fields := strings.Split(line, " ")

	for _, s := range fields {

		if s != "" && s != " " {
			cleanarray = append(cleanarray, s)
		}

	}

	return
}

func featureline1(line string) (reverse bool, class string, startposition int, endposition int, err error) {

	newarray := cleanup(line)

	class = newarray[0]

	for _, s := range newarray {
		if s[0] == '<' {
			s = s[1:]
		}
		if s[0] == '>' {
			s = s[1:]
		}
		var warning error
		if strings.Contains(s, `join`) {
			warning = fmt.Errorf("feature \"%s\" contains join location, adding as one feature only for now", s)
			s = strings.Replace(s, "Join(", "", -1)
			s = strings.Replace(s, ")", "", -1)
			joinhandler := strings.Split(s, `,`)
			split := strings.Split(joinhandler[0], "..")
			startposition, err = strconv.Atoi(split[0])
			if err != nil {
				return
			}
			split = strings.Split(joinhandler[1], "..")
			endposition, err = strconv.Atoi(strings.TrimRight(split[1], "\n"))
			if err != nil {
				return
			}
		} else {
			if strings.Contains(s, "complement") {
				reverse = true
				// though the following line is technically incorrect I'm afraid
				// to fix it because I don't want to cause any accidental bugs,
				// so nolint it (contains duplicate chars in cutset)
				s = strings.TrimLeft(s, `(complement)`) // nolint
				s = strings.TrimRight(s, ")")
				if s[0] == '<' {
					s = s[1:]
				}
				if s[0] == '>' {
					s = s[1:]
				}
			}
			index := strings.Index(s, "..")
			if index != -1 {
				startposition, err = strconv.Atoi(s[0:index])
				if err != nil {
					return
				}
				ss := strings.SplitAfter(s, "..")
				if strings.Contains(ss[1], ")") {
					ss[1] = strings.Replace(ss[1], ")", "", -1)
				}
				if strings.Contains(ss[1], "bp") {
					ss[1] = strings.Replace(ss[1], "bp", "", -1)
				}
				if ss[1][0] == '>' {
					ss[1] = ss[1][1:]
				} else if ss[1][0] == '<' {
					ss[1] = ss[1][1:]
				}
				endposition, err = strconv.Atoi(strings.TrimRight(ss[1], "\n"))
				if err != nil {
					return
				}

			}

		}
		if err == nil {
			err = warning
		}
	}

	return
}
func featureline2(line string) (description string, found bool) {

	fields := strings.Split(line, " ")

	for i, field := range fields {
		if strings.Contains(field, `"`) {
			tempfields := make([]string, i)
			tempfield := strings.Join(fields[i:], " ")
			tempfields = append(tempfields, tempfield)
			fields = tempfields
			break
		}
	}

	newarray := make([]string, 0)
	for _, s := range fields {
		if s != "" && s != " " {
			newarray = append(newarray, s)
		}
	}
	for _, line := range newarray {

		if strings.Contains(line, `/gene`) {

			parts := strings.SplitAfterN(line, `="`, 2)
			if len(parts) == 2 {
				parts[1] = strings.Replace(parts[1], `"`, "", -1)
				description = strings.TrimSpace(parts[1])
				found = true
				return
			}

		}

		if strings.Contains(line, `/label`) {
			if strings.Contains(line, `="`) {
				parts := strings.SplitAfterN(line, `="`, 2)
				if len(parts) == 2 {
					parts[1] = strings.Replace(parts[1], `"`, "", -1)
					description = strings.TrimSpace(parts[1])
					found = true
					return
				}
			} else {
				parts := strings.SplitAfterN(line, "=", 2)
				if len(parts) == 2 {
					description = strings.TrimSpace(parts[1])
					found = true
					return
				}

			}

		}

		if strings.Contains(line, `/product`) {

			parts := strings.SplitAfterN(line, `="`, 2)
			if len(parts) == 2 {
				parts[1] = strings.Replace(parts[1], `"`, "", -1)
				description = strings.TrimSpace(parts[1])
				found = true
				return
			}

		}
	}

	return
}

func handleFeature(lines []string) (description string, reverse bool, class string, startposition int, endposition int, err error) {

	if len(lines) > 0 {
		reverse, class, startposition, endposition, err := featureline1(lines[0])

		if err != nil {
			err = fmt.Errorf("Error with Featureline1 func %s", lines[0])
			return description, reverse, class, startposition, endposition, err
		}

		for i := 1; i < len(lines); i++ {
			description, found := featureline2(lines[i])
			if found {
				return description, reverse, class, startposition, endposition, err
			}

		}

	}
	return
}
func detectFeature(lines []string) (detected bool, startlineindex int, endlineindex int) {

	for i := 0; i < len(lines); i++ {

		if string(lines[i][0]) != " " {
			return
		}

		if string(lines[i][7]) != " " {
			startlineindex = i
		}

		_, found := featureline2(lines[i])
		if found {
			endlineindex = i + 1
		}

		if startlineindex != -1 && endlineindex != 0 {
			detected = true
			return
		}
	}

	return
}
func handleFeatures(lines []string, seq string, seqtype string) (features []wtype.Feature, err error) {
	featurespresent := false
	for _, line := range lines {
		if strings.Contains(line, "FEATURES") {
			featurespresent = true
		}
	}
	if !featurespresent {
		return
	}
	features = make([]wtype.Feature, 0)
	var feature wtype.Feature

	for i := 0; i < len(lines); i++ {
		if lines[i][0:8] == "FEATURES" {
			lines = lines[i+1:]
			break
		}
	}

	linesatstart := lines
	for i := 0; i < len(linesatstart); i++ {

		detected, start, end := detectFeature(lines)

		if detected {
			description, reverse, class, startposition, endposition, err := handleFeature(lines[start:end])
			if err != nil {
				return features, err
			}
			rev := ""
			if reverse {
				rev = "Reverse"
			}

			// Warning! this needs to change to handle cases where start and position assignment has failed rather than just ignoring the problem
			if startposition != 0 && endposition != 0 {
				feature = sequences.MakeFeature(description, seq[startposition-1:endposition], startposition, endposition, seqtype, class, rev)
			}

			features = append(features, feature)
			lines = lines[end:]
			if start > end {
				return features, fmt.Errorf("Start position cannot be greater than end position in feature")
			}
		}

	}
	//features = search.DuplicateFeatures(features)

	return

}

var (
	illegal string = "1234567890"
)

func handleSequence(lines []string) (dnaseq string) {
	originallines := len(lines)
	originfound := false

	if len(lines) > 0 {
		for i := 0; i < originallines; i++ {
			if len([]byte(lines[0])) > 0 {
				if !originfound {
					if lines[i][0:6] == "ORIGIN" {
						originfound = true
					}
				}
				if originfound {

					lines = lines[i+1 : originallines]
					seq := strings.Join(lines, "")
					seq = strings.Replace(seq, " ", "", -1)

					for _, character := range illegal {
						seq = strings.Replace(seq, string(character), "", -1)
					}
					seq = strings.Replace(seq, "\n", "", -1)
					seq = strings.Replace(seq, "//", "", -1)
					dnaseq = seq
					return
				}
			}
		}

	}
	return
}
