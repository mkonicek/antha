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
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lhhelpCmd = &cobra.Command{
	Use:   "lhhelp",
	Short: "List available liquid handling policy commands",
	RunE:  lhhelp,
}

func lhhelp(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())
	fmt.Println()
	fmt.Println("Liquid handling policy commands available:")
	fmt.Println()
	fmt.Println("name,         type,   description")
	fmt.Println()
	fmt.Print(wtype.GetPolicyConsequents().TypeList())
	fmt.Println()
	return nil
}

func init() {
	c := lhhelpCmd
	//flags := c.Flags()
	RootCmd.AddCommand(c)
}
