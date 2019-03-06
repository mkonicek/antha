package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/driver/liquidhandling/server"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func assertOutputsEqual(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Errorf("output 'file' length doesn't match.\ne: %v\ng: %v", len(expected), len(got))
	} else {
		wrong := make([]string, 0, len(expected))
		for i, g := range got {
			e := expected[i]
			if e != g {
				wrong = append(wrong, fmt.Sprintf("e: %q\ng: %q", e, g))
			}
		}

		if len(wrong) > 0 {
			t.Errorf("%d lines don't match in output\n%s", len(wrong), strings.Join(wrong, "\n"))
		}
	}
}

type testDriver struct {
	callList []string
	plates   []wtype.LHObject
}

func (td *testDriver) call(s string) driver.CommandStatus {
	td.callList = append(td.callList, s)
	return driver.CommandStatus{ErrorCode: driver.OK, Msg: s}
}

func (td *testDriver) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	if o, ok := plate.(wtype.LHObject); ok {
		td.plates = append(td.plates, o)
	}
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
	return hltd.testDriver.call(fmt.Sprintf("Transfer(%v, %v, %v, %v, %v, %v)", what, platefrom, wellfrom, plateto, wellto, volume))
}

func (hltd *HighLevelTestDriver) DriverType() ([]string, error) {
	return []string{"antha.mixer.v1.Mixer", "HighLevelTestDriver"}, nil
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

func (lltd *LowLevelTestDriver) DriverType() ([]string, error) {
	return []string{"antha.mixer.v1.Mixer", "LowLevelTestDriver"}, nil
}

type HighLevelConnectionTest struct {
	Name           string
	Calls          func(liquidhandling.HighLevelLiquidhandlingDriver)
	ExpectedCalls  []string
	ExpectedPlates []wtype.LHObject
}

func (test *HighLevelConnectionTest) Run(t *testing.T) {

	drv := &HighLevelTestDriver{}
	go func() {
		if srv, err := server.NewHighLevelServer(drv); err != nil {
			t.Error(err)
		} else if err := srv.Listen(3000); err != nil {
			t.Error(err)
		}
	}()

	// give the server a moment to get set up in the thread
	time.Sleep(500 * time.Millisecond)

	c, err := client.NewHighLevelClient(":3000")
	if err != nil {
		t.Error(err)
	}

	test.Calls(c)

	b, _ := c.GetOutputFile()
	got := strings.Split(string(b), "\n")
	assertOutputsEqual(t, test.ExpectedCalls, got)

	if !reflect.DeepEqual(test.ExpectedPlates, drv.plates) {
		t.Errorf("recieved plates don't match sent plates\ne: %v\ng: %v", test.ExpectedPlates, drv.plates)
	}
}

type HighLevelConnectionTests []HighLevelConnectionTest

func (tests HighLevelConnectionTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestHighLevelConnection(t *testing.T) {
	plates := []wtype.LHObject{makePlateForTest(), makePlateForTest()}
	HighLevelConnectionTests{
		{
			Name: "simple",
			Calls: func(drv liquidhandling.HighLevelLiquidhandlingDriver) {
				drv.Initialize()
				drv.GetCapabilities()
				drv.AddPlateTo("position_1", plates[0], "firstPlate")
				drv.AddPlateTo("position_2", plates[1], "secondPlate")
				drv.Transfer([]string{"the crown jewels"}, []string{"london"}, []string{"tower of"}, []string{"me"}, []string{"head"}, []float64{100.0})
				drv.Message(100, "all your joules", "are belong to me", false)
				drv.Finalize()
			},
			ExpectedCalls: []string{
				"Initialize()",
				"GetCapabilities()",
				"AddPlateTo(\"position_1\", *wtype.Plate, \"firstPlate\")",
				"AddPlateTo(\"position_2\", *wtype.Plate, \"secondPlate\")",
				"Transfer([the crown jewels], [london], [tower of], [me], [head], [100])",
				"Message(100, \"all your joules\", \"are belong to me\", false)",
				"Finalize()",
				"GetOutputFile()",
			},
			ExpectedPlates: plates,
		},
	}.Run(t)
}

type LowLevelConnectionTest struct {
	Name           string
	Calls          func(liquidhandling.LowLevelLiquidhandlingDriver)
	ExpectedCalls  []string
	ExpectedPlates []wtype.LHObject
}

func (test *LowLevelConnectionTest) Run(t *testing.T) {

	drv := &LowLevelTestDriver{}
	go func() {
		if srv, err := server.NewLowLevelServer(drv); err != nil {
			t.Error(err)
		} else if err := srv.Listen(3001); err != nil {
			t.Error(err)
		}
	}()

	// give the server a moment to get set up in the thread
	time.Sleep(500 * time.Millisecond)

	c, err := client.NewLowLevelClient(":3001")
	if err != nil {
		t.Error(err)
	}

	test.Calls(c)

	b, _ := c.GetOutputFile()
	got := strings.Split(string(b), "\n")

	assertOutputsEqual(t, test.ExpectedCalls, got)

	if !reflect.DeepEqual(test.ExpectedPlates, drv.plates) {
		t.Errorf("recieved plates don't match sent plates\ne: %v\ng: %v", test.ExpectedPlates, drv.plates)
	}
}

type LowLevelConnectionTests []LowLevelConnectionTest

func (tests LowLevelConnectionTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestLowLevelConnection(t *testing.T) {
	plates := []wtype.LHObject{
		makeTipwasteForTest(),
		makeTipboxForTest(),
		makePlateForTest(),
		makePlateForTest(),
	}
	LowLevelConnectionTests{
		{
			Name: "simple",
			Calls: func(drv liquidhandling.LowLevelLiquidhandlingDriver) {
				drv.Initialize()
				drv.GetCapabilities()
				drv.AddPlateTo("position_1", plates[0], "tipwaste")
				drv.AddPlateTo("position_2", plates[1], "tipbox")
				drv.AddPlateTo("position_3", plates[2], "firstPlate")
				drv.AddPlateTo("position_4", plates[3], "secondPlate")
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
			ExpectedCalls: []string{
				"Initialize()",
				"GetCapabilities()",
				"AddPlateTo(\"position_1\", *wtype.LHTipwaste, \"tipwaste\")",
				"AddPlateTo(\"position_2\", *wtype.LHTipbox, \"tipbox\")",
				"AddPlateTo(\"position_3\", *wtype.Plate, \"firstPlate\")",
				"AddPlateTo(\"position_4\", *wtype.Plate, \"secondPlate\")",
				"Move([position_2], [A1], [1], [0], [0], [5], [tipbox], 0)",
				"LoadTips([0], 0, 1, [tipbox], [position_2], [A1])",
				"Move([position_3], [A1], [1], [0], [0], [5], [tipbox], 0)",
				"Aspirate([100], [false], 0, 1, [plate], [water], [false])",
				"Move([position_4], [A1], [1], [0], [0], [5], [tipbox], 0)",
				"Dispense([100], [false], 0, 1, [plate], [wine], [false])",
				"Move([position_1], [A1], [1], [0], [0], [5], [tipbox], 0)",
				"UnloadTips([0], 0, 1, [tipbox], [position_2], [A1])",
				"Message(100, \"from water\", \"into wine\", false)",
				"Finalize()",
				"GetOutputFile()",
			},
			ExpectedPlates: plates,
		},
	}.Run(t)
}

type LLGetCapabilities struct {
	LowLevelTestDriver
	Props *liquidhandling.LHProperties
}

func (gc *LLGetCapabilities) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	return *gc.Props, driver.CommandOk()
}

func TestGetCapabilities(t *testing.T) {
	expected := MakeGilsonWithPlatesAndTipboxesForTest("")
	for _, tip := range expected.Tips { //tip parents not preserved
		tip.ClearParent()
	}

	go func() {
		if srv, err := server.NewLowLevelServer(&LLGetCapabilities{Props: expected}); err != nil {
			t.Error(err)
		} else if err := srv.Listen(3002); err != nil {
			t.Error(err)
		}
	}()

	// give the server a moment to get set up in the thread
	time.Sleep(500 * time.Millisecond)

	c, err := client.NewLowLevelClient(":3002")
	if err != nil {
		t.Error(err)
	}

	if got, status := c.GetCapabilities(); !status.Ok() {
		t.Errorf("got bad status: %v", status)
	} else if !reflect.DeepEqual(expected, &got) {
		t.Errorf("Proerties changed:\ne: %+v\ng:%+v", expected, got)
	}
}
