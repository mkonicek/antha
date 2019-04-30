package parquet

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/data"
	"github.com/pkg/errors"
	"github.com/xitongsys/parquet-go/ParquetFile"
	"github.com/xitongsys/parquet-go/ParquetReader"
	"github.com/xitongsys/parquet-go/parquet"
)

type readState struct {
	reader      *ParquetReader.ParquetReader
	columnNames []data.ColumnName
}

// ReadOpt sets an optional behavior when reading parquet files.
type ReadOpt func(*readState) error

// Columns returns a ReadOpt that selects a subset of columns to read from the
// source. If no column names are specified, reads all the columns. This is an
// optimization for projections, that should not be needed in future (when lazy
// access patterns are possible).
func Columns(columnNames ...data.ColumnName) ReadOpt {
	return func(r *readState) error {
		r.columnNames = columnNames
		return nil
	}
}

// TableFromReader reads a data.Table eagerly from io.Reader.
func TableFromReader(reader io.Reader, opts ...ReadOpt) (*data.Table, error) {
	// reading into a memory buffer
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "read Parquet file bytes")
	}

	return TableFromBytes(buffer, opts...)
}

// TableFromBytes reads a data.Table eagerly from a memory buffer
func TableFromBytes(buffer []byte, opts ...ReadOpt) (*data.Table, error) {
	// wrap a byte buffer into a ParquetFile object
	file, err := ParquetFile.NewBufferFile(buffer)
	if err != nil {
		panic(errors.Wrap(err, "SHOULD NOT HAPPEN: creating in-memory ParquetFile"))
	}

	return readTable(read(file, opts))
}

// TableFromFile reads a data.Table eagerly from a Parquet file
func TableFromFile(filePath string, opts ...ReadOpt) (*data.Table, error) {
	// opening the file on disk via ParquetFile object
	file, err := ParquetFile.NewLocalFileReader(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "read Parquet file '%s'", filePath)
	}
	defer file.Close() //nolint

	return readTable(read(file, opts))
}

func read(file ParquetFile.ParquetFile, opts []ReadOpt) (*readState, error) {
	// for now, reading Parquet file in 1 thread, 100 rows at once
	// TODO: which parameters to use for really large datasets?
	np := 1

	// parquet reader
	parquetReader, err := ParquetReader.NewParquetReader(file, nil, int64(np))
	if err != nil {
		return nil, errors.Wrapf(err, "create Parquet reader")
	}
	r := &readState{reader: parquetReader}
	for _, o := range opts {
		if err = o(r); err != nil {
			return nil, errors.Wrapf(err, "option %v generated an error state when reading %v", o, file)
		}
	}
	return r, nil
}

// reads a table from an arbitrary source (in the form of ParquetFile)
func readTable(r *readState, err error) (*data.Table, error) {
	if err != nil {
		return nil, err
	}

	// reading Parquet file metadata
	metadata := r.reader.Footer

	// transforming Parquet file metadata into parquetSchema
	schema, err := schemaFromParquetMetadata(metadata, r.columnNames)
	if err != nil {
		return nil, err
	}

	// converting parquetSchame to json understandable by parquet-go
	jsonSchema, err := schema.toJSON()
	if err != nil {
		return nil, err
	}

	// creating a struct type to accomodate single Table row data (a dynamic data type required by parquet-go)
	rowType := rowStructFromSchema(schema)

	// starting building a Table
	builder, err := data.NewTableBuilder(schema.Columns)
	if err != nil {
		return nil, err
	}

	// reading from Parquet
	err = r.readFromParquet(jsonSchema, rowType, func(rowStruct interface{}) {
		// populating row values from a dynamic row struct
		row := rowValuesFromRowStruct(rowStruct, schema)
		// appending the row to the table builder
		builder.Append(row)
	})
	if err != nil {
		return nil, err
	}

	// building a Table
	return builder.Build(), nil
}

// Reads rows from Parquet file
func (r *readState) readFromParquet(jsonSchema string, rowType reflect.Type, onRow func(interface{})) error {
	defer r.reader.ReadStop()

	// !!! hack !!!
	patchWhitespaces(r.reader.Footer)
	// parquet schema
	if err := r.reader.SetSchemaHandlerFromJSON(jsonSchema); err != nil {
		return errors.Wrap(err, "set Parquet schema")
	}
	batchSize := 100

	// total number of rows
	numRows := int(r.reader.GetNumRows())

	// type []rowType
	sliceType := reflect.SliceOf(rowType)
	// var *[]rowType
	slicePtr := reflect.New(sliceType)

	for numRows > 0 {
		rowCount := min(batchSize, numRows)
		numRows -= rowCount

		// make([]rowType, rowCount, rowCount)
		slicePtr.Elem().Set(reflect.MakeSlice(sliceType, rowCount, rowCount))

		// reading
		if err := r.reader.Read(slicePtr.Interface()); err != nil {
			return errors.Wrap(err, "reading data from Parquet")
		}

		// callback
		slice := slicePtr.Elem()
		for i := 0; i < rowCount; i++ {
			onRow(slice.Index(i).Interface())
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// !!! hack !!!
// parquet-go currently does not support whitespaces in field tags; it deletes them (github.com/xitongsys/parquet-go/Common/Common.go, line 57)
// workaround: this function deletes whitespaces in Parquet file footer just as parquet-go does it in tags
// TODO: delete this; if we decide we need column names with whitespaces, then patch parquet-go
func patchWhitespaces(footer *parquet.FileMetaData) {
	for _, schemaElement := range footer.Schema {
		schemaElement.Name = strings.Replace(schemaElement.Name, " ", "", -1)
	}
	for _, rowGroup := range footer.RowGroups {
		for _, chunk := range rowGroup.Columns {
			for i := range chunk.MetaData.PathInSchema {
				chunk.MetaData.PathInSchema[i] = strings.Replace(chunk.MetaData.PathInSchema[i], " ", "", -1)
			}
		}
	}
}

// populates row values from a dynamic row struct
func rowValuesFromRowStruct(rowStruct interface{}, schema *parquetSchema) []interface{} {
	rowStructValue := reflect.ValueOf(rowStruct)
	if rowStructValue.NumField() != len(schema.Columns) {
		panic("unexpected number of rowStruct fields")
	}

	values := make([]interface{}, len(schema.Columns))
	for i := range schema.Columns {
		field := rowStructValue.Field(i)
		// if fields are optional, they are stored as pointers here
		if schema.nullable(i) {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}
		// a workaround for timestamps: parquet-go reads them into int64 field, so converting them from data.TimestampSmth here
		values[i] = field.Convert(schema.Columns[i].Type).Interface()
	}

	return values
}
