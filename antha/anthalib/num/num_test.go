package num

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/floats"
)

func TestAdd(t *testing.T) {
	// Add
	actual := Add([]float64{1, 2, 3}, []float64{4, 5, 6})
	expected := []float64{5, 7, 9}
	assertEqual(t, expected, actual, "Add")

	// Add - wrong inputs sizes
	assert.Panics(t, func() {
		Add([]float64{1, 2, 3}, []float64{1, 2})
	}, "Add - wrong inputs sizes")

	// AddConst
	actual = AddConst([]float64{1, 2, 3}, 1)
	expected = []float64{2, 3, 4}
	assertEqual(t, expected, actual, "AddConst")
}

func TestSub(t *testing.T) {
	// Sub
	actual := Sub([]float64{1, 2, 3}, []float64{4, 5, 6})
	expected := []float64{-3, -3, -3}
	assertEqual(t, expected, actual, "Sub")

	// Sub - wrong inputs sizes
	assert.Panics(t, func() {
		Sub([]float64{1, 2, 3}, []float64{1, 2})
	}, "Sub - wrong inputs sizes")

	// SubConst
	actual = SubConst([]float64{1, 2, 3}, 1)
	expected = []float64{0, 1, 2}
	assertEqual(t, expected, actual, "SubConst")

	// SubFromConst
	actual = SubFromConst(1, []float64{1, 2, 3})
	expected = []float64{0, -1, -2}
	assertEqual(t, expected, actual, "SubFromConst")
}

func TestMul(t *testing.T) {
	// Mul
	actual := Mul([]float64{1, 2, 3}, []float64{4, 5, 6})
	expected := []float64{4, 10, 18}
	assertEqual(t, expected, actual, "Mul")

	// Mul - wrong inputs sizes
	assert.Panics(t, func() {
		Mul([]float64{1, 2, 3}, []float64{1, 2})
	}, "Mul - wrong inputs sizes")

	// MulByConst
	actual = MulByConst([]float64{1, 2, 3}, 2)
	expected = []float64{2, 4, 6}
	assertEqual(t, expected, actual, "MulByConst")
}

func TestDiv(t *testing.T) {
	// Div
	actual := Div([]float64{4, 10, 24}, []float64{4, 5, 6})
	expected := []float64{1, 2, 4}
	assertEqual(t, expected, actual, "Div")

	// Div - wrong inputs sizes
	assert.Panics(t, func() {
		Div([]float64{1, 2, 3}, []float64{1, 2})
	}, "Sub - wrong inputs sizes")

	// DivConst
	actual = DivConst(6, []float64{1, 2, 3})
	expected = []float64{6, 3, 2}
	assertEqual(t, expected, actual, "DivConst")

	// DivByConst
	actual = DivByConst([]float64{2, 4, 6}, 2)
	expected = []float64{1, 2, 3}
	assertEqual(t, expected, actual, "DivByConst")
}

func TestZeroes(t *testing.T) {
	// Zeroes
	actual := Zeroes(3)
	expected := []float64{0, 0, 0}
	assertEqual(t, expected, actual, "Zeroes")

	// Zeroes - negative size
	assert.Panics(t, func() {
		_ = Zeroes(-5)
	}, "Zeroes - negative size")
}

func TestOnes(t *testing.T) {
	// Ones
	actual := Ones(3)
	expected := []float64{1, 1, 1}
	assertEqual(t, expected, actual, "Ones")

	// Ones - negative size
	assert.Panics(t, func() {
		_ = Ones(-5)
	}, "Ones - negative size")
}

func TestLinspace(t *testing.T) {
	// Linspace
	actual := Linspace(1, 4, 4)
	expected := []float64{1, 2, 3, 4}
	assertEqual(t, expected, actual, "Linspace")

	// Linspace - negative steps count
	assert.Panics(t, func() {
		_ = Linspace(1, 4, -4)
	}, "Linspace - negative steps count")
}

func TestConvolve(t *testing.T) {
	// Tests from https://docs.scipy.org/doc/numpy/reference/generated/numpy.convolve.html
	actual := Convolve([]float64{1, 2, 3}, []float64{0, 1, 0.5}, Full)
	expected := []float64{0., 1., 2.5, 4., 1.5}
	assertEqual(t, expected, actual, "Convolve (full)")

	actual = Convolve([]float64{1, 2, 3}, []float64{0, 1, 0.5}, Same)
	expected = []float64{1., 2.5, 4.}
	assertEqual(t, expected, actual, "Convolve (same)")

	actual = Convolve([]float64{1, 2, 3}, []float64{0, 1, 0.5}, Valid)
	expected = []float64{2.5}
	assertEqual(t, expected, actual, "Convolve (valid)")

	// Convolve - empty input
	assert.Panics(t, func() {
		_ = Convolve([]float64{}, []float64{1}, Full)
	}, "Convolve - empty input")
	assert.Panics(t, func() {
		_ = Convolve([]float64{1}, []float64{}, Full)
	}, "Convolve - empty input")
}

func TestSavGolFilter(t *testing.T) {
	// Tests from https://github.com/scipy/scipy/blob/master/scipy/signal/_savitzky_golay.py#L307
	x := []float64{2, 2, 5, 2, 1, 0, 1, 4, 9}
	actual := SavGolFilter(x, 5, 2, 0, 1.0)
	expected := []float64{1.66, 3.17, 3.54, 2.86, 0.66, 0.17, 1., 4., 9.}
	// edge values are inavitably different because the convolution algorithm we use here is different from the Python one;
	// so cutting the edge values off
	truncate := func(vector []float64) []float64 {
		const edge = 2
		return vector[edge:(len(vector) - edge)]
	}
	assertEqualWithTolerance(t, truncate(expected), truncate(actual), 0.1, "SavGolFilter")

	// error cases
	// TODO: add more ones

	assert.Panics(t, func() {
		_ = SavGolFilter(x, 4, 2, 0, 1.0)
	}, "no error on even window size")
	assert.Panics(t, func() {
		_ = SavGolFilter(x, 5, 6, 0, 1.0)
	}, "no error on polyorder > window_size")
}

const defaultTolerance = 1e-14

// Compares two float vectors using a hard-coded max discrepancy value.
func assertEqual(t *testing.T, expected, actual []float64, msg string) {
	assertEqualWithTolerance(t, expected, actual, defaultTolerance, msg)
}

// Compares two float vectors using a hard-coded max discrepancy value.
func assertEqualWithTolerance(t *testing.T, expected, actual []float64, tolerance float64, msg string) {
	if !floats.EqualApprox(expected, actual, tolerance) {
		t.Errorf(msg+": expected %v, actual %v", expected, actual)
	}
}
