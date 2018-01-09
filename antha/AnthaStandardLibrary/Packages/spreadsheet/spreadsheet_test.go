package spreadsheet

import (
	"testing"
)

// Parses an a1 style excel cell coordinate into ints for row and column for use by plotinum library
// note that 1 is subtracted from the column number in accordance with the go convention of counting from 0

type coordinatetest struct {
	a1format string
	row      int
	col      int
}

var coordinatetests = []coordinatetest{
	{a1format: "a1",
		row: 0,
		col: 0,
	},
	{a1format: "b1",
		row: 1,
		col: 0,
	},
	{a1format: "aa1",
		row: 26,
		col: 0,
	},
	{a1format: "a2",
		row: 0,
		col: 1,
	},
}

func TestA1formattorowcolumn(t *testing.T) {
	for _, test := range coordinatetests {
		r, c, _ := A1FormatToRowColumn(test.a1format)
		if c != test.col {
			t.Error(
				"For", test.a1format, "\n",
				"expected", test.col, "\n",
				"got", c, "\n",
			)
		}
		if r != test.row {
			t.Error(
				"For", test.a1format, "\n",
				"expected", test.row, "\n",
				"got", r, "\n",
			)
		}
	}
}
