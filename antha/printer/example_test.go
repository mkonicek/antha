// antha/printer/example_test.go: Part of the Antha language
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

package printer_test

import (
	"testing"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

// Dummy test function so that anthadoc does not use the entire file as example.
func Test(*testing.T) {}

func parseFunc(filename, functionname string) (fun *ast.FuncDecl, fset *token.FileSet) {
	fset = token.NewFileSet()
	if file, err := parser.ParseFile(fset, filename, nil, 0); err == nil {
		for _, d := range file.Decls {
			if f, ok := d.(*ast.FuncDecl); ok && f.Name.Name == functionname {
				fun = f
				return
			}
		}
	}
	panic("function not found")
}

//func ExampleFprint() {
//	// Parse source file and extract the AST without comments for
//	// this function, with position information referring to the
//	// file set fset.
//	funcAST, fset := parseFunc("example_test.go", "ExampleFprint")
//
//	// Print the function body into buffer buf.
//	// The file set is provided to the printer so that it knows
//	// about the original source formatting and can add additional
//	// line breaks where they were present in the source.
//	var buf bytes.Buffer
//	printer.Fprint(&buf, fset, funcAST.Body)
//
//	// Remove braces {} enclosing the function body, unindent,
//	// and trim leading and trailing white space.
//	s := buf.String()
//	s = s[1 : len(s)-1]
//	s = strings.TrimSpace(strings.Replace(s, "\n\t", "\n", -1))
//
//	// Print the cleaned-up body text to stdout.
//	fmt.Println(s)
//
//	// output:
//	// funcAST, fset := parseFunc("example_test.go", "ExampleFprint")
//	//
//	// var buf bytes.Buffer
//	// printer.Fprint(&buf, fset, funcAST.Body)
//	//
//	// s := buf.String()
//	// s = s[1 : len(s)-1]
//	// s = strings.TrimSpace(strings.Replace(s, "\n\t", "\n", -1))
//	//
//	// fmt.Println(s)
//}
