package codegen

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/target"
)

func (a *ir) getMixes() (ret []*target.Mix) {
	for _, insts := range a.output {
		for _, inst := range insts {
			mix, ok := inst.(*target.Mix)
			if !ok {
				continue
			}
			ret = append(ret, mix)
		}
	}
	return
}

func isIncubator(dev ast.Device) bool {
	incubates := dev.CanCompile(ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1ShakerIncubator,
		},
	})

	human := dev.CanCompile(ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1Human,
		},
	})

	return incubates && !human
}

// Hacky function to identify metadata for incubator setup in lieu of better
// device modeling;
func findIncubationPlates(prop *liquidhandling.LHProperties) ([]*wtype.Plate, error) {
	var ret []*wtype.Plate
	for _, plate := range prop.Plates {
		switch {
		case strings.HasSuffix(plate.Type, "bioshake"):
			fallthrough
		case strings.HasSuffix(plate.Type, "bioshake_96well_adaptor"):
			fallthrough
		case strings.HasSuffix(plate.Type, "bioshake_standard_adaptor"):
			ret = append(ret, plate)
		}
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no incubation plates found. Only bioshake plates are currently supported")
	}

	return ret, nil
}

// addImplicitMixNodes adds additional setup for mix nodes
func (a *ir) addImplicitMixInsts() error {
	// TODO: ddn 2017-08-24: Revisit when mixes have initializers and
	// finalizers which can be scoped at the target level. Right now, only
	// device level scoping can be implemented. Other option would be to
	// continue adding more global code generation passes.
	mixes := a.getMixes()
	if len(mixes) == 0 {
		return nil
	}

	if len(mixes) != 0 {
		a.initializers = append(a.initializers,
			&target.Order{
				Mixes: mixes,
			},
			&target.PlatePrep{
				Mixes: mixes,
			},
		)
	}

	seen := make(map[ast.Device]bool)
	for d := range a.output {
		if seen[d.Device] {
			continue
		}
		seen[d.Device] = true
		if !isIncubator(d.Device) {
			continue
		}

		if len(mixes) != 1 {
			return fmt.Errorf("advanced incubator setup not yet supported")
		}

		incPlates, err := findIncubationPlates(mixes[0].Properties)
		if err != nil {
			return err
		}

		a.initializers = append(a.initializers, &target.SetupIncubator{
			Mix:              mixes[0],
			IncubationPlates: incPlates,
		})
	}

	for _, mix := range mixes {
		a.initializers = append(a.initializers, &target.SetupMixer{
			Mixes: []*target.Mix{
				mix,
			},
		})
	}

	return nil
}

// addIzers adds device-specific initializers and finalizers
func (a *ir) addIzers(deviceOrder []*drun) error {
	for _, d := range deviceOrder {
		instGraph := &target.Graph{
			Insts: a.output[d],
		}
		order, err := graph.TopoSort(graph.TopoSortOpt{
			Graph: instGraph,
		})
		if err != nil {
			return err
		}

		for _, node := range order {
			inst := node.(ast.Inst)
			init, ok := inst.(target.Initializer)
			if ok {
				a.initializers = append(a.initializers, init.GetInitializers()...)
			}

			final, ok := inst.(target.Finalizer)
			if ok {
				a.finalizers = append(a.finalizers, final.GetFinalizers()...)
			}
		}
	}

	return nil
}

// addImplicitInstrs is a cleanup pass to add implicit instructions
func (a *ir) addImplicitInsts(deviceOrder []*drun) error {
	if err := a.addImplicitMixInsts(); err != nil {
		return err
	}

	if err := a.addIzers(deviceOrder); err != nil {
		return err
	}

	return nil
}
