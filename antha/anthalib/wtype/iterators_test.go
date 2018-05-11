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

func TestTickingIterator(t *testing.T) {
	addr := &testAddressable{4, 4}

	//A1, A2, A3, A4
	//B1, B2, B3, B4
	//C1, C2, C3, C4
	//D1, D2, D3, D4

	tests := []*addressSliceIteratorTest{
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 1, 1)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 1, 1),
			Expected:  "[A1],[B1],[C1],[D1],[A2],[B2],[C2],[D2],[A3],[B3],[C3],[D3],[A4],[B4],[C4],[D4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 1)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 1),
			Expected:  "[A1],[C1],[A2],[C2],[A3],[C3],[A4],[C4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 2),
			Expected:  "[A1],[A1],[C1],[C1],[A2],[A2],[C2],[C2],[A3],[A3],[C3],[C3],[A4],[A4],[C4],[C4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 1)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 1),
			Expected:  "[A1,C1],[A2,C2],[A3,C3],[A4,C4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 2),
			Expected:  "[A1,A1],[C1,C1],[A2,A2],[C2,C2],[A3,A3],[C3,C3],[A4,A4],[C4,C4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 4, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 4, 2, 2),
			Expected:  "[A1,A1,C1,C1],[A2,A2,C2,C2],[A3,A3,C3,C3],[A4,A4,C4,C4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, true, 1, 1, 1)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, true, 1, 1, 1),
			Expected:  "[A1],[B1],[C1],[D1],[A2],[B2],[C2],[D2],[A3],[B3],[C3],[D3],[A4],[B4],[C4],[D4],[A1],[B1],[C1],[D1],[A2],[B2],[C2],[D2],[A3],[B3],[C3],[D3],[A4],[B4],[C4],[D4]",
			CallLimit: 32,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}

}

func TestTickingIteratorLarger(t *testing.T) {
	tests := []*addressSliceIteratorTest{
		{
			TestName:  "Test NewTickingIterator(&testAddressable{8, 12}, ColumnWise, TopToBottom, LeftToRight, false, 8, 1, 1)",
			It:        NewTickingIterator(&testAddressable{8, 12}, ColumnWise, TopToBottom, LeftToRight, false, 8, 1, 1),
			Expected:  "[A1,B1,C1,D1,E1,F1,G1,H1],[A2,B2,C2,D2,E2,F2,G2,H2],[A3,B3,C3,D3,E3,F3,G3,H3],[A4,B4,C4,D4,E4,F4,G4,H4],[A5,B5,C5,D5,E5,F5,G5,H5],[A6,B6,C6,D6,E6,F6,G6,H6],[A7,B7,C7,D7,E7,F7,G7,H7],[A8,B8,C8,D8,E8,F8,G8,H8],[A9,B9,C9,D9,E9,F9,G9,H9],[A10,B10,C10,D10,E10,F10,G10,H10],[A11,B11,C11,D11,E11,F11,G11,H11],[A12,B12,C12,D12,E12,F12,G12,H12]",
			CallLimit: 106,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}

}
