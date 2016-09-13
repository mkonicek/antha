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

package parser

import (
	"bytes"
	"fmt"
	"log"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	//"io/ioutil"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func GenbanktoSimpleSeq(filename string) (string, error) {
	var line string
	genbanklines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return HandleSequence(genbanklines), nil
}

func GenbanktoFeaturelessDNASequence(filename string) (wtype.DNASequence, error) {
	line := ""
	genbanklines := make([]string, 0)
	var file *os.File
	file, err := os.Open(filename)
	if err != nil {
		return wtype.DNASequence{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	return HandleGenbank(genbanklines)
}

func GenbankFeaturetoDNASequence(filename string, featurename string) (wtype.DNASequence, error) {
	line := ""
	genbanklines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return wtype.DNASequence{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	annotated, err := HandleGenbank(genbanklines)
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

func GenbankContentstoAnnotatedSeq(contentsinbytes []byte) (annotated wtype.DNASequence, err error) {
	line := ""
	genbanklines := make([]string, 0)

	file := bytes.NewBuffer(contentsinbytes)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	annotated, err = HandleGenbank(genbanklines)

	return

}

func GenbanktoAnnotatedSeq(filename string) (annotated wtype.DNASequence, err error) {

	line := ""
	genbanklines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return wtype.DNASequence{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	return HandleGenbank(genbanklines)
}

func ParseGenbankfile(file *os.File) (wtype.DNASequence, error) {
	line := ""
	genbanklines := make([]string, 0)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		return wtype.DNASequence{}, err
	}

	return HandleGenbank(genbanklines)
}

func HandleGenbank(lines []string) (annotatedseq wtype.DNASequence, err error) {
	if lines[0][0:5] == `LOCUS` {
		// // fmt.Println("in Locus")
		name, _, _, circular, _, err := Locusline(lines[0])

		if err != nil {
			return annotatedseq, err
		}
		/*if seqtype != "DNA" {
			err = fmt.Errorf("Can't parse genbank files which are not classified as type DNA at present")
			// fmt.Println(err.Error())
			return annotatedseq, err
		}*/
		seq := HandleSequence(lines)
		// // // fmt.Println("foundout this seq", seq)

		features := HandleFeatures(lines, seq, "DNA")
		// // // fmt.Println("found these features", features)
		annotatedseq, err = wtype.MakeAnnotatedSeq(name, seq, circular, features)
		// // // fmt.Println("annotated", annotatedseq)
	} else {
		err = fmt.Errorf("no LOCUS found on first line")
	}
	return
}
func Locusline(line string) (name string, seqlength int, seqtype string, circular bool, date string, err error) {

	fields := strings.SplitN(line, " ", 2)
	//// // fmt.Println("length of fields", len(fields))

	restofline := fields[1]

	fields = strings.Split(restofline, " ")
	//// // fmt.Println("length of fields", len(fields))

	newarray := make([]string, 0)
	for _, s := range fields {
		if s != "" && s != " " {
			newarray = append(newarray, s)
		}
	}
	fields = newarray
	//// // fmt.Println("length of fields", len(fields))
	//// fmt.Println(fields)
	if len(fields) > 1 {
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
		err = fmt.Errorf("invalid genbank line: ", line)
	}

	return
}
func Cleanup(line string) (cleanarray []string) {
	fields := strings.Split(line, " ")

	for _, s := range fields {

		if s != "" && s != " " {
			cleanarray = append(cleanarray, s)
		}

	}

	return
}

func Featureline1(line string) (reverse bool, class string, startposition int, endposition int, err error) {

	newarray := Cleanup(line)

	class = newarray[0]

	for _, s := range newarray {
		fmt.Println(newarray)
		if s[0] == '<' {
			s = s[1:]
		}
		if s[0] == '>' {
			s = s[1:]
		}
		fmt.Println(s)
		if strings.Contains(s, `join`) {
			err = fmt.Errorf("double position of feature!!", s, "adding as one feature only for now")
			s = strings.Replace(s, "Join(", "", -1)
			s = strings.Replace(s, ")", "", -1)
			//index := strings.Index(s, "..")
			joinhandler := strings.Split(s, `,`)
			split := strings.Split(joinhandler[0], "..")
			startposition, err = strconv.Atoi(split[0])

			split = strings.Split(joinhandler[1], "..")
			endposition, err = strconv.Atoi(strings.TrimRight(split[1], "\n"))

		} else {
			if strings.Contains(s, `complement`) {
				reverse = true
				s = strings.TrimLeft(s, `(complement)`)
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
					// fmt.Println(err.Error())
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
					// fmt.Println("trimmed", ss[1])
				} else if ss[1][0] == '<' {
					ss[1] = ss[1][1:]
					// fmt.Println("trimmed", ss[1])
				}
				endposition, err = strconv.Atoi(strings.TrimRight(ss[1], "\n"))
				// fmt.Println("trimmed", ss[1], "endposition", endposition)
				// fmt.Println("trimmed", s[0:index], "startposition", startposition)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(reverse, class, startposition, endposition)
			}
		}
	}
	fmt.Println(reverse, class, startposition, endposition)
	return
}
func Featureline2(line string) (description string, found bool) {

	fields := strings.Split(line, " ")
	//// // fmt.Println("length of fields", len(fields))

	// reassemble fields to preserve linked items with spaces e.g. "Green fluorescent protein"
	for i, field := range fields {
		if strings.Contains(field, `"`) {
			tempfields := make([]string, i)
			tempfield := strings.Join(fields[i:len(fields)-1], " ")
			tempfields = fields[0 : i-1]
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
				// // fmt.Println("line", line)
				// // fmt.Println("parts", parts)
				// // fmt.Println("len(parts) =2 yes")
				// // fmt.Println("parts[1]", parts[1])
				description = parts[1] //strings.Replace(parts[1], " ", "_", -1)
				// // fmt.Println("Huh!", description)
				found = true
				return
			}

		}

		if strings.Contains(line, `/label`) {

			parts := strings.SplitAfterN(line, "=", 2)
			if len(parts) == 2 {
				description = strings.TrimSpace(parts[1])
				found = true
				return
			}

		}

		if strings.Contains(line, `/product`) {

			parts := strings.SplitAfterN(line, `="`, 2)
			if len(parts) == 2 {
				// // fmt.Println("line", line)
				// // fmt.Println("parts", parts)
				// // fmt.Println("len(parts) =2 yes")
				// // fmt.Println("parts[1]", parts[1])
				description = parts[1] //strings.Replace(parts[1], " ", "_", -1)
				// // fmt.Println("Huh!", description)
				found = true
				return
			}

		}

	}
	/*for _, line := range newarray {
		if strings.Contains(line, `/product`) {
			parts := strings.SplitAfterN(line, `="`, 2)
			if len(parts) == 2 {
				// // fmt.Println("line", line)
				// // fmt.Println("parts", parts)
				// // fmt.Println("len(parts) =2 yes")
				// // fmt.Println("parts[1]", parts[1])
				description = parts[1] //strings.Replace(parts[1], " ", "_", -1)
				// // fmt.Println("Huh!", description)
				found = true
				return
			}

		}

	}*/
	return
}

func HandleFeature(lines []string) (description string, reverse bool, class string, startposition int, endposition int, err error) {

	if len(lines) > 0 {
		reverse, class, startposition, endposition, err := Featureline1(lines[0])
		//	// fmt.Println(reverse, class, startposition, endposition, err)

		if err != nil {
			fmt.Errorf("Error with Featureline1 func", lines[0])
			return description, reverse, class, startposition, endposition, err
		}
		for i := 1; i < len(lines); i++ {

			description, found := Featureline2(lines[i])
			if found {
				return description, reverse, class, startposition, endposition, err
			}

		}
	}
	return
}
func DetectFeature(lines []string) (detected bool, startlineindex int, endlineindex int) {
	for i := 0; i < len(lines); i++ {
		if string(lines[i][0]) != " " {
			return
		}
		if startlineindex != -1 && endlineindex != 0 {
			detected = true
			//		// // fmt.Println("Yay, detected")
			return
		}
		// // fmt.Println("linerz", lines[i])
		if string(lines[i][7]) != " " {
			startlineindex = i
			// // fmt.Println("start:", i, lines[i])
		}

		_, found := Featureline2(lines[i])
		if found {
			endlineindex = i + 1
			//		// // fmt.Println("end:", i, lines[i])
		}
	}

	return
}
func HandleFeatures(lines []string, seq string, seqtype string) (features []wtype.Feature) {

	featurespresent := false
	for _, line := range lines {
		if strings.Contains(line, "FEATURES") {
			featurespresent = true
		}
	}
	if featurespresent != true {
		return
	}
	features = make([]wtype.Feature, 0)
	var feature wtype.Feature

	for i := 0; i < len(lines); i++ { //, line := range lines {
		//	// fmt.Println(lines)
		//	// fmt.Println(line)
		if lines[i][0:8] == "FEATURES" {
			// fmt.Println(lines[i])
			lines = lines[i+1 : len(lines)]
			// // fmt.Println("broken")
			// fmt.Println(lines)
			//// fmt.Println(line)
			//// fmt.Println(lines[i])
			break
		}
	}
	// // fmt.Println("broken again")
	linesatstart := lines

	for i := 0; i < len(linesatstart); i++ {

		//jumpout := false

		if string(lines[0][0]) != " " {

			return
		}

		detected, start, end := DetectFeature(lines)
		// // fmt.Println("start", start, "end", end)
		if detected {
			// // fmt.Println("detected!!!!!!!!!!!!!", lines[start:end])

			description, reverse, class, startposition, endposition, err := HandleFeature(lines[start:end])
			// // fmt.Println("featuredectected: ", description, reverse, class, startposition, endposition, err)
			if err != nil {
				panic(err.Error())
			}
			rev := ""
			if reverse {
				rev = "Reverse"
			}

			// fmt.Println("seq,start,end = ", seq, startposition, endposition)

			// Warning! this needs to change to handle cases where start and position assignment has failed rather than just ignoring the problem
			if startposition != 0 && endposition != 0 {

				feature = sequences.MakeFeature(description, seq[startposition-1:endposition], seqtype, class, rev)
			}
			if start > end {
				return

				//	fmt.Println(startposition)
				//		if len(seq) > 0 /*&& startposition > 0 /*&& endposition < len(seq) */ {
				/*			feature = sequences.MakeFeature(description, seq[startposition-1:endposition], seqtype, class, rev)
								if start > end {
									return
								}

							} else {
								// fmt.Println("sequence", description, seq, "startposition", startposition, "endposition", endposition, " not valid")
				*/

			}
			features = append(features, feature)
			lines = lines[end:len(lines)]

		}

	}
	return

}

var (
	illegal string = "1234567890"
)

func HandleSequence(lines []string) (dnaseq string) {
	originallines := len(lines)
	originfound := false
	// fmt.Println(originallines)
	if len(lines) > 0 {
		for i := 0; i < originallines; i++ {

			// // fmt.Println("lines", lines[i])
			if len([]byte(lines[0])) > 0 {
				if originfound == false {
					if lines[i][0:6] == "ORIGIN" {
						originfound = true
					}
				}
				if originfound {

					// // fmt.Println("i+1", i, len(lines))
					// fmt.Println(lines[i+1])
					lines = lines[i+1 : originallines]
					seq := strings.Join(lines, "")
					seq = strings.Replace(seq, " ", "", -1)

					for _, character := range illegal {
						seq = strings.Replace(seq, string(character), "", -1)
					}
					seq = strings.Replace(seq, "\n", "", -1)
					seq = strings.Replace(seq, "//", "", -1)
					dnaseq = seq

					// // fmt.Println("dnaseq:", dnaseq)
					return
				}
			}
		}

	}
	return
}
