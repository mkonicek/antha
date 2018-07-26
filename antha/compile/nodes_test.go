// /compile/nodes_test.go: Part of the Antha language
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
	"testing"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

var fset = token.NewFileSet()

// TestLineComments, using a simple test case, checks that consequtive line
// comments are properly terminated with a newline even if the AST position
// information is incorrect.
//
func TestLineComments(t *testing.T) {
	const src = `// comment 1
	// comment 2
	// comment 3
	package main
	`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err) // error in test
	}

	var buf bytes.Buffer
	fset = token.NewFileSet() // use the wrong file set
	if err := Fprint(&buf, fset, f); err != nil {
		t.Fatal(err)
	}

	nlines := 0
	for _, ch := range buf.Bytes() {
		if ch == '\n' {
			nlines++
		}
	}

	const expected = 3
	if nlines < expected {
		t.Errorf("got %d, expected %d\n", nlines, expected)
		t.Errorf("result:\n%s", buf.Bytes())
	}
}

// Verify that the printer can be invoked during initialization.
func init() {
	const name = "foobar"
	var buf bytes.Buffer
	if err := Fprint(&buf, fset, &ast.Ident{Name: name}); err != nil {
		panic(err) // error in test
	}
	// in debug mode, the result contains additional information;
	// ignore it
	if s := buf.String(); !debug && s != name {
		panic("got " + s + ", want " + name)
	}
}

// Verify that the printer doesn't crash if the AST contains BadXXX nodes.

// ddn: test disabled because additional boilerplate throws off comparisons
func TestIllegalProgram(t *testing.T) {
	t.Skip("disabled due to boilerplate")

	const src = "package p\n("
	//	const res = "package p\nBadDecl\n"
	const res = "package p\n\nimport \"github.com/antha-lang/antha/antha/execute\"\nimport \"github.com/Synthace/goflow\"\nimport \"sync\"\nimport \"encoding/json\"\n//import \"log\"\n//import \"bytes\"\n//import \"io\"\n\n\n\nBadDecl\n// AsyncBag functions\nfunc (e *P) Complete(params interface{}) {\n\tp := params.(PParamBlock)\n\tif p.Error {\n\n\t\treturn\n\t}\n\tr := new(PResultBlock)\n\te.startup.Do(func() { e.setup(p) })\n\te.steps(p, r)\n\tif r.Error {\n\n\t\treturn\n\t}\n\n\te.analysis(p, r)\n\t\tif r.Error {\n\n\n\t\treturn\n\t}\n\n\te.validation(p, r)\n\t\tif r.Error {\n\n\t\treturn\n\t}\n\n}\n\n// empty function for interface support\nfunc (e *P) anthaElement() {}\n\n// init function, read characterization info from seperate file to validate ranges?\nfunc (e *P) init() {\n\te.params = make(map[execute.ThreadID]*execute.AsyncBag)\n}\n\nfunc NewP() interface{} {//*P {\n\te := new(P)\n\te.init()\n\treturn e\n}\n\n// Mapper function\nfunc (e *P) Map(m map[string]interface{}) interface{} {\n\tvar res PParamBlock\n\tres.Error = false \n\n\n\treturn res\n}\n\n\ntype P struct {\n\tflow.Component                    // component \"superclass\" embedded\n\tlock           sync.Mutex\n\tstartup        sync.Once\n\tparams         map[execute.ThreadID]*execute.AsyncBag\n}\n\ntype PParamBlock struct{\n\tID\t\texecute.ThreadID\n\tError\tbool\n}\ntype PResultBlock struct{\n\tID\t\texecute.ThreadID\n\tError\tbool\n}\ntype PJSONBlock struct{\n\tID\t\t\t*execute.ThreadID\n\tError\t\t*bool\n}\n" //TODO review the intention of this test
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err == nil {
		t.Error("expected illegal program") // error in test
	}
	var buf bytes.Buffer
	if err := Fprint(&buf, fset, f); err != nil {
		t.Fatal(err)
	}
	if buf.String() != res {
		t.Errorf("got %q, expected %q", buf.String(), res)
	}
}

// Verify that the printer doesn't crash if the AST contains BadXXX nodes.
func TestBadNodes(t *testing.T) {
	const src = "package p\n("
	_, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err == nil {
		t.Error("expected illegal program") // error in test
	}
}

// testComment verifies that f can be parsed again after printing it
// with its first comment set to comment at any possible source offset.
func testComment(t *testing.T, f *ast.File, srclen int, comment *ast.Comment) {
	f.Comments[0].List[0] = comment
	var buf bytes.Buffer
	for offs := 0; offs <= srclen; offs++ {
		buf.Reset()
		// Printing f should result in a correct program no
		// matter what the (incorrect) comment position is.
		if err := Fprint(&buf, fset, f); err != nil {
			t.Error(err)
		}
		if _, err := parser.ParseFile(fset, "", buf.Bytes(), 0); err != nil {
			t.Fatalf("incorrect program for pos = %d:\n%s", comment.Slash, buf.String())
		}
		// Position information is just an offset.
		// Move comment one byte down in the source.
		comment.Slash++
	}
}

// Verify that the printer produces always produces a correct program
// even if the position information of comments introducing newlines
// is incorrect.
func TestBadComments(t *testing.T) {
	t.Skip("not supported in antha")
	const src = `
// first comment - text and position changed by test
package p
import "fmt"
const pi = 3.14 // rough circle
var (
	x, y, z int = 1, 2, 3
	u, v float64
)
func fibo(n int) {
	if n < 2 {
		return n /* seed values */
	}
	return fibo(n-1) + fibo(n-2)
}
`

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Error(err) // error in test
	}

	comment := f.Comments[0].List[0]
	pos := comment.Pos()
	if fset.Position(pos).Offset != 1 {
		t.Error("expected offset 1") // error in test
	}

	testComment(t, f, len(src), &ast.Comment{Slash: pos, Text: "//-style comment"})
	testComment(t, f, len(src), &ast.Comment{Slash: pos, Text: "/*-style comment */"})
	testComment(t, f, len(src), &ast.Comment{Slash: pos, Text: "/*-style \n comment */"})
	testComment(t, f, len(src), &ast.Comment{Slash: pos, Text: "/*-style comment \n\n\n */"})
}

type visitor chan *ast.Ident

func (v visitor) Visit(n ast.Node) (w ast.Visitor) {
	if ident, ok := n.(*ast.Ident); ok {
		v <- ident
	}
	return v
}

// idents is an iterator that returns all idents in f via the result channel.
func idents(f *ast.File) <-chan *ast.Ident {
	v := make(visitor)
	go func() {
		ast.Walk(v, f)
		close(v)
	}()
	return v
}

// identCount returns the number of identifiers found in f.
func identCount(f *ast.File) int {
	n := 0
	for range idents(f) {
		n++
	}
	return n
}

// NB(ddn): test disabled because additional boilerplate throws off comparisons
func TestSourcePos(t *testing.T) {
	t.Skip("disabled due to boilerplate")

	const src = `
package p
import ( "go/printer"; "math" )
const pi = 3.14; var x = 0
type t struct{ x, y, z int; u, v, w float32 }
func (t *t) foo(a, b, c int) int {
	return a*t.x + b*t.y +
		// two extra lines here
		// ...
		c*t.z
}
`

	// parse original
	f1, err := parser.ParseFile(fset, "src", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	// pretty-print original
	var buf bytes.Buffer
	_, err = (&Config{Mode: UseSpaces | SourcePos, Tabwidth: 8}).Fprint(&buf, fset, f1)
	if err != nil {
		t.Fatal(err)
	}

	// parse pretty printed original
	// (//line comments must be interpreted even w/o parser.ParseComments set)
	f2, err := parser.ParseFile(fset, "", buf.Bytes(), 0)
	if err != nil {
		t.Fatalf("%s\n%s", err, buf.Bytes())
	}

	// At this point the position information of identifiers in f2 should
	// match the position information of corresponding identifiers in f1.

	// number of identifiers must be > 0 (test should run) and must match
	n1 := identCount(f1)
	n2 := identCount(f2)
	if n1 == 0 {
		t.Fatal("got no idents")
	}
	if n2 != n1 {
		t.Errorf("got %d idents; want %d", n2, n1)
	}

	// verify that all identifiers have correct line information
	i2range := idents(f2)
	for i1 := range idents(f1) {
		i2 := <-i2range

		if i2.Name != i1.Name {
			t.Errorf("got ident %s; want %s", i2.Name, i1.Name)
		}

		l1 := fset.Position(i1.Pos()).Line
		l2 := fset.Position(i2.Pos()).Line
		if l2 != l1 {
			t.Errorf("got line %d; want %d for %s", l2, l1, i1.Name)
		}
	}

	if t.Failed() {
		t.Logf("\n%s", buf.Bytes())
	}
}

// TextX is a skeleton test that can be filled in for debugging one-off cases.
// Do not remove.
func TestX(t *testing.T) {
	t.Skip()
	const src = `
package p
func _() {}
`
	// parse original
	f, err := parser.ParseFile(fset, "src", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	// pretty-print original
	var buf bytes.Buffer
	if _, err = (&Config{Mode: UseSpaces, Tabwidth: 8}).Fprint(&buf, fset, f); err != nil {
		t.Fatal(err)
	}

	// parse pretty printed original
	if _, err := parser.ParseFile(fset, "", buf.Bytes(), 0); err != nil {
		t.Fatalf("%s\n%s", err, buf.Bytes())
	}

}
