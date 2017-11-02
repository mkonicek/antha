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
	"fmt"

	"bufio"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
)

var (
	alphabet = sequences.WobbleMap
)

func SnapGenetoSimpleSeq(filename string) (string, error) {

	var line string
	var snapgenelines []string
	if !strings.Contains(filename, ".dna") {
		return "", fmt.Errorf("wrong file type, must have file extension .dna")
	}
	contents, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())

		snapgenelines = append(snapgenelines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return Handlesnapgenelines(snapgenelines), nil
}

func Handlesnapgenelines(lines []string) (dnaseq string) {
	originallines := len(lines)
	//startfound := false
	seqlines := make([]string, 0)
	// fmt.Println(originallines)
	if len(lines) > 0 {
		for i := 3; i < originallines; i++ {

			// // fmt.Println("lines", lines[i])
			/*	if len([]byte(lines[0])) > 0 {
				if startfound == false {
					if len([]byte(lines[i])) > 0 {
						startfound = true
					}
				}*/
			//if startfound {

			// // fmt.Println("i+1", i, len(lines))
			// fmt.Println(lines[i+1])
			if len([]byte(lines[i])) > 0 {
				//// fmt.Println("line:", lines[i])
				ok, _, _ := noillegalshere(lines[i])

				if ok != true {
					templine := removeweirdthings(lines[i])
					//	// fmt.Println("original line:", lines[i], "templine:", templine, string(badchars))
					seqlines = append(seqlines, templine)
					break
				} else {
					seqlines = append(seqlines, lines[i])
				}
			} else {
				break
			}
		}

		seq := strings.Join(lines, "")
		seq = strings.Replace(seq, " ", "", -1)

		ok, _, _ := noillegalshere(seq)
		//// fmt.Println("Handling this", ok, badpositions, string(badcharacters))
		if !ok {
			//	// fmt.Println("seq1", seq)
			seq = removeweirdthings(seq)
			//	// fmt.Println("seq2", seq)
		}
		dnaseq = seq

		// // fmt.Println("dnaseq:", dnaseq)

	}
	return
}

/*
func somethingweirdatstart(seq string) (found bool, weirdthing string) {

	var badposition int
	var badpositions = make([]int, 0)
	var badcharacter rune

	for i, letter := range seq {
		for _, valid := range alphabet {
			if letter == valid {
				badposition = i
				badpositions = append(badpositions, badposition)
				return
			}
		}
	}

	if len(badpositions) == 0 {
		nothingfound = true
	}
	return
}

func somethingweirdatend(seq string) (found bool, weirdthing string) {

	return
}
*/
func removeweirdthings(seq string) (weirdthingfreeseq string) {

	if len(seq) == 1 {
		_, _, badcharacters := noillegalshere(string(seq))

		if len(badcharacters) > 0 {
			weirdthingfreeseq = ""
			return
		}
	}

	temp := seq

	_, _, badcharacters := noillegalshere(temp)

	for _, badcharacter := range badcharacters {
		temp = strings.Replace(temp, string(badcharacter), "", -1)
		//	// fmt.Println("temp =", temp, "bad character:", string(badcharacter), "all bad characters:", badcharacters)

	}

	weirdthingfreeseq = temp
	//	// fmt.Println("weirdthingfreeseq", weirdthingfreeseq)
	return
}

func noillegalshere(line string) (nothingfound bool, badpositions []int, badcharacters []rune) {

	/*dnaseq := wtype.MakeSingleStrandedDNASequence("temp", line)
	_, illegals, _ := sequences.Illegalnucleotides(dnaseq)

	if len(illegals) == 0 {
		nothingfound = true
	}*/
	var badposition int
	badpositions = make([]int, 0)
	badcharacters = make([]rune, 0)

	for i, letter := range line {
		//for _, valid := range alphabet {

		_, foundinmap := alphabet[string(letter)]

		if foundinmap == false {

			//if letter == valid {
			badposition = i
			badpositions = append(badpositions, badposition)
			badcharacters = append(badcharacters, letter)

		}
	}

	if len(badpositions) == 0 {
		nothingfound = true
	}
	return
}
