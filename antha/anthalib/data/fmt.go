package data

import (
	"fmt"
	"strings"
)

// Pretty printing Rows and Schema
var _ fmt.Stringer = Rows{}
var _ fmt.Stringer = Schema{}

// String formats the rows as an ascii art table
func (r Rows) String() string {
	const sep = "|"
	const line = "-"
	const hdrSize = 2
	// row-major array with offset for header rows
	cellVals := make([][]interface{}, len(r.Data)+hdrSize)
	colMaxes := make([]int, len(r.Schema.Columns)+1)

	add := func(rownumWithOffset, colnumWithOffset int, value interface{}) {
		// TODO an appropriate format could be determined from the schema
		str := fmt.Sprintf("%+v", value)
		cellVals[rownumWithOffset] = append(cellVals[rownumWithOffset], str)
		if len(str) > colMaxes[colnumWithOffset] {
			colMaxes[colnumWithOffset] = len(str)
		}
	}
	// headers
	add(0, 0, "")
	add(1, 0, "")

	for c, col := range r.Schema.Columns {
		add(0, c+1, col.Name)
		add(1, c+1, col.Type)
	}

	for rownum, rr := range r.Data {
		add(rownum+hdrSize, 0, rr.Index())
		for c, v := range rr.Values() {
			add(rownum+hdrSize, c+1, v.Interface())
		}
	}
	fmtStrBuilder := strings.Builder{}

	for _, colMax := range colMaxes {
		// TODO pad string columns (only) on the right. etc
		fmtStrBuilder.WriteString(fmt.Sprintf("%s%%%ds", sep, colMax)) //nolint
	}
	fmtStrBuilder.WriteString(sep + "\n") //nolint
	fmtStr := fmtStrBuilder.String()
	builder := strings.Builder{}
	for idx, cells := range cellVals {
		rowStr := fmt.Sprintf(fmtStr, cells...)
		builder.WriteString(rowStr) //nolint
		if idx == 1 {
			// divider
			builder.WriteString(strings.Repeat(line, len(rowStr)-1)) //nolint
			builder.WriteString("\n")                                //nolint
		}
	}
	return fmt.Sprintf("%d Row(s):\n%s", len(cellVals)-hdrSize, builder.String())
}

// String formats the Row as the Rows containing a single row
func (r Row) String() string {
	return Rows{
		Data:   []Row{r},
		Schema: r.Schema(),
	}.String()
}

// String formats the schema as a list of columns names and types (each from a new line).
func (s Schema) String() string {
	builder := strings.Builder{}
	for _, column := range s.Columns {
		builder.WriteString(fmt.Sprintf("%s, %+v\n", column.Name, column.Type)) //nolint
	}
	return builder.String()
}

// TODO: method to scan a table literal string into Rows
