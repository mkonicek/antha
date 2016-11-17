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
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lhcomponentsCmd = &cobra.Command{
	Use:   "lhcomponents",
	Short: "List available antha lhcomponents",
	RunE:  lhComponents,
}

func lhComponents(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	cs := factory.GetComponentList()
	/*if err != nil {
		return err
	}*/

	switch viper.GetString("output") {
	case jsonOutput:
		if bs, err := json.Marshal(cs); err != nil {
			return err
		} else {
			_, err = fmt.Println(string(bs))
			return err
		}
	default:

		prettystrings := make([]string, 0)

		prettystrings = append(prettystrings, text.Print("Component Name", "LiquidTypeName"))

		for i := range cs {
			prettystrings = append(prettystrings, text.Print(cs[i], factory.GetComponentByType(cs[i]).TypeName()))
		}
		_, err := fmt.Println(strings.Join(prettystrings, ""))
		return err
	}
}

func init() {
	c := lhcomponentsCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String(
		"output",
		stringOutput,
		fmt.Sprintf("Output format: one of {%s}", strings.Join([]string{stringOutput, jsonOutput}, ",")))
}
