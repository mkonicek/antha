package workflow

import (
	"testing"
)

func TestValidateBasicId(t *testing.T) {
	type testCase struct {
		sid            BasicId
		shouldValidate bool
	}
	testCases := []testCase{
		{BasicId(""), true},
		{BasicId("abc"), true},
		{BasicId("ABC"), true},
		{BasicId("a b"), false},
		{BasicId("a"), true},
		{BasicId("a0"), true},
		{BasicId("0a"), true},
		{BasicId("*"), false},
		{BasicId("åbc"), false},
		{BasicId("🤡"), false},
	}

	for _, testCase := range testCases {
		err := testCase.sid.Validate(true)
		if err != nil && testCase.shouldValidate {
			t.Errorf("BasicId %v expected to be valid but failed validation", testCase.sid)
		} else if err == nil && !testCase.shouldValidate {
			t.Errorf("BasicId %v expected to be invalid but passed validation", testCase.sid)
		}
	}
}
