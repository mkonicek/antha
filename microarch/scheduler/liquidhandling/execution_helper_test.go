package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func testOutputSort(t *testing.T) {
	fmt.Println("TOHOHO9O")
	rq := GetLHRequestForTest()
	cmpIn := GetComponentForTest(fmt.Sprintf("water"), wunit.NewVolume(100.0, "ul"))

	for k := 0; k < 5; k++ {
		//	ins.AddProduct(GetComponentForTest("water", wunit.NewVolume(17.0, "ul")))
		cmpOut := GetComponentForTest(fmt.Sprintf("water"), wunit.NewVolume(100.0-float64(k)*3.0, "ul"))
		cmpOut.SetSample(true)
		ins := wtype.NewLHMixInstruction()
		ins.AddComponent(cmpIn)
		ins.AddProduct(cmpOut)
		rq.Add_instruction(ins)
		cmpIn = cmpOut
	}

	lh := GetLiquidHandlerForTest()
	rq.ConfigureYourself()

	err := lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got an error planning with no inputs: ", err))
	}

	fmt.Println("HOTOHO9O")
}
