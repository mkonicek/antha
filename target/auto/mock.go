package auto

import (
	"io/ioutil"
	"github.com/ghodss/yaml"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/platereader"
	"fmt"
)

// MockTargetConfig defines a mock-target
type MockTargetConfig struct {
	MockDevices []MockDevice  `json:"devices"`
}


// MockDevice defines a mock-device
type MockDevice struct {
	DeviceClass	string `json:"class"`
	DeviceName  string `json:"name"`
	Properties  map[string]string `json:"properties"`
}


// UnmarshalMockTargetConfig parses the --target file to get a list of TargetConfig
func UnmarshalMockTargetConfig(targetConfigFilePath string) (*MockTargetConfig, error) {

	// There was no config
	if targetConfigFilePath == "" {
		return nil, nil
	}

	bTargetConfig, err := ioutil.ReadFile(targetConfigFilePath)
	if err != nil {
		return nil, err
	}
	v := new(MockTargetConfig)
	err = yaml.Unmarshal(bTargetConfig, v)
	return v, err
}


// ToDevice makes a Device from a MockDevice
func (a *MockDevice) ToDevice() (target.Device, error) {
	if a == nil {
		return nil, fmt.Errorf("no device given")
	}

	// Very basic for now
	switch a.DeviceClass {
	case "antha_platereader_v1":
		return &platereader.PlateReader{}, nil
	}

	return nil, fmt.Errorf("unknown mock device class: '%s'", a.DeviceClass)
}
