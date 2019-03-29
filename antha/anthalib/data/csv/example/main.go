package main

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/data/csv"
	"github.com/pkg/errors"
)

func Example() {
	// read a Table from CSV
	table, err := csv.TableFromFile("sample.csv")
	if err != nil {
		panic(errors.Wrapf(err, "read table"))
	}

	// print top 10 rows
	fmt.Println("input\n", table.Head(10).ToRows())
}

func main() {
	Example()
}
