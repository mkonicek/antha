// new_element.go: Part of the Antha language
// Copyright (C) 2016 The Antha authors. All rights reserved.
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

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
	"unicode"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newElementCmd = &cobra.Command{
	Use:   "element <name> [output directory]",
	Short: "Create template element in directory",
	RunE:  newElement,
}

const elementTemplate = `// Protocol {{.Name}} performs something.
//
// All of this text should be used to describe what this protocol does.  It
// should begin with a one sentence summary begining with "Protocol X...". If
// neccessary, a empty line with a detailed description can follow (like this
// description does).
//
// Spend some time thinking of a good protocol name as this is the name by
// which this protocol will be referred. It should convey the purpose and scope
// of the protocol to an outsider and should suggest an obvious
// parameterization. 
//
// Protocol names are also case-sensitive, so try to use a consistent casing
// scheme.
//
// Examples of bad names:
//   - MyProtocol
//   - GeneAssembly
//   - WildCAPSsmallANDLARGE
//
// Better names:
//   - Aliquot
//   - TypeIIsConstructAssembly
protocol {{.Name}}

// Place golang packages to import here
import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
)

// Parameters to this protocol
Parameters {
{{range .Parameters}}	{{.Name}} {{.Type}}
{{end}}}

// Output data of this protocol
Data {
{{range .Data}}	{{.Name}} {{.Type}}
{{end}}}

// Physical inputs to this protocol
Inputs {
{{range .Inputs}}	{{.Name}} {{.Type}}
{{end}}}

// Physical outputs to this protocol
Outputs {
{{range .Outputs}}	{{.Name}} {{.Type}}
{{end}}}

// Conditions to run on startup
Setup {

}

// The core process for this protocol. These steps are executed for each input.
Steps {
{{range .Steps}}{{.}}
{{end}}}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
Analysis {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
Validation {

}
`

// Is string a valid element name?
func validElementName(name string) error {
	// Elements are just go types so they must follow go identifier naming:
	//
	//   identifier     = letter { letter | unicode_digit }
	//   letter         = unicode_letter | "_"
	//   unicode_digit  = /* Unicode code point classification Nd */
	//   unicode_letter = /* Unicode code point classification Lu, Ll, Lt, Lm, Lo */
	//
	// Elements are exported types so they must begin with a letter in Unicode
	// class Lu.

	if len(name) == 0 {
		return fmt.Errorf("empty string")
	}

	for idx, runeValue := range name {
		if idx == 0 && !unicode.In(runeValue, unicode.Lu) {
			return fmt.Errorf("must begin with upper case character")
		}
		if runeValue == '_' {
			continue
		}
		if !unicode.In(runeValue, unicode.Lu, unicode.Ll, unicode.Lm, unicode.Lo, unicode.Nd) {
			return fmt.Errorf("character %q is not a letter or digit", runeValue)
		}
	}

	return nil
}

func writeAn(outputDir string, steps []string, element *component.Component) error {
	type parg struct {
		Name string
		Type string
	}
	type targ struct {
		Name       string
		Steps      []string
		Parameters []parg
		Data       []parg
		Inputs     []parg
		Outputs    []parg
	}
	arg := targ{
		Name:  element.Name,
		Steps: steps,
	}
	for _, p := range element.Description.Params {
		switch p.Kind {
		case "Parameters":
			arg.Parameters = append(arg.Parameters, parg{Name: p.Name, Type: p.Type})
		case "Inputs":
			arg.Inputs = append(arg.Inputs, parg{Name: p.Name, Type: p.Type})
		case "Data":
			arg.Data = append(arg.Data, parg{Name: p.Name, Type: p.Type})
		case "Outputs":
			arg.Outputs = append(arg.Outputs, parg{Name: p.Name, Type: p.Type})
		default:
			return fmt.Errorf("unknown kind %q", p.Kind)
		}
	}

	var out bytes.Buffer
	if err := template.Must(template.New("").Parse(elementTemplate)).Execute(&out, arg); err != nil {
		return err
	}

	fn := filepath.Join(outputDir, element.Name+".an")
	if err := ioutil.WriteFile(fn, out.Bytes(), 0666); err != nil {
		return fmt.Errorf("cannot write file %q: %s", fn, err)
	}
	return nil
}

func writeBundle(outputDir string, element *component.Component, params map[string]interface{}) error {
	desc := workflow.Desc{
		Processes: map[string]workflow.Process{
			"Process1": {
				Component: element.Name,
			},
		},
	}
	p := execute.Params{
		Parameters: map[string]map[string]interface{}{
			"Process1": params,
		},
	}

	type bundle struct {
		workflow.Desc
		execute.Params
	}

	b := bundle{
		Desc:   desc,
		Params: p,
	}

	bs, err := json.MarshalIndent(&b, "", "  ")
	if err != nil {
		return err
	}

	fn := filepath.Join(outputDir, element.Name+".bundle.json")
	if err := ioutil.WriteFile(fn, bs, 0600); err != nil {
		return fmt.Errorf("cannot write file %q: %s", fn, err)
	}
	return nil
}

func newElement(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	var name, outputDir string

	switch len(args) {
	case 0:
		return fmt.Errorf("missing element name")
	case 1:
		name = args[0]
		outputDir = name
	default:
		name = args[0]
		outputDir = args[1]
	}

	if err := validElementName(name); err != nil {
		return fmt.Errorf("element name %q is invalid: %s", name, err)
	}

	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("cannot make directory %q: %s", outputDir, err)
	}

	// TODO: link params, input, output, and element
	type input struct {
		A           float64
		B           float64
		ComponentA  *wtype.Liquid
		ComponentB  *wtype.Liquid
		Option      bool
		String      string
		StringArray []string
		VolumeA     wunit.Volume
		VolumeB     wunit.Volume
	}

	type output struct {
		MixedComponent *wtype.Liquid
		Sum            float64
	}

	params := map[string]interface{}{
		"A":           2.99,
		"B":           -1.0,
		"ComponentA":  "water",
		"ComponentB":  "dna",
		"Option":      false,
		"String":      "Example",
		"StringArray": []string{"A", "B", "C"},
		"VolumeA":     "1ul",
		"VolumeB":     "2ul",
	}

	var in input
	var out output

	typeName := func(i interface{}) string {
		return fmt.Sprint(reflect.TypeOf(i))
	}
	element := &component.Component{
		Name: name,
		Constructor: func() interface{} {
			return &inject.CheckedRunner{
				In:  &input{},
				Out: &output{},
			}
		},
		Description: component.Description{
			Params: []component.ParamDesc{
				{Name: "A", Kind: "Parameters", Type: typeName(in.A)},
				{Name: "B", Kind: "Parameters", Type: typeName(in.B)},
				{Name: "ComponentA", Kind: "Inputs", Type: typeName(in.ComponentA)},
				{Name: "ComponentB", Kind: "Inputs", Type: typeName(in.ComponentB)},
				{Name: "MixedComponent", Kind: "Outputs", Type: typeName(out.MixedComponent)},
				{Name: "Option", Kind: "Parameters", Type: typeName(in.Option)},
				{Name: "String", Kind: "Parameters", Type: typeName(in.String)},
				{Name: "StringArray", Kind: "Parameters", Type: typeName(in.StringArray)},
				{Name: "Sum", Kind: "Data", Type: typeName(out.Sum)},
				{Name: "VolumeA", Kind: "Parameters", Type: typeName(in.VolumeA)},
				{Name: "VolumeB", Kind: "Parameters", Type: typeName(in.VolumeB)},
			},
		},
	}

	steps := []string{
		"\tSum = A + B",
		"\tsampleA := mixer.Sample(ComponentA, VolumeA)",
		"\tsampleB := mixer.Sample(ComponentB, VolumeB)",
		"\tMixedComponent = Mix(sampleA, sampleB)",
	}

	if err := writeAn(outputDir, steps, element); err != nil {
		return err
	}

	if err := writeBundle(outputDir, element, params); err != nil {
		return err
	}

	return nil
}

func init() {
	c := newElementCmd

	newCmd.AddCommand(c)
}
