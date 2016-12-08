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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/antha-lang/antha/cmd/antharun/comp"
	"github.com/antha-lang/antha/cmd/antharun/pretty"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	jsonOutput   = "json"
	yamlOutput   = "yaml"
	stringOutput = "string"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available antha components",
	RunE:  listComponents,
}

func listComponents(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	paths := make(map[string]string)
	for _, comp := range library {
		p, seen := paths[comp.Name]
		if seen {
			return fmt.Errorf("protocol %q defined in more than one file %q and %q", comp.Name, p, comp.Desc.Path)
		}
		paths[comp.Name] = comp.Desc.Path
	}

	cs, err := comp.New(library)
	if err != nil {
		return err
	}

	output := viper.GetString("output")
	switch output {
	case jsonOutput:
		bs, err := json.Marshal(cs)
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
		_, err = fmt.Println(string(bs))
		return err
	case stringOutput:
		return pretty.Components(os.Stdout, cs)
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String(
		"output",
		stringOutput,
		fmt.Sprintf("Output format: one of {%s}", strings.Join([]string{stringOutput, yamlOutput, jsonOutput}, ",")))
}
