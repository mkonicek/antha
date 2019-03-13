package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
)

// WriteTable writes a data.Table to a CSV file
func WriteTable(table *data.Table, filePath string) error {
	// CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "creating CSV file '%s'", filePath)
	}
	defer file.Close() //nolint

	// CSV writer
	writer := csv.NewWriter(bufio.NewWriter(file))

	// writing CSV header
	schema := table.Schema()
	err = writeHeader(&schema, writer)
	if err != nil {
		return errors.Wrap(err, "writing a CSV header")
	}

	rows, done := table.Iter()
	defer done()

	// writing data
	for line := 1; ; line++ {
		// reading a row from the table
		row, ok := <-rows
		if !ok {
			break
		}

		// making a CSV record
		record := rowToCsvRecord(row, &schema)

		// writing to buffer
		if err := writer.Write(record); err != nil {
			return errors.Wrapf(err, "writing a CSV header at line %d", line)
		}
	}

	writer.Flush()
	return nil
}

// writes a CSV header in the form `"columnName1,columnType1","columnName2,columnType2"..."`
func writeHeader(schema *data.Schema, writer *csv.Writer) error {
	header := make([]string, len(schema.Columns))
	for i, column := range schema.Columns {
		csvType, err := csvTypeByReflectType(column.Type)
		if err != nil {
			return errors.Wrapf(err, "writing CSV header column %d", i)
		}
		header[i] = string(column.Name) + "," + csvType.name
	}

	return writer.Write(header)
}

func rowToCsvRecord(row data.Row, schema *data.Schema) []string {
	if len(row.Values) != len(schema.Columns) {
		panic("wrong number of columns")
	}

	record := make([]string, len(row.Values))
	for i := range record {
		if value := row.Values[i].Interface(); value != nil {
			// for now, using just %v to print values: this suits all our types including timestamps
			// in case we need to store more complex types we should instead extend csvType with a print func
			record[i] = fmt.Sprintf("%v", value)
		}
	}

	return record
}
