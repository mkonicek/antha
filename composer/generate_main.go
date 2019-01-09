package composer

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"
)

type mainRenderer struct {
	composer *Composer
	varCount uint64
	varMemo  map[ElementInstanceName]string
}

func newMainRenderer(c *Composer) *mainRenderer {
	return &mainRenderer{
		composer: c,
		varMemo:  make(map[ElementInstanceName]string),
	}
}

func (mr *mainRenderer) render(w io.Writer) error {
	funcs := template.FuncMap{
		"elementTypes": mr.elementTypes,
		"varName":      mr.varName,
	}
	if t, err := template.New("main").Funcs(funcs).Parse(tpl); err != nil {
		return err
	} else {
		return t.Execute(w, mr.composer.Workflow)
	}
}

func (mr *mainRenderer) varName(name ElementInstanceName) string {
	if res, found := mr.varMemo[name]; found {
		return res
	}

	res := make([]rune, 0, len(name))
	ensureUpper := false
	for _, r := range []rune(name) {
		switch {
		case 'a' <= r && r <= 'z' && ensureUpper:
			ensureUpper = false
			res = append(res, unicode.ToUpper(r))
		case 'a' <= r && r <= 'z':
			res = append(res, r)
		case 'A' <= r && r <= 'Z' && len(res) == 0:
			res = append(res, unicode.ToLower(r))
		case 'A' <= r && r <= 'Z':
			res = append(res, r)
			ensureUpper = false
		case strings.ContainsRune(" -_", r):
			ensureUpper = true
		}
	}
	resStr := fmt.Sprintf("%s%d", string(res), mr.varCount)
	mr.varCount++
	mr.varMemo[name] = resStr
	return resStr
}

func (mr *mainRenderer) elementTypes() map[ElementTypeName]*ElementType {
	return mr.composer.elementTypes
}

var tpl = `// Code generated by antha composer. DO NOT EDIT.
package main

import (
	"encoding/json"

	"github.com/antha-lang/antha/laboratory"

{{range elementTypes}}{{if .IsAnthaElement}}	{{printf "%q" .ImportPath}}
{{end}}{{end}})

func main() {
	labBuild := laboratory.NewLaboratoryBuilder({{printf "%q" .JobId}})
	// Register line maps for the elements we're using
{{range elementTypes}}{{if .IsAnthaElement}}	{{.Name}}.RegisterLineMap(labBuild)
{{end}}{{end}}
	// Create the elements
{{range $name, $type := .ElementInstances}}	{{varName $name}} := {{$type.ElementTypeName}}.New{{$type.ElementTypeName}}(labBuild, {{printf "%q" $name}})
{{end}}
	// Add wiring
{{range .ElementInstancesConnections}}	labBuild.AddLink({{varName .Source.ElementInstance}}, {{varName .Target.ElementInstance}}, func () { {{varName .Target.ElementInstance}}.Inputs.{{.Target.ParameterName}} = {{varName .Source.ElementInstance}}.Outputs.{{.Source.ParameterName}} })
{{end}}
	// Set parameters
{{range $name, $params := .ElementInstancesParameters}}{{range $param, $value := $params}}	if err := json.Unmarshal([]byte({{printf "%q" $value}}), &{{varName $name}}.Parameters.{{$param}}); err != nil {
		labBuild.Fatal(err)
	}
{{end}}{{end}}
	// Run!
	errRun := labBuild.RunElements()
	errSave := labBuild.Save()
	if errRun != nil {
		labBuild.Fatal(errRun)
	}
	if errSave != nil {
		labBuild.Fatal(errSave)
	}
}
`
