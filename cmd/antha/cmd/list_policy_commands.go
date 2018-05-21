// lhhelp.go: Part of the Antha language
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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listPolicyCommandsCmd = &cobra.Command{
	Use:   "policyCommands",
	Short: "List available liquid handling policy commands",
	RunE:  listPolicyCommands,
}

type simplePolicyCommand struct {
	Name string
	Type string
	Desc string
}

type simplePolicyCommands []simplePolicyCommand

func (a simplePolicyCommands) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a simplePolicyCommands) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a simplePolicyCommands) Len() int {
	return len(a)
}

func listPolicyCommands(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	var cs simplePolicyCommands
	for _, c := range wtype.GetPolicyConsequents() {
		cs = append(cs, simplePolicyCommand{
			Name: c.Name,
			Type: c.TypeName(),
			Desc: c.Desc,
		})
	}

	sort.Sort(cs)

	output := viper.GetString("output")
	switch output {
	case jsonOutput:
		bs, err := json.MarshalIndent(cs, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Println(string(bs))
		return err
	case yamlOutput:
		bs, err := yaml.Marshal(cs)
		if err != nil {
			return err
		}
		_, err = fmt.Print(string(bs))
		return err
	case csvOutput:
		var lines [][]string
		lines = append(lines, []string{"name", "type", "description"})

		for _, c := range cs {
			lines = append(lines, []string{c.Name, c.Type, c.Desc})
		}

		w := csv.NewWriter(os.Stdout)
		err := w.WriteAll(lines)

		return err
	case textOutput:
		var lines []string
		lines = append(lines, "name,type,description")

		for _, c := range cs {
			lines = append(lines, fmt.Sprintf("%s, %s, %s", c.Name, c.Type, c.Desc))
		}

		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listPolicyCommandsCmd
	listCmd.AddCommand(c)
}
