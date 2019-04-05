package workflow

import (
	"testing"
)

func TestValidateSimpleId(t *testing.T) {
	type testCase struct {
		sid            SimpleId
		shouldValidate bool
	}
	testCases := []testCase{
		{SimpleId(""), true},
		{SimpleId("abc"), true},
		{SimpleId("ABC"), true},
		{SimpleId("a b"), false},
		{SimpleId("a"), true},
		{SimpleId("a0"), true},
		{SimpleId("0a"), false},
		{SimpleId("*"), false},
		{SimpleId("åbc"), false},
		{SimpleId("🤡"), false},
	}

	for _, testCase := range testCases {
		err := testCase.sid.Validate()
		if err != nil && testCase.shouldValidate {
			t.Errorf("SimpleId %v expected to be valid but failed validation", testCase.sid)
		} else if err == nil && !testCase.shouldValidate {
			t.Errorf("SimpleId %v expected to be invalid but passed validation", testCase.sid)
		}
	}
}
