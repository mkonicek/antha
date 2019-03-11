// list.go: Part of the Antha language
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

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportTemplateCmd = &cobra.Command{
	Use:   "reportTemplate",
	Short: "produce a report template in Markdown which records commits of dependent Antha repositories",
	RunE:  readme,
}

func init() {
	c := reportTemplateCmd
	newCmd.AddCommand(c)
}

var synthaceElements = `cd $GOPATH/src/github.com/Synthace/elements
git fetch
git checkout ***COMMIT****
cd -`

func readme(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	switch viper.GetString("output") {
	case jsonOutput:
		_, err := fmt.Println("json not valid for report templates")
		return err

	default:
		var err error

		anthaCommit, err := gitCommit(filepath.Join(gopath(), "github.com/antha-lang/antha"))
		if err != nil {
			anthaCommit = fmt.Sprintln("error getting git commit for antha-lang/antha:", err.Error())
		}
		fmt.Println("github.com/antha-lang/antha", anthaCommit)

		elementsCommit, err := gitCommit(filepath.Join(gopath(), "github.com/antha-lang/elements"))
		if err != nil {
			elementsCommit = fmt.Sprintln("error getting git commit for antha-lang/elements:", err.Error())
		}
		fmt.Println("github.com/antha-lang/elements", elementsCommit)

		synthaceElementsCommit, err := gitCommit(filepath.Join(gopath(), "github.com/Synthace/elements"))
		if err != nil {
			fmt.Println("error getting git commit for Synthace/elements:", err.Error())
			synthaceElementsCommit = fmt.Sprintln("error getting git commit for Synthace/elements:", err.Error())
		}
		fmt.Println("github.com/Synthace/elements", synthaceElementsCommit)

		pipetmaxCommit, err := gitCommit(filepath.Join(gopath(), "github.com/Synthace/instruction-plugins/PipetMax"))
		if err != nil {
			pipetmaxCommit = fmt.Sprintln("error getting git commit for Synthace/instruction-plugins/PipetMax:", err.Error())
		}
		fmt.Println("github.com/Synthace/instruction-plugins/PipetMax", pipetmaxCommit)

		cybioCommit, err := gitCommit(filepath.Join(gopath(), "github.com/Synthace/instruction-plugins/CyBio"))
		if err != nil {
			cybioCommit = fmt.Sprintln("error getting git commit for Synthace/instruction-plugins/CyBio:", err.Error())
		}
		fmt.Println("github.com/Synthace/instruction-plugins/CyBio", cybioCommit)

		otherDependencies := strings.Replace(synthaceElements, "***COMMIT****", synthaceElementsCommit, 1)

		type targ struct {
			TripleQuote          string
			ANTHACOMMIT          string
			ELEMENTSCOMMIT       string
			OTHERDEPENDENCIES    string
			PIPETMAXDRIVERCOMMIT string
			CYBIODRIVERCOMMIT    string
		}
		arg := targ{
			TripleQuote:          "```",
			ANTHACOMMIT:          anthaCommit,
			ELEMENTSCOMMIT:       elementsCommit,
			OTHERDEPENDENCIES:    otherDependencies,
			PIPETMAXDRIVERCOMMIT: pipetmaxCommit,
			CYBIODRIVERCOMMIT:    cybioCommit,
		}

		var out bytes.Buffer
		if err := template.Must(template.New("").Parse(reportTemplate)).Execute(&out, arg); err != nil {
			return err
		}
		fn := fmt.Sprintf("report%s.md", time.Now().Format("20060102150405"))
		if err := ioutil.WriteFile(fn, out.Bytes(), 0666); err != nil {
			return fmt.Errorf("cannot write file %q: %s", fn, err)
		}

		return err
	}
}

var reportTemplate = `
## Aim:



## Status
 


## Next steps:



## Execution instructions:


#### Get required repos

1. branch of antha-lang/antha :

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/antha-lang/antha

git fetch 
git checkout {{.ANTHACOMMIT}}
cd -
{{.TripleQuote}}


2.  branch of antha-lang/elements


{{.TripleQuote}}bash
cd $GOPATH/src/github.com/antha-lang/elements
git fetch
git checkout {{.ELEMENTSCOMMIT}}
cd -
{{.TripleQuote}}

3. Other Dependencies:

{{.TripleQuote}}bash
{{.OTHERDEPENDENCIES}}
{{.TripleQuote}}

4. (A) Pipetmaxdriver

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/instruction-plugins/PipetMax
git fetch
git checkout {{.PIPETMAXDRIVERCOMMIT}}
{{.TripleQuote}}

Or

4. (B) CybioDriver

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/instruction-plugins/CyBio
git fetch
git checkout {{.CYBIODRIVERCOMMIT}}
cd -
{{.TripleQuote}}

#### Run whenever any source code is changed  (e.g. plate definitions, antha element changes, liquid class changes)

5. Build 

{{.TripleQuote}}bash
make current -C $GOPATH/src/github.com/antha-lang/elements
{{.TripleQuote}}

or

{{.TripleQuote}}
anthabuild
{{.TripleQuote}}


#### Run when parameters or workflow is changed

5. run


PipetMax:


{{.TripleQuote}}bash
antharun --driver  go://github.com/Synthace/instruction-plugins/PipetMax
{{.TripleQuote}}


Cybio:


{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/instruction-plugins/CyBio
go build ./...
./server -machine felix
{{.TripleQuote}}


{{.TripleQuote}}bash
antharun --driver localhost:50051 --inputPlateTypes pcrplate_skirted
{{.TripleQuote}}

6. Rename output file

e.g.

{{.TripleQuote}}bash
mv generated.sqlite pipetmaxday1.sqlite
{{.TripleQuote}}

or 

{{.TripleQuote}}bash
mv cybio.xml felixday1.xml
{{.TripleQuote}}

`
