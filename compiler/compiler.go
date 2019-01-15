package compiler

import (
	"fmt"
	"path/filepath"

	pm_driver "github.com/Synthace/PipetMaxDriver/driver"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/mixer"
)

func Compile(path string) error {
	labBuild := laboratory.NewLaboratoryBuilder("")
	if err := labBuild.SetEffectsFromPath(path); err != nil {
		labBuild.Fatal(err)
	}

	pmPath := filepath.Join(filepath.Dir(path), "PipetMax")

	driver, err := pm_driver.New(labBuild.LaboratoryEffects, pmPath, labBuild.JobId)
	if err != nil {
		labBuild.Fatal(err)
	}
	pipetmax, err := mixer.New(mixer.DefaultOpt, driver)
	if err != nil {
		labBuild.Fatal(err)
	}

	tgt := target.New()
	tgt.AddDevice(pipetmax)

	nodes, err := labBuild.Maker.MakeNodes(labBuild.Trace.Instructions())
	if err != nil {
		return err
	}

	instrs, err := codegen.Compile(labBuild.LaboratoryEffects, tgt, nodes)
	if err != nil {
		return err
	}

	fmt.Println("instrs:", instrs)

	return nil
}
