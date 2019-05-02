package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
)

// TableToFile writes a data.Table to a CSV file
func TableToFile(table *data.Table, filePath string) error {
	// CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "creating CSV file '%s'", filePath)
	}
	defer file.Close() //nolint

	// writing the table
	return TableToWriter(table, bufio.NewWriter(file))
}

// TableToBytes writes a data.Table to a memory buffer
func TableToBytes(table *data.Table) ([]byte, error) {
	// a memory buffer writer
	buffer := bytes.NewBuffer(nil)

	// writing the table
	if err := TableToWriter(table, buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// TableToWriter writes a data.Table to io.Writer
func TableToWriter(table *data.Table, writer io.Writer) error {
	// CSV writer
	csvWriter := csv.NewWriter(writer)

	// writing CSV header
	schema := table.Schema()
	err := writeHeader(&schema, csvWriter)
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
		record := rowToCsvRecord(row)

		// writing to buffer
		if err := csvWriter.Write(record); err != nil {
			return errors.Wrapf(err, "writing a CSV header at line %d", line)
		}
	}

	csvWriter.Flush()
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

func rowToCsvRecord(row data.Row) []string {
	record := make([]string, len(row.Schema().Columns))
	for i, value := range row.Values() {
		if !value.IsNull() {
			// for now, using just %v to print values: this suits all our types including timestamps
			// in case we need to store more complex types we should instead extend csvType with a print func
			record[i] = fmt.Sprintf("%v", value.Interface())
		}
	}

	return record
}
