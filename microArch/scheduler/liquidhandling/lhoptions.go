package liquidhandling

type LHOptions struct {
	ModelEvaporation        bool
	OutputSort              bool
	ExecutionPlannerVersion string
}

func NewLHOptions() LHOptions {
	var lho LHOptions
	return lho
}
