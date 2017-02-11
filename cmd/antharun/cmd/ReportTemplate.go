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
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var reportTemplateCmd = &cobra.Command{
	Use:   "reporttemplate",
	Short: "produce a report template",
	RunE:  readme,
}

var synthacelements string = `cd $GOPATH/src/github.com/Synthace/elements
git fetch
git checkout ***COMMIT****
cd -`

func readme(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	switch viper.GetString("output") {
	case jsonOutput:
		_, err := fmt.Println("Json not valid for report templates")
		return err

	default:
		file := "report" + fmt.Sprint(time.Now().Format("20060102150405")) + ".md"
		var err error  
		//sdf

		anthacommit, err := GitCommit(filepath.Join(gopath(),"github.com/antha-lang/antha"))
		if err != nil {
			anthacommit = fmt.Sprintln("error getting git commit for antha-lang/antha:",err.Error())
		}
		fmt.Println("github.com/antha-lang/antha",anthacommit)
		
		elementscommit, err := GitCommit(filepath.Join(gopath(),"github.com/antha-lang/elements"))
		if err != nil {
			elementscommit = fmt.Sprintln("error getting git commit for antha-lang/elements:",err.Error())
		}
		fmt.Println("github.com/antha-lang/elements", elementscommit)
		
		synthaceelementscommit, err := GitCommit(filepath.Join(gopath(),"github.com/Synthace/elements"))
		if err != nil {
			fmt.Println("error getting git commit for Synthace/elements:",err.Error())
			synthaceelementscommit = fmt.Sprintln("error getting git commit for Synthace/elements:",err.Error())
		}
		fmt.Println("github.com/Synthace/elements", synthaceelementscommit)
		
		pipetmaxcommit, err := GitCommit(filepath.Join(gopath(),"github.com/Synthace/PipetMaxDriver"))
		if err != nil {
			pipetmaxcommit = fmt.Sprintln("error getting git commit for Synthace/PipetMaxDriver:",err.Error())
		}
		fmt.Println("github.com/Synthace/PipetMaxDriver", pipetmaxcommit)
		
		cybiocommit, err := GitCommit(filepath.Join(gopath(),"github.com/Synthace/CybioXMLDriver"))
		if err != nil {
			cybiocommit = fmt.Sprintln("error getting git commit for Synthace/CybioXMLDriver:",err.Error())
		}
		fmt.Println("github.com/Synthace/CybioXMLDriver", cybiocommit)
			
		if _, err = os.Stat(file); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
				return err
			}
			
			originalfile, err := os.Open(filepath.Join(gopath(),"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Templates/reporttemplate.md"))
			if err != nil {
				return err
			}
			
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			defer f.Close()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, originalfile); err != nil {
				return err
			}
			readme := string(buf.Bytes())

			newreadme := strings.Replace(readme, "***ANTHACOMMIT****", anthacommit, 1)
			newreadme = strings.Replace(newreadme, "***ELEMENTSCOMMIT****", elementscommit, 1)
			
			otherdependencies := strings.Replace(synthacelements, "***COMMIT****", synthaceelementscommit, 1)
			newreadme = strings.Replace(newreadme, "***OTHERDEPENDENCIES***", otherdependencies, 1)
			newreadme = strings.Replace(newreadme, "***PIPETMAXDRIVERCOMMIT****", pipetmaxcommit, 1)
			newreadme = strings.Replace(newreadme, "***CYBIODRIVERCOMMIT****", cybiocommit, 1)
			if err := ioutil.WriteFile(file, []byte(newreadme), 0666); err != nil {
				return err
			}
			return nil
		}

		return err
	}
}

func init() {
	c := reportTemplateCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String(
		"output",
		textOutput,
		fmt.Sprintf("Output format: one of {%s}", strings.Join([]string{textOutput, jsonOutput}, ",")))
}
