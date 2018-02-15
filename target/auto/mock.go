package auto

import (
	"io/ioutil"
	"github.com/ghodss/yaml"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/platereader"
	"errors"
	"fmt"
)


type FakeCall struct {
	Method string 	`json:"method"`
	Response interface{}	`json:"response"`
}


type MockTargetConfig struct {
	MockDevices []MockDevice  `json:"devices"`
}


type MockDevice struct {
	DeviceClass	string `json:"class"`
	DeviceName  string `json:"name"`
	Properties  map[string]string `json:"properties"`
	FakeCalls  []FakeCall  `json:"fake_calls"`
}


// Parse the --target file
// To get a list of TargetConfig
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


func (a *MockDevice) ToDevice() (target.Device, error) {
	if a == nil {
		return nil, errors.New("no device given")
	}

	// TODO: Better way of doing this?

	// Very basic for now
	switch a.DeviceClass {
	case "antha_platereader_v1":
		return &platereader.PlateReader{}, nil
	}

	return nil, errors.New(fmt.Sprintf("unknown mock device class: '%s'", a.DeviceClass))
}
