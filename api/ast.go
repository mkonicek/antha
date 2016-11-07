package api

type AstNode struct {
	Id           string
	Inputs       []*AstNode    `json:"inputs"`
	Outputs      []*AstNode    `json:"outputs"`
	RawInst      *RawInst      `json:"handle_inst,omitempty"`
	IncubateInst *IncubateInst `json:"incubate_inst,omitempty"`
	MixInst      *MixInst      `json:"mix_inst,omitempty"`
}

type RawInst struct {
	Group    string            `json:"group"`
	Selector map[string]string `json:"selector"`
	Calls    []*GrpcCall       `json:"calls"`
}

type IncubateInst struct {
	Time Measurement `json:"time"`
	Temp Measurement `json:"temp"`
}

type MixInst struct {
	// TBD
}
