package wtype

import (
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
			It:       GetAddressIterator(addr, RowWise, TopToBottom, LeftToRight, false),
			Expected: "A1,A2,A3,B1,B2,B3",
		},
		{
			TestName: "RightLeft-TopBottom",
			It:       GetAddressIterator(addr, RowWise, TopToBottom, RightToLeft, false),
			Expected: "A3,A2,A1,B3,B2,B1",
		},
		{
			TestName: "RightLeft-BottomTop",
			It:       GetAddressIterator(addr, RowWise, BottomToTop, RightToLeft, false),
			Expected: "B3,B2,B1,A3,A2,A1",
		},
		{
			TestName: "LeftRight-BottomTop",
			It:       GetAddressIterator(addr, RowWise, BottomToTop, LeftToRight, false),
			Expected: "B1,B2,B3,A1,A2,A3",
		},
		{
			TestName: "TopBottom-LeftRight",
			It:       GetAddressIterator(addr, ColumnWise, TopToBottom, LeftToRight, false),
			Expected: "A1,B1,A2,B2,A3,B3",
		},
		{
			TestName: "TopBottom-RightLeft",
			It:       GetAddressIterator(addr, ColumnWise, TopToBottom, RightToLeft, false),
			Expected: "A3,B3,A2,B2,A1,B1",
		},
		{
			TestName: "BottomTop-RightLeft",
			It:       GetAddressIterator(addr, ColumnWise, BottomToTop, RightToLeft, false),
			Expected: "B3,A3,B2,A2,B1,A1",
		},
		{
			TestName: "BottomTop-LeftRight",
			It:       GetAddressIterator(addr, ColumnWise, BottomToTop, LeftToRight, false),
			Expected: "B1,A1,B2,A2,B3,A3",
		},
		{
			TestName: "LeftRight-TopBottom",
			It:       GetAddressIterator(addr, RowWise, TopToBottom, LeftToRight, true),
			Expected: "A1,A2,A3,B1,B2,B3,A1,A2,A3,B1,B2,B3",
		},
		{
			TestName: "RightLeft-TopBottom",
			It:       GetAddressIterator(addr, RowWise, TopToBottom, RightToLeft, true),
			Expected: "A3,A2,A1,B3,B2,B1,A3,A2,A1,B3,B2,B1",
		},
		{
			TestName: "RightLeft-BottomTop",
			It:       GetAddressIterator(addr, RowWise, BottomToTop, RightToLeft, true),
			Expected: "B3,B2,B1,A3,A2,A1,B3,B2,B1,A3,A2,A1",
		},
		{
			TestName: "LeftRight-BottomTop",
			It:       GetAddressIterator(addr, RowWise, BottomToTop, LeftToRight, true),
			Expected: "B1,B2,B3,A1,A2,A3,B1,B2,B3,A1,A2,A3",
		},
		{
			TestName: "TopBottom-LeftRight",
			It:       GetAddressIterator(addr, ColumnWise, TopToBottom, LeftToRight, true),
			Expected: "A1,B1,A2,B2,A3,B3,A1,B1,A2,B2,A3,B3",
		},
		{
			TestName: "TopBottom-RightLeft",
			It:       GetAddressIterator(addr, ColumnWise, TopToBottom, RightToLeft, true),
			Expected: "A3,B3,A2,B2,A1,B1,A3,B3,A2,B2,A1,B1",
		},
		{
			TestName: "BottomTop-RightLeft",
			It:       GetAddressIterator(addr, ColumnWise, BottomToTop, RightToLeft, true),
			Expected: "B3,A3,B2,A2,B1,A1,B3,A3,B2,A2,B1,A1",
		},
		{
			TestName: "BottomTop-LeftRight",
			It:       GetAddressIterator(addr, ColumnWise, BottomToTop, LeftToRight, true),
			Expected: "B1,A1,B2,A2,B3,A3,B1,A1,B2,A2,B3,A3",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}

}
