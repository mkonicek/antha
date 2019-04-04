package workflow

import (
	"testing"
)

func TestValidateJobId(t *testing.T) {
	type testCase struct {
		jobID          JobId
		shouldValidate bool
	}
	testCases := []testCase{
		{JobId("abc"), true},
		{JobId("ABC"), true},
		{JobId("a b"), true},
		{JobId("a"), true},
		{JobId("a0"), true},
		{JobId("0a"), false},
		{JobId("*"), false},
		{JobId("åbc"), false},
		{JobId("🤡"), false},
	}

	for _, testCase := range testCases {
		err := testCase.jobID.Validate()
		if err != nil && testCase.shouldValidate {
			t.Errorf("JobId %v expected to be valid but failed validation", testCase.jobID)
		} else if err == nil && !testCase.shouldValidate {
			t.Errorf("JobId %v expected to be invalid but passed validation", testCase.jobID)
		}
	}
}
