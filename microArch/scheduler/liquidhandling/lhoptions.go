package liquidhandling

type LHOptions struct {
	ModelEvaporation        bool
	OutputSort              bool
	ExecutionPlannerVersion string
	PrintInstructions       bool
	LegacyVolume            bool
}

func NewLHOptions() LHOptions {
	var lho LHOptions
	return lho
}
