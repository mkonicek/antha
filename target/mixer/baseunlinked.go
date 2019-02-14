// +build !linkedDrivers

package mixer

import "github.com/antha-lang/antha/workflow"

func (bm *BaseMixer) maybeLinkedDriver(wf *workflow.Workflow, data []byte) error {
	return nil
}
