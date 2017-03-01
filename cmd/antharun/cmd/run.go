// run.go: Part of the Antha language
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
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/cmd/antharun/frontend"
	"github.com/antha-lang/antha/cmd/antharun/pretty"
	"github.com/antha-lang/antha/cmd/antharun/spawn"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/auto"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultPort = 50051
)

var runCmd = &cobra.Command{
	Use:           "antharun",
	Short:         "Run an antha workflow",
	RunE:          runWorkflow,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func makeMixerOpt() (mixer.Opt, error) {
	opt := mixer.Opt{}
	if i := viper.GetInt("maxPlates"); i != 0 {
		f := float64(i)
		opt.MaxPlates = &f
	}
	if i := viper.GetInt("maxWells"); i != 0 {
		f := float64(i)
		opt.MaxWells = &f
	}
	if i := viper.GetFloat64("residualVolumeWeight"); i != 0 {
		f := float64(i)
		opt.ResidualVolumeWeight = &f
	}
	opt.InputPlateType = GetStringSlice("inputPlateType")
	opt.OutputPlateType = GetStringSlice("outputPlateType")
	opt.TipType = GetStringSlice("tipType")

	for _, fn := range GetStringSlice("inputPlates") {
		p, err := mixer.ParseInputPlateFile(fn)
		if err != nil {
			return opt, err
		}
		opt.InputPlates = append(opt.InputPlates, p)
	}

	return opt, nil
}

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range library {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, fmt.Errorf("error adding protocol %q: %s", desc.Name, err)
		}
	}
	return ctx, nil
}

type runOpt struct {
	MixerOpt               mixer.Opt
	Drivers                []string
	BundleFile             string
	ParametersFile         string
	WorkflowFile           string
	MixInstructionFileName string
}

func (a *runOpt) Run() error {
	var wdata, pdata, bdata []byte
	var err error

	if len(a.BundleFile) != 0 {
		bdata, err = ioutil.ReadFile(a.BundleFile)
		if err != nil {
			return err
		}
	} else {
		wdata, err = ioutil.ReadFile(a.WorkflowFile)
		if err != nil {
			return err
		}

		pdata, err = ioutil.ReadFile(a.ParametersFile)
		if err != nil {
			return err
		}
	}

	wdesc, params, err := executeutil.Unmarshal(executeutil.UnmarshalOpt{
		WorkflowData: wdata,
		BundleData:   bdata,
		ParamsData:   pdata,
	})
	if err != nil {
		return err
	}

	mixerOpt := mixer.DefaultOpt.Merge(params.Config).Merge(&a.MixerOpt)
	opt := auto.Opt{
		MaybeArgs: []interface{}{mixerOpt},
	}
	for _, uri := range a.Drivers {
		opt.Endpoints = append(opt.Endpoints, auto.Endpoint{Uri: uri})
	}
	t, err := auto.New(opt)
	if err != nil {
		return err
	}

	fe, err := frontend.New(frontend.Opt{})
	if err != nil {
		return err
	}
	defer fe.Shutdown()

	ctx, err := makeContext()
	if err != nil {
		return err
	}

	rout, err := execute.Run(ctx, execute.Opt{
		Target:   t.Target,
		Workflow: wdesc,
		Params:   params,
		TransitionalReadLocalFiles: true,
	})
	if err != nil {
		return err
	}

	// if option is set, add  liquid handling instruction output

	if a.MixInstructionFileName != "" {
		fmt.Println("OUTPUTING YOUR FILES ")
		countFiles := 1
		for _, inst := range rout.Insts {
			mi, ok := inst.(*target.Mix)

			if ok {
				fn := fmt.Sprintf("%s-%d.txt", a.MixInstructionFileName, countFiles)
				countFiles += 1

				ba := []byte(mi.Request.InstructionText)
				fb := org_antha_lang_antha_v1.FromBytes{Bytes: ba}
				bb := org_antha_lang_antha_v1.Blob_Bytes{Bytes: &fb}
				blb := org_antha_lang_antha_v1.Blob{Name: fn, From: &bb}
				var f wtype.File
				e := f.UnmarshalBlob(&blb)
				if e != nil {
					panic(fmt.Sprintf("error making file: %s", e.Error))
				}
				prcname := fmt.Sprintf("MixInstruction%d", countFiles)
				port := workflow.Port{Process: prcname, Port: "InstructionText"}

				rout.Workflow.Outputs[port] = f
			}
		}
	}

	if err := pretty.SaveFiles(os.Stdout, rout); err != nil {
		return err
	}

	if err := pretty.Timeline(os.Stdout, t, rout); err != nil {
		return err
	}

	if err := pretty.Run(os.Stdout, os.Stdin, t, rout); err != nil {
		return err
	}

	return nil
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var drivers []string
	for idx, uri := range GetStringSlice("driver") {
		u, err := url.Parse(uri)
		if err != nil {
			return err
		}

		switch u.Scheme {
		case "go":
			p := u.Host + u.Path
			s, err := spawn.GoPackage(p, fmt.Sprintf("%d %s", idx, path.Base(u.Path)))
			if s != nil {
				defer s.Close()
			}
			if err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			} else if err := s.Start(); err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			} else if uri, err := s.Uri(); err != nil {
				return fmt.Errorf("cannot parse port for package %s: %s", p, err)
			} else {
				drivers = append(drivers, uri)
			}
		case "tcp":
			drivers = append(drivers, u.Host)
		default:
			drivers = append(drivers, u.String())
		}
	}

	mopt, err := makeMixerOpt()
	if err != nil {
		return err
	}

	opt := &runOpt{
		MixerOpt:               mopt,
		Drivers:                drivers,
		BundleFile:             viper.GetString("bundle"),
		ParametersFile:         viper.GetString("parameters"),
		WorkflowFile:           viper.GetString("workflow"),
		MixInstructionFileName: viper.GetString("mixInstructionFileName"),
	}

	return opt.Run()
}

func init() {
	c := runCmd
	flags := c.Flags()

	//RootCmd.AddCommand(c)
	flags.String("parameters", "parameters.json", "Parameters to workflow")
	flags.String("workflow", "workflow.json", "Workflow definition file")
	flags.String("bundle", "", "Input bundle with parameters and workflow together (overrides parameter and workflow arguments)")
	flags.StringSlice("driver", nil, "Uris of remote drivers ({tcp,go}://...); use multiple flags for multiple drivers")
	flags.StringSlice("component", nil, "Uris of remote components ({tcp,go}://...); use multiple flags for multiple components")
	flags.Int("maxPlates", 0, "Maximum number of plates")
	flags.Int("maxWells", 0, "Maximum number of wells on a plate")
	flags.Float64("residualVolumeWeight", 0.0, "Residual volume weight")
	flags.StringSlice("inputPlateType", nil, "Default input plate types (in order of preference)")
	flags.StringSlice("outputPlateType", nil, "Default output plate types (in order of preference)")
	flags.StringSlice("inputPlates", nil, "File containing input plates")
	flags.StringSlice("tipType", nil, "Names of permitted tip types")
	flags.String("mixInstructionFileName", "", "Name of instructions files to output to for mixes")
}
