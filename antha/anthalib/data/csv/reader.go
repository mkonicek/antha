package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
)

// TableFromFile reads a data.Table eagerly from a CSV file
func TableFromFile(filePath string) (*data.Table, error) {
	// CSV file
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, errors.Wrapf(err, "opening CSV file '%s'", filePath)
	}
	defer file.Close() //nolint

	// reading a table from an io.Reader on the top of the CSV file
	return TableFromReader(bufio.NewReader(file))
}

// TableFromBytes reads a data.Table eagerly from a memory buffer
func TableFromBytes(buffer []byte) (*data.Table, error) {
	// reading a table from an io.Reader on the top of a memory buffer
	return TableFromReader(bytes.NewReader(buffer))
}

// TableFromReader reads a data.Table eagerly from io.Reader
func TableFromReader(reader io.Reader) (*data.Table, error) {
	// CSV reader
	csvReader := csv.NewReader(reader)

	// reading schema
	schema, err := readSchema(csvReader)
	if err != nil {
		return nil, errors.Wrapf(err, "reading CSV file schema")
	}

	// starting building a Table
	builder, err := data.NewTableBuilder(schema.Columns)
	if err != nil {
		return nil, err
	}

	// reading data
	for line := 1; ; line++ {
		// reading a record
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "reading CSV record")
		}
		// parsing a record into a row ([]interface{})
		row, err := parseCSVRecord(schema, record, line)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing CSV record")
		}
		// appending the row
		builder.Append(row)
	}

	// building a Table
	return builder.Build(), nil
}

// reads a schema from a CSV file first line
func readSchema(reader *csv.Reader) (*data.Schema, error) {
	record, err := reader.Read()
	if err != nil {
		return nil, errors.Wrapf(err, "reading CSV file header")
	}

	columns := make([]data.Column, len(record))
	for i, header := range record {
		name, typeName, err := parseColumnHeader(header)
		if err != nil {
			return nil, err
		}

		csvType, err := csvTypeByName(typeName)
		if err != nil {
			return nil, err
		}

		columns[i] = data.Column{Name: data.ColumnName(name), Type: csvType.typ}
	}

	return data.NewSchema(columns), nil
}

// parses CSV column header
func parseColumnHeader(header string) (string, string, error) {
	// header is assumed to contain column name and column type separated by a comma
	headerParts := strings.Split(header, ",")
	if len(headerParts) != 2 {
		return "", "", errors.Errorf("Column header '%s' does not contain column name and type separated by a comma", header)
	}
	return headerParts[0], headerParts[1], nil
}

// parses column values in a single record
func parseCSVRecord(schema *data.Schema, record []string, line int) ([]interface{}, error) {
	if len(record) != len(schema.Columns) {
		panic("wrong number of fields") // unreachable since csv.Writer checks the number of fields
	}

	row := make([]interface{}, len(record))
	for i := range record {
		valueText := record[i]
		column := schema.Columns[i]

		// for now, treating empty text as null
		// TODO: some better way to distinguish nulls and empty strings
		if len(valueText) == 0 {
			row[i] = nil
			continue
		}

		csvType, err := csvTypeByReflectType(column.Type)
		if err != nil {
			panic(err)
		}

		value, err := csvType.parse(valueText)
		if err != nil {
			return nil, errors.Wrapf(err, "line %d, column %d", line, i)
		}
		row[i] = value
	}

	return row, nil
}
