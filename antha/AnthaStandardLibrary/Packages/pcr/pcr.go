// This is a package that contains types and function used for PCR reactions
package pcr

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type Reaction struct {
	ReactionName string
	Template     wtype.DNASequence
	PrimerPair   [2]wtype.DNASequence
}

func (r Reaction) TemplateName() string {
	return r.Template.Nm
}
func (r Reaction) PrimerNames() [2]string {
	var primerpair [2]string
	primerpair[0] = r.PrimerPair[0].Nm
	primerpair[1] = r.PrimerPair[1].Nm
	return primerpair
}
