package liquidhandling

type LHOptions struct {
	ModelEvaporation        bool
	OutputSort              bool
	ExecutionPlannerVersion string
	PrintInstructions       bool
	LegacyVolume            bool
	FixVolumes              bool
}

func NewLHOptions() LHOptions {
	var lho LHOptions
	return lho
}
