// This is a package that contains types and function used for PCR reactions
package pcr

type Reaction struct {
	ReactionName string
	Template     string
	PrimerPair   [2]string
}
