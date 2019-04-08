// antha/ast/example_test.go: Part of the Antha language
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

package ast_test

import (
	"bytes"
	"fmt"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/format"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

// This example demonstrates how to inspect the AST of a Go/Antha program.
func ExampleInspect() {
	// src is the input for which we want to inspect the AST.
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", []byte(src), 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})

	// output:
	// src.go:2:9:	p
	// src.go:3:7:	c
	// src.go:3:11:	1.0
	// src.go:4:5:	X
	// src.go:4:9:	f
	// src.go:4:11:	3.14
	// src.go:4:17:	2
	// src.go:4:21:	c
}

// This example shows what an AST looks like when printed for debugging.
func ExamplePrint() {
	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
	println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", []byte(src), 0)
	if err != nil {
		panic(err)
	}
	// Print the AST.
	ast.Print(fset, f)

	// output:
	//0  *ast.File {
	//      1  .  Package: 2:1
	//      2  .  Tok: package
	//      3  .  Name: *ast.Ident {
	//      4  .  .  NamePos: 2:9
	//      5  .  .  Name: "main"
	//      6  .  }
	//      7  .  Decls: []ast.Decl (len = 1) {
	//      8  .  .  0: *ast.FuncDecl {
	//      9  .  .  .  Name: *ast.Ident {
	//     10  .  .  .  .  NamePos: 3:6
	//     11  .  .  .  .  Name: "main"
	//     12  .  .  .  .  Obj: *ast.Object {
	//     13  .  .  .  .  .  Kind: func
	//     14  .  .  .  .  .  Name: "main"
	//     15  .  .  .  .  .  Decl: *(obj @ 8)
	//     16  .  .  .  .  }
	//     17  .  .  .  }
	//     18  .  .  .  Type: *ast.FuncType {
	//     19  .  .  .  .  Func: 3:1
	//     20  .  .  .  .  Params: *ast.FieldList {
	//     21  .  .  .  .  .  Opening: 3:10
	//     22  .  .  .  .  .  Closing: 3:11
	//     23  .  .  .  .  }
	//     24  .  .  .  }
	//     25  .  .  .  Body: *ast.BlockStmt {
	//     26  .  .  .  .  Lbrace: 3:13
	//     27  .  .  .  .  List: []ast.Stmt (len = 1) {
	//     28  .  .  .  .  .  0: *ast.ExprStmt {
	//     29  .  .  .  .  .  .  X: *ast.CallExpr {
	//     30  .  .  .  .  .  .  .  Fun: *ast.Ident {
	//     31  .  .  .  .  .  .  .  .  NamePos: 4:2
	//     32  .  .  .  .  .  .  .  .  Name: "println"
	//     33  .  .  .  .  .  .  .  }
	//     34  .  .  .  .  .  .  .  Lparen: 4:9
	//     35  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
	//     36  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
	//     37  .  .  .  .  .  .  .  .  .  ValuePos: 4:10
	//     38  .  .  .  .  .  .  .  .  .  Kind: STRING
	//     39  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
	//     40  .  .  .  .  .  .  .  .  }
	//     41  .  .  .  .  .  .  .  }
	//     42  .  .  .  .  .  .  .  Ellipsis: -
	//     43  .  .  .  .  .  .  .  Rparen: 4:25
	//     44  .  .  .  .  .  .  }
	//     45  .  .  .  .  .  }
	//     46  .  .  .  .  }
	//     47  .  .  .  .  Rbrace: 5:1
	//     48  .  .  .  }
	//     49  .  .  }
	//     50  .  }
	//     51  .  Scope: *ast.Scope {
	//     52  .  .  Objects: map[string]*ast.Object (len = 1) {
	//     53  .  .  .  "main": *(obj @ 12)
	//     54  .  .  }
	//     55  .  }
	//     56  .  Unresolved: []*ast.Ident (len = 1) {
	//     57  .  .  0: *(obj @ 30)
	//     58  .  }
	//     59  }
}

// This example illustrates how to remove a variable declaration
// in a Go/Antha program while maintaining correct comment association
// using an ast.CommentMap.
func ExampleCommentMap() {
	// src is the input for which we create the AST that we
	// are going to manipulate.
	src := `
// This is the package comment.
package main

// This comment is associated with the hello constant.
const hello = "Hello, World!" // line comment 1

// This comment is associated with the foo variable.
var foo = hello // line comment 2 

// This comment is associated with the main function.
func main() {
	fmt.Println(hello) // line comment 3
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", []byte(src), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Create an ast.CommentMap from the ast.File's comments.
	// This helps keeping the association between comments
	// and AST nodes.
	cmap := ast.NewCommentMap(fset, f, f.Comments)

	// Remove the first variable declaration from the list of declarations.
	f.Decls = removeFirstVarDecl(f.Decls)

	// Use the comment map to filter comments that don't belong anymore
	// (the comments associated with the variable declaration), and create
	// the new comments list.
	f.Comments = cmap.Filter(f).Comments()

	// Print the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	// output:
	// // This is the package comment.
	// package main
	//
	// // This comment is associated with the hello constant.
	// const hello = "Hello, World!" // line comment 1
	//
	// // This comment is associated with the main function.
	// func main() {
	// 	fmt.Println(hello) // line comment 3
	// }
}

func removeFirstVarDecl(list []ast.Decl) []ast.Decl {
	for i, decl := range list {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
			copy(list[i:], list[i+1:])
			return list[:len(list)-1]
		}
	}
	panic("variable declaration not found")
}
