package main

import (
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/data"
)

func example() {
	// Creating a bounded table from scalar slices.
	column1 := data.Must().NewSeriesFromSlice("measure", []int64{10, 10, 30, 0, 5}, []bool{true, true, true, false, true})
	column2 := data.Must().NewSeriesFromSlice("label", []string{"", "aa", "xx", "aa", ""}, nil)
	tab := data.NewTable([]*data.Series{column1, column2})

	// just print Table as a whole.
	fmt.Println("input\n", tab.ToRows())

	// iterate over the entire Table.
	for record := range tab.IterAll() {
		// s := &struct{ x int64}
		// err := record.ToStruct(&s)
		// s.x
		m, _ := record.Observation("measure") //nolint
		if m.IsNull() {
			fmt.Println("measure=null at index", record.Index)
		} else {
			fmt.Println("measure=", m.MustInt64())
		}
	}

	// subset of rows
	fmt.Println("tab.Slice(2,4)\n", tab.Slice(2, 4).ToRows())

	// produce a new Table by filtering
	smallerTab := tab.Must().Filter().On("label").Interface(data.Eq("aa"))
	fmt.Println("after filter\n", smallerTab.ToRows())

	mult := func(r data.Row) interface{} {
		m, _ := r.Observation("measure") //nolint
		if m.IsNull() {
			return nil
		}
		return float64(m.MustInt64()) * float64(2.5)
	}
	extended := tab.Must().
		Extend("multiplied").By(mult, reflect.TypeOf(float64(0)))
	extendedAndFiltered := extended.Must().Filter().On("multiplied").Interface(data.Eq(25))
	fmt.Println("extended and filtered\n", extendedAndFiltered.ToRows())

	// equivalent extension using static types
	projected := tab.
		Must().Convert("measure", reflect.TypeOf(float64(0))).
		Must().Extend("multiplied").On("measure").Float64(
		func(vals ...float64) float64 {
			return vals[0] * 2.5
		}).
		Must().Project("label", "multiplied")
	fmt.Println("extended and projected\n", projected.ToRows())

	alternateProjected := extended.ProjectAllBut("measure")
	fmt.Printf("alternateProjected.Equal(projected): %v\n", alternateProjected.Equal(projected))
}

func main() {
	example()
}
