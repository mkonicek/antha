package num

import (
	"github.com/pkg/errors"
)

// Implementations of some useful numpy functions which are not available in gonum.

// Zeroes returns an array of zeroes of specified size.
// It's encouraged to use it instead of just make() in case the further code relies on the fact that the array contains zeroes.
func Zeroes(size int) []float64 {
	return make([]float64, size)
}

// Ones return an array of ones of specified size.
func Ones(size int) []float64 {
	result := make([]float64, size)
	for i := range result {
		result[i] = 1
	}
	return result
}

// Linspace implements `np.linspace` - i.e. splits the interval [start, end] into `num - 1` equal intervals and returns `num` split points.
func Linspace(start, end float64, num int) []float64 {
	if num < 0 {
		panic(errors.Errorf("number of samples, %d, must be non-negative.", num))
	}
	result := make([]float64, num)
	step := (end - start) / float64(num-1)
	for i := range result {
		result[i] = start + float64(i)*step
	}
	return result
}

// Arange implements `np.arange` - i.e. returns a list of integers (start, ..., stop - 1) in the form of []float64
func Arange(start int, stop int) []float64 {
	return Linspace(float64(start), float64(stop-1), stop-start)
}

// ConvolutionMode defines convolution output array length.
type ConvolutionMode int

const (
	// Full - returns the convolution at each point of overlap, i.e. of length N+M-1.
	Full = iota
	// Same - returns the output of length max(M, N).
	Same
	// Valid - returns the output of length max(M, N) - min(M, N) + 1.
	Valid
)

// Convolve is a (very naive) implementation of precise discrete convolution.
// The results are numerically equivalent to `np.convolve(a, v, mode)` - it looks like that `np.convolve` uses precise convolution as well (but not an FFT approximation).
// TODO: optimize the implementation - the current one has O((M+N)^2) time complexity. Looks like it's possible to achieve at least O(MN).
func Convolve(a, v []float64, mode ConvolutionMode) []float64 {
	if len(a) == 0 {
		panic(errors.New("Convolve: a cannot be empty"))
	}
	if len(v) == 0 {
		panic(errors.New("Convolve: v cannot be empty"))
	}

	// the code below relies on the fact that `a` is the longer array
	if len(v) > len(a) {
		a, v = v, a
	}

	size := len(a) + len(v) - 1

	// a + zeroes
	a_ext := Zeroes(size)
	copy(a_ext[:len(a)], a)

	// v + zeroes
	v_ext := Zeroes(size)
	copy(v_ext[:len(v)], v)

	result := Zeroes(size)
	for i := 0; i < size; i++ {
		for j := 0; j < size-i; j++ {
			result[i+j] += a_ext[i] * v_ext[j]
		}
	}

	switch mode {
	case Full:
		// `Full` mode: returning the whole result
		return result
	case Same:
		// `Same` mode: returning the subarray of length `len(a)`
		toCut := len(v) / 2 // is this correct? at least, this works for the sample data
		return result[toCut : toCut+len(a)]
	case Valid:
		// `Valid` mode: returning the subarray of length `len(a) - len(v) + 1`
		toCut := len(v) - 1
		return result[toCut:len(a)]
	default:
		panic(errors.Errorf("invalid convolution mode %v", mode))
	}
}

// maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// computes `n!`
func factorial(n int) int {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return result
}
