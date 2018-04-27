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

func TestDNAPolicy(t *testing.T) {
	pft, _ := wtype.GetLHPolicyForTest()

	tp := TransferParams{
		What:    "dna",
		Volume:  wunit.NewVolume(2.0, "ul"),
		Channel: getChannelForTest(),
	}

	ins1 := NewSuckInstruction()
	ins1.AddTransferParams(tp)

	p, err := GetPolicyFor(pft, ins1)

	if err != nil {
		t.Error(err)
	}

	m, ok := p["POST_MIX"]

	if ok && m.(int) > 0 {
		t.Fatal("DNA must not post mix at volumes > 2 ul")
	}

	tp.Volume = wunit.NewVolume(1.99, "ul")

	ins2 := NewSuckInstruction()
	ins2.AddTransferParams(tp)
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

	tp := TransferParams{
		What:    "PEG",
		Volume:  wunit.NewVolume(190.0, "ul"),
		Channel: getChannelForTest(),
	}

	ins1 := NewSuckInstruction()
	ins1.AddTransferParams(tp)

	p, err := GetPolicyFor(pft, ins1)

	if err != nil {
		t.Error(err)
	}

	if p["ASPZOFFSET"].(float64) != 1.0 {
		t.Fatal("ASPZOFFSET for PEG must be 1.0")
	}
	if p["DSPZOFFSET"].(float64) != 1.0 {
		t.Fatal("DSPZOFFSET for PEG must be 1.0")
	}
	if p["POST_MIX_Z"].(float64) != 1.0 {
		t.Fatal("POST_MIX_Z for PEG must be 1.0")
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

	tp := TransferParams{
		What:    "Protoplasts",
		Volume:  wunit.NewVolume(10.0, "ul"),
		Channel: getChannelForTest(),
	}

	ins1 := NewBlowInstruction()
	ins1.AddTransferParams(tp)

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
		ins := NewSingleChannelBlockInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewMultiChannelBlockInstruction()
		for i := 0; i < 8; i++ {
			ins.What = append(ins.What, waters)
		}
		ret = append(ret, ins)
	}

	{
		ins := NewSingleChannelTransferInstruction()
		ins.What = "water"
		ret = append(ret, ins)
	}

	{
		ins := NewMultiChannelTransferInstruction()
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

	{
		ins := NewBlowInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewMoveRawInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewSuckInstruction()
		ins.What = waters
		ret = append(ret, ins)
	}

	{
		ins := NewBlowInstruction()
		ins.What = waters
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
				InstructionTypeName(ins),
				ins.GetParameter("LIQUIDCLASS"))
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

	tp := TransferParams{
		What:    "SmartMix",
		Volume:  wunit.NewVolume(25.0, "ul"),
		TVolume: wunit.NewVolume(1000.0, "ul"),
		FVolume: wunit.NewVolume(1000.0, "ul"),
		Channel: getChannelForTest(),
	}

	ins1 := NewSuckInstruction()
	ins1.AddTransferParams(tp)

	p, err := GetPolicyFor(pft, ins1)

	if err != nil {
		t.Error(err)
	}

	m, ok := p["POST_MIX_VOLUME"]

	if !ok || m.(float64) != 200.0 {
		t.Error("SmartMix must have post mix volume of 200 ul when pipetting into 1000ul \n",
			"found ", m, " ul \n")
	}

	tp.Volume = wunit.NewVolume(25, "ul")
	tp.TVolume = wunit.NewVolume(50, "ul")

	ins2 := NewSuckInstruction()
	ins2.AddTransferParams(tp)
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
