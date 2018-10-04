// /compile/generator_test.go: Part of the Antha language
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

package compile

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

func TestTypeSugaring(t *testing.T) {
	nodeSizes := make(map[ast.Node]int)
	cfg := &Config{}
	compiler := &compiler{}
	fset := token.NewFileSet()
	compiler.init(cfg, fset, nodeSizes)

	root := NewAnthaRoot("")
	antha := NewAntha(root)
	expr, err := parser.ParseExpr("func(x Volume) Concentration { x := Volume }")
	if err != nil {
		t.Fatal(err)
	}
	desired, err := parser.ParseExpr("func(x wunit.Volume) wunit.Concentration { x := Volume }")
	if err != nil {
		t.Fatal(err)
	}

	ast.Inspect(expr, antha.inspectTypes)
	var buf1, buf2 bytes.Buffer
	if _, err := compiler.Fprint(&buf1, fset, expr); err != nil {
		t.Fatal(err)
	}
	if _, err := compiler.Fprint(&buf2, fset, desired); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Errorf("wanted\n'''%s'''\ngot\n'''%s'''\n", buf2.String(), buf1.String())
	}
}

func TestRelativeToGoPath(t *testing.T) {
	type TestCase struct {
		GoPath   []string
		Name     string
		Expected string
		HasError bool
	}

	cases := []TestCase{
		{GoPath: []string{"/xx"}, Name: "/xx/file", Expected: "file"},
		{GoPath: []string{"/xx"}, Name: "file", Expected: "file"},
		{GoPath: []string{"/noxx"}, Name: "/xx/file", Expected: "/xx/file"},
		{GoPath: []string{"/xx", "/xx/deeper"}, Name: "/xx/file", Expected: "file"},
		{GoPath: []string{"/xx", "/xx/deeper"}, Name: "/xx/deeper/file", Expected: "file"},
	}

	for idx, c := range cases {
		var goPath []string
		for _, v := range c.GoPath {
			goPath = append(goPath, filepath.FromSlash(v))
		}
		name := filepath.FromSlash(c.Name)
		expected := filepath.FromSlash(c.Expected)

		f, err := relativeTo(goPath, name)
		if c.HasError && err == nil {
			t.Errorf("%d: %+v: expected error but found success", idx, c)
		} else if err != nil {
			t.Errorf("%d: %+v: %s", idx, c, err)
		} else if expected != f {
			t.Errorf("%d: %+v: expected %q found %q", idx, c, expected, f)
		}
	}
}
