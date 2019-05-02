/*
Package num provides tools for floats arithmetic.

Overview

The aim of this package is to provide Antha users with a library containing basic float vectors operations.

The package supplements the gonum package and provides a simpler interface for some of its parts.
While gonum focuses on efficiency and tends to avoid allocations inside its functions,
the num package allows to do multiple element-wise operations more concisely (though at the cost of some additional allocations).

Example:
	x := []float64{1, 2, 3}
	linear_transformed_x := num.Add(num.Mul(a, x), b)

Also num is aimed at implementing at least some basic numpy functionality currently missing in gonum.
*/
package num
