// +build !protobuf

package laboratory

import (
	"github.com/antha-lang/antha/instructions"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func export(idGen *id.IDGenerator, inDir, outDir string, instrs instructions.Insts, err error) error {
	return nil
}
