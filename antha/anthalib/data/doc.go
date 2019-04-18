/*
Package data provides lazy data tables with file format support.

Overview

The aim of this package is to make it possible to efficiently manipulate tabular data with arbitrary schemas
in Go (and Antha).  It is intended to handle large datasets efficiently by operating lazily where
possible.

An example usecase is to load a dataset, project it dynamically to some subset of columns, calculate a new column
and write the combined  output to a new file.

Tables

The main data type is Table.  Typically we will start by reading from a data source such as a Parquet file:
	import 	"github.com/antha-lang/antha/antha/anthalib/data/parquet"
	table, err := parquet.TableFromBytes(myBuffer)

Tables are lazy and will not typically load all data until requested to.

We can materialize the Table into a slice of Row objects.  This is useful for printing it, so that:
	fmt.Println(table.ToRows())

might print:
	5 Row(s):
	| |Quantity|Capacity| label|
	| | float64| float64|string|
	----------------------------
	|0|      11|      30|     A|
	|1|      12|      30|     B|
	|2|    10.5|      50|     A|
	|3|   <nil>|      30|     C|
	|4|     5.5|      30|     C|

A Table collects Series values as columns.  These are not really useful on their own but their types are exposed
via the table Schema, so that:
	fmt.Printf("%+v\n", table.Schema().Columns)

might print:
	[{Name:quantity Type:float64} {Name:label Type:string}]

Each column in a Table has a Go datatype, but may also be nil, even for scalar types such as string.  nil is not
thre same as zero or NaN.

Data import

In addition to reading tables from Parquet files (see above), there are several other ways to create tables.

It is possible to construct a table "column-wise" by creating Series instances from a Go slice and a nullability mask, like this:
	column1 := data.Must().NewSeriesFromSlice("quantity", []float64{11.0, 12.0, 10.5, 0, 5.5}, []bool{true, true, true, false, true})
	column2 := data.Must().NewSeriesFromSlice("label", []string{"A", "B", "A", "C", "C"}, nil)
	table := data.NewTable(column1, column2)

If it is needeed to create a table "row-wise" rather than "column-wise", the easiest way is to read it from a slice of structs:
	type myType struct {
		Capacity float64
		Label    string
	}
	myData := []myType{{30, "A"}, {50, "B"}}
	table := data.Must().NewTableFromStructs(myData)

Another way to create a table "row-wise", especially when its schema is not known at compile time, is building a table with a Builder:
	columns := []data.Column{
		{Name: "Id", Type: reflect.TypeOf(0)},
		{Name: "Label", Type: reflect.TypeOf("")},
	}
	builder := data.Must().NewTableBuilder(columns)
	builder.Append([]interface{}{30, "A"})
	builder.Append([]interface{}{40, nil})
	table := builder.Build()

Data manipulation

Tables can be transformed with methods that return a new Table.  Typically the returned Table is lazy.  This means that an expression like:
	callback := func(vals ...float64) float64 { return 100 * vals[0] / vals[1] }
	table.Extend("quantity_as_percent").On("quantity", "Capacity").Float64(callback)

does not evaluate the callback.  It will be evaluated when the quantity_as_percent value is needed by a dependent Table.

Here the value of the Table would be:

	5 Row(s):
	| |quantity|Capacity| Label|quantity_as_percent|
	| | float64| float64|string|            float64|
	------------------------------------------------
	|0|      11|      30|     A| 36.666666666666664|
	|1|      12|      30|     B|                 40|
	|2|    10.5|      50|     A|                 21|
	|3|   <nil>|      30|     C|              <nil>|
	|4|     5.5|      30|     C| 18.333333333333332|

Sort is a special case, as it is eager: it materializes the entire table.

Other data manipulation methods include Project, Slice, Sort, Filter, Distinct, Pivot, and Join.

Calling Cache on a Table returns a fully materialized copy, which is useful if the data needs to be used for more than one
subsequent operation.

Data export

We can write data back to files, such as Parquet, using parquet.TableToBytes.

Also it is possible to populate a slice of structs with a table data:
	type myType struct {
		Capacity float64
		Label    string
	}
	structs := []myType{}
	table.ToStructs(&structs)
	fmt.Printf("%+v\n", structs)

This might print:
	[{Capacity:30 Label:A} {Capacity:30 Label:B} {Capacity:50 Label:A} {Capacity:30 Label:C} {Capacity:30 Label:C}]

Or, alternatively, table data can be exported manually using iteration.

Iteration

If necessary the whole table can be processed by a callback:
	for record := range table.IterAll() {
		m, _ := record.Value("Quantity")
		if m.IsNull() {
			fmt.Println("quantity=null at index", record.Index)
		} else {
			fmt.Println("quantity=", m.MustFloat64())
		}
	}
which here might print:
	quantity= 11
	quantity= 12
	quantity= 10.5
	quantity=null at index 3
	quantity= 5.5


Future features

The main features still to implement include:

Aggregation, such as group by, along with aggregate functions such as average, max, etc.

Relational operations - concat, union, join.

Complete Parquet and CSV read/write support.

There is the possibility of extending this to an event processing dataflow framework
(like Apache Flink) for more 'realtime' usecases on unbounded datasets.

*/
package data
