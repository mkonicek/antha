package liquidhandling

type LHOptions struct {
	OutputSort               bool
	PrintInstructions        bool
	IgnorePhysicalSimulation bool
}

func NewLHOptions() LHOptions {
	return LHOptions{}
}
