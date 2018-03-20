// microArch/frontend/socket/socket_test.go: Part of the Antha language
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

package socket

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalMessageType(t *testing.T) {
	mt := MessageType(1)
	out, err := json.Marshal(mt)
	if err != nil {
		t.Fatal(err)
	}
	exp := `"display"`
	if string(out) != exp {
		t.Fatal(fmt.Sprintf("Expecting %s. Got %s.", mt, out))
	}
}

func TestBuildMessageSocketMessage(t *testing.T) {
	testSuite := make(map[string]string)

	testSuite[`{"Display":[{"Locale":"","Message":"Hello World"}],"Request":{"Type":"display","Options":null}}`] = `Hello World`

	for expected, message := range testSuite {
		res, err := json.Marshal(*buildMessageSocketMessage(message))
		if err != nil {
			t.Fatal(err)
		}
		if string(res) != expected {
			t.Fatal(fmt.Sprintf("Expecting %s. Got %s.", expected, string(res)))
		}
	}
}
