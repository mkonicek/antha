package data

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

var output *Table

func BenchmarkDistinctNumeric(b *testing.B) {
	b.StopTimer()
	input := genericInputTable()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		output = input.Must().Distinct().On("IntCol", "FloatCol").Must().Cache()
	}
}

func genericInputTable() *Table {
	return generateTable(10000, newInt64Gen("IntCol", 0.1), newFloat64Gen("FloatCol", 0.1), newStringGen("StringCol", 0.1, 15))
}

func BenchmarkDistinctHeterogenous(b *testing.B) {
	b.StopTimer()
	input := genericInputTable()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		output = input.Must().Distinct().On("IntCol", "FloatCol", "StringCol").Must().Cache()
	}
}

func BenchmarkSort(b *testing.B) {
	b.StopTimer()
	input := genericInputTable()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		output = input.Must().Sort(Key{{"IntCol", true}, {"FloatCol", false}, {"StringCol", true}}).Must().Cache()
	}
}

func BenchmarkPivot(b *testing.B) {
	b.StopTimer()
	input := pivotInputTable()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		output = input.Must().Pivot().Key("Key").Columns("FieldName", "FieldValue").Must().Cache()
	}
}

func pivotInputTable() *Table {
	size := 10000
	numOfCols := 60
	stringLen := 15
	return generateTable(size,
		newRestrictedStringGen("Key", size/numOfCols, stringLen),
		newRestrictedStringGen("FieldName", numOfCols, stringLen),
		newStringGen("FieldValue", 0.1, stringLen),
	)
}

func BenchmarkJoin(b *testing.B) {
	b.StopTimer()
	left, right := joinInputTables()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		output = left.Must().Join().NaturalInner(right).Must().Cache()
	}
}

func joinInputTables() (*Table, *Table) {
	leftSize := 10000
	rightSize := 20
	stringLen := 15

	// key is a string column with limited number of distinct values
	keyColGen := newRestrictedStringGen("Key", rightSize, stringLen)

	return generateTable(leftSize,
			keyColGen,
			newInt64Gen("IntCol", 0.1),
			newFloat64Gen("FloatCol", 0.1),
		), generateTable(rightSize,
			keyColGen,
			newInt64Gen("OtherIntCol", 0.1),
		)
}

/*
 * table generation tools
 */

// the only random generator with a fixed seed used for all generations (for the sake of determinism)
var r *rand.Rand = rand.New(rand.NewSource(0))

func generateTable(size int, generators ...seriesGen) *Table {
	// creating columns descriptions
	columns := make([]Column, len(generators))
	for i, g := range generators {
		columns[i] = Column{
			Name: g.name,
			Type: g.typ,
		}
	}

	// initializing the table builder
	builder := Must().NewTableBuilder(columns)
	builder.Reserve(size)

	// filling the table
	row := make([]interface{}, len(generators))
	for i := 0; i < size; i++ {
		for j, g := range generators {
			row[j] = g.next()
		}
		builder.Append(row)
	}

	// building the table
	return builder.Build()
}

// random series generator
type seriesGen struct {
	name ColumnName
	typ  reflect.Type
	next func() interface{}
}

// generates a series containing random int64s
func newInt64Gen(name ColumnName, nullProbability float64) seriesGen {
	return seriesGen{
		name: name,
		typ:  reflect.TypeOf(int64(0)),
		next: func() interface{} {
			if genIsNull(nullProbability) {
				return nil
			}
			return r.Int63()
		},
	}
}

// generates a series containing random float64s
func newFloat64Gen(name ColumnName, nullProbability float64) seriesGen {
	return seriesGen{
		name: name,
		typ:  reflect.TypeOf(float64(0)),
		next: func() interface{} {
			if genIsNull(nullProbability) {
				return nil
			}
			return r.Float64()
		},
	}
}

// generates a series containing random strings
func newStringGen(name ColumnName, nullProbability float64, maxLength int) seriesGen {
	return seriesGen{
		name: name,
		typ:  reflect.TypeOf(""),
		next: func() interface{} {
			if genIsNull(nullProbability) {
				return nil
			}
			return randomString(maxLength)
		},
	}
}

// generates a series containing numOfDistinct distinct strings only
func newRestrictedStringGen(name ColumnName, numOfDistinct int, maxLength int) seriesGen {
	strings := make([]string, numOfDistinct)
	for i := range strings {
		strings[i] = randomString(maxLength)
	}

	return seriesGen{
		name: name,
		typ:  reflect.TypeOf(""),
		next: func() interface{} {
			return strings[r.Intn(numOfDistinct)]
		},
	}
}

// 64 letters
const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

// generates random string
func randomString(maxLength int) string {
	// generate a slice of bytes of length [0, maxLength - 1]
	length := r.Intn(maxLength)
	bytes := make([]byte, length)
	_, err := rand.Read(bytes) //nolint
	handle(err)
	// transform into string
	for i, b := range bytes {
		bytes[i] = letters[b%64]
	}
	return string(bytes)
}

func genIsNull(nullProbability float64) bool {
	return float64(r.Int31()) < math.MaxInt32*nullProbability
}
