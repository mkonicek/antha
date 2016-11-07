package api

// Configuration for a device.
//
// swagger:model ReportTaskConfig
type DeviceConfig struct {
	// Device to execute
	DeviceId string `json:"device_id"`
	// Inputs
	DeviceInputs []*DeviceItem `json:"device_inputs"`
	// Outputs
	DeviceOutputs []*DeviceItem `json:"device_outputs"`
	// AST
	AstNodes []*AstNode `json:"ast_nodes"`
	// Raw calls to execute on device
	Calls []*GrpcCall `json:"calls"`
}

// Extract items
func (a *DeviceConfig) DeviceItems() (rs []*DeviceItem) {
	rs = append(rs, a.DeviceInputs...)
	rs = append(rs, a.DeviceOutputs...)
	return
}

type DeviceItem struct {
	Id string `json:"id"`
	// Placement of input in parent
	DeviceCoord DeviceCoord `json:"device_coord,omitempty"`
}

// Extract parent relationship
func Parent(items []*DeviceItem) map[*DeviceItem]*DeviceItem {
	parent := make(map[*DeviceItem]*DeviceItem)
	seen := make(map[string]*DeviceItem)
	for _, item := range items {
		seen[item.Id] = item
	}
	for _, item := range items {
		parent[item] = seen[item.DeviceCoord.Parent]
	}
	return parent
}

// Extract child relationship
func Children(items []*DeviceItem) map[*DeviceItem][]*DeviceItem {
	children := make(map[*DeviceItem][]*DeviceItem)
	seen := make(map[string]*DeviceItem)
	for _, item := range items {
		seen[item.Id] = item
	}
	for _, item := range items {
		p := seen[item.DeviceCoord.Parent]
		children[p] = append(children[p], seen[item.Id])
	}
	return children
}
