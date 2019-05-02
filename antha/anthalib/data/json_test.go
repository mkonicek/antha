package data

import (
	"encoding/json"
	"fmt"
	"testing"
)

//antha interop

func TestTableJSON(t *testing.T) {
	// Antha parameter values must currently roundtrip via json.
	tab := NewTable(
		makeArrowSeries("float_measure", []float64{0, 2, 1000}, nil),
		makeArrowSeries("int_measure", []int{1, 3, 100}, nil),
		makeArrowSeries("time", []TimestampMillis{TimestampMillis(0), TimestampMillis(0), TimestampMillis(1)}, nil),
		makeArrowSeries("label", []string{"", "abcdef", "abcd"}, []bool{false, true, true}),
	)
	js, err := json.Marshal(tab)
	if err != nil {
		t.Fatal(err)
	}
	un := new(Table)
	err = json.Unmarshal(js, un)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, tab, un, fmt.Sprintf("roundtrip: %s", js))
}

func BenchmarkJsonMarshal(b *testing.B) {
	b.StopTimer()
	tab := genericInputTable()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, err := json.Marshal(tab)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	b.StopTimer()
	tab := genericInputTable()
	js, err := json.Marshal(tab)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		un := new(Table)
		err = json.Unmarshal(js, un)
		if err != nil {
			b.Fatal(err)
		}
	}
}
