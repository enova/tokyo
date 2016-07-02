package dbl

// Epsilon Threshold
var e = 1e-10

// SetEpsilon sets the epsilon threshold value for comparing float64's.
// This value is global and cannot be set per thread.
func SetEpsilon(epsilon float64) {
	e = epsilon
}

// GetEpsilon returns the epsilon theshold value
func GetEpsilon() float64 {
	return e
}

// LT x < y
func LT(x float64, y float64) bool {
	return (x < y-e)
}

// GT x > y
func GT(x float64, y float64) bool {
	return (x > y+e)
}

// LE x <= y
func LE(x float64, y float64) bool {
	return (x < y+e)
}

// GE x >= y
func GE(x float64, y float64) bool {
	return (x > y-e)
}

// EQ x == y
func EQ(x float64, y float64) bool {
	return (y-e < x) && (x < y+e)
}

// NE x != y
func NE(x float64, y float64) bool {
	return (x < y-e) || (y+e < x)
}

// IsZero x == 0
func IsZero(x float64) bool {
	return EQ(x, 0)
}

// IsPos x == 0
func IsPos(x float64) bool {
	return GT(x, 0)
}

// IsNeg x == 0
func IsNeg(x float64) bool {
	return LT(x, 0)
}

// Floor floors the given number to the lext lowest round unit
func Floor(x float64, unit float64) float64 {
	if IsZero(unit) {
		return x
	}

	units := int64((x + unit*e) / unit)
	return float64(units) * unit
}

// Rectify brings the supplied value closer to the precision.
// It should be used when a value is expected to be an integer
// multiple of the unit, but due to repeated computations the value
// may have diverged.
//
// (35.4501, 0.01) => (35.45 +- 1e-8)
//
func Rectify(x *float64, unit float64) {
	if IsZero(unit) {
		return
	}

	switch {
	case *x < 0:
		units := uint64((-(*x) + unit*0.01) / unit)
		*x = float64(units) * unit * -1.0

	default:
		units := uint64((*x + unit*0.01) / unit)
		*x = float64(units) * unit
	}
}

// SafeDiv performs safe division of two numbers i.e. if the divisor
// is zero, it returns zero
func SafeDiv(x float64, y float64) float64 {
	if IsZero(y) {
		return 0
	}

	return x / y
}
