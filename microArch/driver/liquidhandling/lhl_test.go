package liquidhandling

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func getTestBlowout(ch *wtype.LHChannelParameter, multi int, tipType string) RobotInstruction {
	bi := NewBlowInstruction()
	bi.Multi = multi
	bi.TipType = tipType
	for i := 0; i < multi; i++ {
		bi.What = append(bi.What, "soup")
		bi.PltTo = append(bi.PltTo, "position_4")
		bi.WellTo = append(bi.WellTo, "A1")
		bi.Volume = append(bi.Volume, wunit.NewVolume(10.0, "ul"))
		bi.TPlateType = append(bi.TPlateType, "pcrplate_skirted_riser40")
		bi.TVolume = append(bi.TVolume, wunit.NewVolume(5.0, "ul"))
	}
	bi.Prms = ch
	bi.Head = ch.Head
	return bi
}

func getTestSuck(ch *wtype.LHChannelParameter, multi int, tipType string) RobotInstruction {
	ret := NewSuckInstruction()
	ret.Multi = multi
	ret.TipType = tipType
	for i := 0; i < multi; i++ {
		ret.What = append(ret.What, "soup")
		ret.PltFrom = append(ret.PltFrom, "position_4")
		ret.WellFrom = append(ret.WellFrom, "A1")
		ret.Volume = append(ret.Volume, wunit.NewVolume(10.0, "ul"))
		ret.FPlateType = append(ret.FPlateType, "pcrplate_skirted_riser40")
		ret.FVolume = append(ret.FVolume, wunit.NewVolume(20.0, "ul"))
	}
	ret.Prms = ch
	ret.Head = ch.Head
	return ret
}

func TestBlowMixing(t *testing.T) {

	tenUl := wunit.NewVolume(10.0, "ul")

	tests := []*PolicyTest{
		{
			Name: "blow no tip change single channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
					},
				},
			},
			Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MIX,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 5, //the Mix
					Values: map[string]interface{}{
						"CYCLES": []int{5},
						"VOLUME": []wunit.Volume{tenUl},
					},
				},
			},
		},
		{
			Name: "blow no tip change eight channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
					},
				},
			},
			Instruction:          getTestBlowout(getLVConfig(), 8, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MIX,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 5, //the Mix
					Values: map[string]interface{}{
						"CYCLES": []int{5, 5, 5, 5, 5, 5, 5, 5},
						"VOLUME": []wunit.Volume{tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestBlowNoMixing(t *testing.T) {

	tests := []*PolicyTest{
		{
			Name: "z-offset single channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{
						"DSPZOFFSET": 5.4,
					},
				},
			},
			Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[string]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 1,
					Values: map[string]interface{}{
						"DRIVE": "Z",
						"SPEED": defaultZSpeed,
					},
				},
				{
					Instruction: 2, //the move before the dispense
					Values: map[string]interface{}{
						"OFFSETZ": []float64{5.4},
					},
				},
			},
		},
		{
			Name: "entry speed single channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{
						"DSPENTRYSPEED": 50.0,
					},
				},
			},
			Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,SDS,MOV,DSP,SDS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[string]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 1,
					Values: map[string]interface{}{
						"DRIVE": "Z",
						"SPEED": defaultZSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[string]interface{}{
						"DRIVE": "Z",
						"SPEED": 50.0,
					},
				},
				{
					Instruction: 4,
					Values: map[string]interface{}{
						"OFFSETZ": []float64{defaultZOffset},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestDSPPipetSpeed(t *testing.T) {

	tests := []*PolicyTest{
		{
			Name: "blow pipette speed unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{},
				},
			},
			Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[string]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		/*	{
				Name: "blow pipette speed OK",
				Rules: []*Rule{
					{
						Name: "soup",
						Conditions: []Condition{
							&CategoryCondition{
								Attribute: "LIQUIDCLASS",
								Value:     "soup",
							},
						},
						Policy: map[string]interface{}{
							"DSPSPEED": 1.5,
						},
					},
				},
				Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
				Robot:                nil,
				ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
				Assertions: []*InstructionAssertion{
					{
						Instruction: 0,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
					{
						Instruction: 3,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   1.5,
						},
					},
					{
						Instruction: 5,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
				},
			},
			{
				Name: "blow pipette speed too large, made safe",
				Rules: []*Rule{
					{
						Name: "soup",
						Conditions: []Condition{
							&CategoryCondition{
								Attribute: "LIQUIDCLASS",
								Value:     "soup",
							},
						},
						Policy: map[string]interface{}{
							"DSPSPEED":             LVMaxRate + 1.0,
							"OVERRIDEPIPETTESPEED": true,
						},
					},
				},
				Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
				Robot:                nil,
				ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
				Assertions: []*InstructionAssertion{
					{
						Instruction: 0,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
					{
						Instruction: 3,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   LVMaxRate,
						},
					},
					{
						Instruction: 5,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
				},
			},
			{
				Name: "blow pipette speed too small, made safe",
				Rules: []*Rule{
					{
						Name: "soup",
						Conditions: []Condition{
							&CategoryCondition{
								Attribute: "LIQUIDCLASS",
								Value:     "soup",
							},
						},
						Policy: map[string]interface{}{
							"DSPSPEED":             LVMinRate * 0.5,
							"OVERRIDEPIPETTESPEED": true,
						},
					},
				},
				Instruction:          getTestBlowout(getLVConfig(), 1, "Gilson20"),
				Robot:                nil,
				ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
				Assertions: []*InstructionAssertion{
					{
						Instruction: 0,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
					{
						Instruction: 3,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   LVMinRate,
						},
					},
					{
						Instruction: 5,
						Values: map[string]interface{}{
							"HEAD":    1,
							"CHANNEL": -1,
							"SPEED":   defaultPipetteSpeed,
						},
					},
				},
			},
			{
				Name: "blow pipette speed too large, error",
				Rules: []*Rule{
					{
						Name: "soup",
						Conditions: []Condition{
							&CategoryCondition{
								Attribute: "LIQUIDCLASS",
								Value:     "soup",
							},
						},
						Policy: map[string]interface{}{
							"DSPSPEED":             4.75,
							"OVERRIDEPIPETTESPEED": false,
						},
					},
				},
				Instruction: getTestBlowout(getLVConfig(), 1, "Gilson20"),
				Error:       "setting pipette dispense speed: value 4.750000 out of range 0.022500 - 3.750000",
			},
			{
				Name: "blow pipette speed too small, error",
				Rules: []*Rule{
					{
						Name: "soup",
						Conditions: []Condition{
							&CategoryCondition{
								Attribute: "LIQUIDCLASS",
								Value:     "soup",
							},
						},
						Policy: map[string]interface{}{
							"DSPSPEED":             0.01,
							"OVERRIDEPIPETTESPEED": false,
						},
					},
				},
				Instruction: getTestBlowout(getLVConfig(), 1, "Gilson20"),
				Error:       "setting pipette dispense speed: value 0.010000 out of range 0.022500 - 3.750000",
			},*/
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestASPPipetSpeed(t *testing.T) {

	tests := []*PolicyTest{
		{
			Name: "suck pipette speed unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[string]interface{}{},
				},
			},
			Instruction:          getTestSuck(getLVConfig(), 1, "Gilson20"),
			Robot:                nil,
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[string]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		/*		{
					Name: "suck pipette speed OK",
					Rules: []*Rule{
						{
							Name: "soup",
							Conditions: []Condition{
								&CategoryCondition{
									Attribute: "LIQUIDCLASS",
									Value:     "soup",
								},
							},
							Policy: map[string]interface{}{
								"ASPSPEED": 1.5,
							},
						},
					},
					Instruction:          getTestSuck(getLVConfig(), 1, "Gilson20"),
					Robot:                nil,
					ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
					Assertions: []*InstructionAssertion{
						{
							Instruction: 0,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
						{
							Instruction: 3,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   1.5,
							},
						},
						{
							Instruction: 5,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
					},
				},
				{
					Name: "suck pipette speed too large, made safe",
					Rules: []*Rule{
						{
							Name: "soup",
							Conditions: []Condition{
								&CategoryCondition{
									Attribute: "LIQUIDCLASS",
									Value:     "soup",
								},
							},
							Policy: map[string]interface{}{
								"ASPSPEED":             LVMaxRate + 1.0,
								"OVERRIDEPIPETTESPEED": true,
							},
						},
					},
					Instruction:          getTestSuck(getLVConfig(), 1, "Gilson20"),
					Robot:                nil,
					ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
					Assertions: []*InstructionAssertion{
						{
							Instruction: 0,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
						{
							Instruction: 3,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   LVMaxRate,
							},
						},
						{
							Instruction: 5,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
					},
				},
				{
					Name: "suck pipette speed too small, made safe",
					Rules: []*Rule{
						{
							Name: "soup",
							Conditions: []Condition{
								&CategoryCondition{
									Attribute: "LIQUIDCLASS",
									Value:     "soup",
								},
							},
							Policy: map[string]interface{}{
								"ASPSPEED":             LVMinRate * 0.5,
								"OVERRIDEPIPETTESPEED": true,
							},
						},
					},
					Instruction:          getTestSuck(getLVConfig(), 1, "Gilson20"),
					Robot:                nil,
					ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
					Assertions: []*InstructionAssertion{
						{
							Instruction: 0,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
						{
							Instruction: 3,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   LVMinRate,
							},
						},
						{
							Instruction: 5,
							Values: map[string]interface{}{
								"HEAD":    1,
								"CHANNEL": -1,
								"SPEED":   defaultPipetteSpeed,
							},
						},
					},
				},
				{
					Name: "suck pipette speed too large, error",
					Rules: []*Rule{
						{
							Name: "soup",
							Conditions: []Condition{
								&CategoryCondition{
									Attribute: "LIQUIDCLASS",
									Value:     "soup",
								},
							},
							Policy: map[string]interface{}{
								"ASPSPEED":             4.75,
								"OVERRIDEPIPETTESPEED": false,
							},
						},
					},
					Instruction: getTestSuck(getLVConfig(), 1, "Gilson20"),
					Error:       "setting pipette aspirate speed: value 4.750000 out of range 0.022500 - 3.750000",
				},
				{
					Name: "suck pipette speed too small, error",
					Rules: []*Rule{
						{
							Name: "soup",
							Conditions: []Condition{
								&CategoryCondition{
									Attribute: "LIQUIDCLASS",
									Value:     "soup",
								},
							},
							Policy: map[string]interface{}{
								"ASPSPEED":             0.01,
								"OVERRIDEPIPETTESPEED": false,
							},
						},
					},
					Instruction: getTestSuck(getLVConfig(), 1, "Gilson20"),
					Error:       "setting pipette aspirate speed: value 0.010000 out of range 0.022500 - 3.750000",
				},*/
	}

	for _, test := range tests {
		test.Run(t)
	}
}
