package data

import (
	"testing"
	// TODO "github.com/stretchr/testify/assert"
)

func TestFmt(t *testing.T) {
	tab := NewTable(
		Must().NewSeriesFromSlice("measure", []int64{1, 1000}, nil),
		Must().NewSeriesFromSlice("label", []string{"abcdef", "abcd"}, nil),
	)
	formatted := tab.ToRows().String()
	if formatted != `2 Row(s):
| |measure| label|
| |  int64|string|
------------------
|0|      1|abcdef|
|1|   1000|  abcd|
` {
		t.Errorf("fmt: %s", formatted)
	}
}

func TestFmtEmpty(t *testing.T) {
	tab := NewTable(
		Must().NewSeriesFromSlice("A", []float64{}, nil),
	)
	formatted := tab.ToRows().String()
	expected := `0 Row(s):
||      A|
||float64|
----------
`
	if formatted != expected {
		t.Errorf("fmt: %s", formatted)
	}
}
