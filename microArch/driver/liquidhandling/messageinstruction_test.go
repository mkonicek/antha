// messageinstruction_test
package liquidhandling

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
)

// Some potentially confusing logic in need of testing
// 5 test scenarios:
// 1. A prompt with a wait time:
//    - should generate a wait instruction followed by a message instruction
// 2. A prompt with no wait time:
//    - should just generate a message instruction
// 3. A wait time with no prompt:
//    - should generate an empty prompt instruction followed by a wait instruction
// 	  - the generation of the empty prompt is critical to correct generation of the correct wait instruction in trilution
// 4. A magic barrier prompt with no wait time:
//    - should not generate either instruction (the prompt is just there to split the workflow)
// 5. A magic barrier prompt with a wait time:
//    - should generate an error
func TestPreOutput(t *testing.T) {

	type messageTest struct {
		MSG                *MessageInstruction
		ExpectedOutputMSG  *MessageInstruction
		ExpectedOutputWait *WaitInstruction
		ExpectedErr        bool
	}

	waitFor10s := NewWaitInstruction()
	waitFor10s.Time = 10.0

	noPromptWithWait := NewMessageInstruction(
		&wtype.LHInstruction{
			Message:  "",
			WaitTime: 1e10,
		},
	)

	aPrompt := NewMessageInstruction(
		&wtype.LHInstruction{
			Message: "hello everybody",
		},
	)

	aPromptWithWait := NewMessageInstruction(
		&wtype.LHInstruction{
			Message:  "hello everybody",
			WaitTime: 1e10,
		},
	)

	var tests = []messageTest{
		// 1. A prompt with a wait time:
		//    - should generate a wait instruction followed by a message instruction
		{
			MSG: NewMessageInstruction(
				&wtype.LHInstruction{
					Message:  "hello everybody",
					WaitTime: 1e10,
				},
			),
			ExpectedOutputMSG:  aPromptWithWait,
			ExpectedOutputWait: waitFor10s,
		},
		// 2. A prompt with no wait time:
		//    - should just generate a message instruction
		{
			MSG: NewMessageInstruction(
				&wtype.LHInstruction{
					Message:  "hello everybody",
					WaitTime: 0,
				},
			),
			ExpectedOutputMSG:  aPrompt,
			ExpectedOutputWait: nil,
		},
		// 3. A wait time with no prompt:
		//    - should generate an empty prompt instruction followed by a wait instruction
		// 	  - the generation of the empty prompt is critical to correct generation of the correct wait instruction in trilution
		{
			MSG: NewMessageInstruction(
				&wtype.LHInstruction{
					Message:  "",
					WaitTime: 1e10,
				},
			),
			ExpectedOutputMSG:  noPromptWithWait,
			ExpectedOutputWait: waitFor10s,
		},
		// 4. A magic barrier prompt with no wait time:
		//    - should not generate either instruction (the prompt is just there to split the workflow)
		{
			MSG: NewMessageInstruction(
				&wtype.LHInstruction{
					Message:  wtype.MAGICBARRIERPROMPTSTRING,
					WaitTime: 0,
				},
			),
			ExpectedOutputMSG:  nil,
			ExpectedOutputWait: nil,
		},
		// 5. A magic barrier prompt with a wait time:
		//    - should generate an error
		{
			MSG: NewMessageInstruction(
				&wtype.LHInstruction{
					Message:  wtype.MAGICBARRIERPROMPTSTRING,
					WaitTime: 1e10,
				},
			),
			ExpectedErr:        true,
			ExpectedOutputMSG:  nil,
			ExpectedOutputWait: waitFor10s,
		},
	}

	for _, test := range tests {
		messageWait, err := test.MSG.PreOutput()

		if (err != nil) != test.ExpectedErr {
			t.Errorf(
				"expected error does not match error obtained (%v) for test %s", err, text.PrettyPrint(test),
			)
		}

		if !reflect.DeepEqual(messageWait.MessageInstruction, test.ExpectedOutputMSG) {
			t.Errorf(
				"found unexpected diffs between expected message (%+v) and actual result (%+v) for test %s./n", test.ExpectedOutputMSG, messageWait.MessageInstruction, text.PrettyPrint(test),
			)
		}
		if !reflect.DeepEqual(messageWait.WaitInstruction, test.ExpectedOutputWait) {
			t.Errorf(
				"found unexpected diffs between expected wait (%+v) and actual result (%+v) for test %s./n", test.ExpectedOutputWait, messageWait.WaitInstruction, text.PrettyPrint(test),
			)
		}
	}
}
