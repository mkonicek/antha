package data

import (
	"reflect"

	"github.com/pkg/errors"
)

/*
 * join interfaces
 */

// JoinSelection contains a left table to apply join.
type JoinSelection struct {
	t *Table
}

// NaturalInner performs inner join between js.t and t by columns having the same names.
func (js *JoinSelection) NaturalInner(t *Table) (*Table, error) {
	cols := js.commonCols(t)
	return hashJoin(js.t, cols, t, cols, inner)
}

// NaturalLeftOuter performs left outer join between js.t and t by columns having the same names.
func (js *JoinSelection) NaturalLeftOuter(t *Table) (*Table, error) {
	cols := js.commonCols(t)
	return hashJoin(js.t, cols, t, cols, leftOuter)
}

// finds columns with the same name in two tables in order to perform natural join
func (js *JoinSelection) commonCols(t *Table) []ColumnName {
	common := []ColumnName{}
	for _, column := range js.t.schema.Columns {
		if _, err := t.schema.ColIndex(column.Name); err == nil {
			common = append(common, column.Name)
		}
	}
	return common
}

// On sets left table columns for a join.
func (js *JoinSelection) On(cols ...ColumnName) *JoinOn {
	return &JoinOn{t: js.t, cols: cols}
}

// JoinOn contains a left table and its columns to apply join.
type JoinOn struct {
	t    *Table
	cols []ColumnName
}

// Inner sets a right table and columns for inner join and performs the join itself.
func (jo *JoinOn) Inner(t *Table, cols ...ColumnName) (*Table, error) {
	return hashJoin(jo.t, jo.cols, t, cols, inner)
}

// LeftOuter sets a right table and columns for left outer join and performs the join itself.
func (jo *JoinOn) LeftOuter(t *Table, cols ...ColumnName) (*Table, error) {
	return hashJoin(jo.t, jo.cols, t, cols, leftOuter)
}

/*
 * join internals
 */

// Hash join itself. Currently only inner and left outer joins are supported.
func hashJoin(left *Table, leftCols []ColumnName, right *Table, rightCols []ColumnName, typ joinType) (*Table, error) {
	jq := &joinQuery{
		left:  joinQueryTable{t: left, cols: leftCols},
		right: joinQueryTable{t: right, cols: rightCols},
		typ:   typ,
	}

	// checking input tables
	if err := jq.check(); err != nil {
		return nil, err
	}

	// creating the right table index (if the key columns are not indexable, then unable to continue)
	indexKeyType, err := makeIndexKeyType(right.schema.MustProject(rightCols...))
	if err != nil {
		return nil, err
	}

	// partially creating an output table (for the time being, without series iterators and metadata)
	outputTable := jq.newJointTable()

	// joint table series iterators common state generator
	group := newSeriesGroup(func() seriesGroupStateImpl {
		return jq.newJoinState(indexKeyType)
	})

	// creating an iterator generator and metadata for each series
	for i, series := range outputTable.series {
		series.read = group.read(i)
		series.meta = jq.newJointSeriesMeta()
	}

	return outputTable, nil
}

// join input data
type joinQuery struct {
	left  joinQueryTable
	right joinQueryTable
	typ   joinType
}

type joinQueryTable struct {
	t    *Table
	cols []ColumnName
}

// join type
type joinType int

const (
	inner joinType = iota
	leftOuter
	//rightOuter  - not implemented yet
	//fullOuter - not implemented yet
)

// check checks input tables and columns correctness
func (jq *joinQuery) check() error {
	// checking the left table data
	if err := jq.left.check(); err != nil {
		return errors.Wrap(err, "left table")
	}
	// checking the right table data
	if err := jq.right.check(); err != nil {
		return errors.Wrap(err, "right table")
	}
	// checking that the numbers of columns are equal
	if len(jq.left.cols) != len(jq.right.cols) {
		return errors.Errorf("different number of columns for joining: %d on the left, %d on the right", len(jq.left.cols), len(jq.right.cols))
	}
	// checking that the types of columns are equal
	for i := range jq.left.cols {
		leftCol := jq.left.cols[i]
		leftType := jq.left.t.schema.MustCol(leftCol).Type
		rightCol := jq.right.cols[i]
		rightType := jq.right.t.schema.MustCol(rightCol).Type
		if leftType != rightType {
			return errors.Errorf("different types of columns for joining: column %s is of type %d, while column %s is of type %d",
				leftCol, leftType, rightCol, rightType)
		}
	}
	return nil
}

// check checks input table
func (jqt *joinQueryTable) check() error {
	// the table should be bounded
	if !isBounded(jqt.t.series) {
		return errors.New("unable to join unbounded tables")
	}
	// for now, joining on an empty columns list is not supported
	if len(jqt.cols) == 0 {
		return errors.New("unable to join on empty columns list")
	}
	// all the columns should exist
	if err := jqt.t.schema.CheckColumnsExist(jqt.cols...); err != nil {
		return errors.Wrap(err, "join table columns lookup")
	}
	return nil
}

// newJointTable partially creates a joint table (for the time being, without series iterators and metadata)
func (jq *joinQuery) newJointTable() *Table {
	// joint table contains all columns of both tables
	columns := append(append([]Column{}, jq.left.t.schema.Columns...), jq.right.t.schema.Columns...)
	// current implementation of join preserves the left table sort key
	key := jq.left.t.sortKey

	return newFromSchema(NewSchema(columns), key...)
}

// newJointSeriesMeta creates metadata for columns of joint table
func (jq *joinQuery) newJointSeriesMeta() *jointSeriesMeta {
	// calculating resulting series max size
	_, leftMaxSize, _ := seriesSize(jq.left.t.series)
	_, rightMaxSize, _ := seriesSize(jq.right.t.series)
	return &jointSeriesMeta{maxSize: leftMaxSize * rightMaxSize}
}

// joinState stores the common state of the joined table series iterators
type joinState struct {
	jq           *joinQuery
	leftColsNum  int
	rightColsNum int

	left           *tableIterator // the left (presumably large) table iterator
	leftProjection projection     // projection to calculate left index value
	leftKey        reflect.Value  // current left key

	right      indexMap // the right (presumably small) table in the form of a reflectively created index: indexKey -> []row (rows with such key)
	rightRows  []raw    // the right table rows corresponding to the current left table row
	rightIndex Index

	currRow raw // current row of the output table: |leftCol1|...|leftColN|rightCol1|...|rightColM|
}

func (jq *joinQuery) newJoinState(indexKeyType reflect.Type) *joinState {
	// indexing the right table
	rightIndex := indexTable(jq.right.t, jq.right.cols, indexKeyType)

	leftColsNum := len(jq.left.t.schema.Columns)
	rightColsNum := len(jq.right.t.schema.Columns)

	return &joinState{
		jq:             jq,
		leftColsNum:    leftColsNum,
		rightColsNum:   rightColsNum,
		left:           newTableIterator(jq.left.t.series),
		leftProjection: mustNewProjection(jq.left.t.schema, jq.left.cols...),
		leftKey:        reflect.New(indexKeyType).Elem(),
		right:          rightIndex,
		rightIndex:     -1,
		currRow:        newRaw(leftColsNum + rightColsNum),
	}
}

// creates a reflective index of a table by specified key columns: key -> []raw (rows having this key)
func indexTable(t *Table, keyCols []ColumnName, indexKeyType reflect.Type) indexMap {
	keyProjection := mustNewProjection(t.schema, keyCols...)
	iter := t.read(t.series)
	index := newIndexMap(indexKeyType, reflect.TypeOf([]raw{}))
	key := reflect.New(indexKeyType).Elem() // creating the reflective key once in order not to call reflect.New multiple times
	for iter.Next() {
		// load a row
		row := iter.rawValue()
		// row key columns
		rowProj := row.project(keyProjection)
		// load key column values into the reflective key
		loadIndexKeyFromRow(rowProj, key)
		// adding the row to the index
		var keyRows []raw // rows having this key
		if keyRowsValue, ok := index.lookup(key); ok {
			keyRows = keyRowsValue.Interface().([]raw)
		}
		keyRows = append(keyRows, row)
		index.set(key, reflect.ValueOf(keyRows))
	}
	return index
}

func (js *joinState) Next() bool {
	for {
		// advance rightRows if possible
		if js.rightIndex+1 < Index(len(js.rightRows)) {
			js.rightIndex++
			js.setRightRow(js.rightRows[js.rightIndex])
			return true
		}

		// rightRows is exhausted => advance left table if possible
		if !js.left.Next() {
			return false
		}

		// fetch left row
		leftRow := js.left.rawValue()
		copy(js.currRow[:js.leftColsNum], leftRow)

		// fetch rightRows corresponding to leftRow
		leftRowProj := leftRow.project(js.leftProjection)
		loadIndexKeyFromRow(leftRowProj, js.leftKey) // load current key from left row
		if rightRawsValue, ok := js.right.lookup(js.leftKey); ok {
			js.rightRows = rightRawsValue.Interface().([]raw)
		} else {
			js.rightRows = nil
		}
		js.rightIndex = -1

		// in case of left join, return empty rightRow if rightRows is empty
		if js.jq.typ == leftOuter && len(js.rightRows) == 0 {
			js.setRightRow(newRaw(js.rightColsNum))
			return true
		}
	}
}

// setRightRow sets the right row data and indicates that the next value is fetched
func (js *joinState) setRightRow(rightRow raw) {
	copy(js.currRow[js.leftColsNum:], rightRow)
}

func (js *joinState) Value(colIndex int) interface{} {
	return js.currRow[colIndex]
}

// jointSeriesMeta is a metadata type for joint table series
type jointSeriesMeta struct {
	maxSize int
}

func (m *jointSeriesMeta) IsMaterialized() bool { return false }
func (m *jointSeriesMeta) ExactSize() int       { return -1 }
func (m *jointSeriesMeta) MaxSize() int         { return m.maxSize }
