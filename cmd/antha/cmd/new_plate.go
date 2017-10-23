// new_plate.go: Part of the Antha language
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
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newPlateCmd = &cobra.Command{
	Use:   "plate <plate type> <filename.csv>",
	Short: "Create template plate file",
	RunE:  newPlate,
}

func newPlate(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	switch len(args) {
	case 0:
		return fmt.Errorf("no plate type given")
	case 1:
		return fmt.Errorf("no output file given")
	}

	ptype := args[0]
	file := args[1]

	ctx := testinventory.NewContext(context.Background())
	plate, err := inventory.NewPlate(ctx, ptype)
	if err != nil {
		return fmt.Errorf("cannot make plate %q %s", ptype, err)
	}
	_, err = wtype.AutoExportPlateCSV(file, plate)
	return err
}

func init() {
	c := newPlateCmd
	newCmd.AddCommand(c)
}
