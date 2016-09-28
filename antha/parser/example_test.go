// antha/parser/example_test.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
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

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

//func ExampleParseFile() {
//	fset := token.NewFileSet() // positions are relative to fset
//
//	// Parse the file containing this very example
//	// but stop after processing the imports.
//	f, err := parser.ParseFile(fset, "example_test.go", nil, parser.ImportsOnly)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// Print the imports from the file's AST.
//	for _, s := range f.Imports {
//		fmt.Println(s.Path.Value)
//	}
//
//	// output:
//	//
//	// "fmt"
//	// "github.com/antha-lang/antha/antha/parser"
//	// "github.com/antha-lang/antha/antha/token"
//}
