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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wtype/liquidtype"
	"github.com/antha-lang/antha/cmd/antha/pretty"
	"github.com/antha-lang/antha/cmd/antha/spawn"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/auto"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflowtest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:           "run",
	Short:         "Run an antha workflow",
	RunE:          runWorkflow,
	SilenceErrors: true,
}

func makeMixerOpt(ctx context.Context) (mixer.Opt, error) {
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
		f := i
		opt.ResidualVolumeWeight = &f
	}
	opt.InputPlateTypes = GetStringSlice("inputPlateTypes")
	opt.OutputPlateTypes = GetStringSlice("outputPlateTypes")
	opt.TipTypes = GetStringSlice("tipTypes")

	for _, fn := range GetStringSlice("inputPlates") {
		p, err := mixer.ParseInputPlateFile(ctx, fn)
		if err != nil {
			return opt, err
		}
		opt.InputPlates = append(opt.InputPlates, p)
	}

	policyFileName := viper.GetString("policyFile")

	if policyFileName != "" {
		data, err := ioutil.ReadFile(policyFileName)
		if err != nil {
			return opt, err
		}
		opt.CustomPolicyData, err = liquidtype.PolicyMakerFromBytes(data, wtype.PolicyName(liquidtype.BASEPolicy))
		if err != nil {
			return opt, err
		}
	}

	opt.OutputSort = viper.GetBool("outputSort")

	executionPlannerVersion := ""
	if viper.GetBool("withMulti") {
		executionPlannerVersion = "ep3"
	}
	opt.PlanningVersion = executionPlannerVersion

	opt.PrintInstructions = viper.GetBool("printInstructions")

	opt.UseDriverTipTracking = viper.GetBool("useDriverTipTracking")
	opt.IgnorePhysicalSimulation = viper.GetBool("ignorePhysicalSimulation")
	opt.LegacyVolume = viper.GetBool("legacyVolumeTracking")

	opt.FixVolumes = viper.GetBool("fixVolumes")

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
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name, Stage: desc.Stage}, runner); err != nil {
			return nil, fmt.Errorf("adding protocol %q: %s", desc.Name, err)
		}
	}
	ctx = testinventory.NewContext(ctx)
	return ctx, nil
}

type runOpt struct {
	MixerOpt               mixer.Opt
	Drivers                []string
	BundleFile             string
	ParametersFile         string
	WorkflowFile           string
	MixInstructionFileName string
	TestBundleFileName     string
	LayoutSummaryFile      string
	MixSummaryFile         string
	RunTest                bool
}

func (a *runOpt) Run() error {
	bundle, err := executeutil.UnmarshalSingle(a.BundleFile, a.WorkflowFile, a.ParametersFile)
	if err != nil {
		return err
	}

	mixerOpt := mixer.DefaultOpt.Merge(bundle.RawParams.Config).Merge(&a.MixerOpt)

	opt := auto.Opt{
		MaybeArgs: []interface{}{mixerOpt},
	}
	for _, uri := range a.Drivers {
		opt.Endpoints = append(opt.Endpoints, auto.Endpoint{URI: uri})
	}

	// Auto detect gRPC devices on network interfaces
	t, err := auto.New(opt)
	if err != nil {
		return err
	}

	ctx, err := makeContext()
	if err != nil {
		return err
	}

	rout, err := execute.Run(ctx, execute.Opt{
		Target:                     t.Target,
		Workflow:                   &bundle.Desc,
		Params:                     &bundle.RawParams,
		TransitionalReadLocalFiles: true,
	})
	if err != nil {
		return err
	}

	mixes := make([]*target.Mix, 0, len(rout.Insts))
	for _, inst := range rout.Insts {
		if mix, ok := inst.(*target.Mix); ok {
			mixes = append(mixes, mix)
		}
	}

	if a.LayoutSummaryFile != "" {
		for i, mix := range mixes {
			outFile := a.LayoutSummaryFile
			if len(mixes) > 1 {
				outFile = fmt.Sprintf("%s.%d", a.LayoutSummaryFile, i)
			}

			if bs, err := mix.SummarizeLayout(); err != nil {
				return err
			} else if err := ioutil.WriteFile(outFile, bs, 0644); err != nil {
				return err
			}
		}
	}

	if a.MixSummaryFile != "" {
		for i, mix := range mixes {
			outFile := a.MixSummaryFile
			if len(mixes) > 1 {
				outFile = fmt.Sprintf("%s.%d", a.MixSummaryFile, i)
			}

			if bs, err := mix.SummarizeActions(); err != nil {
				return err
			} else if err := ioutil.WriteFile(outFile, bs, 0644); err != nil {
				return err
			}
		}
	}

	// if option is set, add liquid handling instruction output
	if a.MixInstructionFileName != "" {
		countFiles := 1
		for _, inst := range rout.Insts {
			mi, ok := inst.(*target.Mix)

			if !ok {
				continue
			}

			fn := fmt.Sprintf("%s-%d.txt", a.MixInstructionFileName, countFiles)
			countFiles++
			if err := ioutil.WriteFile(fn, []byte(mi.Request.InstructionText), 0666); err != nil {
				return err
			}
		}
	}

	// if option is set, cache outputs for testing

	if a.TestBundleFileName != "" {
		expected := workflowtest.SaveTestOutputs(rout, "")
		bundleWithOutputs := *bundle
		bundleWithOutputs.TestOpt = expected

		if bundleWithOutputs.Version == "" {
			bundleWithOutputs.Version = "1.2.0"
		}

		bundleWithOutputs = addRun1s(bundleWithOutputs)

		serializedOutputs, err := json.MarshalIndent(bundleWithOutputs, "", "  ")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(a.TestBundleFileName, serializedOutputs, 0666); err != nil {
			return err
		}
	}

	if a.RunTest {
		err := workflowtest.CompareTestResults(rout, bundle.TestOpt)
		if err != nil {
			return err
		}
		fmt.Println("TEST BUNDLE COMPARISON OK")
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
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	ctx := testinventory.NewContext(context.Background())

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
				defer s.Close() // nolint: errcheck
			}
			if err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			} else if err := s.Start(); err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			}
			uri, err := s.URI()
			if err != nil {
				return fmt.Errorf("cannot parse port for package %s: %s", p, err)
			}
			drivers = append(drivers, uri)
		case "tcp":
			drivers = append(drivers, u.Host)
		default:
			drivers = append(drivers, u.String())
		}
	}

	mopt, err := makeMixerOpt(ctx)
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
		TestBundleFileName:     viper.GetString("makeTestBundle"),
		RunTest:                viper.GetBool("runTest"),
		LayoutSummaryFile:      viper.GetString("layoutSummary"),
		MixSummaryFile:         viper.GetString("mixSummary"),
	}

	return opt.Run()
}

func init() {
	c := runCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)
	flags.Bool("legacyVolumeTracking", false, "Do not track volumes for intermediate components")
	flags.Bool("outputSort", false, "Sort execution by output - improves tip usage")
	flags.Bool("printInstructions", false, "Output the raw instructions sent to the driver")
	flags.Bool("useDriverTipTracking", false, "If the driver has tip tracking available, use it")
	flags.Bool("ignorePhysicalSimulation", false, "Ignore errors when physically simulating the workflow - use to suppress issues caused by bugs in physical simulations")
	flags.Bool("withMulti", false, "Allow use of new multichannel planning - deprecated")
	flags.Float64("residualVolumeWeight", 0.0, "Residual volume weight")
	flags.Int("maxPlates", 0, "Maximum number of plates")
	flags.Int("maxWells", 0, "Maximum number of wells on a plate")
	flags.String("bundle", "", "Input bundle with parameters and workflow together (overrides parameter and workflow arguments)")
	flags.String("makeTestBundle", "", "Generate json format bundle for testing and put it here")
	flags.String("mixInstructionFileName", "", "Name of instructions files to output to for mixes")
	flags.String("parameters", "", "Parameters to workflow")
	flags.String("workflow", "", "Workflow definition file")
	flags.String("mixSummary", "", "save a summary of the generated liquidhandling actions to the given filename")
	flags.String("layoutSummary", "", "save a summary of the generated deck layout to the given filename")
	flags.StringSlice("component", nil, "Uris of remote components ({tcp,go}://...); use multiple flags for multiple components")
	flags.StringSlice("driver", nil, "Uris of remote drivers ({tcp,go}://...); use multiple flags for multiple drivers")
	flags.StringSlice("inputPlateTypes", nil, "Default input plate types (in order of preference)")
	flags.StringSlice("inputPlates", nil, "File containing input plates")
	flags.StringSlice("outputPlateTypes", nil, "Default output plate types (in order of preference)")
	flags.StringSlice("tipTypes", nil, "Names of permitted tip types")
	flags.Bool("runTest", false, "compare mix instructions and time estimates with results previously generated by using the makeTestBundle flag. ")
	flags.Bool("fixVolumes", true, "Make all volumes sufficient for later uses")
	flags.String("policyFile", "", "Design file of custom liquid policies in format of .xlsx JMP file")
}

func idempotentRun1Addition(name string) string {
	if !strings.HasSuffix(name, "_run1") {
		name = name + "_run1"
	}

	return name
}

// the workflow editor refuses to recognise any element without _run1 at the end
// this adds _run1 iff it is not already present as a suffix to the name of an element
func addRun1s(bin executeutil.Bundle) executeutil.Bundle {
	for name, value := range bin.Parameters {
		delete(bin.Parameters, name)
		name = idempotentRun1Addition(name)
		bin.Parameters[name] = value
	}

	for name, value := range bin.Processes {
		delete(bin.Processes, name)
		name = idempotentRun1Addition(name)
		bin.Processes[name] = value
	}

	for i, conn := range bin.Connections {
		conn.Src.Process = idempotentRun1Addition(conn.Src.Process)
		conn.Tgt.Process = idempotentRun1Addition(conn.Tgt.Process)
		bin.Connections[i] = conn
	}

	return bin
}
