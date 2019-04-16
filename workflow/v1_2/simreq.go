package v1_2

import (
	"github.com/antha-lang/antha/workflow/protobuf"
)

func FromSimulateRequest(simReq protobuf.SimulateRequest) (*workflowv1_2, error) {
	wf := &workflowv1_2{}
	return wf, nil
}
