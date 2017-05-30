// This is a package that contains types and function used for PCR reactions
package PCR

type PCRReaction struct {
	ReactionName string
	Template     string
	PrimerPair   [2]string
}
