// Part of the Antha language
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

// Package for processing sequencing results
package seqtools

import (
	"fmt"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Npercent(NumberofN int, FullSeq string) float64 {
	a := float64(NumberofN)
	b := float64(len(FullSeq))
	var undet float64 = (a / b * 100)

	return undet
}

func CountN(Seq string) int {
	n := strings.Count(Seq, "N")
	return n
}

func ReadtoCsvfromcurrentdir(filenametocreate string) {

	//Make an output file
	//filename := "Sequencing_results.csv"

	file, err := os.Create(filenametocreate)

	if err != nil {
		fmt.Println(err)
	}

	//Write headers to output file
	// fmt.Println("Sequence data written to file: " + filenametocreate)
	var headers string = "Filename,	Sequence length (nts),	Proprotion undetermined (%),	Sequence"
	n, err1 := io.WriteString(file, headers)

	if err1 != nil {
		fmt.Println(n, err1)
	}

	file.Close() // nolint

	//Search for files within current directory
	dirname := "." + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close() // nolint

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println("Reading " + dirname)
	var skip bool
	filesdone := make([]string, 0)
	//Determine if file extension is ".seq"
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".seq" {
			if len(filesdone) != 0 {
				for _, filedone := range filesdone {

					if file.Name() == filedone {
						skip = true
					}

				}
			}
			if !skip {
				//Read file containing sequencing result
				bs, err := ioutil.ReadFile(file.Name())
				if err != nil {
					return
				}

				var name1 string = file.Name()

				//Print the filename
				// fmt.Println("Sequence filename:", name1)

				//assign sequence result to a variable
				seq1 := string(bs)
				Seq := strings.Replace(seq1, "\n", "", -1)

				//print the sequence result
				// fmt.Println("Sequencing result =", Seq)

				//print the length of the sequencing run
				// fmt.Println("Sequencing run length =", len(Seq))

				//count the number of nucleotides which have been designated "n" (undetermined) and print the amount
				N := CountN(Seq)
				// fmt.Println("Nucleotides not sequenced =", N)

				//print the proportion of undetermined nucleotides in the sequence
				fmt.Printf("Proportion of sequence not determined = %0.2f %% \n", Npercent(N, Seq))

				//Append sequence data to the file
				f, err := os.OpenFile("Sequencing_results.csv", os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err)
				}

				defer f.Close() // nolint

				if _, err = f.WriteString(fmt.Sprintf("\n %s, %v, %0.2f, %s", name1, len(Seq), Npercent(N, Seq), Seq)); err != nil {
					panic(err)

				}
				filesdone = append(filesdone, file.Name())

			}
		}
	}
}
