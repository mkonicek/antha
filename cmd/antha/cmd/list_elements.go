// list.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"os"

	"github.com/antha-lang/antha/cmd/antha/comp"
	"github.com/antha-lang/antha/cmd/antha/pretty"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listElementsCmd = &cobra.Command{
	Use:   "elements",
	Short: "List available antha elements",
	RunE:  listElements,
}

func listElements(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	paths := make(map[string]string)
	comps := runComponents()
	for _, c := range comps {
		p, seen := paths[c.Name]
		if seen {
			return fmt.Errorf("protocol %q defined in more than one file %q and %q", c.Name, p, c.Description.Path)
		}
		paths[c.Name] = c.Description.Path
	}

	cs, err := comp.New(comps)
	if err != nil {
		return err
	}

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
		return pretty.Components(os.Stdout, cs)
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listElementsCmd
	listCmd.AddCommand(c)
}
