package data

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

type jsonHeader struct {
	Schema  Schema
	SortKey Key `json:",omitempty"`
	Size    int
}

type jsonCol struct {
	Name string
	Type string
}

type jsonRead struct {
	Header jsonHeader
	Data   []*json.RawMessage
}

// MarshalJSON encodes type info only for well-known types.
func (c *Column) MarshalJSON() ([]byte, error) {
	typName, found := typeSupport[c.Type]
	if !found {
		return nil, errors.Errorf("serializing type %+v for column %q is unsupported", c.Type, c.Name)
	}
	return json.Marshal(jsonCol{string(c.Name), typName})
}

// UnmarshalJSON decodes type info only for well-known types.
func (c *Column) UnmarshalJSON(bytes []byte) error {
	jc := new(jsonCol)
	err := json.Unmarshal(bytes, jc)
	if err != nil {
		return errors.Wrap(err, "unable to read column")
	}
	var found bool
	// TODO this type registry could be extended dynamically using tokenizing and reflection
	c.Type, found = typeSupportByName[jc.Type]
	if !found {
		return errors.Errorf("deserializing type %+v for column %q is unsupported", jc.Type, jc.Name)
	}
	c.Name = ColumnName(jc.Name)
	return nil
}

// MarshalJSON saves the table as an object with a header array (for the schema)
// followed by a row-major data array. This is supported for scalar types only.
func (t *Table) MarshalJSON() ([]byte, error) {
	// TODO error if data is unbounded
	b := new(bytes.Buffer)
	b.WriteString(`{`)
	// printing the header first to allow potentially seeking through the dataset lazily
	b.WriteString(`"Header":`)
	hdr, err := json.Marshal(jsonHeader{Schema: t.schema, SortKey: t.sortKey, Size: t.Size()})
	if err != nil {
		return nil, errors.Wrap(err, "unable to emit header")
	}
	b.Write(hdr)
	b.WriteString(",\n" + `"Data":[`)
	// writing data
	firstRow := true

	for row := range t.IterAll() {
		// TODO use Encoder here?
		if !firstRow {
			b.WriteString(",\n")
		}
		b.WriteString(`[`)
		for i, o := range row.Values {
			if i > 0 {
				b.WriteString(",")
			}
			if o.IsNull() {
				b.WriteString(`null`)
			} else {
				rowBytes, err := json.Marshal(o.value)
				if err != nil {
					return nil, errors.Wrapf(err, "unable to emit value of %q at index %d", o.ColumnName(), row.Index)
				}
				b.Write(rowBytes)
			}
		}
		firstRow = false
		b.WriteString("]")
	}

	b.WriteString(`]}`)
	return b.Bytes(), nil
}

// UnmarshalJSON eagerly deserializes to a Table.
// This is supported for scalar types only.
func (t *Table) UnmarshalJSON(b []byte) (err error) {
	// TODO lazy decoding!
	partiallyParsed := new(jsonRead)
	err = json.Unmarshal(b, partiallyParsed)
	if err != nil {
		return errors.Wrap(err, "unable to parse table header")
	}
	builder, err := NewTableBuilder(partiallyParsed.Header.Schema.Columns)
	if err != nil {
		return
	}
	schema := builder.schema()
	var idx int
	var values *json.RawMessage
	defer func() {
		if errRec := recover(); errRec != nil {
			err = errors.Errorf("failed to build table with schema %v: %+v (last handled row %d: %s)", schema, errRec, idx, *values)
		}
	}()
	row := make([]interface{}, schema.NumColumns())
	colBaseVals := make([]reflect.Value, schema.NumColumns())
	for cIdx, c := range schema.Columns {
		// ptr
		colBaseVals[cIdx] = reflect.New(c.Type)
	}
	for idx, values = range partiallyParsed.Data {
		// if err := json.Unmarshal(*values, row); err != nil {
		// 	return errors.Wrapf(err, "parsing table at row %d", idx)
		// }
		// TODO the encoder could be constructed just once outside loop.
		dec := json.NewDecoder(bytes.NewReader(*values))
		_, err = dec.Token() // == delim
		checkErr(err)

		// convert values back to reqd types using reflective call
		for cIdx := range schema.Columns {
			// note: this won't work
			// 			val := reflect.NewAt(c.Type, unsafe.Pointer(uintptr(0)))

			// do a sneaky hack to peek at the next value and see if it is nil.  this is horrific because of brittle whitespace.
			// However it seems to be a relatively (?) quick way of doing this
			hackBuf := make([]byte, 2)
			_, err := dec.Buffered().Read(hackBuf)
			checkErr(err)
			isNull := hackBuf[0] == 'n' || (hackBuf[0] == ',' && hackBuf[1] == 'n')
			if isNull {
				row[cIdx] = nil
				continue
			}
			err = dec.Decode(colBaseVals[cIdx].Interface())
			if err != nil {
				return errors.Wrapf(err, "failed to build table with schema %v: parsing table at row %d: %s", schema, idx, *values)
			}
			row[cIdx] = colBaseVals[cIdx].Elem().Interface()

		}
		builder.Append(row)
	}
	*t = *builder.Build()
	t.sortKey = partiallyParsed.Header.SortKey
	return
}

func checkErr(err error) {
	if err != nil {
		panic(errors.Wrap(err, "SHOULD NOT HAPPEN"))
	}
}
