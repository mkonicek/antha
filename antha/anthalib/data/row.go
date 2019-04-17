package data

// Row represents a materialized record.
type Row struct {
	index  Index
	schema *Schema
	values []interface{}
}

// Rows are materialized table data, suitable for printing for example.
type Rows struct {
	Data   []Row
	Schema Schema
}

// Index returns a row index.
func (r Row) Index() Index {
	return r.index
}

// Schema returns a schema of the table which contains the Row.
func (r Row) Schema() Schema {
	return *r.schema
}

// Interface returns the row values in the form of []interface{}.
func (r Row) Interface() []interface{} {
	return r.values
}

// Value accesses a column value by column name.
func (r Row) Value(c ColumnName) (Value, error) {
	colIndex, err := r.schema.ColIndex(c)
	if err != nil {
		return Value{}, err
	}
	return Value{
		colIndex: colIndex,
		schema:   r.schema,
		value:    r.values[colIndex],
	}, nil
}

// MustValue accesses a column value by column name. Panics on errors
func (r Row) MustValue(c ColumnName) Value {
	v, err := r.Value(c)
	handle(err)
	return v
}

// ValueAt accesses a column value by column index.
func (r Row) ValueAt(colIndex int) Value {
	return Value{
		colIndex: colIndex,
		schema:   r.schema,
		value:    r.values[colIndex],
	}
}

// Values retreivies all the columns values from a Row.
func (r Row) Values() []Value {
	values := make([]Value, len(r.values))
	for i := range values {
		values[i] = r.ValueAt(i)
	}
	return values
}

// Value holds an arbitrary, nullable column value.
type Value struct {
	colIndex int
	schema   *Schema
	value    interface{}
}

// Column returns the column of the value.
func (v Value) ColIndex() int {
	return v.colIndex
}

// Column returns the column of the value.
func (v Value) Column() Column {
	return v.schema.Columns[v.colIndex]
}

// IsNull returns true if the value is null.
func (v Value) IsNull() bool {
	return v.value == nil
}

// Interface returns the underlying representation of the value.
func (v Value) Interface() interface{} {
	return v.value
}

// raw is a compact internal representation of Row (just values without metadata)
type raw []interface{}

func newRaw(size int) raw {
	return make([]interface{}, size)
}

func (r raw) row(index Index, schema *Schema) Row {
	return Row{
		index:  index,
		schema: schema,
		values: r,
	}
}

// project projects a row according to p
func (r raw) project(p projection) raw {
	newR := newRaw(len(p.newToOld))
	for new, old := range p.newToOld {
		newR[new] = r[old]
	}
	return newR
}

// projection efficiently (without dealing with strings) defines an ordered subset of the set of columns
type projection struct {
	newToOld []int // new column index -> old column index
}

func newProjection(s Schema, cols ...ColumnName) (projection, error) {
	newToOld := make([]int, len(cols))
	for new, col := range cols {
		old, err := s.ColIndex(col)
		if err != nil {
			return projection{}, err
		}
		newToOld[new] = old
	}
	return projection{newToOld: newToOld}, nil
}

func mustNewProjection(s Schema, cols ...ColumnName) projection {
	p, err := newProjection(s, cols...)
	handle(err)
	return p
}
