// antha/compiler/generator.go: Part of the Antha language
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

// This file implements the generation of go code from the
// antha blocks. It is primarily used in nodes.go and compile.go

package compile

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/token"
)

const (
	runStepsIntrinsic = "RunSteps"
)

var (
	unknownToken = errors.New("unknown token")
	badRun       = errors.New("bad run instruction")
)

// An input or an output
type param struct {
	Type string      // String representation of type
	Desc string      // Freeform description
	Kind token.Token // One of token.{DATA, PARAMETERS, OUTPUT, INPUT}
}

// Augmentation to compiler state to parse Antha files
type antha struct {
	element     string                          // Element name
	desc        string                          // Description of this element
	path        string                          // Normalized path to element
	inputs      map[string]param                // Inputs of an element
	inputOrder  []string                        // Canonical order of inputs to ensure deterministic output
	outputs     map[string]param                // Outputs of an element
	outputOrder []string                        // Canonical order of outputs to ensure deterministic output
	reuseMap    map[token.Token]map[string]bool // Map use of data through execution phases
	intrinsics  map[string]string               // Replacements for identifiers in expressions in functions
	types       map[string]string               // Replacement for type names in type expressions and types and type lists
	imports     map[string]string               // Additional imports with expression to suppress unused import
}

func (p *compiler) anthaInit() {
	p.inputs = make(map[string]param)
	p.outputs = make(map[string]param)
	p.reuseMap = make(map[token.Token]map[string]bool)
	for _, tok := range []token.Token{token.ANALYSIS, token.VALIDATION, token.STEPS, token.SETUP, token.REQUIREMENTS} {
		p.reuseMap[tok] = make(map[string]bool)
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
		"ReadEM":        "execute.ReadEM",
		"SetInputPlate": "execute.SetInputPlate",
	}
	p.types = map[string]string{
		"Amount":               "wunit.Amount",
		"Angle":                "wunit.Angle",
		"AngularVelocity":      "wunit.AngularVelocity",
		"Area":                 "wunit.Area",
		"Capacitance":          "wunit.Capacitance",
		"Component":            "wtype.LHComponent",
		"Concentration":        "wunit.Concentration",
		"Density":              "wunit.Density",
		"Energy":               "wunit.Energy",
		"File":                 "wtype.File",
		"FlowRate":             "wunit.FlowRate",
		"Force":                "wunit.Force",
		"HandleOpt":            "execute.HandleOpt",
		"Length":               "wunit.Length",
		"Mass":                 "wunit.Mass",
		"Plate":                "wtype.LHPlate",
		"Pressure":             "wunit.Pressure",
		"Rate":                 "wunit.Rate",
		"Resistance":           "wunit.Resistance",
		"Solution":             "wtype.LHSolution",
		"SpecificHeatCapacity": "wunit.SpecificHeatCapacity",
		"SubstanceQuantity":    "wunit.SubstanceQuantity",
		"Temperature":          "wunit.Temperature",
		"Time":                 "wunit.Time",
		"Tip":                  "wtype.LHTip",
		"Tipbox":               "wtype.LHTipbox",
		"Velocity":             "wunit.Velocity",
		"Voltage":              "wunit.Voltage",
		"Volume":               "wunit.Volume",
		"Well":                 "wtype.LHWell",
	}
	p.imports = map[string]string{
		"github.com/antha-lang/antha/antha/anthalib/wtype": "wtype.FALSE",
		"github.com/antha-lang/antha/antha/anthalib/wunit": "wunit.Make_units",
		"github.com/antha-lang/antha/execute":              "execute.MixInto",
		"github.com/antha-lang/antha/inject":               "",
		"github.com/antha-lang/antha/component":            "",
		"context": "",
	}
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

// Merges multiple import blocks and then adds paths
func (p *compiler) addAnthaImports(file *ast.File, paths []string) {

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

	for _, p := range paths {
		specs = append(specs,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(p),
					ValuePos: insertPos,
				}})
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
	// NB(ddn): tried to sort here, but the following needs proper token.Pos,
	// which are annoying to generate now. A gofmt on the generated file should
	// be just as good.
	//ast.SortImports(token.NewFileSet(), file)
}

// Return appropriate go type string for an antha (type) expr
func (p *compiler) getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		if v, ok := p.types[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}
	case *ast.SelectorExpr:
		res = p.getTypeString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		// note: array types can use a param as the length, so must be
		// allocated and treated as a slice since they can be dynamic
		res = "[]" + p.getTypeString(t.Elt)
	case *ast.StarExpr:
		res = "*" + p.getTypeString(t.X)
	case *ast.MapType:
		res = fmt.Sprintf("map[%s]%s", p.getTypeString(t.Key), p.getTypeString(t.Value))
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
}

// Return appropriate go configuration type for an antha (type) expr
func (p *compiler) getConfigTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		if v, ok := p.types[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}
		return
	case *ast.SelectorExpr:
		res = p.getConfigTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType:
		// note: array types can use a param as the length, so must be
		// allocated and treated as a slice since they can be dynamic
		res = "[]" + p.getConfigTypeString(t.Elt)
		return
	case *ast.StarExpr:
		res = "wtype.FromFactory"
		return
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
}

// Remove antha param declarations since they are fully captured in indexParams()
func (p *compiler) removeParamDecls(src *ast.File) {
	var newd []ast.Decl
	for _, decl := range src.Decls {
		toadd := true
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.PARAMETERS, token.DATA, token.INPUTS, token.OUTPUTS:
				toadd = false
			}
		}
		if toadd {
			newd = append(newd, decl)
		}
	}
	src.Decls = newd
}

// Modify AST
func (p *compiler) transform(src *ast.File) {
	if len(p.Package) != 0 {
		src.Name.Name = p.Package
	}

	p.sugarAST(src.Decls)
	p.removeParamDecls(src)

	var imports []string
	for i := range p.imports {
		imports = append(imports, i)
	}
	p.addAnthaImports(src, imports)
}

// Print out additional go code for each antha file
func (p *compiler) generate() {
	p.genFunctions()
	p.genUses()
	if err := p.genStructs(); err != nil {
		panic(err)
	}
}

func sortKeys(m map[string]param) []string {
	sorted := make([]string, len(m))
	i := 0
	for k, _ := range m {
		sorted[i] = k
		i += 1
	}
	sort.Strings(sorted)
	return sorted
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
	for idx := len(prefixes) - 1; idx >= 0; idx -= 1 {
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

// Collect information needed in downstream generation passes
func (p *compiler) analyze(src *ast.File) {
	p.desc = src.Doc.Text()
	f := p.fset.File(src.Package)
	if f != nil {
		name := f.Name()
		name, _ = relativeTo(getGoPath(), name)
		p.path = filepath.ToSlash(name)
	}
	p.element = strings.Title(src.Name.Name)

	p.recordParams(src.Decls)
	p.recordOutputUse(src.Decls)

	p.inputOrder = sortKeys(p.inputs)
	p.outputOrder = sortKeys(p.outputs)
}

// index all the spec definitions for inputs and outputs to element
func (p *compiler) recordParams(decls []ast.Decl) {
	//TODO check that names are not lowercase
	addParam := func(m *map[string]param, decl *ast.GenDecl) {
		for _, spec := range decl.Specs {
			spec := spec.(*ast.ValueSpec)
			var descs []string
			if spec.Doc != nil {
				descs = append(descs, spec.Doc.Text())
			}
			if spec.Comment != nil {
				descs = append(descs, spec.Comment.Text())
			}
			for _, name := range spec.Names {
				(*m)[name.String()] = param{Type: p.getTypeString(spec.Type), Desc: strings.Join(descs, "\n"), Kind: decl.Tok}
			}
		}
	}

	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.PARAMETERS, token.INPUTS:
				addParam(&p.inputs, decl)
			case token.OUTPUTS, token.DATA:
				addParam(&p.outputs, decl)
			}
		}
	}
}

// Track use of output parameters in antha blocks
func (p *compiler) recordOutputUse(decls []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.AnthaDecl:
			tok := decl.Tok
			ast.Inspect(decl.Body, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.Ident:
					if _, is := p.outputs[n.Name]; is {
						varName := n.Name
						p.reuseMap[tok][varName] = true
					}
				}
				return true
			})
		}
	}
}

// Generate synthesized antha functions
func (p *compiler) genFunctions() {
	// TODO(ddn): reduce boilerplate by just generating basic functions and
	// types and then implementing higher-order functionality in Go directly
	// (rather than via template generation)
	tmpl := `func _{{.Element}}Run(_ctx context.Context, input *{{.Element}}Input) *{{.Element}}Output {
	output := &{{.Element}}Output{}
	_{{.Element}}Setup(_ctx, input)
	_{{.Element}}Steps(_ctx, input, output)
	_{{.Element}}Analysis(_ctx, input, output)
	_{{.Element}}Validation(_ctx, input, output)
	return output
}

func {{.Element}}RunSteps(_ctx context.Context, input *{{.Element}}Input) *{{.Element}}SOutput {
	soutput := &{{.Element}}SOutput{}
	output := _{{.Element}}Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func {{.Element}}New() interface{} {
	return &{{.Element}}Element{
		inject.CheckedRunner {
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &{{.Element}}Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _{{.Element}}Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In: 	 &{{.Element}}Input{},
			Out: 	 &{{.Element}}Output{},
		},
	}
}
`
	params := struct {
		Element string
	}{Element: p.element}

	var buf bytes.Buffer
	template.Must(template.New("").Parse(tmpl)).Execute(&buf, params)
	p.print(buf)
}

func (p *compiler) genUses() {
	tmpl := `var (
	{{range .Uses}} _ = {{.}}
	{{end}}
)
`
	var params struct {
		Uses []string
	}

	for _, use := range p.imports {
		if len(use) > 0 {
			params.Uses = append(params.Uses, use)
		}
	}
	sort.Strings(params.Uses)

	var buf bytes.Buffer
	template.Must(template.New("").Parse(tmpl)).Execute(&buf, params)
	p.print(buf)
}

// helper function to generate the paramblock and primary element structs
func (p *compiler) genStructs() error {
	var tmpl = `type {{.Element}}Element struct {
	inject.CheckedRunner
}

type {{.Element}}Input struct {
	{{range .Inputs}}{{.Name}} {{.Type}}
	{{end}}
}

type {{.Element}}Output struct {
	{{range .Outputs}}{{.Name}} {{.Type}}
	{{end}}
}

type {{.Element}}SOutput struct {
	Data struct {
		{{range .Data}}{{.Name}} {{.Type}}
		{{end}}
	}
	Outputs struct {
		{{range .SOutputs}}{{.Name}} {{.Type}}
		{{end}}
	}
}

func init() {
	if err := addComponent(component.Component{Name: "{{.Element}}",
		Constructor: {{.Element}}New, 
		Desc: component.ComponentDesc{
			Desc: {{.Desc}},
			Path: {{.Path}},
			Params: []component.ParamDesc{
				{{range .PDesc}}component.ParamDesc{Name: {{.Name}}, Desc: {{.Desc}}, Kind: {{.Kind}}},
				{{end}}
			},
		},
	}); err != nil {
		panic(err)
	}
}
`

	type field struct {
		Name, Type string
	}

	type pdesc struct {
		Name, Desc, Kind string
	}

	params := struct {
		Element  string
		Desc     string
		Path     string
		Inputs   []field
		Outputs  []field
		Data     []field
		SOutputs []field
		PDesc    []pdesc
	}{
		Element: p.element,
		Desc:    strconv.Quote(p.desc),
		Path:    strconv.Quote(p.path),
	}

	add := func(param param, name string) error {
		f := field{Name: name, Type: param.Type}
		switch param.Kind {
		case token.INPUTS, token.PARAMETERS:
			params.Inputs = append(params.Inputs, f)
		case token.DATA:
			params.Outputs = append(params.Outputs, f)
			params.Data = append(params.Data, f)
		case token.OUTPUTS:
			params.Outputs = append(params.Outputs, f)
			params.SOutputs = append(params.SOutputs, f)
		default:
			return unknownToken
		}
		params.PDesc = append(params.PDesc, pdesc{
			Name: strconv.Quote(name),
			Desc: strconv.Quote(param.Desc),
			Kind: strconv.Quote(param.Kind.String()),
		})
		return nil
	}
	for _, name := range p.inputOrder {
		if err := add(p.inputs[name], name); err != nil {
			return err
		}
	}
	for _, name := range p.outputOrder {
		if err := add(p.outputs[name], name); err != nil {
			return err
		}
	}

	var b bytes.Buffer
	template.Must(template.New("").Parse(tmpl)).Execute(&b, params)
	p.print(b)

	return nil
}

// Update AST for antha semantics
func (p *compiler) sugarAST(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.AnthaDecl:
			ast.Inspect(decl.Body, p.inspectForIntrinsics)
			ast.Inspect(decl.Body, p.inspectForParams)
			p.sugarForTypes(decl.Body)
		}
	}
}

// Return appropriate nested SelectorExpr for the replacement for Identifier
func (p *compiler) sugarForIdent(t *ast.Ident) ast.Expr {
	if v, ok := p.types[t.Name]; ok {
		cs := strings.Split(v, ".")
		var base ast.Expr = &ast.Ident{NamePos: t.NamePos, Name: cs[0]}
		for _, c := range cs[1:] {
			base = &ast.SelectorExpr{X: base, Sel: &ast.Ident{NamePos: t.NamePos, Name: c}}
		}
		return base
	} else {
		return t
	}
}

// Return appropriate go type for an antha (type) expr
func (p *compiler) sugarForType(t ast.Expr) ast.Expr {
	switch t := t.(type) {
	case nil:
	case *ast.Ident:
		return p.sugarForIdent(t)
	case *ast.ParenExpr:
		return &ast.ParenExpr{Lparen: t.Lparen, X: p.sugarForType(t.X), Rparen: t.Rparen}
	case *ast.SelectorExpr:
		return t
	case *ast.StarExpr:
		return &ast.StarExpr{Star: t.Star, X: p.sugarForType(t.X)}
	case *ast.ArrayType:
		return &ast.ArrayType{Lbrack: t.Lbrack, Len: t.Len, Elt: p.sugarForType(t.Elt)}
	case *ast.StructType:
		return &ast.StructType{Struct: t.Struct, Fields: p.sugarForFieldList(t.Fields), Incomplete: t.Incomplete}
	case *ast.FuncType:
		return &ast.FuncType{Func: t.Func, Params: p.sugarForFieldList(t.Params), Results: p.sugarForFieldList(t.Results)}
	case *ast.InterfaceType:
		return &ast.InterfaceType{Interface: t.Interface, Methods: p.sugarForFieldList(t.Methods), Incomplete: t.Incomplete}
	case *ast.MapType:
		return &ast.MapType{Map: t.Map, Key: p.sugarForType(t.Key), Value: p.sugarForType(t.Value)}
	case *ast.ChanType:
		return &ast.ChanType{Begin: t.Begin, Arrow: t.Arrow, Dir: t.Dir, Value: p.sugarForType(t.Value)}
	default:
		log.Panicf("unexpected expression %s of type %s", t, reflect.TypeOf(t))
	}

	return t
}

func (p *compiler) sugarForFieldList(t *ast.FieldList) *ast.FieldList {
	if t == nil {
		return nil
	}
	var fields []*ast.Field
	for _, f := range t.List {
		fields = append(fields, &ast.Field{Doc: f.Doc, Names: f.Names, Type: p.sugarForType(f.Type), Tag: f.Tag, Comment: f.Comment})
	}
	return &ast.FieldList{Opening: t.Opening, List: fields, Closing: t.Closing}
}

// Replace bare antha types with go qualified names
func (p *compiler) sugarForTypes(root ast.Node) ast.Node {
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		case nil:
			return false
		case *ast.FuncLit:
			n.Type = p.sugarForType(n.Type).(*ast.FuncType)
			return false
		case *ast.CompositeLit:
			n.Type = p.sugarForType(n.Type)
			return false
		case *ast.TypeAssertExpr:
			n.Type = p.sugarForType(n.Type)
			return false
		case *ast.ValueSpec:
			n.Type = p.sugarForType(n.Type)
			return false
		default:
			return true
		}
	})
	return root
}

// Replace bare antha identifiers with go qualified names
func (p *compiler) inspectForParams(node ast.Node) bool {
	// sugar if it is a known param
	rewriteIdent := func(node *ast.Ident) {
		if _, ok := p.inputs[node.Name]; ok {
			node.Name = "_input." + node.Name
		} else if _, ok := p.outputs[node.Name]; ok {
			node.Name = "_output." + node.Name
		}
	}

	rewriteAssignLhs := func(node *ast.AssignStmt) {
		for j := range node.Lhs {
			if ident, ok := node.Lhs[j].(*ast.Ident); ok {
				isUpper := ident.Name[0:1] == strings.ToUpper(ident.Name[0:1])
				if _, fixme := p.outputs[ident.Name]; fixme && isUpper {
					ident.Name = "_output." + ident.Name
				}
			}
		}
	}

	switch n := node.(type) {
	case nil:
		return false
	case *ast.AssignStmt:
		rewriteAssignLhs(n)
	case *ast.KeyValueExpr:
		if _, identKey := n.Key.(*ast.Ident); identKey {
			// Skip identifiers that are keys
			ast.Inspect(n.Value, p.inspectForParams)
			return false
		}
	case *ast.Ident:
		rewriteIdent(n)
	case *ast.SelectorExpr:
		// Skip identifiers that are field accesses
		ast.Inspect(n.X, p.inspectForParams)
		return false
	}
	return true
}

// Replace bare antha function names with go qualified names
func (p *compiler) inspectForIntrinsics(node ast.Node) bool {
	// Transform
	//  Run(Fun, _{A: v}, _{B: v}
	// to
	//  FunRun(_ctx, FunInputs_{A: v, B: v})
	rewriteRun := func(call *ast.CallExpr) error {
		if len(call.Args) != 3 {
			return badRun
		} else if fun, ok := call.Args[0].(*ast.Ident); !ok {
			return badRun
		} else if params, ok := call.Args[1].(*ast.CompositeLit); !ok {
			return badRun
		} else if inputs, ok := call.Args[2].(*ast.CompositeLit); !ok {
			return badRun
		} else {
			call.Fun = ast.NewIdent(fun.Name + runStepsIntrinsic)
			call.Args = []ast.Expr{
				ast.NewIdent("_ctx"),
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: ast.NewIdent(fun.Name + "Input"),
						Elts: append(params.Elts, inputs.Elts...),
					},
				},
			}
		}
		return nil
	}

	switch n := node.(type) {
	case *ast.CallExpr:
		if ident, direct := n.Fun.(*ast.Ident); !direct {
		} else if ident.Name == runStepsIntrinsic {
			if err := rewriteRun(n); err != nil {
				p.internalError(err)
			}
		} else if sugar, ok := p.intrinsics[ident.Name]; !ok {
		} else {
			ident.Name = sugar
			n.Args = append([]ast.Expr{ast.NewIdent("_ctx")}, n.Args...)
		}
	}
	return true
}
