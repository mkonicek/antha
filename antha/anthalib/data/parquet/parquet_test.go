package parquet

import (
	"bytes"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/data"
)

func TestParquet(t *testing.T) {
	// creating a Table
	table := data.NewTable([]*data.Series{
		data.Must().NewSeriesFromSlice("bool_column", []bool{true, true, false, false, true}, nil),
		data.Must().NewSeriesFromSlice("int64_column", []int64{10, 10, 30, -1, 5}, []bool{true, true, true, false, true}),
		data.Must().NewSeriesFromSlice("float32_column", []float64{1.5, 2.5, 3.5, math.NaN(), 5.5}, []bool{true, true, true, false, true}),
		data.Must().NewSeriesFromSlice("string_column", []string{"", "aa", "xx", "aa", ""}, nil),
		data.Must().NewSeriesFromSlice("timestamp_millis_column", []data.TimestampMillis{1, 2, 3, 4, 5}, nil),
		data.Must().NewSeriesFromSlice("timestamp_micros_column", []data.TimestampMicros{1000, 2000, 3000, 4000, 5000}, nil),
	})

	// file: write + read
	fileName := parquetFileName(t)
	defer os.Remove(fileName)

	if err := TableToFile(table, fileName); err != nil {
		t.Errorf("TableToFile: %s", err)
	}

	readTable, err := TableFromFile(fileName)
	if err != nil {
		t.Errorf("TableFromFile: %s", err)
	}

	assertEqual(t, table, readTable, "tables are different after serialization to a file")

	// bytes: write + read
	blob, err := TableToBytes(table)
	if err != nil {
		t.Errorf("TableToBytes: %s", err)
	}

	readTable, err = TableFromBytes(blob)
	if err != nil {
		t.Errorf("TableFromBytes: %s", err)
	}

	assertEqual(t, table, readTable, "tables are different after serialization to a memory buffer")

	// write to io.Writer + read from io.Reader
	buffer := bytes.NewBuffer(nil)

	if err := TableToWriter(table, buffer); err != nil {
		t.Errorf("TableToWriter: %s", err)
	}

	readTable, err = TableFromReader(buffer)
	if err != nil {
		t.Errorf("TableFromReader: %s", err)
	}

	assertEqual(t, table, readTable, "tables are different after serialization to io.Writer")
}

func parquetFileName(t *testing.T) string {
	f, err := ioutil.TempFile("", "table*.parquet")
	if err != nil {
		t.Errorf("create temp file: %s", err)
	}
	defer f.Close() //nolint
	return f.Name()
}

func assertEqual(t *testing.T, expected, actual *data.Table, msg string) {
	if !actual.Equal(expected) {
		t.Error(msg)
		t.Log("actual", actual.Head(20).ToRows())
	}
}
