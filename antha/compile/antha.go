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
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	rdebug "runtime/debug"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"github.com/pkg/errors"
)

const (
	runStepsIntrinsic   = "RunSteps"
	lineNumberConstName = "_lineNumber"
)

const (
	tabWidth    = 8
	printerMode = UseSpaces | TabIndent
)

var (
	errUnknownToken = errors.New("unknown token")
	errNotAnthaFile = errors.New("not antha file")
)

type parseError interface {
	error
	Pos() token.Pos
}

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

// Antha is a preprocessing pass from antha file to go file
type Antha struct {
	SourceSHA256 []byte

	// Description of this element
	Desc string
	// Package of element
	Package string
	// Messages of an element as well as inputs and outputs
	Messages []*Message

	// inputs or outputs of an element but not messages
	tokenByParamName map[string]token.Token
	blocksUsed       map[token.Token]bool
	// Replacements for identifiers in expressions in functions
	intrinsics map[string]string
	// Replacement for go type names in type expressions and types and type lists
	types map[string]string
	// Additional imports to add
	importReqs   []*importReq
	importByName map[string]*importReq
	importProtos []string
	// externalCallPackages
	externalCallPackages []string
}

// NewAntha creates a new antha pass
func NewAntha() *Antha {
	p := &Antha{
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
		//	"Wait":          "execute.Wait",
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
		"Energy":               "wunit.Energy",
		"File":                 "wtype.File",
		"FlowRate":             "wunit.FlowRate",
		"Force":                "wunit.Force",
		"HandleOpt":            "execute.HandleOpt",
		"IncubateOpt":          "execute.IncubateOpt",
		"LHComponent":          "wtype.LHComponent",
		"LHPlate":              "wtype.LHPlate",
		"LHTip":                "wtype.LHTip",
		"LHTipbox":             "wtype.LHTipbox",
		"LHWell":               "wtype.LHWell",
		"Length":               "wunit.Length",
		"LiquidType":           "wtype.LiquidType",
		"Mass":                 "wunit.Mass",
		"PolicyName":           "wtype.PolicyName",
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
		Path: "fmt",
	})
	p.addImportReq(&importReq{
		Path: "errors",
	})
	p.addImportReq(&importReq{
		Path:    "github.com/antha-lang/antha/antha/anthalib/wtype",
		UseExpr: "wtype.FALSE",
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

// addImports merges multiple import blocks and then adds paths
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
		for _, s := range gd.Specs {
			specs = append(specs, s)
		}
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

// getTypeString return appropriate go type string for an antha (type) expr
func (p *Antha) getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {

	case nil:
		res = ""

	case *ast.Ident:
		if v, ok := p.types[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}

	case *ast.SelectorExpr:
		res = p.getTypeString(t.X) + "." + t.Sel.Name

	case *ast.BasicLit:
		res = t.Value

	case *ast.ArrayType:
		bound := p.getTypeString(t.Len)
		res = "[" + bound + "]" + p.getTypeString(t.Elt)

	case *ast.StarExpr:
		res = "*" + p.getTypeString(t.X)

	case *ast.MapType:
		res = fmt.Sprintf("map[%s]%s", p.getTypeString(t.Key), p.getTypeString(t.Value))

	default:
		throwErrorf(e.Pos(), "invalid type spec to get type of: %T", t)
	}

	return
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
	// TODO: remove below
	rdebug.PrintStack()
	return fmt.Errorf("%s: %s", p.Filename, msg)
}

func getGoPackage(filename string) string {
	filename, _ = relativeTo(getGoPath(), filename)
	dir := filepath.ToSlash(filepath.Dir(filename))
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
		r := part
		r = strings.Map(func(r rune) rune {
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

	protocolName := src.Name.Name

	src.Name.Name = "main"
	src.Tok = token.PACKAGE

	p.Desc = src.Doc.Text()

	file := fileSet.File(src.Package)
	p.Package = getGoPackage(file.Name())

	if e, f := protocolName, path.Base(p.Package); e != f {
		return fmt.Errorf("%s: expecting protocol %s to be in directory %s", file.Name(), protocolName, e)
	}

	p.recordImports(src.Decls)
	p.recordBlocks(src.Decls)
	p.recordMessages(src.Decls)
	if err := p.validateMessages(p.Messages); err != nil {
		return err
	}

	p.desugar(fileSet, src)

	for _, pkg := range p.externalCallPackages {
		p.addImportReq(&importReq{
			Name: manglePackageName(pkg),
			Path: pkg,
		})
	}

	p.addImportReq(&importReq{
		Path: p.Package,
		Name: "protobuf",
	})

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
	for _, v := range bases {
		prefixes = append(prefixes, v)
	}

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

			req := &importReq{
				Added: true,
				Path:  im.Path.Value,
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

			p.Messages = append(p.Messages, &Message{
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

func (p *Antha) generateMain(fileSet *token.FileSet, file *ast.File) ([]byte, error) {
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
	io.Copy(&out, bytes.NewReader(main))

	if err := p.printFunctions(&out); err != nil {
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

func (p *Antha) generateProto() ([]byte, error) {
	var tmpl = `syntax = "proto3";
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

	merge := func(messages []*Message, merged, aname, bname string) (ret []*Message) {
		var a Message
		var b Message
		for _, msg := range messages {
			if msg.Name == aname {
				a = *msg
				continue
			}
			if msg.Name == bname {
				b = *msg
				continue
			}
			ret = append(ret, msg)
		}
		ret = append(ret, &Message{
			Name:   merged,
			Doc:    a.Doc + b.Doc,
			Fields: append(a.Fields, b.Fields...),
		})

		return
	}

	type field struct {
		Desc string
		Type string
		Name string
		Tag  int
	}

	type message struct {
		Name   string
		Desc   string
		Fields []field
	}

	type tvars struct {
		PackageName string
		Imports     []string
		Messages    []message
	}

	tv := tvars{
		PackageName: getProtobufPackage(p.Package),
	}

	messages := p.Messages
	messages = merge(messages, "Request", "Parameters", "Inputs")
	messages = merge(messages, "Response", "Data", "Outputs")

	used := make(map[*importReq]bool)
	for _, msg := range messages {
		m := message{
			Name: msg.Name,
			Desc: fmtDoc("", msg.Doc),
		}
		for idx, f := range msg.Fields {
			m.Fields = append(m.Fields, field{
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

		tv.Imports = append(tv.Imports, path.Join(req.Path, "element.proto"))
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
func (p *Antha) Generate(fileSet *token.FileSet, file *ast.File) (map[string][]byte, error) {
	mainBs, err := p.generateMain(fileSet, file)
	if err != nil {
		return nil, err
	}

	protoBs, err := p.generateProto()
	if err != nil {
		return nil, err
	}

	defaultGo := []byte("//go:generate protoc -I. -Ivendor --go_out=plugins=grpc:. element.proto\npackage protobuf\n")

	return map[string][]byte{
		"server/main.go": mainBs,
		"element.proto":  protoBs,
		"element.go":     defaultGo,
	}, nil
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
	var tmpl = `
type Element struct {}

func (Element) Run(_ctx context.Context, request *protobuf.Request) (response *protobuf.Response, err error) {
	defer func() {
		if res := recover(); res != nil {
			e, ok := res.(error)
			if ok {
				err = e
			} else {
				err = errors.New(fmt.Sprint(res))
			}
		}
	}()
	response = &protobuf.Response{}
	{{if .HasSetup}}_Setup(_ctx, request, response){{end}}
	{{if .HasSteps}}_Steps(_ctx, request, response){{end}}
	{{if .HasAnalysis}}_Analysis(_ctx, request, response){{end}}
	{{if .HasValidation}}_Validation(_ctx, request, response){{end}}
	return
}

func (Element) Metadata(_ctx context.Context, request *api.Empty) (*api.ElementMetadata, error) {
	return _metadata, nil
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
	type field struct {
		Name string
		Type string
	}

	type port struct {
		Name     string
		RawName  string
		PortType string
		Doc      string
	}

	type call struct {
		Name       string
		InputType  string
		OutputType string
	}

	type dependency struct {
		Package   string
		GitCommit string
		Dirty     bool
	}

	type tvars struct {
		SHA256        string
		CommentDoc    string
		Package       string
		Ports         []*port
		Calls         []*call
		Dependencies  []*dependency
		Inputs        []field
		Outputs       []field
		HasSteps      bool
		HasValidation bool
		HasSetup      bool
		HasAnalysis   bool
	}

	tv := tvars{
		SHA256:        encodeByteArray(p.SourceSHA256),
		CommentDoc:    strconv.Quote(p.Desc),
		Package:       strconv.Quote(p.Package),
		HasSteps:      p.blocksUsed[token.STEPS],
		HasValidation: p.blocksUsed[token.VALIDATION],
		HasSetup:      p.blocksUsed[token.SETUP],
		HasAnalysis:   p.blocksUsed[token.ANALYSIS],
	}

	seen := make(map[string]bool)
	for _, pkg := range p.externalCallPackages {
		if seen[pkg] {
			continue
		}
		seen[pkg] = true

		mangledPackage := manglePackageName(pkg)
		mangledFunc := "_" + runStepsIntrinsic + mangledPackage

		tv.Calls = append(tv.Calls, &call{
			Name:       mangledFunc,
			InputType:  mangledPackage + ".Input",
			OutputType: mangledPackage + ".Output",
		})
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
		panic(errors.Wrap(err, x))
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
				&ast.Field{
					Names: identList("_ctx"),
					Type:  mustParseExpr("context.Context"),
				},
				&ast.Field{
					Names: identList("_input"),
					Type:  mustParseExpr("*protobuf.Request"),
				},
				&ast.Field{
					Names: identList("_output"),
					Type:  mustParseExpr("*protobuf.Response"),
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

// rewriteRunSteps transforms
//  RunSteps("Fun", _{A: v}, _{B: v}
// to
//  _RunStepsxxxxxx(_ctx, &xxxxxx.Inputs{A: v, B: v})
func (p *Antha) rewriteRunSteps(call *ast.CallExpr) {
	if len(call.Args) != 3 {
		throwErrorf(call.Pos(), "%s takes three arguments", runStepsIntrinsic)
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		throwErrorf(call.Pos(), "first argument of %s must be a string literal", runStepsIntrinsic)
	} else if lit.Kind != token.STRING {
		throwErrorf(call.Pos(), "first argument of %s must be a string literal", runStepsIntrinsic)
	}

	params, ok := call.Args[1].(*ast.CompositeLit)
	if !ok {
		throwErrorf(call.Pos(), "second argument of %s must be a struct literal", runStepsIntrinsic)
	}
	inputs, ok := call.Args[2].(*ast.CompositeLit)
	if !ok {
		throwErrorf(call.Pos(), "third argument of %s must be a struct literal", runStepsIntrinsic)
	}

	pkg, err := strconv.Unquote(lit.Value)
	if err != nil {
		throwErrorf(call.Pos(), err.Error())
	}

	p.externalCallPackages = append(p.externalCallPackages, pkg)

	mangledPackage := manglePackageName(pkg)
	mangledFunc := "_" + runStepsIntrinsic + mangledPackage

	call.Fun = ast.NewIdent(mangledFunc)
	call.Args = []ast.Expr{
		ast.NewIdent("_ctx"),
		&ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: mustParseExpr(mangledPackage + ".Input"),
				Elts: append(params.Elts, inputs.Elts...),
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
		} else if desugar, ok := p.intrinsics[ident.Name]; ok {
			ident.Name = desugar
			n.Args = append([]ast.Expr{ast.NewIdent("_ctx")}, n.Args...)
		}
	}
	return true
}
