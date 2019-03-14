// +build linkedDrivers

package mixer

import (
	"fmt"

	cb_driver "github.com/Synthace/instruction-plugins/CyBio/driver"
	lc_driver "github.com/Synthace/instruction-plugins/LabcyteEcho/driver"
	pm_driver "github.com/Synthace/instruction-plugins/PipetMax/driver"
	tecan_driver "github.com/Synthace/instruction-plugins/TecanScript/driver"
	lhdriver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

var linkedDriverFuns = map[target.MixerDriverSubType](func() lhdriver.LiquidhandlingDriver){
	target.GilsonPipetmaxSubType: func() lhdriver.LiquidhandlingDriver {
		return pm_driver.New(false)
	},
	target.CyBioSubType: func() lhdriver.LiquidhandlingDriver {
		return cb_driver.New()
	},
	target.LabcyteSubType: func() lhdriver.LiquidhandlingDriver {
		return lc_driver.New(false)
	},
	target.TecanSubType: func() lhdriver.LiquidhandlingDriver {
		return tecan_driver.New()
	},
}

func (bm *BaseMixer) maybeLinkedDriver(wf *workflow.Workflow, data []byte) error {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	if bm.connection.ExecFile == "" && bm.connection.HostPort == "" && bm.properties == nil {
		if fun, found := linkedDriverFuns[bm.expectedSubType]; !found {
			return fmt.Errorf("Unable to find linked driver function for mixer subtype %v", bm.expectedSubType)
		} else {
			bm.logger.Log("msg", "Using linked driver")
			driver := fun()
			if props, status := driver.Configure(wf.JobId, wf.Meta.Name, bm.id, data); !status.Ok() {
				return status.GetError()
			} else {
				props.Driver = driver
				bm.properties = props
				return nil
			}
		}
	}
	return nil
}
