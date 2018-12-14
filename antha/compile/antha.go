// antha.go: Part of the Antha language
// Copyright (C) 2017 The Antha authors. All rights reserved.
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
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

const (
	elementFilename = "element.go"
)

const (
	tabWidth    = 8
	printerMode = UseSpaces | TabIndent
)

var (
	errNotAnthaFile = errors.New("not antha file")
)

type posError struct {
	message string
	pos     token.Pos
}

func (e posError) Error() string {
	return e.message
}

func (e posError) Pos() token.Pos {
	return e.pos
}

func throwErrorf(pos token.Pos, format string, args ...interface{}) {
	panic(posError{
		message: fmt.Sprintf(format, args...),
		pos:     pos,
	})
}

// A Field is a field of a message
type Field struct {
	Name string
	Type ast.Expr // Fully qualified go type name
	Doc  string
	Tag  string
}

// A Message is an input or an output or user defined type
type Message struct {
	Name   string
	Doc    string
	Fields []*Field
	Kind   token.Token // One of token.{DATA, PARAMETERS, OUTPUTS, INPUTS}
}

func (m *Message) getFields() []*Field {
	if m == nil {
		return nil
	}
	return m.Fields
}

func isOutput(tok token.Token) bool {
	switch tok {
	case token.OUTPUTS, token.DATA:
		return true
	default:
		return false
	}
}

func isInput(tok token.Token) bool {
	switch tok {
	case token.INPUTS, token.PARAMETERS:
		return true
	default:
		return false
	}
}

func isAnthaGenDeclToken(tok token.Token) bool {
	switch tok {
	case token.OUTPUTS, token.DATA, token.PARAMETERS, token.INPUTS:
		return true
	default:
		return false
	}
}

func manglePackageName(pkg string) string {
	return "_" + hex.EncodeToString([]byte(pkg))
}

// An importReq is a request to add an import
type importReq struct {
	Path         string // Package path
	Name         string // Package identifier
	UseExpr      string // Dummy expression to supress unused imports
	Added        bool   // Has the import already been added?
	ProtoPackage string // Protobuf package this import cooresponds to
}

func (r *importReq) ImportName() string {
	if len(r.Name) != 0 {
		return r.Name
	}

	return path.Base(r.Path)
}

// Antha is a preprocessing pass from antha file to go file
type Antha struct {
	// Description of this element
	description string
	// Path to element file
	elementPath string
	// messages of an element as well as inputs and outputs
	messages []*Message
	// Protocol name as given in Antha file
	protocolName string

	root *AnthaRoot

	// inputs or outputs of an element but not messages
	tokenByParamName map[string]token.Token
	blocksUsed       map[token.Token]bool
	// Replacements for identifiers in expressions in functions
	intrinsics map[string]string
	// Replacement for go type names in type expressions and types and type lists
	types map[string]string
	// Imports in protocol and imports to add
	importReqs   []*importReq
	importByName map[string]*importReq
}

// NewAntha creates a new antha pass
func NewAntha(root *AnthaRoot) *Antha {
	p := &Antha{
		root:         root,
		importByName: make(map[string]*importReq),
	}

	p.intrinsics = map[string]string{
		"Centrifuge":    "execute.Centrifuge",
		"Electroshock":  "execute.Electroshock",
		"ExecuteMixes":  "execute.ExecuteMixes",
		"Handle":        "execute.Handle",
		"Incubate":      "execute.Incubate",
		"Mix":           "execute.Mix",
		"MixInto":       "execute.MixInto",
		"MixNamed":      "execute.MixNamed",
		"MixTo":         "execute.MixTo",
		"MixerPrompt":   "execute.MixerPrompt",
		"NewComponent":  "execute.NewComponent",
		"NewPlate":      "execute.NewPlate",
		"Prompt":        "execute.Prompt",
		"ReadEM":        "execute.ReadEM",
		"Sample":        "execute.Sample",
		"SetInputPlate": "execute.SetInputPlate",
		"SplitSample":   "execute.SplitSample",
	}

	p.types = map[string]string{
		"Amount":               "wunit.Amount",
		"Angle":                "wunit.Angle",
		"AngularVelocity":      "wunit.AngularVelocity",
		"Area":                 "wunit.Area",
		"Capacitance":          "wunit.Capacitance",
		"Concentration":        "wunit.Concentration",
		"DNASequence":          "wtype.DNASequence",
		"Density":              "wunit.Density",
		"DeviceMetadata":       "api.DeviceMetadata",
		"Energy":               "wunit.Energy",
		"File":                 "wtype.File",
		"FlowRate":             "wunit.FlowRate",
		"Force":                "wunit.Force",
		"HandleOpt":            "execute.HandleOpt",
		"JobID":                "jobfile.JobID",
		"IncubateOpt":          "execute.IncubateOpt",
		"LHComponent":          "wtype.Liquid",
		"LHPlate":              "wtype.LHPlate",
		"LHTip":                "wtype.LHTip",
		"LHTipbox":             "wtype.LHTipbox",
		"LHWell":               "wtype.LHWell",
		"Length":               "wunit.Length",
		"Liquid":               "wtype.Liquid",
		"LiquidType":           "wtype.LiquidType",
		"Mass":                 "wunit.Mass",
		"Moles":                "wunit.Moles",
		"PolicyName":           "wtype.PolicyName",
		"Plate":                "wtype.Plate",
		"Pressure":             "wunit.Pressure",
		"Rate":                 "wunit.Rate",
		"Resistance":           "wunit.Resistance",
		"SpecificHeatCapacity": "wunit.SpecificHeatCapacity",
		"SubstanceQuantity":    "wunit.SubstanceQuantity",
		"Temperature":          "wunit.Temperature",
		"Time":                 "wunit.Time",
		"Velocity":             "wunit.Velocity",
		"Voltage":              "wunit.Voltage",
		"Volume":               "wunit.Volume",
		"Warning":              "wtype.Warning",
	}

	// TODO: add usage tracking to replace UseExpr
	p.addImportReq(&importReq{
		Path: "github.com/antha-lang/antha/laboratory",
	})
	p.addImportReq(&importReq{
		Path:    "github.com/antha-lang/antha/antha/anthalib/wtype",
		UseExpr: "wtype.FALSE",
	})
	p.addImportReq(&importReq{
		Path:    "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/jobfile",
		UseExpr: "jobfile.DefaultClient",
	})
	p.addImportReq(&importReq{
		Path:    "github.com/antha-lang/antha/antha/anthalib/wunit",
		UseExpr: "wunit.GetGlobalUnitRegistry",
	})
	p.addImportReq(&importReq{
		Path:    "github.com/antha-lang/antha/execute",
		UseExpr: "execute.MixInto",
	})
	p.addImportReq(&importReq{
		Path:         "github.com/antha-lang/antha/api/v1",
		Name:         "api",
		UseExpr:      "api.State_CREATED",
		ProtoPackage: "org.antha_lang.antha.v1",
	})

	return p
}

func filterDupSpecs(specs []ast.Spec) []ast.Spec {
	type pair struct {
		name, path string
	}
	seen := make(map[pair]bool)
	var keep []ast.Spec
	for _, spec := range specs {
		ispec, ok := spec.(*ast.ImportSpec)
		if !ok {
			keep = append(keep, spec)
			continue
		}
		key := pair{
			name: ispec.Name.String(),
			path: ispec.Path.Value,
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		keep = append(keep, spec)
	}

	return keep
}

// getImportInsertPos returns position of last import decl or last decl if no
// import decl is present.
func getImportInsertPos(decls []ast.Decl) token.Pos {
	var lastNode ast.Node
	for _, d := range decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			if lastNode == nil {
				lastNode = d
			}
			continue
		}
		lastNode = gd
	}

	if lastNode == nil {
		return token.NoPos
	}
	return lastNode.Pos()
}

// addImports merges multiple import blocks and then adds paths; returns
// merged imports
func (p *Antha) addImports(file *ast.File) {
	var specs []ast.Spec
	var restDecls []ast.Decl
	insertPos := getImportInsertPos(file.Decls)

	for _, d := range file.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			restDecls = append(restDecls, d)
			continue
		}
		specs = append(specs, gd.Specs...)
	}

	for _, req := range p.importReqs {
		if req.Added {
			continue
		}

		imp := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:     token.STRING,
				Value:    strconv.Quote(req.Path),
				ValuePos: insertPos,
			},
		}
		if len(req.Name) != 0 {
			imp.Name = ast.NewIdent(req.Name)
		}
		specs = append(specs, imp)
	}

	if len(specs) == 0 {
		if len(restDecls) != len(file.Decls) {
			// Clean up empty imports
			file.Decls = restDecls
		}
		return
	}

	merged := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: insertPos,
		Rparen: insertPos,
		Specs:  filterDupSpecs(specs),
	}

	file.Decls = append([]ast.Decl{merged}, restDecls...)
}

// getTypeString return appropriate go type string for an antha (type) expr;
// mark used package selectors
func (p *Antha) getTypeString(e ast.Expr, used map[string]bool) string {
	switch t := e.(type) {

	case nil:
		return ""

	case *ast.Ident:
		v, ok := p.types[t.Name]
		if ok {
			return v
		}
		return t.Name

	case *ast.SelectorExpr:
		sel := p.getTypeString(t.X, used)
		if used != nil {
			used[sel] = true
		}
		return sel + "." + t.Sel.Name

	case *ast.BasicLit:
		return t.Value

	case *ast.ArrayType:
		bound := p.getTypeString(t.Len, used)
		return "[" + bound + "]" + p.getTypeString(t.Elt, used)

	case *ast.StarExpr:
		return "*" + p.getTypeString(t.X, used)

	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.getTypeString(t.Key, used), p.getTypeString(t.Value, used))

	default:
		throwErrorf(e.Pos(), "invalid type spec to get type of: %T", t)
	}

	return ""
}

// Transform rewrites AST to go standard primitives
func (p *Antha) Transform(fileSet *token.FileSet, src *ast.File) (err error) {
	if src.Tok != token.PROTOCOL {
		return errNotAnthaFile
	}

	p.protocolName = src.Name.Name

	src.Tok = token.PACKAGE

	p.description = src.Doc.Text()

	file := fileSet.File(src.Package)
	p.elementPath = file.Name()
	p.root.addProtocolDirectory(p.protocolName, filepath.Dir(p.elementPath))

	// Case-insensitive comparison because some filesystems are
	// case-insensitive
	packagePath := filepath.Base(filepath.Dir(p.elementPath))
	if e, f := p.protocolName, packagePath; strings.ToLower(e) != strings.ToLower(f) {
		return fmt.Errorf("%s: expecting protocol %s to be in directory %s", file.Name(), e, f)
	}

	p.recordImports(src.Decls)
	p.recordMessages(src.Decls)
	if err := p.validateMessages(p.messages); err != nil {
		return err
	}

	p.desugar(fileSet, src)
	p.addImports(src)
	p.addUses(src)

	return
}

func (p *Antha) addImportReq(req *importReq) {
	name := req.Name
	if len(name) == 0 {
		name = path.Base(req.Path)
	}

	p.importReqs = append(p.importReqs, req)
	p.importByName[name] = req
}

// Return name relative to a base if possible
func relativeTo(bases []string, name string) (string, error) {
	absName, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}

	var prefixes []string
	prefixes = append(prefixes, bases...)

	// In reverse alphabetical to ensure longest match first
	sort.Strings(prefixes)
	for idx := len(prefixes) - 1; idx >= 0; idx-- {
		p := prefixes[idx]
		if !strings.HasPrefix(absName, p) {
			continue
		} else if rp, err := filepath.Rel(p, absName); err != nil {
			return "", err
		} else {
			return rp, nil
		}
	}

	return name, nil
}

func (p *Antha) recordImports(decls []ast.Decl) {
	for _, decl := range decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}

		for _, spec := range gd.Specs {
			im := spec.(*ast.ImportSpec)

			value, _ := strconv.Unquote(im.Path.Value)
			req := &importReq{
				Added: true,
				Path:  value,
			}
			if im.Name != nil {
				req.Name = im.Name.String()
			}

			p.addImportReq(req)
		}
	}
}

// recordMessages records all the spec definitions for inputs and outputs to element
func (p *Antha) recordMessages(decls []ast.Decl) {
	join := func(xs ...string) string {
		var ret []string
		for _, x := range xs {
			if len(x) == 0 {
				continue
			}
			ret = append(ret, x)
		}
		return strings.Join(ret, "\n")
	}

	for _, decl := range decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if !isAnthaGenDeclToken(decl.Tok) {
			continue
		}

		for _, spec := range decl.Specs {
			spec, ok := spec.(*ast.TypeSpec)
			if !ok {
				throwErrorf(spec.Pos(), "expecting type")
			}
			typ, ok := spec.Type.(*ast.StructType)
			if !ok {
				throwErrorf(spec.Pos(), "expecting struct type")
			}

			var fields []*Field
			for _, field := range typ.Fields.List {
				for _, name := range field.Names {
					f := &Field{
						Name: name.String(),
						Type: p.desugarTypeExpr(field.Type),
						Doc:  join(field.Comment.Text(), field.Doc.Text()),
					}
					if field.Tag != nil {
						f.Tag = field.Tag.Value
					}
					fields = append(fields, f)
				}
			}

			p.messages = append(p.messages, &Message{
				Name:   spec.Name.String(),
				Fields: fields,
				Doc:    join(decl.Doc.Text(), spec.Comment.Text(), spec.Doc.Text()),
				Kind:   decl.Tok,
			})
		}
	}
}

func uniqueFields(fields []*Field) error {
	seen := make(map[string]bool)
	for _, f := range fields {
		if seen[f.Name] {
			return fmt.Errorf("%s already declared", f.Name)
		}
		seen[f.Name] = true
	}

	return nil
}

func (p *Antha) validateMessages(messages []*Message) error {
	p.tokenByParamName = make(map[string]token.Token)

	seen := make(map[string]*Message)

	for _, msg := range messages {
		name := msg.Name

		for _, field := range msg.Fields {
			if _, seen := p.tokenByParamName[field.Name]; seen {
				return fmt.Errorf("%s already declared", name)
			}
			p.tokenByParamName[field.Name] = msg.Kind

			if !ast.IsExported(field.Name) {
				return fmt.Errorf("field %s must begin with an upper case letter", name)
			}
		}

		if _, seen := seen[name]; seen {
			return fmt.Errorf("%s already declared", name)
		}

		seen[name] = msg

		if err := uniqueFields(msg.Fields); err != nil {
			return err
		}
	}

	return nil
}

func (p *Antha) generateElement(fileSet *token.FileSet, file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	compiler := &Config{
		Mode:     printerMode,
		Tabwidth: tabWidth,
	}
	lineMap, err := compiler.Fprint(&buf, fileSet, file)
	if err != nil {
		return nil, err
	}

	if err := p.printFunctions(&buf, lineMap); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Generate returns files with slash names to complete antha to go
// transformation
func (p *Antha) Generate(fileSet *token.FileSet, file *ast.File) (*AnthaFiles, error) {
	elementBs, err := p.generateElement(fileSet, file)
	if err != nil {
		return nil, err
	}

	elementName := path.Join(p.protocolName, elementFilename)

	files := NewAnthaFiles()
	files.addFile(elementName, elementBs)

	return files, nil
}

func (p *Antha) addUses(src *ast.File) {
	for _, req := range p.importReqs {
		if len(req.UseExpr) == 0 {
			continue
		}
		decl := &ast.GenDecl{
			Tok: token.VAR,
		}

		decl.Specs = append(decl.Specs, &ast.ValueSpec{
			Names:  identList("_"),
			Values: []ast.Expr{mustParseExpr(req.UseExpr)},
		})
		src.Decls = append(src.Decls, decl)
	}
}

func encodeByteArray(bs []byte) string {
	var buf bytes.Buffer
	buf.WriteString("[]byte {\n")
	for len(bs) > 0 {
		n := 16
		if n > len(bs) {
			n = len(bs)
		}

		for _, c := range bs[:n] {
			buf.WriteString("0x")
			buf.WriteString(hex.EncodeToString([]byte{c}))
			buf.WriteString(",")
		}

		buf.WriteString("\n")

		bs = bs[n:]
	}

	buf.WriteString("}")

	return buf.String()
}

// printFunctions generates synthetic antha functions and data stuctures
func (p *Antha) printFunctions(out io.Writer, lineMap map[int]int) error {
	var tmpl = `
type {{.ElementTypeName}} struct {
	Inputs     Inputs
	Outputs    Outputs
	Parameters Parameters
	Data       Data
}

func New{{.ElementTypeName}}(lab *laboratory.Laboratory) *{{.ElementTypeName}} {
	element := &{{.ElementTypeName}}{}
	lab.InstallElement(element)
	return element
}

var LineMap = map[int]int{
	{{range $key, $value := .LineMap}}{{$key}}: {{$value}}, {{end}}
}
`
	type TVars struct {
		ElementTypeName string
		LineMap         map[int]int
	}

	tv := TVars{
		ElementTypeName: p.protocolName,
		LineMap:         lineMap,
	}

	return template.Must(template.New("").Parse(tmpl)).Execute(out, tv)
}

// desugar updates AST for antha semantics
func (p *Antha) desugar(fileSet *token.FileSet, src *ast.File) {
	for idx, d := range src.Decls {
		switch d := d.(type) {

		case *ast.GenDecl:
			ast.Inspect(d, p.inspectTypes)
			p.desugarGenDecl(d)

		case *ast.AnthaDecl:
			ast.Inspect(d.Body, p.inspectIntrinsics)
			ast.Inspect(d.Body, p.inspectParamUses)
			ast.Inspect(d.Body, p.inspectTypes)
			src.Decls[idx] = p.desugarAnthaDecl(fileSet, src, d)

		default:
			ast.Inspect(d, p.inspectTypes)
		}
	}
}

func identList(name string) []*ast.Ident {
	return []*ast.Ident{ast.NewIdent(name)}
}

func mustParseExpr(x string) ast.Expr {
	r, err := parser.ParseExpr(x)
	if err != nil {
		panic(fmt.Errorf("cannot parse %s: %s", x, err))
	}
	return r
}

// desugarGenDecl returns standard go ast for antha GenDecls
func (p *Antha) desugarGenDecl(d *ast.GenDecl) {
	if !isAnthaGenDeclToken(d.Tok) {
		return
	}

	d.Tok = token.TYPE
}

// desugarAnthaDecl returns standard go ast for antha decl.
//
// E.g.,
//   Validation
// to
//  func (element *ElementTypeName) Validation(lab *laboratory.Laboratory)
func (p *Antha) desugarAnthaDecl(fileSet *token.FileSet, src *ast.File, d *ast.AnthaDecl) ast.Decl {
	f := &ast.FuncDecl{
		Doc: d.Doc,
		Recv: &ast.FieldList{
			Opening: d.Pos(),
			List: []*ast.Field{
				{
					Names: identList("element"),
					Type:  mustParseExpr("*" + p.protocolName),
				},
			},
		},
		Name: ast.NewIdent(d.Tok.String()),
		Body: d.Body,
	}

	f.Type = &ast.FuncType{
		Func: d.Pos(),
		Params: &ast.FieldList{
			Opening: d.Pos(),
			List: []*ast.Field{
				{
					Names: identList("lab"),
					Type:  mustParseExpr("*laboratory.Laboratory"),
				},
			},
		},
	}

	return f
}

// desugarTypeIdent returns appropriate nested SelectorExpr for the replacement for
// Identifier
func (p *Antha) desugarTypeIdent(t *ast.Ident) ast.Expr {
	v, ok := p.types[t.Name]
	if !ok {
		return t
	}

	return mustParseExpr(v)
}

// desugarTypeExpr returns appropriate go type for an antha (type) expr
func (p *Antha) desugarTypeExpr(t ast.Node) ast.Expr {
	switch t := t.(type) {
	case nil:
		return nil

	case *ast.Ident:
		return p.desugarTypeIdent(t)

	case *ast.ParenExpr:
		t.X = p.desugarTypeExpr(t.X)

	case *ast.SelectorExpr:

	case *ast.StarExpr:
		t.X = p.desugarTypeExpr(t.X)

	case *ast.ArrayType:
		t.Elt = p.desugarTypeExpr(t.Elt)

	case *ast.StructType:
		ast.Inspect(t, p.inspectTypes)

	case *ast.FuncType:
		ast.Inspect(t, p.inspectTypes)

	case *ast.InterfaceType:
		ast.Inspect(t, p.inspectTypes)

	case *ast.MapType:
		t.Key = p.desugarTypeExpr(t.Key)
		t.Value = p.desugarTypeExpr(t.Value)

	case *ast.ChanType:
		t.Value = p.desugarTypeExpr(t.Value)

	case *ast.Ellipsis:

	default:
		throwErrorf(t.Pos(), "unexpected expression %s of type %T", t, t)
	}

	return t.(ast.Expr)
}

func inspectExprList(exprs []ast.Expr, w func(ast.Node) bool) {
	for _, expr := range exprs {
		ast.Inspect(expr, w)
	}
}

// inspectTypes replaces bare antha types with go qualified names.
//
// Changing all idents blindly would be simpler but opt instead with only
// replacing idents that appear in types.
func (p *Antha) inspectTypes(n ast.Node) bool {
	switch n := n.(type) {
	case nil:

	case *ast.Field:
		n.Type = p.desugarTypeExpr(n.Type)

	case *ast.TypeSpec:
		n.Type = p.desugarTypeExpr(n.Type)

	case *ast.MapType:
		n.Key = p.desugarTypeExpr(n.Key)
		n.Value = p.desugarTypeExpr(n.Value)

	case *ast.ArrayType:
		n.Elt = p.desugarTypeExpr(n.Elt)

	case *ast.ChanType:
		n.Value = p.desugarTypeExpr(n.Value)

	case *ast.FuncLit:
		n.Type = p.desugarTypeExpr(n.Type).(*ast.FuncType)
		ast.Inspect(n.Body, p.inspectTypes)

	case *ast.CompositeLit:
		n.Type = p.desugarTypeExpr(n.Type)
		inspectExprList(n.Elts, p.inspectTypes)

	case *ast.TypeAssertExpr:
		n.Type = p.desugarTypeExpr(n.Type)

	case *ast.ValueSpec:
		n.Type = p.desugarTypeExpr(n.Type)
		inspectExprList(n.Values, p.inspectTypes)

	default:
		return true
	}

	return false
}

// inspectParamUses replaces bare antha identifiers with go qualified names
func (p *Antha) inspectParamUses(node ast.Node) bool {
	// desugar if it is a known param
	rewriteIdent := func(node *ast.Ident) {
		tok, ok := p.tokenByParamName[node.Name]
		if !ok {
			return
		}

		node.Name = "element." + tok.String() + "." + node.Name
	}

	rewriteAssignLHS := func(node *ast.AssignStmt) {
		for _, lhs := range node.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}

			tok, found := p.tokenByParamName[ident.Name]
			if !found || !isOutput(tok) {
				continue
			}

			ident.Name = "element." + tok.String() + "." + ident.Name
		}
	}

	switch n := node.(type) {

	case nil:
		return false

	case *ast.AssignStmt:
		rewriteAssignLHS(n)

	case *ast.KeyValueExpr:
		if _, identKey := n.Key.(*ast.Ident); identKey {
			// Skip identifiers that are keys
			ast.Inspect(n.Value, p.inspectParamUses)
			return false
		}
	case *ast.Ident:
		rewriteIdent(n)

	case *ast.SelectorExpr:
		// Skip identifiers that are field accesses
		ast.Inspect(n.X, p.inspectParamUses)
		return false
	}
	return true
}

// inspectIntrinsics replaces bare antha function names with go qualified
// names
func (p *Antha) inspectIntrinsics(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.CallExpr:
		ident, direct := n.Fun.(*ast.Ident)
		if !direct {
			break
		}

		if desugar, ok := p.intrinsics[ident.Name]; ok {
			ident.Name = desugar
			n.Args = append([]ast.Expr{ast.NewIdent("lab")}, n.Args...) // only for now.
		}
	}
	return true
}
