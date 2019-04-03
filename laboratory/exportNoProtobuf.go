// +build !protobuf

package laboratory

import (
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func export(idGen *id.IDGenerator, inDir, outDir string, instrs effects.Insts) error {
	return nil
}
