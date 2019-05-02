package num

import (
	"gonum.org/v1/gonum/floats"
)

// `gonum.org/v1/gonum/floats` extensions - mostly float vector element-wise operations.

// Add adds two vectors.
func Add(a, b []float64) []float64 {
	return floats.AddTo(make([]float64, len(a)), a, b)
}

// AddConst adds a scalar to a vector.
func AddConst(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b
	}
	return result
}

// Sub subtracts one vector from another.
func Sub(a, b []float64) []float64 {
	return floats.SubTo(make([]float64, len(a)), a, b)
}

// SubConst subtracts a scalar from a vector.
func SubConst(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] - b
	}
	return result
}

// SubFromConst subtracts a vector from a scalar.
func SubFromConst(a float64, b []float64) []float64 {
	result := make([]float64, len(b))
	for i := range b {
		result[i] = a - b[i]
	}
	return result
}

// Mul multiplies two vectors.
func Mul(a, b []float64) []float64 {
	return floats.MulTo(make([]float64, len(a)), a, b)
}

// MulByConst multiplies a vector by a scalar.
func MulByConst(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] * b
	}
	return result
}

// Div divides one vector by another.
func Div(a, b []float64) []float64 {
	return floats.DivTo(make([]float64, len(a)), a, b)
}

// DivByConst divides a vector by a scalar.
func DivByConst(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] / b
	}
	return result
}

// DivConst divides a scalar by a vector.
func DivConst(a float64, b []float64) []float64 {
	result := make([]float64, len(b))
	for i := range b {
		result[i] = a / b[i]
	}
	return result
}
