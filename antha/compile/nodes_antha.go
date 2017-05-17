// antha/compile/nodes_antha.go: Part of the Antha language
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

// Package compile implements the antha source-to-source transformation
package compile

import (
	"bytes"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/token"
)

// Called from node.go walker to print antha code blocks
func (p *compiler) anthaDecl(d *ast.AnthaDecl) {
	switch d.Tok {
	case token.STEPS, token.SETUP, token.REQUIREMENTS, token.ANALYSIS, token.VALIDATION:
		ctx := new(anthaContext)
		ctx.init(p.element, d.Tok)
		p.anthaSig(d, ctx)
	default:
		panic("Bad anthaDecl")

	}
}

// Outputs the Antha block signatures
func (p *compiler) anthaSig(d *ast.AnthaDecl, c *anthaContext) {
	p.setComment(d.Doc)

	p.print(d.Pos(), c.getSignature())

	// adjust the padding to the body
	p.adjBlock(p.distanceFrom(d.Pos()), vtab, d.Body)
}

// antha generator context structure
// Note: this code is not thread safe
type anthaContext struct {
	PkgName   string
	emitStart bool
	emitEnd   bool
	tok       token.Token
}

// Special case blocks
//	PARAMETERS
//	DATA
//	INPUTS
//	OUTPUTS

var anthaSigs = map[token.Token]string{
	token.STEPS:        `func _{{.PkgName}}Steps(_ctx context.Context, _input *{{.PkgName}}Input, _output *{{.PkgName}}Output)`,
	token.REQUIREMENTS: `func _{{.PkgName}}Requirements()`,
	token.SETUP:        `func _{{.PkgName}}Setup(_ctx context.Context, _input *{{.PkgName}}Input)`,
	token.ANALYSIS:     `func _{{.PkgName}}Analysis(_ctx context.Context, _input *{{.PkgName}}Input, _output *{{.PkgName}}Output)`,
	token.VALIDATION:   `func _{{.PkgName}}Validation(_ctx context.Context, _input *{{.PkgName}}Input, _output *{{.PkgName}}Output)`,
}

// templates for any start lines in block
var anthaStarts = map[token.Token]string{
	token.STEPS:        ``,
	token.REQUIREMENTS: ``,
	token.SETUP:        ``,
	token.ANALYSIS:     ``,
	token.VALIDATION:   ``,
}

// templates for any end lines in block
var anthaEnds = map[token.Token]string{
	token.STEPS:        ``,
	token.REQUIREMENTS: ``,
	token.SETUP:        ``,
	token.ANALYSIS:     ``,
	token.VALIDATION:   ``,
}

// init the context to the appropriate type of antha block
func (c *anthaContext) init(name string, tok token.Token) {
	c.emitStart = true
	c.emitEnd = true
	c.tok = tok
	c.PkgName = name
}

// simple helper function to generate the appropriate block
// function signatures
func (c *anthaContext) getSignature() (b bytes.Buffer) {
	t := template.Must(template.New(c.tok.String() + "_sig").Parse(anthaSigs[c.tok]))
	t.Execute(&b, c)
	return
}

// simple helper function to generate the beginning of an antha block
// depending on the block type
func (c *anthaContext) beginBlock() (b bytes.Buffer, empty bool) {
	if !c.emitStart {
		empty = true
		return
	}
	t := template.Must(template.New(c.tok.String() + "_start").Parse(anthaStarts[c.tok]))
	t.Execute(&b, c)
	empty = b.Len() > 0
	c.emitStart = false
	return
}

// simple helper function to generate the ending of an antha block
// depending on the block type
func (c *anthaContext) endBlock() (b bytes.Buffer, empty bool) {
	if !c.emitEnd {
		empty = true
		return
	}
	t := template.Must(template.New(c.tok.String() + "_end").Parse(anthaEnds[c.tok]))
	t.Execute(&b, c)
	empty = b.Len() > 0
	c.emitEnd = false
	return
}
