package liquidhandling

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
)

func TestTransferBlock(t *testing.T) {
	inss := make([]*wtype.LHInstruction, 8)
	dstp := factory.GetPlateByType("pcrplate_skirted_riser40")

	for i := 0; i < 8; i++ {
		c := factory.GetComponentByType("water")
		c.Vol = 100.0
		c.Vunit = "ul"
		c.SetSample(true)
		ins := wtype.NewLHInstruction()
		ins.Components = append(ins.Components, c)

		c2 := factory.GetComponentByType("tartrazine")
		c2.Vol = 24.0
		c2.Vunit = "ul"
		c2.SetSample(true)

		ins.Components = append(ins.Components, c2)

		c3 := c.Dup()
		c3.Mix(c2)
		ins.Result = c3

		ins.Platetype = "pcrplate_skirted_riser40"
		ins.Welladdress = wutil.NumToAlpha(i+1) + "1"
		ins.SetPlateID(dstp.ID)
		inss[i] = ins
	}

	tb := NewTransferBlockInstruction(inss)

	rbt := makeTestGilson()

	// make a couple of plates

	// src

	p := factory.GetPlateByType("pcrplate_skirted_riser40")

	c := factory.GetComponentByType("water")
	c.Vol = 1600
	c.Vunit = "ul"
	p.AddComponent(c, true)

	c = factory.GetComponentByType("tartrazine")
	c.Vol = 1600
	c.Vunit = "ul"

	p.AddComponent(c, true)

	rbt.AddPlate("position_4", p)

	// dst
	rbt.AddPlate("position_8", dstp)

	pol, err := GetLHPolicyForTest()
	tb.Generate(pol, rbt)

	if err != nil {
		t.Fatal(err)
	}
}
