// list_plates.go: Part of the Antha language
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

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/ghodss/yaml"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listPlatesCmd = &cobra.Command{
	Use:   "plates",
	Short: "List available antha plates",
	RunE:  listPlates,
}

// TODO: replace with api definition
type simplePlate struct {
	Type          string
	WellsX        int
	WellsY        int
	WellShape     string
	WellBottom    string
	MaxWellVolume wunit.Volume
}

type simplePlates []simplePlate

func (a simplePlates) Less(i, j int) bool {
	return a[i].Type < a[j].Type
}

func (a simplePlates) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a simplePlates) Len() int {
	return len(a)
}

func listPlates(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	red := func(x string) string {
		return ansi.Color(x, "red")
	}

	ctx := testinventory.NewContext(context.Background())

	var ps simplePlates
	for _, p := range testinventory.GetPlates(ctx) {
		ps = append(ps, simplePlate{
			Type:          p.Type,
			WellsX:        p.WellsX(),
			WellsY:        p.WellsY(),
			WellShape:     p.Welltype.Shape().Type.String(),
			WellBottom:    p.Welltype.Bottom.String(),
			MaxWellVolume: p.Welltype.MaxVolume(),
		})
	}

	sort.Sort(ps)

	output := viper.GetString("output")
	switch output {
	case jsonOutput:
		bs, err := json.MarshalIndent(ps, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Println(string(bs))
		return err
	case yamlOutput:
		bs, err := yaml.Marshal(ps)
		if err != nil {
			return err
		}
		_, err = fmt.Print(string(bs))
		return err
	case textOutput:
		var lines []string
		lines = append(lines, red("PlateName")+" Properties")

		for _, p := range ps {
			prop := fmt.Sprintf("%d by %d %s %s shaped %s wells",
				p.WellsX, p.WellsY, p.WellShape, p.WellBottom, p.MaxWellVolume)

			lines = append(lines, red(p.Type)+" "+prop)
		}

		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	case csvOutput:
		var lines []string
		lines = append(lines, "PlateName,Properties")
		for _, p := range ps {
			prop := fmt.Sprintf("%d by %d,%s,%s shaped %s wells",
				p.WellsX, p.WellsY, p.WellShape, p.WellBottom, p.MaxWellVolume)
			lines = append(lines, p.Type+","+prop)
		}
		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listPlatesCmd
	listCmd.AddCommand(c)
}
