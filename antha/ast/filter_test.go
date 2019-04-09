// antha/ast/filter_test.go: Part of the Antha language
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

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// To avoid a cyclic dependency with antha/parser, this file is in a separate package.

package ast_test

import (
	"bytes"
	"testing"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/format"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

const input = `package p

type t1 struct{}
type t2 struct{}

func f1() {}
func f1() {}
func f2() {}

func (*t1) f1() {}
func (t1) f1() {}
func (t1) f2() {}

func (t2) f1() {}
func (t2) f2() {}
func (x *t2) f2() {}
`

// Calling ast.MergePackageFiles with ast.FilterFuncDuplicates
// keeps a duplicate entry with attached documentation in favor
// of one without, and it favors duplicate entries appearing
// later in the source over ones appearing earlier. This is why
// (*t2).f2 is kept and t2.f2 is eliminated in this test case.
//
const golden = `package p

type t1 struct{}
type t2 struct{}

func f1() {}
func f2() {}

func (t1) f1() {}
func (t1) f2() {}

func (t2) f1() {}

func (x *t2) f2() {}
`

func TestFilterDuplicates(t *testing.T) {
	// parse input
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", []byte(input), 0)
	if err != nil {
		t.Fatal(err)
	}

	// create package
	files := map[string]*ast.File{"": file}
	pkg, err := ast.NewPackage(fset, files, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// filter
	merged := ast.MergePackageFiles(pkg, ast.FilterFuncDuplicates)

	// pretty-print
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, merged); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	if output != golden {
		t.Errorf("incorrect output:\n%s", output)
	}
}
