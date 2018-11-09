package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/driver/liquidhandling/server"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type testDriver struct {
	callList []string
}

func (td *testDriver) call(s string) driver.CommandStatus {
	td.callList = append(td.callList, s)
	return driver.CommandStatus{OK: true, Msg: s}
}

func (td *testDriver) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	return td.call(fmt.Sprintf("AddPlateTo(%q, %T, %q)", position, plate, name))
}

func (td *testDriver) RemoveAllPlates() driver.CommandStatus {
	return td.call("RemoveAllPlates()")
}

func (td *testDriver) RemovePlateAt(position string) driver.CommandStatus {
	return td.call(fmt.Sprintf("RemovePlateAt(%q)", position))
}

func (td *testDriver) Initialize() driver.CommandStatus {
	return td.call("Initialize()")
}

func (td *testDriver) Finalize() driver.CommandStatus {
	return td.call("Finalize()")
}

func (td *testDriver) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	return td.call(fmt.Sprintf("Message(%d, %q, %q, %t)", level, title, text, showcancel))
}

func (td *testDriver) GetOutputFile() ([]byte, driver.CommandStatus) {
	r := td.call("GetOutputFile()")
	return []byte(strings.Join(td.callList, "\n")), r
}

func (td *testDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	return liquidhandling.LHProperties{}, td.call("GetCapabilities()")
}

type HighLevelTestDriver struct {
	testDriver
}

func (hltd *HighLevelTestDriver) Transfer(what, platefrom, wellfrom, plateto, wellto []string, volume []float64) driver.CommandStatus {
	return hltd.testDriver.call(fmt.Sprintf("Trasfer(%v, %v, %v, %v, %v, %v)", what, platefrom, wellfrom, plateto, wellto, volume))
}

type LowLevelTestDriver struct {
	testDriver
}

func (lltd *LowLevelTestDriver) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("Move(%v, %v, %v, %v, %v, %v, %v, %v)", deckposition, wellcoords, reference, offsetX, offsetY, offsetZ, plate_type, head))
}

func (lltd *LowLevelTestDriver) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("Aspirate(%v, %v, %d, %d, %v, %v, %v)", volume, overstroke, head, multi, platetype, what, llf))
}

func (lltd *LowLevelTestDriver) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("Dispense(%v, %v, %d, %d, %v, %v, %v)", volume, blowout, head, multi, platetype, what, llf))
}

func (lltd *LowLevelTestDriver) LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("LoadTips(%v, %d, %d, %v, %v, %v)", channels, head, multi, platetype, position, well))
}

func (lltd *LowLevelTestDriver) UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("UnloadTips(%v, %d, %d, %v, %v, %v)", channels, head, multi, platetype, position, well))
}

func (lltd *LowLevelTestDriver) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("SetPipetteSpeed(%d, %d, %f)", head, channel, rate))
}

func (lltd *LowLevelTestDriver) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("SetDriveSpeed(%q, %f)", drive, rate))
}

func (lltd *LowLevelTestDriver) Wait(time float64) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("Wait(%f)", time))
}

func (lltd *LowLevelTestDriver) Mix(head int, volume []float64, platetype []string, cycles []int, multi int, what []string, blowout []bool) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("Mix(%d, %v, %v, %v, %d, %v, %v)", head, volume, platetype, cycles, multi, what, blowout))
}

func (lltd *LowLevelTestDriver) ResetPistons(head, channel int) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("ResetPistons(%d, %d)", head, channel))
}

func (lltd *LowLevelTestDriver) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
	return lltd.testDriver.call(fmt.Sprintf("UpdateMetaData(props)")) //props serialisation should be tested in liquidhandlign package
}

type HighLevelConnectionTest struct {
	Name     string
	Calls    func(liquidhandling.HighLevelLiquidhandlingDriver)
	Expected []string
}

func (test *HighLevelConnectionTest) Run(t *testing.T) {

	go func() {
		if srv, err := server.NewHighLevelServer(&HighLevelTestDriver{}); err != nil {
			t.Error(err)
		} else {
			srv.Listen(3000)
		}
	}()

	c, err := client.NewHighLevelClient(":3000")
	if err != nil {
		t.Error(err)
	}

	test.Calls(c)

	b, _ := c.GetOutputFile()
	got := strings.Split(string(b), "\n")

	if !reflect.DeepEqual(test.Expected, got) {
		t.Errorf("output 'file' doesn't match.\nexpected: %v\ngot: %v", test.Expected, got)
	}
}

type HighLevelConnectionTests []HighLevelConnectionTest

func (tests HighLevelConnectionTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestHighLevelConnection(t *testing.T) {
	HighLevelConnectionTests{
		{
			Name: "simple",
			Calls: func(drv liquidhandling.HighLevelLiquidhandlingDriver) {
				drv.Initialize()
				drv.GetCapabilities()
				drv.AddPlateTo("position_1", makePlateForTest(), "firstPlate")
				drv.AddPlateTo("position_2", makePlateForTest(), "secondPlate")
				drv.Transfer([]string{"the crown jewels"}, []string{"london"}, []string{"tower of"}, []string{"me"}, []string{"head"}, []float64{100.0})
				drv.Message(100, "all your joules", "are belong to me", false)
				drv.Finalize()
			},
			Expected: []string{
				"Initialize()",
				"GetCapabilities()",
				"AddPlateTo(\"position_1\", *wtype.LHPlate, \"firstPlate\")",
				"AddPlateTo(\"position_2\", *wtype.LHPlate, \"secondPlate\")",
				"Transfer([]string{the crown jewels}, []string{london}, []string{tower of}, []string{me}, []string{head}, []float64{100.0})",
				"Message(100, \"all your joules\", \"are belong to me\", false)",
				"Finalize()",
			},
		},
	}.Run(t)
}

type LowLevelConnectionTest struct {
	Name     string
	Calls    func(liquidhandling.LowLevelLiquidhandlingDriver)
	Expected []string
}

func (test *LowLevelConnectionTest) Run(t *testing.T) {

	go func() {
		if srv, err := server.NewLowLevelServer(&LowLevelTestDriver{}); err != nil {
			t.Error(err)
		} else {
			srv.Listen(3000)
		}
	}()

	c, err := client.NewLowLevelClient(":3000")
	if err != nil {
		t.Error(err)
	}

	test.Calls(c)

	b, _ := c.GetOutputFile()
	got := strings.Split(string(b), "\n")

	if !reflect.DeepEqual(test.Expected, got) {
		t.Errorf("output 'file' doesn't match.\nexpected: %v\ngot: %v", test.Expected, got)
	}
}

type LowLevelConnectionTests []LowLevelConnectionTest

func (tests LowLevelConnectionTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestLowLevelConnection(t *testing.T) {
	LowLevelConnectionTests{
		{
			Name: "simple",
			Calls: func(drv liquidhandling.LowLevelLiquidhandlingDriver) {
				drv.Initialize()
				drv.GetCapabilities()
				drv.AddPlateTo("position_1", makeTipwasteForTest(), "tipwaste")
				drv.AddPlateTo("position_2", makeTipboxForTest(), "tipbox")
				drv.AddPlateTo("position_3", makePlateForTest(), "firstPlate")
				drv.AddPlateTo("position_4", makePlateForTest(), "secondPlate")
				drv.Move([]string{"position_2"}, []string{"A1"}, []int{1}, []float64{0.0}, []float64{0.0}, []float64{5.0}, []string{"tipbox"}, 0)
				drv.LoadTips([]int{0}, 0, 1, []string{"tipbox"}, []string{"position_2"}, []string{"A1"})
				drv.Move([]string{"position_3"}, []string{"A1"}, []int{1}, []float64{0.0}, []float64{0.0}, []float64{5.0}, []string{"tipbox"}, 0)
				drv.Aspirate([]float64{100.0}, []bool{false}, 0, 1, []string{"plate"}, []string{"water"}, []bool{false})
				drv.Move([]string{"position_4"}, []string{"A1"}, []int{1}, []float64{0.0}, []float64{0.0}, []float64{5.0}, []string{"tipbox"}, 0)
				drv.Dispense([]float64{100.0}, []bool{false}, 0, 1, []string{"plate"}, []string{"wine"}, []bool{false})
				drv.Move([]string{"position_1"}, []string{"A1"}, []int{1}, []float64{0.0}, []float64{0.0}, []float64{5.0}, []string{"tipbox"}, 0)
				drv.UnloadTips([]int{0}, 0, 1, []string{"tipbox"}, []string{"position_2"}, []string{"A1"})
				drv.Message(100, "from water", "into wine", false)
				drv.Finalize()
			},
			Expected: []string{
				"Initialize()",
				"GetCapabilities()",
				"AddPlateTo(\"position_1\", *wtype.LHPlate, \"firstPlate\")",
				"AddPlateTo(\"position_2\", *wtype.LHPlate, \"secondPlate\")",
				"Transfer([]string{the crown jewels}, []string{london}, []string{tower of}, []string{me}, []string{head}, []float64{100.0})",
				"Message(100, \"all your joules\", \"are belong to me\", false)",
				"Finalize()",
			},
		},
	}.Run(t)
}
