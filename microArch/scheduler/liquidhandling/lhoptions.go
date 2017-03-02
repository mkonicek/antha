package liquidhandling

type LHOptions struct {
	ModelEvaporation bool
	OutputSort       bool
}

func NewLHOptions() LHOptions {
	var lho LHOptions
	return lho
}
