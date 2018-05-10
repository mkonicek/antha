package wtype

import (
	"fmt"
	"strings"
	"testing"
)

type testAddressable struct {
	rows int
	cols int
}

func (self *testAddressable) AddressExists(wc WellCoords) bool {
	return wc.X >= 0 && wc.X < self.cols && wc.Y >= 0 && wc.Y < self.rows
}

func (self *testAddressable) NRows() int {
	return self.rows
}

func (self *testAddressable) NCols() int {
	return self.cols
}

func (self *testAddressable) GetChildByAddress(WellCoords) LHObject {
	return (LHObject)(nil)
}

func (self *testAddressable) CoordsToWellCoords(c Coordinates) (WellCoords, Coordinates) {
	return WellCoords{}, Coordinates{}
}

func (self *testAddressable) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
	return Coordinates{}, true
}

type addressIteratorTest struct {
	TestName  string
	It        AddressIterator
	CallLimit int
	Expected  string
}

func (self *addressIteratorTest) Run(t *testing.T) {
	callLimit := self.CallLimit
	if callLimit <= 0 {
		callLimit = 12
	}

	got := make([]string, 0, callLimit)
	it := self.It
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		got = append(got, wc.FormatA1())
		if len(got) >= callLimit {
			break
		}
	}

	gotS := strings.Join(got, ",")
	if gotS != self.Expected {
		t.Errorf("%s: \n\twant : \"%s\"\n\t got : \"%s\"", self.TestName, self.Expected, gotS)
	}
}

func TestAddressIterator(t *testing.T) {
	addr := &testAddressable{2, 3}

	tests := []*addressIteratorTest{
		{
			TestName: "LeftRight-TopBottom",
			It:       NewAddressIterator(addr, RowWise, TopToBottom, LeftToRight, false),
			Expected: "A1,A2,A3,B1,B2,B3",
		},
		{
			TestName: "RightLeft-TopBottom",
			It:       NewAddressIterator(addr, RowWise, TopToBottom, RightToLeft, false),
			Expected: "A3,A2,A1,B3,B2,B1",
		},
		{
			TestName: "RightLeft-BottomTop",
			It:       NewAddressIterator(addr, RowWise, BottomToTop, RightToLeft, false),
			Expected: "B3,B2,B1,A3,A2,A1",
		},
		{
			TestName: "LeftRight-BottomTop",
			It:       NewAddressIterator(addr, RowWise, BottomToTop, LeftToRight, false),
			Expected: "B1,B2,B3,A1,A2,A3",
		},
		{
			TestName: "TopBottom-LeftRight",
			It:       NewAddressIterator(addr, ColumnWise, TopToBottom, LeftToRight, false),
			Expected: "A1,B1,A2,B2,A3,B3",
		},
		{
			TestName: "TopBottom-RightLeft",
			It:       NewAddressIterator(addr, ColumnWise, TopToBottom, RightToLeft, false),
			Expected: "A3,B3,A2,B2,A1,B1",
		},
		{
			TestName: "BottomTop-RightLeft",
			It:       NewAddressIterator(addr, ColumnWise, BottomToTop, RightToLeft, false),
			Expected: "B3,A3,B2,A2,B1,A1",
		},
		{
			TestName: "BottomTop-LeftRight",
			It:       NewAddressIterator(addr, ColumnWise, BottomToTop, LeftToRight, false),
			Expected: "B1,A1,B2,A2,B3,A3",
		},
		{
			TestName: "LeftRight-TopBottom",
			It:       NewAddressIterator(addr, RowWise, TopToBottom, LeftToRight, true),
			Expected: "A1,A2,A3,B1,B2,B3,A1,A2,A3,B1,B2,B3",
		},
		{
			TestName: "RightLeft-TopBottom",
			It:       NewAddressIterator(addr, RowWise, TopToBottom, RightToLeft, true),
			Expected: "A3,A2,A1,B3,B2,B1,A3,A2,A1,B3,B2,B1",
		},
		{
			TestName: "RightLeft-BottomTop",
			It:       NewAddressIterator(addr, RowWise, BottomToTop, RightToLeft, true),
			Expected: "B3,B2,B1,A3,A2,A1,B3,B2,B1,A3,A2,A1",
		},
		{
			TestName: "LeftRight-BottomTop",
			It:       NewAddressIterator(addr, RowWise, BottomToTop, LeftToRight, true),
			Expected: "B1,B2,B3,A1,A2,A3,B1,B2,B3,A1,A2,A3",
		},
		{
			TestName: "TopBottom-LeftRight",
			It:       NewAddressIterator(addr, ColumnWise, TopToBottom, LeftToRight, true),
			Expected: "A1,B1,A2,B2,A3,B3,A1,B1,A2,B2,A3,B3",
		},
		{
			TestName: "TopBottom-RightLeft",
			It:       NewAddressIterator(addr, ColumnWise, TopToBottom, RightToLeft, true),
			Expected: "A3,B3,A2,B2,A1,B1,A3,B3,A2,B2,A1,B1",
		},
		{
			TestName: "BottomTop-RightLeft",
			It:       NewAddressIterator(addr, ColumnWise, BottomToTop, RightToLeft, true),
			Expected: "B3,A3,B2,A2,B1,A1,B3,A3,B2,A2,B1,A1",
		},
		{
			TestName: "BottomTop-LeftRight",
			It:       NewAddressIterator(addr, ColumnWise, BottomToTop, LeftToRight, true),
			Expected: "B1,A1,B2,A2,B3,A3,B1,A1,B2,A2,B3,A3",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}

}

type addressSliceIteratorTest struct {
	TestName  string
	It        AddressSliceIterator
	CallLimit int
	Expected  string
}

func (self *addressSliceIteratorTest) Run(t *testing.T) {
	callLimit := self.CallLimit
	if callLimit <= 0 {
		callLimit = 6
	}

	got := make([]string, 0, callLimit)
	it := self.It
	for wellCoords := it.Curr(); it.Valid(); wellCoords = it.Next() {
		wcS := make([]string, 0, len(wellCoords))
		for _, wc := range wellCoords {
			wcS = append(wcS, wc.FormatA1())
		}
		got = append(got, strings.Join(wcS, ","))
		if len(got) >= callLimit {
			break
		}
	}

	gotS := fmt.Sprintf("[%s]", strings.Join(got, "],["))
	if gotS != self.Expected {
		t.Errorf("%s: \n\twant : \"%s\"\n\t got : \"%s\"", self.TestName, self.Expected, gotS)
	}
}

func TestAddressSliceIterator(t *testing.T) {
	addr := &testAddressable{2, 3}

	//A1, A2, A3
	//B1, B2, B3

	tests := []*addressSliceIteratorTest{
		{
			TestName: "Cols:TopToBottom, LeftToRight, repeat = false",
			It:       NewColumnIterator(addr, TopToBottom, LeftToRight, false),
			Expected: "[A1,B1],[A2,B2],[A3,B3]",
		},
		{
			TestName: "Cols:BottomToTop, LeftToRight, repeat = false",
			It:       NewColumnIterator(addr, BottomToTop, LeftToRight, false),
			Expected: "[B1,A1],[B2,A2],[B3,A3]",
		},
		{
			TestName: "Cols:TopToBottom, RightToLeft, repeat = false",
			It:       NewColumnIterator(addr, TopToBottom, RightToLeft, false),
			Expected: "[A3,B3],[A2,B2],[A1,B1]",
		},
		{
			TestName: "Cols:BottomToTop, RightToLeft, repeat = false",
			It:       NewColumnIterator(addr, BottomToTop, RightToLeft, false),
			Expected: "[B3,A3],[B2,A2],[B1,A1]",
		},
		{
			TestName: "Rows:TopToBottom, LeftToRight, repeat = false",
			It:       NewRowIterator(addr, TopToBottom, LeftToRight, false),
			Expected: "[A1,A2,A3],[B1,B2,B3]",
		},
		{
			TestName: "Rows:BottomToTop, LeftToRight, repeat = false",
			It:       NewRowIterator(addr, BottomToTop, LeftToRight, false),
			Expected: "[B1,B2,B3],[A1,A2,A3]",
		},
		{
			TestName: "Rows:TopToBottom, RightToLeft, repeat = false",
			It:       NewRowIterator(addr, TopToBottom, RightToLeft, false),
			Expected: "[A3,A2,A1],[B3,B2,B1]",
		},
		{
			TestName: "Rows:BottomToTop, RightToLeft, repeat = false",
			It:       NewRowIterator(addr, BottomToTop, RightToLeft, false),
			Expected: "[B3,B2,B1],[A3,A2,A1]",
		},
		{
			TestName:  "Cols:TopToBottom, LeftToRight, repeat = true",
			It:        NewColumnIterator(addr, TopToBottom, LeftToRight, true),
			Expected:  "[A1,B1],[A2,B2],[A3,B3],[A1,B1],[A2,B2],[A3,B3]",
			CallLimit: 6,
		},
		{
			TestName:  "Cols:BottomToTop, LeftToRight, repeat = true",
			It:        NewColumnIterator(addr, BottomToTop, LeftToRight, true),
			Expected:  "[B1,A1],[B2,A2],[B3,A3],[B1,A1],[B2,A2],[B3,A3]",
			CallLimit: 6,
		},
		{
			TestName:  "Cols:TopToBottom, RightToLeft, repeat = true",
			It:        NewColumnIterator(addr, TopToBottom, RightToLeft, true),
			Expected:  "[A3,B3],[A2,B2],[A1,B1],[A3,B3],[A2,B2],[A1,B1]",
			CallLimit: 6,
		},
		{
			TestName:  "Cols:BottomToTop, RightToLeft, repeat = true",
			It:        NewColumnIterator(addr, BottomToTop, RightToLeft, true),
			Expected:  "[B3,A3],[B2,A2],[B1,A1],[B3,A3],[B2,A2],[B1,A1]",
			CallLimit: 6,
		},
		{
			TestName:  "Rows:TopToBottom, LeftToRight, repeat = true",
			It:        NewRowIterator(addr, TopToBottom, LeftToRight, true),
			Expected:  "[A1,A2,A3],[B1,B2,B3],[A1,A2,A3],[B1,B2,B3]",
			CallLimit: 4,
		},
		{
			TestName:  "Rows:BottomToTop, LeftToRight, repeat = true",
			It:        NewRowIterator(addr, BottomToTop, LeftToRight, true),
			Expected:  "[B1,B2,B3],[A1,A2,A3],[B1,B2,B3],[A1,A2,A3]",
			CallLimit: 4,
		},
		{
			TestName:  "Rows:TopToBottom, RightToLeft, repeat = true",
			It:        NewRowIterator(addr, TopToBottom, RightToLeft, true),
			Expected:  "[A3,A2,A1],[B3,B2,B1],[A3,A2,A1],[B3,B2,B1]",
			CallLimit: 4,
		},
		{
			TestName:  "Rows:BottomToTop, RightToLeft, repeat = true",
			It:        NewRowIterator(addr, BottomToTop, RightToLeft, true),
			Expected:  "[B3,B2,B1],[A3,A2,A1],[B3,B2,B1],[A3,A2,A1]",
			CallLimit: 4,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}

}
