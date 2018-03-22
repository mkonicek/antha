package spreadsheet

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// Parses an a1 style excel cell coordinate into ints for row and column for use by plotinum library
// note that 1 is subtracted from the column number in accordance with the go convention of counting from 0

type coordinatetest struct {
	a1format string
	row      int
	col      int
}

var coordinatetests = []coordinatetest{
	{a1format: "a1",
		row: 0,
		col: 0,
	},
	{a1format: "b1",
		row: 1,
		col: 0,
	},
	{a1format: "aa1",
		row: 26,
		col: 0,
	},
	{a1format: "a2",
		row: 0,
		col: 1,
	},
}

type spreadSheetTest struct {
	fileName    string
	testSheet   int
	testRows    map[int][]interface{}
	testColumns map[int][]interface{}
}

var tests = []spreadSheetTest{
	{
		fileName:  "xlsxParserTestFile.xlsx",
		testSheet: 0,
		testRows: map[int][]interface{}{
			0: {"Well", "A String Header", "A Number Header"},
			1: {"A1", "High", 1},
		},
		testColumns: map[int][]interface{}{
			0: {"Well", "A1", "A2", "A3", "D1"},
			1: {"A String Header", "High", "High", "low", "low"},
		},
	},
}

func TestA1formattorowcolumn(t *testing.T) {
	for _, test := range coordinatetests {
		r, c, _ := A1FormatToRowColumn(test.a1format)
		if c != test.col {
			t.Error(
				"For", test.a1format, "\n",
				"expected", test.col, "\n",
				"got", c, "\n",
			)
		}
		if r != test.row {
			t.Error(
				"For", test.a1format, "\n",
				"expected", test.row, "\n",
				"got", r, "\n",
			)
		}
	}
}

func TestOpenXLSXBinary(t *testing.T) {
	for _, test := range tests {
		data, err := ioutil.ReadFile(test.fileName)
		if err != nil {
			t.Fatal(err)
		}
		xlsx, err := OpenXLSXBinary(data)
		if err != nil {
			t.Error(err.Error())
			xlsx, err = OpenXLSXFromFileName(test.fileName)
			if err != nil {
				t.Error(err.Error())
				break
			}
		}
		sheet, err := Sheet(xlsx, test.testSheet)
		if err != nil {
			t.Error(err.Error())
			break
		}

		for rowIndex, values := range test.testRows {
			rowValues, err := Row(sheet, rowIndex)
			if err != nil {
				t.Error("test error ", err.Error())
			}
			for i := range rowValues {
				if i >= len(values) {
					t.Error("test error ")
				}
				if fmt.Sprint(rowValues[i]) != fmt.Sprint(values[i]) {
					t.Error(
						"For", test.fileName, "row", rowIndex, "\n",
						"expected:", values[i], "\n",
						"got", rowValues[i], "\n",
					)
				}
			}
		}

		for columnIndex, values := range test.testColumns {
			colValues, err := Column(sheet, columnIndex)
			if err != nil {
				t.Error(err.Error())
			}

			for i := range colValues {
				if i >= len(values) {
					t.Error("test error ")
				}
				if fmt.Sprint(colValues[i]) != fmt.Sprint(values[i]) {
					t.Error(
						"For", test.fileName, "column", columnIndex, "\n",
						"expected:", values[i], "\n",
						"got", colValues[i], "\n",
					)
				}
			}
		}

		dataMap, err := ToHeaderDataMap(sheet, 0)
		if err != nil {
			t.Error(err.Error())
		}

		fmt.Println(SheetToCSV(sheet))
		fmt.Println("dataMap: ", dataMap)
		// Output:
		// [[Well A String Header A Number Header] [A1 High 1] [A2 High 1.000000e+04] [A3 low 1] [D1 low 0.4]]
		// dataMap:  map[A String Header:[High High low low] A Number Header:[1 1.000000e+04 1 0.4] Well:[A1 A2 A3 D1]]

	}
}
