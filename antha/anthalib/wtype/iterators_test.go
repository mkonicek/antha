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

func (self *testAddressable) CoordsToWellCoords(c Coordinates3D) (WellCoords, Coordinates3D) {
	return WellCoords{}, Coordinates3D{}
}

func (self *testAddressable) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates3D, bool) {
	return Coordinates3D{}, true
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
			Expected:  "[A1],[C1],[B1],[D1],[A2],[C2],[B2],[D2],[A3],[C3],[B3],[D3],[A4],[C4],[B4],[D4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 1, 2, 2),
			Expected:  "[A1],[A1],[C1],[C1],[B1],[B1],[D1],[D1],[A2],[A2],[C2],[C2],[B2],[B2],[D2],[D2],[A3],[A3],[C3],[C3],[B3],[B3],[D3],[D3],[A4],[A4],[C4],[C4],[B4],[B4],[D4],[D4]",
			CallLimit: 64,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 1)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 1),
			Expected:  "[A1,C1],[B1,D1],[A2,C2],[B2,D2],[A3,C3],[B3,D3],[A4,C4],[B4,D4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 2, 2, 2),
			Expected:  "[A1,A1],[C1,C1],[B1,B1],[D1,D1],[A2,A2],[C2,C2],[B2,B2],[D2,D2],[A3,A3],[C3,C3],[B3,B3],[D3,D3],[A4,A4],[C4,C4],[B4,B4],[D4,D4]",
			CallLimit: 32,
		},
		{
			TestName:  "NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 4, 2, 2)",
			It:        NewTickingIterator(addr, ColumnWise, TopToBottom, LeftToRight, false, 4, 2, 2),
			Expected:  "[A1,A1,C1,C1],[B1,B1,D1,D1],[A2,A2,C2,C2],[B2,B2,D2,D2],[A3,A3,C3,C3],[B3,B3,D3,D3],[A4,A4,C4,C4],[B4,B4,D4,D4]",
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
		{
			TestName:  "Test NewTickingIterator(&testAddressable{16, 24}, ColumnWise, TopToBottom, LeftToRight, false, 8, 2, 1)",
			It:        NewTickingIterator(&testAddressable{16, 24}, ColumnWise, TopToBottom, LeftToRight, false, 8, 2, 1),
			Expected:  "[A1,C1,E1,G1,I1,K1,M1,O1],[B1,D1,F1,H1,J1,L1,N1,P1],[A2,C2,E2,G2,I2,K2,M2,O2],[B2,D2,F2,H2,J2,L2,N2,P2],[A3,C3,E3,G3,I3,K3,M3,O3],[B3,D3,F3,H3,J3,L3,N3,P3],[A4,C4,E4,G4,I4,K4,M4,O4],[B4,D4,F4,H4,J4,L4,N4,P4],[A5,C5,E5,G5,I5,K5,M5,O5],[B5,D5,F5,H5,J5,L5,N5,P5],[A6,C6,E6,G6,I6,K6,M6,O6],[B6,D6,F6,H6,J6,L6,N6,P6],[A7,C7,E7,G7,I7,K7,M7,O7],[B7,D7,F7,H7,J7,L7,N7,P7],[A8,C8,E8,G8,I8,K8,M8,O8],[B8,D8,F8,H8,J8,L8,N8,P8],[A9,C9,E9,G9,I9,K9,M9,O9],[B9,D9,F9,H9,J9,L9,N9,P9],[A10,C10,E10,G10,I10,K10,M10,O10],[B10,D10,F10,H10,J10,L10,N10,P10],[A11,C11,E11,G11,I11,K11,M11,O11],[B11,D11,F11,H11,J11,L11,N11,P11],[A12,C12,E12,G12,I12,K12,M12,O12],[B12,D12,F12,H12,J12,L12,N12,P12],[A13,C13,E13,G13,I13,K13,M13,O13],[B13,D13,F13,H13,J13,L13,N13,P13],[A14,C14,E14,G14,I14,K14,M14,O14],[B14,D14,F14,H14,J14,L14,N14,P14],[A15,C15,E15,G15,I15,K15,M15,O15],[B15,D15,F15,H15,J15,L15,N15,P15],[A16,C16,E16,G16,I16,K16,M16,O16],[B16,D16,F16,H16,J16,L16,N16,P16],[A17,C17,E17,G17,I17,K17,M17,O17],[B17,D17,F17,H17,J17,L17,N17,P17],[A18,C18,E18,G18,I18,K18,M18,O18],[B18,D18,F18,H18,J18,L18,N18,P18],[A19,C19,E19,G19,I19,K19,M19,O19],[B19,D19,F19,H19,J19,L19,N19,P19],[A20,C20,E20,G20,I20,K20,M20,O20],[B20,D20,F20,H20,J20,L20,N20,P20],[A21,C21,E21,G21,I21,K21,M21,O21],[B21,D21,F21,H21,J21,L21,N21,P21],[A22,C22,E22,G22,I22,K22,M22,O22],[B22,D22,F22,H22,J22,L22,N22,P22],[A23,C23,E23,G23,I23,K23,M23,O23],[B23,D23,F23,H23,J23,L23,N23,P23],[A24,C24,E24,G24,I24,K24,M24,O24],[B24,D24,F24,H24,J24,L24,N24,P24]",
			CallLimit: 394,
		},
		{
			TestName:  "Test NewTickingIterator(&testAddressable{1, 12}, ColumnWise, TopToBottom, LeftToRight, false, 8, 1, 8)",
			It:        NewTickingIterator(&testAddressable{1, 12}, ColumnWise, TopToBottom, LeftToRight, false, 8, 1, 8),
			Expected:  "[A1,A1,A1,A1,A1,A1,A1,A1],[A2,A2,A2,A2,A2,A2,A2,A2],[A3,A3,A3,A3,A3,A3,A3,A3],[A4,A4,A4,A4,A4,A4,A4,A4],[A5,A5,A5,A5,A5,A5,A5,A5],[A6,A6,A6,A6,A6,A6,A6,A6],[A7,A7,A7,A7,A7,A7,A7,A7],[A8,A8,A8,A8,A8,A8,A8,A8],[A9,A9,A9,A9,A9,A9,A9,A9],[A10,A10,A10,A10,A10,A10,A10,A10],[A11,A11,A11,A11,A11,A11,A11,A11],[A12,A12,A12,A12,A12,A12,A12,A12]",
			CallLimit: 22,
		},
	}

	//writeTITest(make384platefortest(), 8, 2, 1)
	//writeTITest(maketroughfortest(), 8, 1, 8)

	for _, test := range tests {
		test.Run(t)
	}

}
