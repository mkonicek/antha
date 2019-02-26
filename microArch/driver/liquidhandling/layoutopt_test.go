package liquidhandling

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type LayoutOptTest struct {
	Name          string
	Driver        *LayoutOpt
	User          *LayoutOpt
	ExpectedError string
	Expected      *LayoutOpt
}

func (test *LayoutOptTest) expecting(err error) bool {
	if err == nil {
		return test.ExpectedError == ""
	} else {
		return test.ExpectedError == err.Error()
	}
}

func (test *LayoutOptTest) Run(t *testing.T) {
	if got, err := test.Driver.ApplyUserPreferences(test.User); !test.expecting(err) {
		t.Errorf("errors don't match:\ne: %v\ng: %v", test.ExpectedError, err)
	} else if err == nil {
		assert.Equal(t, test.Expected, got, "output of merge didn't match expected")
	}
}

type LayoutOptTests []*LayoutOptTest

func (tests LayoutOptTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestLayoutOpt(t *testing.T) {
	LayoutOptTests{
		{
			Name: "basic example",
			Driver: &LayoutOpt{
				Tipboxes:  Addresses{"a", "b", "c", "d"},
				Tipwastes: Addresses{"e"},
			},
			User: &LayoutOpt{
				Tipboxes: Addresses{"b", "c"},
			},
			Expected: &LayoutOpt{
				Tipboxes:  Addresses{"b", "c"},
				Tipwastes: Addresses{"e"},
			},
		},
		{
			Name: "error",
			Driver: &LayoutOpt{
				Tipboxes:  Addresses{"a", "b", "c", "d"},
				Tipwastes: Addresses{"e"},
			},
			User: &LayoutOpt{
				Tipboxes: Addresses{"b", "c", "the moon"},
			},
			ExpectedError: "invalid layout preferences: cannot place tipboxes at: \"the moon\"",
		},
	}.Run(t)
}
