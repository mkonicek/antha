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
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

const (
	runStepsIntrinsic   = "RunSteps"
	awaitDataIntrinsic  = "AwaitData"
	lineNumberConstName = "_lineNumber"
	elementProto        = "element.proto"
	elementPackage      = "element"
	elementFilename     = "element.go"
	modelPackage        = "model"
	modelFilename       = "model.go"
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
	Kind   token.Token // One of token.{DATA, PARAMETERS, OUTPUTS, INPUTS, MESSAGE}
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
	case token.OUTPUTS, token.DATA, token.PARAMETERS, token.INPUTS, token.MESSAGE:
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
	SourceSHA256 []byte

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
	importProtos []string
}

// NewAntha creates a new antha pass
func NewAntha(root *AnthaRoot) *Antha {
	p := &Antha{
		root:         root,
		importByName: make(map[string]*importReq),
	}

	p.importProtos = []string{
		"github.com/antha-lang/antha/api/v1/blob.proto",
		"github.com/antha-lang/antha/api/v1/coord.proto",
		"github.com/antha-lang/antha/api/v1/element.proto",
		"github.com/antha-lang/antha/api/v1/empty.proto",
		"github.com/antha-lang/antha/api/v1/inventory.proto",
		"github.com/antha-lang/antha/api/v1/measurement.proto",
		"github.com/antha-lang/antha/api/v1/message.proto",
		"github.com/antha-lang/antha/api/v1/polynomial.proto",
		"github.com/antha-lang/antha/api/v1/state.proto",
		"github.com/antha-lang/antha/api/v1/task.proto",
		"github.com/antha-lang/antha/api/v1/workflow.proto",
	}

	p.intrinsics = map[string]string{
		"Centrifuge":    "execute.Centrifuge",
		"Electroshock":  "execute.Electroshock",
		"Errorf":        "execute.Errorf",
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
		Path: "context",
	})
	p.addImportReq(&importReq{
		Path: "encoding/json",
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
		UseExpr: "wunit.Make_units",
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
	p.addImportReq(&importReq{
		Path: "github.com/antha-lang/antha/component",
	})
	p.addImportReq(&importReq{
		Path: "github.com/antha-lang/antha/inject",
	})

	p.addImportReq(&importReq{
		Path:    "encoding/json",
		UseExpr: "json.Unmarshal",
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

func (p *Antha) addExternalPackage(pkgPath string) *importReq {
	req := &importReq{
		Path: pkgPath,
		Name: manglePackageName(pkgPath),
	}
	p.importReqs = append(p.importReqs, req)
	return req
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

func unwrapRecover(res interface{}, fileSet *token.FileSet, pos token.Pos) error {
	perr, ok := res.(posError)
	msg := res.(error).Error()
	if ok {
		pos = perr.Pos()
	}
	p := fileSet.Position(pos)

	if ok {
		return fmt.Errorf("%s:%d: %s", p.Filename, p.Line, msg)
	}
	return fmt.Errorf("%s: %s", p.Filename, msg)
}

// normalizePath takes a filepath and returns a slash path relative to GOPATH
func normalizePath(filename string) string {
	filename, _ = relativeTo(getGoPath(), filename)
	dir := filepath.ToSlash(filename)
	if strings.HasPrefix(dir, "src/") {
		return dir[len("src/"):]
	}
	return dir
}

func reverse(xs []string) (ret []string) {
	for idx := len(xs) - 1; idx >= 0; idx-- {
		ret = append(ret, xs[idx])
	}
	return
}

func getProtobufPackage(goPackage string) string {
	var parts []string
	idx := strings.Index(goPackage, "/")
	if idx >= 0 {
		parts = append(parts, reverse(strings.Split(goPackage[:idx], "."))...)
		goPackage = goPackage[idx+1:]
	}
	parts = append(parts, strings.Split(goPackage, "/")...)

	var ret []string
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		r := strings.Map(func(r rune) rune {
			switch r {
			case '-', '.':
				return '_'
			default:
				return r
			}
		}, part)
		ret = append(ret, r)
	}

	return strings.Join(ret, ".")
}

// Transform rewrites AST to go standard primitives
func (p *Antha) Transform(fileSet *token.FileSet, src *ast.File) (err error) {
	defer func() {
		if res := recover(); res != nil {
			err = unwrapRecover(res, fileSet, src.Package)
		}
	}()

	if src.Tok != token.PROTOCOL {
		return errNotAnthaFile
	}

	p.protocolName = src.Name.Name
	src.Name.Name = elementPackage

	src.Tok = token.PACKAGE

	p.description = src.Doc.Text()

	file := fileSet.File(src.Package)
	p.elementPath = file.Name()
	p.root.addProtocolDirectory(p.protocolName, filepath.Dir(p.elementPath))

	p.addImportReq(&importReq{
		Path: path.Join(p.root.outputPackageBase, p.protocolName, modelPackage),
	})

	// Case-insensitive comparison because some filesystems are
	// case-insensitive
	packagePath := filepath.Base(filepath.Dir(p.elementPath))
	if e, f := p.protocolName, packagePath; strings.ToLower(e) != strings.ToLower(f) {
		return fmt.Errorf("%s: expecting protocol %s to be in directory %s", file.Name(), e, f)
	}

	p.recordImports(src.Decls)
	p.recordBlocks(src.Decls)
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

// Usually $GOPATH but if not set, future versions of go will assume $HOME/go
func getGoPath() []string {
	ps := filepath.SplitList(os.Getenv("GOPATH"))
	if len(ps) == 0 {
		usr, err := user.Current()
		if err == nil {
			ps = append(ps, filepath.Join(usr.HomeDir, "go"))
		}
	}

	return ps
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

// recordBlocks records all blocks used
func (p *Antha) recordBlocks(decls []ast.Decl) {
	p.blocksUsed = make(map[token.Token]bool)

	for _, decl := range decls {
		decl, ok := decl.(*ast.AnthaDecl)
		if !ok {
			continue
		}
		p.blocksUsed[decl.Tok] = true
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

		switch k := msg.Kind; k {

		case token.MESSAGE:
			if !ast.IsExported(name) {
				return fmt.Errorf("%s %s must begin with an upper case letter", k, name)
			}

		default:
			for _, field := range msg.Fields {

				if _, seen := p.tokenByParamName[field.Name]; seen {
					return fmt.Errorf("%s already declared", name)
				}
				p.tokenByParamName[field.Name] = msg.Kind
			}
		}

		for _, field := range msg.Fields {
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
	if err := compiler.Fprint(&buf, fileSet, file); err != nil {
		return nil, err
	}

	pat := regexp.MustCompile(fmt.Sprintf(`const %s = "([^"\n\r]+)"`, lineNumberConstName))
	main := pat.ReplaceAll(buf.Bytes(), []byte(`//line $1`))
	var out bytes.Buffer
	if _, err := io.Copy(&out, bytes.NewReader(main)); err != nil {
		return nil, err
	}

	if err := p.printFunctions(&out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// generateModel generates json structures that implement parsed messages
func (p *Antha) generateModel() ([]byte, error) {
	const tmpl = `package model 

import (
{{ range .Imports}}{{.Name}} {{.Value}}
{{ end }}
)

type Input struct {
	{{ range .Inputs }}{{ .Name }} {{ .Value }}
	{{ end }}
}

type Output struct {
	{{ range .MergedOutputs }}{{ .Name }} {{ .Value }}
	{{ end }}
}

type RunStepsOutput struct {
	Data struct {
		{{ range .Data }}{{ .Name }} {{ .Value }}
		{{ end }}
	}
	Outputs struct {
		{{ range .Outputs }}{{ .Name }} {{ .Value }}
		{{ end }}
	}
}
`
	type Field struct {
		Name  string
		Value string
	}

	usedSelectors := make(map[string]bool)

	makeFields := func(tok token.Token) (ret []Field) {
		for _, msg := range p.getMessage(tok).getFields() {
			ret = append(ret, Field{
				Name:  msg.Name,
				Value: p.getTypeString(msg.Type, usedSelectors),
			})
		}
		return
	}

	type TVars struct {
		Imports       []Field
		Inputs        []Field
		Outputs       []Field
		Data          []Field
		MergedOutputs []Field
	}

	tv := TVars{}

	tv.Inputs = append(tv.Inputs, makeFields(token.INPUTS)...)
	tv.Inputs = append(tv.Inputs, makeFields(token.PARAMETERS)...)

	tv.Outputs = makeFields(token.OUTPUTS)
	tv.Data = makeFields(token.DATA)
	tv.MergedOutputs = append(tv.MergedOutputs, tv.Outputs...)
	tv.MergedOutputs = append(tv.MergedOutputs, tv.Data...)

	seen := make(map[string]bool)
	for _, req := range p.importReqs {
		iname := req.ImportName()
		if !usedSelectors[iname] {
			continue
		}

		if seen[iname] {
			continue
		}

		seen[iname] = true

		tv.Imports = append(tv.Imports, Field{
			Name:  req.Name,
			Value: strconv.Quote(req.Path),
		})
	}

	var out bytes.Buffer
	if err := template.Must(template.New("").Parse(tmpl)).Execute(&out, tv); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (p *Antha) getProtobufType(e ast.Expr, used map[*importReq]bool) string {
	switch t := e.(type) {

	case nil:
		return ""

	case *ast.Ident:
		switch t.Name {
		case "int", "int32":
			return "int32"
		case "int64":
			return "int64"
		case "float":
			return "float"
		case "float64":
			return "double"
		}
		return t.Name

	case *ast.SelectorExpr:
		if lhs, ok := t.X.(*ast.Ident); ok {
			importReq, ok := p.importByName[lhs.Name]
			if !ok {
				throwErrorf(e.Pos(), "unknown package in protobuf: %s", lhs.Name)
			}
			pkg := importReq.ProtoPackage
			if len(pkg) == 0 {
				pkg = getProtobufPackage(importReq.Path)
			}
			used[importReq] = true
			return pkg + "." + t.Sel.Name
		}

	case *ast.BasicLit:
		return t.Value

	case *ast.ArrayType:
		return fmt.Sprintf("repeated %s", p.getProtobufType(t.Elt, used))

	case *ast.StarExpr:
		return p.getProtobufType(t.X, used)

	case *ast.MapType:
		return fmt.Sprintf("map<%s,%s>", p.getProtobufType(t.Key, used), p.getProtobufType(t.Value, used))

	}

	throwErrorf(e.Pos(), "invalid type in protobuf: %T", e)

	return ""
}

func (p *Antha) getMessage(tok token.Token) *Message {
	for _, msg := range p.messages {
		if msg.Name == tok.String() {
			return msg
		}
	}
	return nil
}

// generateProto generates protobuf file that implements parsed
// messages
func (p *Antha) generateProto() ([]byte, error) {
	const tmpl = `syntax = "proto3";
package {{ .PackageName }};
option go_package = "protobuf";

{{ range .Imports }}import "{{ . }}";
{{ end }}
{{ range .Messages }}{{ if .Desc }}{{ .Desc }}{{ end }}
message {{ .Name }} {
{{ range .Fields }}{{ if .Desc }}{{ .Desc }}{{ end }}
  {{ .Type }} {{ .Name }} = {{ .Tag }};
{{ end }}
}
{{ end }}
service Element {
  rpc Run(Request) returns (Response);
  rpc Run(org.antha_lang.antha.v1.Empty) returns (org.antha_lang.antha.v1.ElementMetadata);
}
`
	fmtDoc := func(indent, s string) string {
		if len(s) == 0 {
			return s
		}

		comment := indent + "// "
		s = strings.TrimSpace(s)
		s = strings.Replace(s, "\n", "\n"+comment, -1)
		return comment + s
	}

	merge := func(messages []*Message, merged string, aTok, bTok token.Token) (ret []*Message) {
		a := p.getMessage(aTok)
		b := p.getMessage(bTok)

		for _, msg := range messages {
			if msg == a || msg == b {
				continue
			}
			ret = append(ret, msg)
		}

		newMsg := &Message{
			Name: merged,
		}

		if a != nil {
			newMsg.Doc += a.Doc
			newMsg.Fields = append(newMsg.Fields, a.Fields...)
		}

		if b != nil {
			newMsg.Doc += b.Doc
			newMsg.Fields = append(newMsg.Fields, b.Fields...)
		}

		ret = append(ret, newMsg)

		return
	}

	type Field struct {
		Desc string
		Type string
		Name string
		Tag  int
	}

	type Message struct {
		Name   string
		Desc   string
		Fields []Field
	}

	type TVars struct {
		PackageName string
		Imports     []string
		Messages    []Message
	}

	tv := TVars{
		PackageName: getProtobufPackage(normalizePath(filepath.Dir(p.elementPath))),
	}

	messages := p.messages
	messages = merge(messages, "Request", token.PARAMETERS, token.INPUTS)
	messages = merge(messages, "Response", token.DATA, token.OUTPUTS)

	used := make(map[*importReq]bool)
	for _, msg := range messages {
		m := Message{
			Name: msg.Name,
			Desc: fmtDoc("", msg.Doc),
		}
		for idx, f := range msg.Fields {
			m.Fields = append(m.Fields, Field{
				Desc: fmtDoc("  ", f.Doc),
				Type: p.getProtobufType(f.Type, used),
				Name: f.Name,
				Tag:  idx + 1,
			})
		}
		tv.Messages = append(tv.Messages, m)
	}

	tv.Imports = append(tv.Imports, p.importProtos...)
	for req := range used {
		// TODO: we assume that only system packages have this field set
		if len(req.ProtoPackage) != 0 {
			continue
		}

		tv.Imports = append(tv.Imports, path.Join(req.Path, elementProto))
	}

	sort.Strings(tv.Imports)

	var out bytes.Buffer
	if err := template.Must(template.New("").Parse(tmpl)).Execute(&out, tv); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// Generate returns files with slash names to complete antha to go
// transformation
func (p *Antha) Generate(fileSet *token.FileSet, file *ast.File) (*AnthaFiles, error) {
	elementBs, err := p.generateElement(fileSet, file)
	if err != nil {
		return nil, err
	}

	var modelBs []byte

	if true {
		modelBs, err = p.generateModel()
	} else {
		modelBs, err = p.generateProto()
	}

	if err != nil {
		return nil, err
	}

	modelName := path.Join(p.protocolName, modelPackage, modelFilename)
	elementName := path.Join(p.protocolName, elementPackage, elementFilename)

	files := NewAnthaFiles()
	files.addFile(modelName, modelBs)
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
func (p *Antha) printFunctions(out io.Writer) error {
	// TODO: put the recover handler here when we get to multi-address space
	// execution. In single-address space, the caller assumes that execeptions
	// will bubble up through Run.

	// NB: serialize in Run to enforce serialization barrier between element
	// calls
	var tmpl = `
type Element struct {
}

func (Element) Run(_ctx context.Context, request *{{ .ModelPackage }}.Input) (response *{{ .ModelPackage }}.Output, err error) {
	_ctx = execute.WithElementName(_ctx, {{ .ElementName }})
	bs, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var in *{{ .ModelPackage }}.Input
	if err := json.Unmarshal(bs, &in); err != nil {
		return nil, err
	}

	response = &{{ .ModelPackage }}.Output{}
	{{if .HasSetup}}_Setup(_ctx, in, response){{end}}
	{{if .HasSteps}}_Steps(_ctx, in, response){{end}}
	{{if .HasAnalysis}}_Analysis(_ctx, in, response){{end}}
	{{if .HasValidation}}_Validation(_ctx, in, response){{end}}
	return
}

func (Element) RunAnalysisValidation(_ctx context.Context, request *{{ .ModelPackage }}.Input) (response *{{ .ModelPackage }}.Output, err error) {
	response = &{{ .ModelPackage }}.Output{}
	{{if .HasAnalysis}}_Analysis(_ctx, request, response){{end}}
	{{if .HasValidation}}_Validation(_ctx, request, response){{end}}
	return
}

func (Element) Metadata(_ctx context.Context, request *api.Empty) (*api.ElementMetadata, error) {
	return _metadata, nil
}

func _newAVRunner() interface{} {
	elem := &Element{}
	return &inject.CheckedRunner {
		RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
			request := &{{ .ModelPackage }}.Input{}
			if err := inject.Assign(value, request); err != nil {
				return nil, err
			}
			resp, err := elem.RunAnalysisValidation(_ctx, request)
			if err != nil {
				return nil, err
			}
			return inject.MakeValue(resp), nil
		},
		In: &{{ .ModelPackage }}.Input{},
		Out: &{{ .ModelPackage }}.Output{},
	}
}
func _newRunner() interface{} {
	elem := &Element{}
	return &inject.CheckedRunner {
		RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
			request := &{{ .ModelPackage }}.Input{}
			if err := inject.Assign(value, request); err != nil {
				return nil, err
			}
			resp, err := elem.Run(_ctx, request)
			if err != nil {
				return nil, err
			}
			return inject.MakeValue(resp), nil
		},
		In: &{{ .ModelPackage }}.Input{},
		Out: &{{ .ModelPackage }}.Output{},
	}
}

func RunSteps(_ctx context.Context, request *{{ .ModelPackage }}.Input) (*{{ .ModelPackage }}.RunStepsOutput) {
	elem := &Element{}
	resp, err := elem.Run(_ctx, request)
	if err != nil {
		panic(err)
	}

	var output {{ .ModelPackage }}.RunStepsOutput
	if err := inject.AssignSome(resp, &output.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(resp, &output.Outputs); err != nil {
		panic(err)
	}
	return &output
}

func GetComponent() []*component.Component{
	return []*component.Component{ 
			&component.Component{
			Name: {{ .ElementName }},
			Stage: api.ElementStage_STEPS,
			Constructor: _newRunner,
			Description: component.Description{
				Desc: {{ .Desc }},
				Path: {{ .Path }},
				Params: []component.ParamDesc{
					{{range .Params}}component.ParamDesc{
						Name: {{ .Name }},
						Desc: {{ .Desc }},
						Kind: {{ .Kind }},
					},
					{{end}}
				},
			},
		},
			&component.Component{
			Name: {{ .ElementName }},
			Stage: api.ElementStage_ANALYSIS,
			Constructor: _newAVRunner,
			Description: component.Description{
				Desc: {{ .Desc }},
				Path: {{ .Path }},
			},
		},
	}
}

var (
	_metadata *api.ElementMetadata
)

func init() {
	_metadata = &api.ElementMetadata {
		SourceSha256: {{.SHA256}},
	}
}
`
	type Param struct {
		Name     string
		BareName string
		Desc     string
		Kind     string
		BareKind string
	}

	type TVars struct {
		ModelPackage  string
		ElementName   string
		SHA256        string
		Desc          string
		Path          string
		Params        []Param
		HasSteps      bool
		HasValidation bool
		HasSetup      bool
		HasAnalysis   bool
	}

	elementPath := normalizePath(p.elementPath)

	tv := TVars{
		ModelPackage:  modelPackage,
		ElementName:   strconv.Quote(p.protocolName),
		SHA256:        encodeByteArray(p.SourceSHA256),
		Desc:          strconv.Quote(p.description),
		Path:          strconv.Quote(elementPath),
		HasSteps:      p.blocksUsed[token.STEPS],
		HasValidation: p.blocksUsed[token.VALIDATION],
		HasSetup:      p.blocksUsed[token.SETUP],
		HasAnalysis:   p.blocksUsed[token.ANALYSIS],
	}

	for _, msg := range p.messages {
		for _, field := range msg.Fields {
			tv.Params = append(tv.Params, Param{
				Name:     strconv.Quote(field.Name),
				BareName: field.Name,
				Desc:     strconv.Quote(field.Doc),
				Kind:     strconv.Quote(msg.Name),
				BareKind: msg.Name,
			})
		}
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
//   _Validation(_ctx context.Context, _input *protobuf.Request, _output *protobuf.Response)
func (p *Antha) desugarAnthaDecl(fileSet *token.FileSet, src *ast.File, d *ast.AnthaDecl) ast.Decl {
	f := &ast.FuncDecl{
		Doc:  d.Doc,
		Name: ast.NewIdent("_" + d.Tok.String()),
		Body: d.Body,
	}

	f.Type = &ast.FuncType{
		Func: d.Pos(),
		Params: &ast.FieldList{
			Opening: d.Pos(),
			List: []*ast.Field{
				{
					Names: identList("_ctx"),
					Type:  mustParseExpr("context.Context"),
				},
				{
					Names: identList("_input"),
					Type:  mustParseExpr("*" + modelPackage + ".Input"),
				},
				{
					Names: identList("_output"),
					Type:  mustParseExpr("*" + modelPackage + ".Output"),
				},
			},
		},
	}

	// HACK: all the ast rewriting invalidates positions, so we insert a dummy
	// decl to hang comment on and turn it into a comment in the adjustment
	// function.
	if len(d.Body.List) > 0 {
		pos := fileSet.Position(d.Body.Lbrace)

		line := fmt.Sprintf("%s:%d", pos.Filename, pos.Line)

		dummy := &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names:  identList(lineNumberConstName),
						Values: []ast.Expr{mustParseExpr(strconv.Quote(line))},
					},
				},
			},
		}
		d.Body.List = append([]ast.Stmt{dummy}, d.Body.List...)
	}

	// HACK: unanchored comments can interrupt regexp replacement of HACK nodes
	// above, remove all unanchored comments to fix.
	src.Comments = nil

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

		if isOutput(tok) {
			node.Name = "_output." + node.Name
		} else if isInput(tok) {
			node.Name = "_input." + node.Name
		}
	}

	rewriteAssignLHS := func(node *ast.AssignStmt) {
		for _, lhs := range node.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}

			param, ok := p.tokenByParamName[ident.Name]
			if !ok || !isOutput(param) {
				continue
			}

			ident.Name = "_output." + ident.Name
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

type awaitDataPrototype struct {
	Annotatee     *ast.Ident
	Metadata      ast.Expr
	NextElement   *ast.Ident
	ReplacedParam *ast.Ident
	Params        *ast.CompositeLit
	Inputs        *ast.CompositeLit
}

func parseAwaitData(call *ast.CallExpr) *awaitDataPrototype {
	if len(call.Args) != 6 {
		return nil
	}

	annotatee, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return nil
	}

	metadata := call.Args[1]

	nextElement, ok := call.Args[2].(*ast.Ident)
	if !ok {
		return nil
	}

	replacedParam, ok := call.Args[3].(*ast.Ident)
	if !ok {
		return nil
	}

	params, ok := call.Args[4].(*ast.CompositeLit)
	if !ok {
		return nil
	}

	inputs, ok := call.Args[5].(*ast.CompositeLit)
	if !ok {
		return nil
	}

	return &awaitDataPrototype{
		Annotatee:     annotatee,
		Metadata:      metadata,
		NextElement:   nextElement,
		ReplacedParam: replacedParam,
		Params:        params,
		Inputs:        inputs,
	}
}

// rewriteAwaitData transforms
//   AwaitData(
//   	annotatee,
//      metadata,
// 		nextElement,
//      replacedParam,
// 		nextParameters,
// 		nextInput)
// to
//   execute.AwaitData(
//		_ctx,
//		annotatee,
// 		metadata,
//		nextElement,
//      replacedParam,
// 		inject.MakeValue(nextParams + nextInput),
//      inject.MakeValue(_output))
func (p *Antha) rewriteAwaitData(call *ast.CallExpr) {
	proto := parseAwaitData(call)
	if proto == nil {
		var p awaitDataPrototype
		throwErrorf(call.Pos(),
			"expecting %s(%s) found %s(%s)",
			awaitDataIntrinsic,
			strings.Join(typesToString(p.Annotatee, p.Metadata, p.NextElement, p.ReplacedParam, p.Params, p.Inputs), ","),
			awaitDataIntrinsic,
			strings.Join(typesToString(call.Args...), ","),
		)
	}

	nextElement := ""
	nextElementArgs := &ast.CompositeLit{Type: mustParseExpr("model.Input"), Elts: []ast.Expr{}}

	// indirect over next element stuff
	if proto.NextElement.Name != "nil" {
		modelPkg := path.Join(p.root.outputPackageBase, proto.NextElement.Name, modelPackage)
		modelReq := p.addExternalPackage(modelPkg)
		nextElement = proto.NextElement.Name
		nextElementArgs = &ast.CompositeLit{
			Type: mustParseExpr(modelReq.Name + ".Input"),
			Elts: append(proto.Params.Elts, proto.Inputs.Elts...),
		}
	}

	call.Fun = mustParseExpr("execute." + awaitDataIntrinsic)

	call.Args = []ast.Expr{
		ast.NewIdent("_ctx"),
		proto.Annotatee,
		proto.Metadata,
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(nextElement),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(proto.ReplacedParam.Name),
		},
		mustParseExpr("inject.MakeValue(_output)"),
		&ast.CallExpr{
			Fun: mustParseExpr("inject.MakeValue"),
			Args: []ast.Expr{
				nextElementArgs,
			},
		},
	}
}

type runStepsPrototype struct {
	Callee *ast.Ident
	Params *ast.CompositeLit
	Inputs *ast.CompositeLit
}

func parseRunSteps(call *ast.CallExpr) *runStepsPrototype {
	if len(call.Args) != 3 {
		return nil
	}
	ident, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return nil
	}

	params, ok := call.Args[1].(*ast.CompositeLit)
	if !ok {
		return nil
	}
	inputs, ok := call.Args[2].(*ast.CompositeLit)
	if !ok {
		return nil
	}

	return &runStepsPrototype{
		Callee: ident,
		Params: params,
		Inputs: inputs,
	}
}

// typesToString returns the type strings of a set of expressions
func typesToString(objs ...ast.Expr) (ret []string) {
	for _, obj := range objs {
		ret = append(ret, fmt.Sprintf("%T", obj))
	}
	return
}

// rewriteRunSteps transforms
//  RunSteps(Fun, _{A: v}, _{B: v}
// to
//  target.RunSteps(_ctx, &model.Inputs{A: v, B: v})
func (p *Antha) rewriteRunSteps(call *ast.CallExpr) {
	proto := parseRunSteps(call)
	if proto == nil {
		var p runStepsPrototype
		throwErrorf(call.Pos(),
			"expecting %s(%s) found %s(%s)",
			runStepsIntrinsic,
			strings.Join(typesToString(p.Callee, p.Params, p.Inputs), ","),
			runStepsIntrinsic,
			strings.Join(typesToString(call.Args...), ","),
		)
	}

	elementPkg := path.Join(p.root.outputPackageBase, proto.Callee.Name, elementPackage)
	elementReq := p.addExternalPackage(elementPkg)
	modelPkg := path.Join(p.root.outputPackageBase, proto.Callee.Name, modelPackage)
	modelReq := p.addExternalPackage(modelPkg)

	call.Fun = mustParseExpr(elementReq.Name + "." + runStepsIntrinsic)

	call.Args = []ast.Expr{
		ast.NewIdent("_ctx"),
		&ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: mustParseExpr(modelReq.Name + ".Input"),
				Elts: append(proto.Params.Elts, proto.Inputs.Elts...),
			},
		},
	}
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

		if ident.Name == runStepsIntrinsic {
			p.rewriteRunSteps(n)
		} else if ident.Name == awaitDataIntrinsic {
			p.rewriteAwaitData(n)
		} else if desugar, ok := p.intrinsics[ident.Name]; ok {
			ident.Name = desugar
			n.Args = append([]ast.Expr{ast.NewIdent("_ctx")}, n.Args...)
		}
	}
	return true
}
