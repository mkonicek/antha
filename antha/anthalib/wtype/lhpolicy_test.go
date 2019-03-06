package wtype

import (
	"fmt"
	"testing"
)

func TestComponentPolicy(t *testing.T) {
	rs := NewLHPolicyRuleSet()
	pol := MakeTestPolicy()
	cnd := MakeTestCondition()
	rs.AddRule(cnd, pol)
	cmp := makeComponent()

	err := cmp.SetPolicies(rs)

	if err != nil {
		t.Fatal(err.Error())
	}

	r2, err := cmp.GetPolicies()

	if err != nil {
		t.Fatal(err.Error())
	}

	if !rs.IsEqualTo(r2) {
		t.Fatal(fmt.Sprintf("Rule set coming out not same as came in! \n%v \n=/= \n%v", *rs, *r2))
	}

}

func MakeTestCondition() LHPolicyRule {
	r := NewLHPolicyRule("verylowvolumerule")
	if err := r.AddNumericConditionOn("VOLUME", 0.0, 1.0); err != nil {
		panic(err)
	}
	return r
}

func MakeTestPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 27)
	// don't set this here -- use defaultpipette speed or there will be inconsistencies
	// defaultpolicy["ASP_SPEED"] = 3.0
	// defaultpolicy["DSP_SPEED"] = 3.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MULTI"] = false
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTOFFSET"] = -5.0
	defaultpolicy["BLOWOUTVOLUME"] = 0.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = true
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	// added to diagnose bubble cause
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["POST_MIX_Z"] = 0.5
	defaultpolicy["PRE_MIX_Z"] = 0.5
	//defaultpolicy["ASP_WAIT"] = 1.0
	//defaultpolicy["DSP_WAIT"] = 1.0
	defaultpolicy["PRE_MIX_VOLUME"] = 10.0
	defaultpolicy["POST_MIX_VOLUME"] = 10.0

	return defaultpolicy
}

func TestPolicyOption(t *testing.T) {
	rs := NewLHPolicyRuleSet()
	err := rs.SetOption("USE_DRIVER_TIP_TRACKING", false)

	if err != nil {
		t.Errorf("Expected nil return when setting USE_DRIVER_TIP_TRACKING, got %v", err)
	}

	err = rs.SetOption("USE_DRIVER_TIP_TRACKING", true)

	if err != nil {
		t.Errorf("Expected nil return when setting USE_DRIVER_TIP_TRACKING, got %v", err)
	}

	err = rs.SetOption("USE_DRIVER_TIP_TRACKING", 3)

	if err == nil {
		t.Errorf("Trying to set a boolean value to an int should fail but did not")
	}
}
