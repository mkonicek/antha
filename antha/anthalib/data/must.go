package data

import (
	"reflect"
)

/*
 * utility for wrapping error functions
 */

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// MustCreate panics on any error when creating tables/series.
type MustCreate struct{}

// Must returns a proxy object that asserts errors are nil.
func Must() MustCreate {
	return MustCreate{}
}

// NewSeriesFromSlice is a errors handling wrapper for NewSeriesFromSlice.
func (m MustCreate) NewSeriesFromSlice(col ColumnName, values interface{}, notNull []bool) *Series {
	ser, err := NewSeriesFromSlice(col, values, notNull)
	handle(err)
	return ser
}

// NewTableFromStructs is a errors handling wrapper for NewTableFromStructs.
func (m MustCreate) NewTableFromStructs(structs interface{}) *Table {
	t, err := NewTableFromStructs(structs)
	handle(err)
	return t
}

// NewTableBuilder is a errors handling wrapper for NewTableBuilder.
func (m MustCreate) NewTableBuilder(columns []Column) *TableBuilder {
	b, err := NewTableBuilder(columns)
	handle(err)
	return b
}

// MustTable panics on any error when creating derived tables.
type MustTable struct {
	*Table
}

// Must returns a proxy object that asserts errors are nil.
func (t *Table) Must() MustTable {
	return MustTable{t}
}

// Cache panics unless *Table.Cache.
func (m MustTable) Cache() *Table {
	t, err := m.Table.Cache()
	handle(err)
	return t
}

// Project panics unless *Table.Project.
func (m MustTable) Project(columns ...ColumnName) *Table {
	t, err := m.Table.Project(columns...)
	handle(err)
	return t
}

// Convert panics unless *Table.Convert.
func (m MustTable) Convert(col ColumnName, typ reflect.Type) *Table {
	t, err := m.Table.Convert(col, typ)
	handle(err)
	return t
}

// Sort panics unless *Table.Sort.
func (m MustTable) Sort(key Key) *Table {
	t, err := m.Table.Sort(key)
	handle(err)
	return t
}

// SortByFunc panics unless *Table.SortByFunc.
func (m MustTable) SortByFunc(f SortFunc) *Table {
	t, err := m.Table.SortByFunc(f)
	handle(err)
	return t
}

// Filter

// Filter returns a proxy for *Table.Filter.
func (m MustTable) Filter() *MustFilterSelection {
	return &MustFilterSelection{m.Table.Filter()}
}

// MustFilterSelection panics on any error when creating derived tables.
type MustFilterSelection struct {
	*FilterSelection
}

// On returns a proxy for *FilterSelection.On.
func (s *MustFilterSelection) On(cols ...ColumnName) *MustFilterOn {
	return &MustFilterOn{s.FilterSelection.On(cols...)}
}

// MustFilterOn panics on any error when creating derived tables.
type MustFilterOn struct {
	*FilterOn
}

// Distinct

// Distinct returns a proxy for *Table.Distinct.
func (m MustTable) Distinct() *MustDistinctSelection {
	return &MustDistinctSelection{m.Table.Distinct()}
}

// MustDistinctSelection panics on any error when creating derived tables.
type MustDistinctSelection struct {
	*DistinctSelection
}

// On selects distinct rows by columns specified. Panics on errors.
func (s *MustDistinctSelection) On(cols ...ColumnName) *Table {
	t, err := s.DistinctSelection.On(cols...)
	handle(err)
	return t
}

// Extend

// Extend returns a proxy for *Table.Extend.
func (m MustTable) Extend(newCol ColumnName) *MustExtension {
	return &MustExtension{m.Table.Extend(newCol)}
}

// MustExtension panics on any error when creating derived tables.
type MustExtension struct {
	*Extension
}

// ConstantType panics unless *Extension.ConstantType.
func (e *MustExtension) ConstantType(value interface{}, typ reflect.Type) *Table {
	t, err := e.Extension.ConstantType(value, typ)
	handle(err)
	return t
}

// On returns a proxy for *Extension.On.
func (e *MustExtension) On(cols ...ColumnName) *MustExtendOn {
	return &MustExtendOn{e.Extension.On(cols...)}
}

// MustExtendOn panics on any error when creating derived tables.
type MustExtendOn struct {
	*ExtendOn
}

// Interface panics on any error when creating derived tables.
func (on *MustExtendOn) Interface(f func(v ...interface{}) interface{}, newType reflect.Type) *Table {
	t, err := on.ExtendOn.Interface(f, newType)
	handle(err)
	return t
}

// Update

// Update returns a proxy for *Table.Update.
func (m MustTable) Update(col ColumnName) *MustUpdateSelection {
	return &MustUpdateSelection{m.Table.Update(col)}
}

// MustUpdateSelection panics on any error when creating derived tables.
type MustUpdateSelection struct {
	*UpdateSelection
}

// By performs Update using the whole table row as input.
// For wide tables this might be inefficient, consider using On(...) instead.
func (us *MustUpdateSelection) By(fn func(r Row) interface{}) *Table {
	t, err := us.UpdateSelection.By(fn)
	handle(err)
	return t
}

// Constant makes column a constant column with the given value, but of the same type.
// Returns error if a non-nil value cannot be converted to the column type.
func (us *MustUpdateSelection) Constant(value interface{}) *Table {
	t, err := us.UpdateSelection.Constant(value)
	handle(err)
	return t
}

// On selects a subset of columns to use as an extension source. If duplicate columns
// exist, the first so named is used.  Note this does not panic yet, even if the
// columns do not exist.  (However subsequent calls to the returned object will
// error.)
func (us *MustUpdateSelection) On(cols ...ColumnName) *MustUpdateOn {
	return &MustUpdateOn{us.UpdateSelection.On(cols...)}
}

// MustUpdateOn panics on any error when creating derived tables.
type MustUpdateOn struct {
	*UpdateOn
}

// Interface updates a column of an arbitrary type using a subset of source columns of arbitrary types.
func (on *MustUpdateOn) Interface(fn func(v ...interface{}) interface{}) *Table {
	t, err := on.UpdateOn.Interface(fn)
	handle(err)
	return t
}

// Pivot

// Pivot returns a proxy for *Table.Pivot.
func (m MustTable) Pivot() *MustPivotSelection {
	return &MustPivotSelection{m.Table.Pivot()}
}

// MustPivotSelection is a proxy for PivotSelection.
type MustPivotSelection struct {
	*PivotSelection
}

// Key returns a proxy for PivotKey.
func (mps *MustPivotSelection) Key(key ...ColumnName) *MustPivotKey {
	return &MustPivotKey{&PivotKey{table: mps.table, key: key}}
}

// MustPivotKey panics on any error when creating derived tables.
type MustPivotKey struct {
	*PivotKey
}

// Columns panics on any error when creating derived tables.
func (mpk *MustPivotKey) Columns(pivot ColumnName, value ColumnName) *Table {
	table, err := mpk.PivotKey.Columns(pivot, value)
	handle(err)
	return table
}

// Join

// Join returns a proxy for *Table.Join.
func (m MustTable) Join() *MustJoinSelection {
	return &MustJoinSelection{m.Table.Join()}
}

// MustJoinSelection is a proxy for JoinSelection.
type MustJoinSelection struct {
	js *JoinSelection
}

// NaturalInner performs inner join between js.t and t by columns having the same names.
func (js *MustJoinSelection) NaturalInner(t *Table) *Table {
	table, err := js.js.NaturalInner(t)
	handle(err)
	return table
}

// NaturalLeftOuter performs left outer join between js.t and t by columns having the same names.
func (js *MustJoinSelection) NaturalLeftOuter(t *Table) *Table {
	table, err := js.js.NaturalLeftOuter(t)
	handle(err)
	return table
}

// On sets left table columns for a join.
func (js *MustJoinSelection) On(cols ...ColumnName) *MustJoinOn {
	return &MustJoinOn{jo: js.js.On(cols...)}
}

// MustJoinOn is a proxy for JoinOn.
type MustJoinOn struct {
	jo *JoinOn
}

// Inner sets a right table and columns for inner join and performs the join itself.
func (jo *MustJoinOn) Inner(t *Table, cols ...ColumnName) *Table {
	table, err := jo.jo.Inner(t, cols...)
	handle(err)
	return table
}

// LeftOuter sets a right table and columns for left outer join and performs the join itself.
func (jo *MustJoinOn) LeftOuter(t *Table, cols ...ColumnName) *Table {
	table, err := jo.jo.LeftOuter(t, cols...)
	handle(err)
	return table
}

// Foreach

// Foreach returns a proxy for *Table.Foreach.
func (m MustTable) Foreach() *MustForeachSelection {
	return &MustForeachSelection{m.Table.Foreach()}
}

// MustForeachSelection is a proxy for ForeachSelection.
type MustForeachSelection struct {
	fs *ForeachSelection
}

// By performs Foreach on the whole table rows.
// For wide tables this might be inefficient, consider using On(...) instead.
func (fs *MustForeachSelection) By(fn func(r Row)) {
	fs.fs.By(fn)
}

// MustForeachOn is a proxy for ForeachOn.
type MustForeachOn struct {
	on *ForeachOn
}

// On selects columns for iterating on.
func (fs *MustForeachSelection) On(cols ...ColumnName) *MustForeachOn {
	return &MustForeachOn{on: fs.fs.On(cols...)}
}

// Interface invokes a user-supplied function passing the named column values as interface{} arguments, including nil.
// If given any SchemaAssertions, they are called in the beginning and may have side effects.
func (on *MustForeachOn) Interface(fn func(v ...interface{}), assertions ...SchemaAssertion) {
	err := on.on.Interface(fn, assertions...)
	handle(err)
}
