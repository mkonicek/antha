// +build protobuf

package laboratory

import (
	runner "github.com/Synthace/antha-runner/export"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func export(idGen *id.IDGenerator, inDir, outDir string, instrs effects.Insts) error {
	return runner.Export(idGen, inDir, outDir, instrs)
}
