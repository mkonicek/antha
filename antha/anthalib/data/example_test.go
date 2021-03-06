package data

import (
	"fmt"
	"reflect"
	"strings"
)

// Example dataset
var pirateBooty *Table

func ExampleTable() {
	// create a table
	pirateBooty = NewTable(
		Must().NewSeriesFromSlice("Name", []string{"doubloon", "grog", "cutlass", "chest"}, nil),
		Must().NewSeriesFromSlice("Price", []float64{1.0, 0, 5.5, 600.0}, []bool{true, false, true, true}),
		Must().NewSeriesFromSlice("Quantity", []int64{1200, 44, 30, 2}, []bool{true, true, true, true}),
	)

	fmt.Println(pirateBooty.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleNewTable() {
	// create a table
	pirateBooty = NewTable(
		Must().NewSeriesFromSlice("Name", []string{"doubloon", "grog", "cutlass", "chest"}, nil),
		Must().NewSeriesFromSlice("Price", []float64{1.0, 0, 5.5, 600.0}, []bool{true, false, true, true}),
		Must().NewSeriesFromSlice("Quantity", []int64{1200, 44, 30, 2}, []bool{true, true, true, true}),
	)

	fmt.Println(pirateBooty.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleTable_Size() {
	fmt.Println(pirateBooty.Size())
	// Output: 4
}

func ExampleTable_Schema() {
	fmt.Println(pirateBooty.Schema())
	// Output:
	// Name, string
	// Price, float64
	// Quantity, int64
}

func ExampleTable_ToRows() {
	rows := pirateBooty.ToRows()
	fmt.Println(rows)
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleTable_Sort() {
	// in ascending order of Name.
	byNameAsc, _ := pirateBooty.Sort(Key{{"Name", true}})
	fmt.Println(byNameAsc.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|   chest|    600|       2|
	// |1| cutlass|    5.5|      30|
	// |2|doubloon|      1|    1200|
	// |3|    grog|  <nil>|      44|
}

func ExampleTable_SortByFunc() {
	// in ascending order of length of Name.
	byNameLenAsc, _ := pirateBooty.SortByFunc(func(r1 Row, r2 Row) bool {
		return len(r1.ValueAt(0).MustString()) < len(r2.ValueAt(0).MustString())
	})
	fmt.Println(byNameLenAsc.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|    grog|  <nil>|      44|
	// |1|   chest|    600|       2|
	// |2| cutlass|    5.5|      30|
	// |3|doubloon|      1|    1200|
}

func ExampleTable_Project() {
	// two columns by name
	projected, _ := pirateBooty.Project("Quantity", "Name")
	fmt.Println(projected.ToRows())
	// Output: 4 Row(s):
	// | |Quantity|    Name|
	// | |   int64|  string|
	// ---------------------
	// |0|    1200|doubloon|
	// |1|      44|    grog|
	// |2|      30| cutlass|
	// |3|       2|   chest|
}

func ExampleTable_ProjectAllBut() {
	// drop some columns.
	justQuantity := pirateBooty.ProjectAllBut("Price", "Name", "No such column")
	fmt.Println(justQuantity.ToRows())
	// Output: 4 Row(s):
	// | |Quantity|
	// | |   int64|
	// ------------
	// |0|    1200|
	// |1|      44|
	// |2|      30|
	// |3|       2|
}

func ExampleTable_Rename() {
	// rename a column.
	renameQuantity := pirateBooty.Rename("Quantity", "QuantityBeforeSinking")
	fmt.Println(renameQuantity.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|QuantityBeforeSinking|
	// | |  string|float64|                int64|
	// ------------------------------------------
	// |0|doubloon|      1|                 1200|
	// |1|    grog|  <nil>|                   44|
	// |2| cutlass|    5.5|                   30|
	// |3|   chest|    600|                    2|
}

func ExampleTable_Head() {
	// first 2 rows
	top := pirateBooty.Head(2)
	fmt.Println(top.ToRows())
	// Output: 2 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
}

func ExampleTable_Slice() {
	// last 2 rows
	bottom := pirateBooty.Slice(2, 4)
	fmt.Println(bottom.ToRows())
	// Output: 2 Row(s):
	// | |   Name|  Price|Quantity|
	// | | string|float64|   int64|
	// ----------------------------
	// |0|cutlass|    5.5|      30|
	// |1|  chest|    600|       2|
}

func ExampleTable_Filter_staticType() {
	// notice that rows with Price=nil are not passed to the filter.
	priceLessThan100, _ := pirateBooty.Filter().On("Price").
		Float64(func(v ...float64) bool {
			price := v[0]
			return price < 100
		})
	fmt.Println(priceLessThan100.ToRows())
	// Output: 2 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1| cutlass|    5.5|      30|
}

func ExampleTable_Filter_dynamicType() {
	// any column value can be filtered on.
	startsWithC := pirateBooty.Filter().
		By(func(r Row) bool {
			name := r.ValueAt(0)
			return !name.IsNull() && strings.HasPrefix(name.MustString(), "c")
		})
	fmt.Println(startsWithC.ToRows())
	// Output: 2 Row(s):
	// | |   Name|  Price|Quantity|
	// | | string|float64|   int64|
	// ----------------------------
	// |0|cutlass|    5.5|      30|
	// |1|  chest|    600|       2|
}

func ExampleTable_Filter_equal() {
	// equality filter is provided as a function.
	grog, _ := pirateBooty.Filter().On("Name").
		Interface(Eq("grog"))
	fmt.Println(grog.ToRows())
	// Output: 1 Row(s):
	// | |  Name|  Price|Quantity|
	// | |string|float64|   int64|
	// ---------------------------
	// |0|  grog|  <nil>|      44|
}

func ExampleTable_Distinct_onSpecifiedColumnsByEq() {
	// distinct by equality of the specified columns.
	distinctByName, _ := pirateBooty.Distinct().On("Name")
	fmt.Println(distinctByName.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleTable_Convert() {
	// convert column values from int to float
	converted, _ := pirateBooty.Convert("Quantity", reflect.TypeOf(float64(0)))
	fmt.Println(converted.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64| float64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleTable_Extend_constant() {
	// add a constant column value.
	withSource := pirateBooty.Extend("Source").
		Constant("BOOTY")
	fmt.Println(withSource.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|Source|
	// | |  string|float64|   int64|string|
	// ------------------------------------
	// |0|doubloon|      1|    1200| BOOTY|
	// |1|    grog|  <nil>|      44| BOOTY|
	// |2| cutlass|    5.5|      30| BOOTY|
	// |3|   chest|    600|       2| BOOTY|
}

func ExampleTable_Extend_dynamicTypes() {
	// calculate new column value on the whole row using dynamic types.
	totals := pirateBooty.Extend("Total").
		By(func(r Row) interface{} {
			q, _ := r.Value("Quantity")
			p, _ := r.Value("Price")
			if p.IsNull() {
				return nil
			}
			return p.MustFloat64() * float64(q.MustInt64())
		}, reflect.TypeOf(float64(0)))
	fmt.Println(totals.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|  Total|
	// | |  string|float64|   int64|float64|
	// -------------------------------------
	// |0|doubloon|      1|    1200|   1200|
	// |1|    grog|  <nil>|      44|  <nil>|
	// |2| cutlass|    5.5|      30|    165|
	// |3|   chest|    600|       2|   1200|
}

func ExampleTable_Extend_onSelectedColsDynamicTypes() {
	// calculate new column value on selected columns using dynamic types.
	salePrices, _ := pirateBooty.Extend("Reduced Price").On("Price").
		Interface(func(v ...interface{}) interface{} {
			// 25% off
			if v[0] == nil {
				return nil
			}
			return v[0].(float64) * 0.75
		}, reflect.TypeOf(float64(0)))
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|Reduced Price|
	// | |  string|float64|   int64|      float64|
	// -------------------------------------------
	// |0|doubloon|      1|    1200|         0.75|
	// |1|    grog|  <nil>|      44|        <nil>|
	// |2| cutlass|    5.5|      30|        4.125|
	// |3|   chest|    600|       2|          450|
}

func ExampleTable_Extend_onSelectedColsMixedTypes() {
	// calculate new column value on selected columns using dynamic types for inputs and static type for output.
	salePrices, _ := pirateBooty.Extend("Reduced Price").On("Price").
		InterfaceFloat64(func(v ...interface{}) (float64, bool) {
			// 25% off
			if v[0] == nil {
				return float64(0), false
			}
			return v[0].(float64) * 0.75, true
		})
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|Reduced Price|
	// | |  string|float64|   int64|      float64|
	// -------------------------------------------
	// |0|doubloon|      1|    1200|         0.75|
	// |1|    grog|  <nil>|      44|        <nil>|
	// |2| cutlass|    5.5|      30|        4.125|
	// |3|   chest|    600|       2|          450|
}

func ExampleTable_Extend_onSelectedColsStaticType() {
	// calculate new column value on selected columns using static type.
	salePrices, _ := pirateBooty.Extend("Reduced Price").On("Price").
		Float64(func(v ...float64) (float64, bool) {
			// 25% off
			return v[0] * 0.75, true
		})
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|Reduced Price|
	// | |  string|float64|   int64|      float64|
	// -------------------------------------------
	// |0|doubloon|      1|    1200|         0.75|
	// |1|    grog|  <nil>|      44|        <nil>|
	// |2| cutlass|    5.5|      30|        4.125|
	// |3|   chest|    600|       2|          450|
}

func ExampleTable_Update_constant() {
	// add a constant column value.
	withSource, _ := pirateBooty.Update("Quantity").
		Constant(int64(0))
	fmt.Println(withSource.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|       0|
	// |1|    grog|  <nil>|       0|
	// |2| cutlass|    5.5|       0|
	// |3|   chest|    600|       0|
}

func ExampleTable_Update_dynamicTypes() {
	// calculate updated column value on the whole row using dynamic types.
	totals, _ := pirateBooty.Update("Quantity").
		By(func(r Row) interface{} {
			q, _ := r.Value("Quantity")
			p, _ := r.Value("Price")
			if p.IsNull() {
				return 0
			}
			return q.MustInt64()
		})
	fmt.Println(totals.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|       0|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
}

func ExampleTable_Update_onSelectedColsDynamicTypes() {
	// calculate updated column value on selected columns using dynamic types.
	salePrices, _ := pirateBooty.Update("Price").On("Price").
		Interface(func(v ...interface{}) interface{} {
			// 25% off
			if v[0] == nil {
				return nil
			}
			return v[0].(float64) * 0.75
		})
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|   0.75|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|  4.125|      30|
	// |3|   chest|    450|       2|
}

func ExampleTable_Update_onSelectedColsMixedTypes() {
	// calculate new column value on selected columns using dynamic types for inputs and static type for output.
	salePrices, _ := pirateBooty.Update("Price").On("Price").
		InterfaceFloat64(func(v ...interface{}) (float64, bool) {
			// 25% off
			if v[0] == nil {
				return float64(0), false
			}
			return v[0].(float64) * 0.75, true
		})
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|   0.75|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|  4.125|      30|
	// |3|   chest|    450|       2|
}

func ExampleTable_Update_onSelectedColsStaticType() {
	// calculate new column value on selected columns using static type.
	salePrices, _ := pirateBooty.Update("Price").On("Price").
		Float64(func(v ...float64) (float64, bool) {
			// 25% off
			return v[0] * 0.75, true
		})
	fmt.Println(salePrices.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|   0.75|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|  4.125|      30|
	// |3|   chest|    450|       2|
}

func ExampleTable_Pivot() {
	// create a table
	pirateBootyNarrow := NewTable(
		Must().NewSeriesFromSlice("Name", []string{"doubloon", "doubloon", "grog", "cutlass", "cutlass", "chest", "chest"}, nil),
		Must().NewSeriesFromSlice("PropertyName", []string{"Price", "Quantity", "Quantity", "Price", "Quantity", "Price", "Quantity"}, nil),
		Must().NewSeriesFromSlice("PropertyValue", []float64{1.0, 1200, 44, 5.5, 30, 600.0, 2}, nil),
	)

	// make a wide table from pirateBootyNarrow.
	pirateBootyWide, _ := pirateBootyNarrow.Pivot().Key("Name").Columns("PropertyName", "PropertyValue")
	fmt.Println(pirateBootyWide.ToRows())
	// Output: 4 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64| float64|
	// -----------------------------
	// |0|   chest|    600|       2|
	// |1| cutlass|    5.5|      30|
	// |2|doubloon|      1|    1200|
	// |3|    grog|  <nil>|      44|
}

func ExampleTable_Join_naturalInner() {
	// create another table to join
	currency := NewTable(
		Must().NewSeriesFromSlice("Name", []string{"doubloon", "piastre"}, nil),
		Must().NewSeriesFromSlice("ExchangeRate", []float64{0.5, 2.3}, nil),
	)

	// natural inner join
	currencyBooty, _ := pirateBooty.Join().NaturalInner(currency)
	fmt.Println(currencyBooty.ToRows())
	// Output: 1 Row(s):
	// | |    Name|  Price|Quantity|    Name|ExchangeRate|
	// | |  string|float64|   int64|  string|     float64|
	// ---------------------------------------------------
	// |0|doubloon|      1|    1200|doubloon|         0.5|
}

func ExampleTable_Join_inner() {
	// create another table to join
	currency := NewTable(
		Must().NewSeriesFromSlice("CurrencyName", []string{"doubloon", "piastre"}, nil),
		Must().NewSeriesFromSlice("CurrencyExchangeRate", []float64{0.5, 2.3}, nil),
	)

	// inner join
	currencyBooty, _ := pirateBooty.Join().On("Name").Inner(currency, "CurrencyName")
	fmt.Println(currencyBooty.ToRows())
	// Output: 1 Row(s):
	// | |    Name|  Price|Quantity|CurrencyName|CurrencyExchangeRate|
	// | |  string|float64|   int64|      string|             float64|
	// ---------------------------------------------------------------
	// |0|doubloon|      1|    1200|    doubloon|                 0.5|
}

func ExampleTable_Append() {
	morePirateBooty, _ := pirateBooty.Append(NewTable(
		Must().NewSeriesFromSlice("Name", []string{"piastre"}, nil),
		Must().NewSeriesFromSlice("Price", []float64{1.5}, nil),
		Must().NewSeriesFromSlice("Quantity", []int64{1000}, nil),
	))
	fmt.Println(morePirateBooty.ToRows())
	// Output: 5 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
	// |4| piastre|    1.5|    1000|
}

func ExampleAppendMany() {
	morePirateBooty, _ := AppendMany(pirateBooty, NewTable(
		Must().NewSeriesFromSlice("Name", []string{"piastre"}, nil),
		Must().NewSeriesFromSlice("Price", []float64{1.5}, nil),
		Must().NewSeriesFromSlice("Quantity", []int64{1000}, nil),
	), NewTable(
		Must().NewSeriesFromSlice("Name", []string{"guinea"}, nil),
		Must().NewSeriesFromSlice("Price", []float64{0.5}, nil),
		Must().NewSeriesFromSlice("Quantity", []int64{100}, nil),
	))
	fmt.Println(morePirateBooty.ToRows())
	// Output: 6 Row(s):
	// | |    Name|  Price|Quantity|
	// | |  string|float64|   int64|
	// -----------------------------
	// |0|doubloon|      1|    1200|
	// |1|    grog|  <nil>|      44|
	// |2| cutlass|    5.5|      30|
	// |3|   chest|    600|       2|
	// |4| piastre|    1.5|    1000|
	// |5|  guinea|    0.5|     100|
}

func ExampleTable_Foreach_wholerow() {
	nonNullCount := 0
	pirateBooty.Foreach().By(func(r Row) {
		for _, v := range r.Values() {
			if !v.IsNull() {
				nonNullCount++
			}
		}
	})
	fmt.Printf("Non null count: %d\n", nonNullCount)
	// Output: Non null count: 11
}

func ExampleTable_Foreach_generic() {
	cost := 0.
	_ = pirateBooty.Foreach().On("Price", "Quantity").Interface(func(v ...interface{}) {
		if v[0] == nil {
			return
		}
		cost += v[0].(float64) * float64(v[1].(int64))
	})
	fmt.Printf("Total cost: %f\n", cost)
	// Output: Total cost: 2565.000000
}

func ExampleTable_Foreach_specialized() {
	quantity := int64(0)
	_ = pirateBooty.Foreach().On("Quantity").Int64(func(v ...int64) {
		quantity += v[0]
	})
	fmt.Printf("Total quantity: %d\n", quantity)
	// Output: Total quantity: 1276
}

func ExampleTable_ForeachKey() {
	_ = pirateBooty.ForeachKey("Name").By(func(key Row, partition *Table) {
		fmt.Printf("Key: %v\n", key)
		fmt.Printf("Partition: %v\n", partition.ToRows())
	})
	// Output: Key: 1 Row(s):
	// | |  Name|
	// | |string|
	// ----------
	// |0|  grog|
	//
	// Partition: 1 Row(s):
	// | |  Price|Quantity|
	// | |float64|   int64|
	// --------------------
	// |0|  <nil>|      44|
	//
	// Key: 1 Row(s):
	// | |    Name|
	// | |  string|
	// ------------
	// |1|doubloon|
	//
	// Partition: 1 Row(s):
	// | |  Price|Quantity|
	// | |float64|   int64|
	// --------------------
	// |0|      1|    1200|
	//
	// Key: 1 Row(s):
	// | |   Name|
	// | | string|
	// -----------
	// |2|cutlass|
	//
	// Partition: 1 Row(s):
	// | |  Price|Quantity|
	// | |float64|   int64|
	// --------------------
	// |0|    5.5|      30|
	//
	// Key: 1 Row(s):
	// | |  Name|
	// | |string|
	// ----------
	// |3| chest|
	//
	// Partition: 1 Row(s):
	// | |  Price|Quantity|
	// | |float64|   int64|
	// --------------------
	// |0|    600|       2|
}

func ExampleTable_ToStructs() {
	type price struct {
		Price float64
		Name  string
	}
	structSlice := make([]*price, 0)
	err := pirateBooty.ToStructs(&structSlice)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	for _, p := range structSlice {
		fmt.Printf("%#v\n", p)
	}
	// Output:
	// &data.price{Price:1, Name:"doubloon"}
	// &data.price{Price:0, Name:"grog"}
	// &data.price{Price:5.5, Name:"cutlass"}
	// &data.price{Price:600, Name:"chest"}
}

func ExampleTable_ToStructs_withTags() {
	type price struct {
		Pr   float64 `table:"Price"`
		Idx  int     `table:",index"`
		Name int     `table:"-"` // ignored
	}
	structSlice := make([]*price, 0)
	err := pirateBooty.ToStructs(&structSlice)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	for _, p := range structSlice {
		fmt.Printf("%#v\n", p)
	}
	// Output:
	// &data.price{Pr:1, Idx:0, Name:0}
	// &data.price{Pr:0, Idx:1, Name:0}
	// &data.price{Pr:5.5, Idx:2, Name:0}
	// &data.price{Pr:600, Idx:3, Name:0}
}

func ExampleNewTableFromStructs() {
	type price struct {
		Pr   float64 `table:"Price"`
		Name string
		Who  string `table:"-"` // ignored
	}
	priceStructs := []*price{
		{Pr: 1.5, Name: "scrimshaw"},
		{Pr: 10000, Name: "yer mortal soul"},
	}
	newPrices := Must().NewTableFromStructs(priceStructs)
	fmt.Println(newPrices.ToRows())
	// Output: 2 Row(s):
	// | |  Price|           Name|
	// | |float64|         string|
	// ---------------------------
	// |0|    1.5|      scrimshaw|
	// |1|  10000|yer mortal soul|
}

func ExampleRow_ToStruct() {
	fmt.Println()
	for row := range pirateBooty.IterAll() {
		// get a struct representation of the row.
		asStruct := &struct {
			Pr   float64 `table:"Price"`
			Idx  int     `table:",index"`
			Name string
		}{}
		err := row.ToStruct(asStruct)
		if err != nil {
			panic(err)
		}
		fmt.Printf("row # %2d has price: %.1f\n", asStruct.Idx, asStruct.Pr)
	}
	// Output:
	// row #  0 has price: 1.0
	// row #  1 has price: 0.0
	// row #  2 has price: 5.5
	// row #  3 has price: 600.0
}
