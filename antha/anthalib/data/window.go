package data

// import (
// 	"time"
// )

// type WindowSpec interface {
// }

// // TODO window semantics
// func (t *Table) Window(spec WindowSpec) *Window {
// 	return nil
// }

// // TimestampWindowSpec implements windowing by equal sized time buckets (Tumble)
// type TimestampWindowSpec struct {
// 	Column   ColumnName
// 	Duration time.Duration
// }

// var _ WindowSpec = (*TimestampWindowSpec)(nil)

// // Window provides aggregate functions over row groups
// type Window struct {
// 	t *Table
// 	w WindowSpec
// }

// // Reduce
// func (w *Window) Reduce(f func(r1, r2 Row) Row) *Table {
// 	return nil
// }

// //Aggregate() *Table

// // Fold
// func (w *Window) Fold(initial Row, f func(r1, r2 Row) Row) *Table {
// 	return nil
// }

// // TODO code gen for the static types!
// func (w *Window) FoldInt32Col(initial int32, f func(r1 []int32) int32) *Table {
// 	return nil
// }
