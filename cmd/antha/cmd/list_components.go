// list_components.go: Part of the Antha language
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
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/ghodss/yaml"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "List available antha liquid component types",
	RunE:  listComponents,
}

// TODO: replace with api definition
type simpleComponent struct {
	Name       string
	LiquidType string
}

type simpleComponents []simpleComponent

func (a simpleComponents) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a simpleComponents) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a simpleComponents) Len() int {
	return len(a)
}

func listComponents(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	red := func(x string) string {
		return ansi.Color(x, "red")
	}

	ctx := testinventory.NewContext(context.Background())

	var cs simpleComponents
	for _, c := range testinventory.GetComponents(ctx) {
		cs = append(cs, simpleComponent{
			Name:       c.CName,
			LiquidType: c.TypeName(),
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
	case textOutput:
		var lines []string
		lines = append(lines, red("ComponentName")+" LiquidTypeName")

		for _, c := range cs {
			lines = append(lines, red(c.Name)+" "+c.LiquidType)
		}

		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	case csvOutput:
		var lines []string
		lines = append(lines, "ComponentNames,LiquidTypeName")
		for _, c := range cs {
			lines = append(lines, c.Name+","+c.LiquidType)
		}

		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listComponentsCmd
	listCmd.AddCommand(c)
}
