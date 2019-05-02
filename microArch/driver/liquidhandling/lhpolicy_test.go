package liquidhandling

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func getChannelForTest() *wtype.LHChannelParameter {
	return wtype.NewLHChannelParameter("ch", "gilson", wunit.NewVolume(20.0, "ul"), wunit.NewVolume(200.0, "ul"), wunit.NewFlowRate(0.0, "ml/min"), wunit.NewFlowRate(100.0, "ml/min"), 8, false, wtype.LHVChannel, 1)
}

func getSingleChannelSuck(what string, volume wunit.Volume) *SuckInstruction {
	return NewSuckInstruction(&ChannelTransferInstruction{
		What:   []string{what},
		Volume: []wunit.Volume{volume},
		Prms:   []*wtype.LHChannelParameter{getChannelForTest()},
		Multi:  1,
	})
}

func getSingleChannelBlow(what string, volume wunit.Volume) *BlowInstruction {
	return NewBlowInstruction(&ChannelTransferInstruction{
		What:   []string{what},
		Volume: []wunit.Volume{volume},
		Prms:   []*wtype.LHChannelParameter{getChannelForTest()},
		Multi:  1,
	})
}

func TestDNAPolicy(t *testing.T) {
	pft, _ := wtype.GetLHPolicyForTest()

	ins1 := getSingleChannelSuck("dna", wunit.NewVolume(2.0, "ul"))

	p, err := GetPolicyFor(pft, ins1)

	if err != nil {
		t.Error(err)
	}

	m, ok := p["POST_MIX"]

	if ok && m.(int) > 0 {
		t.Fatal("DNA must not post mix at volumes > 2 ul")
	}

	ins2 := getSingleChannelSuck("dna", wunit.NewVolume(1.99, "ul"))

	p, err = GetPolicyFor(pft, ins2)

	if err != nil {
		t.Error(err)
	}

	m, ok = p["POST_MIX"]

	if !ok || m.(int) != 1 {
		t.Fatal("DNA must have exactly 1 post mix at volumes < 2 ul")
	}
}

func TestPEGPolicy(t *testing.T) {
	pft, _ := wtype.GetLHPolicyForTest()

	ins1 := getSingleChannelSuck("peg", wunit.NewVolume(190.0, "ul"))

	p, err := GetPolicyFor(pft, ins1)
	if err != nil {
		t.Error(err)
	}

	if p["ASPZOFFSET"].(float64) != 1.0 {
		t.Error("ASPZOFFSET for PEG must be 1.0")
	}
	if p["DSPZOFFSET"].(float64) != 1.0 {
		t.Error("DSPZOFFSET for PEG must be 1.0")
	}
	if p["POST_MIX_Z"].(float64) != 1.0 {
		t.Error("POST_MIX_Z for PEG must be 1.0")
	}

	for i := 0; i < 100; i++ {
		q, err := GetPolicyFor(pft, ins1)

		if err != nil {
			t.Error(err)
		}

		if q["ASPZOFFSET"] != p["ASPZOFFSET"] {
			t.Fatal("Inconsistent Z offsets returned")
		}
	}
}

func TestPPPolicy(t *testing.T) {
	pft, _ := wtype.GetLHPolicyForTest()

	ins1 := getSingleChannelBlow("protoplasts", wunit.NewVolume(10.0, "ul"))

	p, err := GetPolicyFor(pft, ins1)
	if err != nil {
		t.Error(err)
	}

	if p["TIP_REUSE_LIMIT"].(int) != 5 {
		t.Fatal(fmt.Sprintf("Protoplast tip reuse limit is %d, not 5", p["TIP_REUSE_LIMIT"]))
	}

}

func getWaterInstructions() []RobotInstruction {
	var ret []RobotInstruction
	waters := []string{"water", "water", "water", "water", "water", "water", "water", "water"}

	{
		for _, w := range waters {
			ins := NewChannelBlockInstruction()
			tp := NewMultiTransferParams(1)
			tp.Transfers = append(tp.Transfers, TransferParams{What: w})
			ins.AddTransferParams(tp)
			ret = append(ret, ins)
		}
	}

	{
		ins := NewChannelBlockInstruction()
		for i := 0; i < 8; i++ {
			ins.What = append(ins.What, waters)
		}
		ret = append(ret, ins)
	}

	{
		ins := NewChannelTransferInstruction()
		ins.What = append(ins.What, "water")
		ret = append(ret, ins)
	}

	{
		ins := NewChannelTransferInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewAspirateInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewDispenseInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	cti := &ChannelTransferInstruction{
		What: waters,
		Prms: []*wtype.LHChannelParameter{getChannelForTest()},
	}

	{
		ins := NewBlowInstruction(cti)
		ret = append(ret, ins)
	}

	{
		ins := NewMoveRawInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewSuckInstruction(cti)
		ret = append(ret, ins)
	}

	{
		ins := NewBlowInstruction(cti)
		ret = append(ret, ins)
	}

	{
		ins := NewResetInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewMixInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	return ret
}

//TestRobotInstructionCheckLiquidClass tests that the Check method successfully
//matches liquid classes for all instructions
func TestRobotInstructionCheckLiquidClass(t *testing.T) {
	pft, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Fatal(err)
	}

	waterRule, ok := pft.Rules["water"]
	if !ok {
		t.Fatal("Couldn't get water rule")
	}

	waterInstructions := getWaterInstructions()

	for _, ins := range waterInstructions {

		if !ins.Check(waterRule) {
			t.Errorf("Instruction \"%s\" didn't match water rule, LIQUIDCLASS=%s",
				ins.Type().Name,
				ins.GetParameter(LIQUIDCLASS))
		}
	}
}

func TestSmartMixPolicy(t *testing.T) {
	pft, err := wtype.GetLHPolicyForTest()

	if err != nil {
		t.Error(
			err.Error(),
		)
	}

	cti := &ChannelTransferInstruction{
		What:    []string{"SmartMix"},
		Volume:  []wunit.Volume{wunit.NewVolume(25.0, "ul")},
		TVolume: []wunit.Volume{wunit.NewVolume(1000.0, "ul")},
		Prms:    []*wtype.LHChannelParameter{getChannelForTest()},
	}

	ins1 := NewBlowInstruction(cti)

	p, err := GetPolicyFor(pft, ins1)
	if err != nil {
		t.Error(err)
	}

	m, ok := p["POST_MIX_VOLUME"]

	if !ok || m.(float64) != 200.0 {
		t.Error("SmartMix must have post mix volume of 200 ul when pipetting into 1000ul \n",
			"found ", m, " ul \n")
	}

	cti.Volume[0] = wunit.NewVolume(25, "ul")
	cti.TVolume[0] = wunit.NewVolume(50, "ul")

	ins2 := NewBlowInstruction(cti)

	p, err = GetPolicyFor(pft, ins2)
	if err != nil {
		t.Error(err)
	}

	m, ok = p["POST_MIX_VOLUME"]

	if !ok || m.(float64) != 20.0 {
		t.Error("SmartMix must have post mix volume of 20 ul when pipetting into 50ul \n",
			"found ", m, " ul \n")
	}
}
