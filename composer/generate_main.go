package composer

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/antha-lang/antha/workflow"
)

type mainRenderer struct {
	composer *Composer
	varCount uint64
	varMemo  map[workflow.ElementInstanceName]string
}

func newMainRenderer(c *Composer) *mainRenderer {
	return &mainRenderer{
		composer: c,
		varMemo:  make(map[workflow.ElementInstanceName]string),
	}
}

func (mr *mainRenderer) render(w io.Writer) error {
	funcs := template.FuncMap{
		"elementTypes": mr.elementTypes,
		"varName":      mr.varName,
		"token":        mr.token,
	}
	if t, err := template.New("main").Funcs(funcs).Parse(tpl); err != nil {
		return err
	} else {
		return t.Execute(w, mr.composer.Workflow)
	}
}

func (mr *mainRenderer) varName(name workflow.ElementInstanceName) string {
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

func (mr *mainRenderer) elementTypes() map[workflow.ElementTypeName]*TranspilableElementType {
	return mr.composer.elementTypes
}

func (mr *mainRenderer) token(elem workflow.ElementInstanceName, param workflow.ElementParameterName) (string, error) {
	if elemInstance, found := mr.composer.Workflow.Elements.Instances[elem]; !found {
		return "", fmt.Errorf("No such element instance with name '%v'", elem)
	} else if elemType, found := mr.composer.elementTypes[elemInstance.ElementTypeName]; !found {
		return "", fmt.Errorf("No such element type with name '%v' (element instance '%v')",
			elemInstance.ElementTypeName, elem)
	} else if elemType.transpiler == nil {
		return "", fmt.Errorf("The element type '%v' does not appear to contain an Antha element",
			elemInstance.ElementTypeName)
	} else if tok, found := elemType.transpiler.TokenByParamName[string(param)]; !found {
		return "", fmt.Errorf("The element type '%v' has no parameter named '%v' (element instance '%v')",
			elemInstance.ElementTypeName, param, elem)
	} else {
		return tok.String(), nil
	}
}

var tpl = `// Code generated by antha composer.
//go:generate go-bindata data/...

package main

import (
	"bytes"
	"io/ioutil"

	"github.com/antha-lang/antha/laboratory"
	"github.com/ugorji/go/codec"

{{range elementTypes}}{{if .IsAnthaElement}}	{{printf "%q" .ImportPath}}
{{end}}{{end}})

func main() {
	labBuild := laboratory.NewLaboratoryBuilder(ioutil.NopCloser(bytes.NewBuffer(MustAsset("data/workflow.json"))))
	jh := &codec.JsonHandle{}
	labBuild.RegisterJsonExtensions(jh)

	// Register line maps for the elements we're using
{{range elementTypes}}{{if .IsAnthaElement}}	{{.Name}}.RegisterLineMap(labBuild)
{{end}}{{end}}
	// Create the elements
{{range $name, $inst := .Elements.Instances}}	{{varName $name}} := {{$inst.ElementTypeName}}.New{{$inst.ElementTypeName}}(labBuild, {{printf "%q" $name}})
{{end}}
	// Add wiring
{{range .Elements.InstancesConnections}}	labBuild.AddConnection({{varName .Source.ElementInstance}}, {{varName .Target.ElementInstance}}, func() { {{varName .Target.ElementInstance}}.{{token .Target.ElementInstance .Target.ParameterName}}.{{.Target.ParameterName}} = {{varName .Source.ElementInstance}}.{{token .Source.ElementInstance .Source.ParameterName}}.{{.Source.ParameterName}} })
{{end}}
	// Set parameters
{{range $name, $inst := .Elements.Instances}}{{range $param, $value := $inst.Parameters}}	if err := codec.NewDecoderBytes([]byte({{printf "%q" $value}}), jh).Decode(&{{varName $name}}.{{token $name $param}}.{{$param}}); err != nil {
		labBuild.Fatal(err)
	}
{{end}}{{end}}
	// Run!
	errRun := labBuild.RunElements()
	errSave := labBuild.SaveErrors()
	if errRun != nil {
		labBuild.Fatal(errRun)
	}
	if errSave != nil {
		labBuild.Fatal(errSave)
	}
	labBuild.Compile()
	labBuild.Export()
}
`
