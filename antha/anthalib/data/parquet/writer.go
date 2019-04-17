package parquet

import (
	"bytes"
	"io"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
	"github.com/xitongsys/parquet-go/ParquetFile"
	"github.com/xitongsys/parquet-go/ParquetWriter"
	"github.com/xitongsys/parquet-go/parquet"
)

// TableToWriter writes a data.Table to io.Writer
func TableToWriter(table *data.Table, writer io.Writer) error {
	// wrapping io.Writer in a ParquetFile.ParquetFile
	file := ParquetFile.NewWriterFile(writer)

	// writing the table
	return writeTable(table, file)
}

// TableToBytes writes a data.Table to a memory buffer
func TableToBytes(table *data.Table) ([]byte, error) {
	// a memory buffer writer
	buffer := bytes.NewBuffer(nil)

	// a ParquetFile.ParquetFile on the top of the memory buffer writer
	file := ParquetFile.NewWriterFile(buffer)

	// writing the table
	if err := TableToWriter(table, file); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// TableToFile writes a data.Table to a file on disk
func TableToFile(table *data.Table, filePath string) error {
	// opening the file
	file, err := ParquetFile.NewLocalFileWriter(filePath)
	if err != nil {
		return errors.Wrapf(err, "opening Parquet file '%s' for writing", filePath)
	}
	defer file.Close() //nolint

	// writing a table
	return writeTable(table, file)
}

// writes a data.Table to ParquetFile.ParquetFile
func writeTable(table *data.Table, file ParquetFile.ParquetFile) error {
	// parquet schema
	tableSchema := table.Schema()
	schema := newParquetSchema(&tableSchema)

	// converting parquetSchame to json understandable by parquet-go
	jsonSchema, err := schema.toJSON()
	if err != nil {
		return err
	}

	// creating a struct type to accomodate single Table row data (a dynamic data type requered by parquet-go)
	rowType := rowStructFromSchema(schema)

	// starting iterating through the table
	iter, done := table.Iter()
	defer done()

	// writing to Parquet
	return writeToParquet(file, jsonSchema, rowType, func() (interface{}, error) {
		row, ok := <-iter
		if !ok {
			return nil, nil
		}
		return makeRowValue(row, rowType), nil
	})
}

// Writes rows to Parquet file
func writeToParquet(file ParquetFile.ParquetFile, jsonSchema string, rowType reflect.Type, rowIter func() (interface{}, error)) error {
	// Parquet writer and its settings
	writer, err := ParquetWriter.NewParquetWriter(file, nil, 1)
	if err != nil {
		return errors.Wrap(err, "creating Parquet writer")
	}
	writer.RowGroupSize = 128 * 1024 * 1024 //128M
	writer.CompressionType = parquet.CompressionCodec_SNAPPY

	// parquet schema
	if err := writer.SetSchemaHandlerFromJSON(jsonSchema); err != nil {
		return errors.Wrap(err, "set Parquet schema")
	}

	// Writing to Parquet
	for {
		row, err := rowIter()
		if err != nil {
			return err
		}
		if row == nil {
			break
		}
		if err = writer.Write(row); err != nil {
			return err
		}
	}

	// Flush
	return writer.WriteStop()
}

// copies data.Row content into a dynamic data struct suitable for Parquet writer
func makeRowValue(row data.Row, rowType reflect.Type) interface{} {
	rowValue := reflect.New(rowType)
	// filling fields
	for i, obs := range row.Values() {
		field := rowValue.Elem().Field(i)
		if obs.Interface() != nil {
			// creating a pointer (all values are optional, so they are stored in pointer felds)
			field.Set(reflect.New(field.Type().Elem()))
			// a workaround for timestamps: parquet-go uses int64 for them, so converting them from data.TimestampSmth to int64 here
			value := reflect.ValueOf(obs.Interface()).Convert(field.Type().Elem())
			// setting a value
			field.Elem().Set(value)
		}
	}
	return rowValue.Interface()
}
