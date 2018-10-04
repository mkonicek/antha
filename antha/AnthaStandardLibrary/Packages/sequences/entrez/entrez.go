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

// package for querying all of NCBI databases
package entrez

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/genbank"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	biogo "github.com/biogo/ncbi/entrez"
)

var (
	email   = "no-reply@antha-lang.com"
	tool    = "entrez-biogo-antha"
	retries = 5
)

func appendError(err error, err2 error) (newErr error) {

	if err == nil && err2 == nil {
		return nil
	} else if err == nil {
		return err2
	} else if err2 == nil {
		return err
	}
	return fmt.Errorf("Errors: " + err.Error() + ": " + err2.Error())
}

// This queries the selected database saving the record to bytes
// Database options are nucleotide, Protein, Gene. For full list see http://www.ncbi.nlm.nih.gov/books/NBK25497/table/chapter2.T._entrez_unique_identifiers_ui/?report=objectonly
// Return type includes but must match the database type. See http://www.ncbi.nlm.nih.gov/books/NBK25499/table/chapter4.T._valid_values_of__retmode_and/?report=objectonly
// Query can be any string but it is recommended to use GI number if one specific record is required.
func RetrieveRecords(query string, database string, Max int, ReturnType string) (contentsinbytes []byte, err error) {
	// query database

	h := biogo.History{}
	s, err := biogo.DoSearch(database, query, nil, &h, tool, email)
	if err != nil {
		return []byte{}, appendError(fmt.Errorf("Error in biogo.DoSearch: "), err)
	}

	var (
		buf = &bytes.Buffer{}
		p   = &biogo.Parameters{RetMax: Max, RetType: ReturnType, RetMode: "text"}
		bn  int64
	)

	for p.RetStart = 0; p.RetStart < s.Count; p.RetStart += p.RetMax {
		var t int
		for t = 0; t < retries; t++ {
			sleep := time.Duration(1) * time.Second // limit queries to < 3 per second
			time.Sleep(sleep)

			var (
				r   io.ReadCloser
				_bn int64
			)
			r, err = biogo.Fetch(database, p, tool, email, &h)
			if err != nil {
				if r != nil {
					r.Close() // nolint
				}
				continue
			}

			_bn, err = io.Copy(buf, r)
			bn += _bn
			r.Close() // nolint
			if err == nil {
				break
			}
		}
		if err != nil {
			return []byte{}, appendError(fmt.Errorf("Error in fetching record %s", query), err)
		}
	}

	contentsinbytes, err = ioutil.ReadAll(buf)

	if err != nil {
		return contentsinbytes, fmt.Errorf("Error reading record: %s", err.Error())
	}

	if len(contentsinbytes) == 0 {
		return contentsinbytes, fmt.Errorf("no data returned from looking up query %s in %s", query, database)
	}
	return contentsinbytes, nil
}

// This retrieves sequence of any type from any NCBI sequence database
func RetrieveSequence(id string, database string) (seq wtype.DNASequence, err error) {

	contents, err := RetrieveRecords(id, database, 1, "gb")

	if err != nil {
		return wtype.DNASequence{}, err
	}

	seq, err = genbank.GenbankContentsToAnnotatedSeq(contents)
	if err != nil {
		return wtype.DNASequence{}, err
	}
	seq.Seq = strings.ToUpper(seq.Seq)

	return seq, err
}

// This will retrieve vector using fasta or db
func RetrieveVector(id string) (seq wtype.DNASequence, err error) {
	/*//first check if vector sequence is in fasta file
	if seq, err := parser.RetrieveSeqFromFASTA(id, filepath.Join(anthapath.Path(), "vectors.txt")); err != nil {
		// if not in refactor, check db*/
	seq, err = RetrieveSequence(id, "nucleotide")
	return

}
