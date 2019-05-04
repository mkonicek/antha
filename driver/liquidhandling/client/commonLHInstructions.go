package client

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	driver "github.com/antha-lang/antha/microArch/driver"
)

type CommonLHInstructions struct{}

// method common to all liquidhandling instruction drivers - embed in all dependent methods and reimplement as needed
// https://medium.com/@simplyianm/why-gos-structs-are-superior-to-class-based-inheritance-b661ba897c67 !!
func (c *CommonLHInstructions) GetAllowedLocations(plate *wtype.Plate) ([]string, driver.CommandStatus) {
	//self.AddWarning("Pass Through")
	return nil, driver.CommandOk()
}
