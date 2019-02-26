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

package parse

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader/dataset"
)

// ParseSpectraMaxData parses spectra data from an XML file
func ParseSpectraMaxData(xmlFileContents []byte) (dataOutput dataset.SpectraMaxData, err error) {

	s, err := decodeUTF16(xmlFileContents)
	if err != nil {
		panic(err)
	}

	utf8XMLContents := []byte(s)

	err = xml.Unmarshal(utf8XMLContents, &dataOutput)

	if err != nil {
		fmt.Println("error:", err)
	}

	return
}

func decodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		_, err := ret.Write(b8buf[:n])
		if err != nil {
			return "", err
		}
	}

	return ret.String(), nil
}
