package dbl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDbl(t *testing.T) {
	assert := assert.New(t)

	assert.True(LT(0.1, 0.2), "LT failed")
	assert.True(GT(0.2, 0.1), "GT failed")

	assert.True(LE(0.1, 0.1), "LE failed")
	assert.True(LE(0.1, 0.2), "LE failed")

	assert.True(GE(0.1, 0.1), "GE failed")
	assert.True(GE(0.2, 0.1), "GE failed")

	assert.True(EQ(0.1, 0.1), "EQ failed")
	assert.True(NE(0.1, 0.2), "NE failed")
	assert.True(NE(0.2, 0.1), "NE failed")

	assert.True(IsZero(0.0), "IsZero failed")
	assert.False(IsZero(0.1), "IsZero failed")

	assert.True(EQ(Floor(0.130001, 0.01), 0.13), "Floor failed")
	assert.True(EQ(Floor(0.139999, 0.01), 0.13), "Floor failed")

	assert.True(IsZero(0.0), "IsZero failed")
	assert.False(IsZero(0.1), "IsZero failed")

	assert.True(IsPos(0.1), "IsPos failed")
	assert.False(IsPos(0.0), "IsPos failed")
	assert.False(IsPos(-0.1), "IsPos failed")

	assert.True(IsNeg(-0.1), "IsNeg failed")
	assert.False(IsNeg(0.0), "IsNeg failed")
	assert.False(IsNeg(0.1), "IsNeg failed")

	assert.Equal(SafeDiv(7.0, 2.0), 3.5, "SafeDiv failed")
	assert.Equal(SafeDiv(1.0, 0.0), 0.0, "SafeDiv failed")
}

func TestRectify(t *testing.T) {
	assert := assert.New(t)

	x := 45.03
	Rectify(&x, 0.01)
	assert.Equal(x, 45.03)

	// Rectify 10% Error
	x = 45.031
	Rectify(&x, 0.01)
	assert.Equal(x, 45.03)

	x = -45.031
	Rectify(&x, 0.01)
	assert.Equal(x, -45.03)

	// Don't Rectify
	x = 45.03000001
	assert.NotEqual(x, 45.03)

	// Rectify
	Rectify(&x, 0.01)
	assert.Equal(x, 45.03)

	/////////////////////////////////////
	// Exhibit Divergence Then Rectify //
	/////////////////////////////////////

	// Add A Penny 100000 Times
	x = 45.03
	for i := 0; i < 100000; i++ {
		x += 0.01
	}

	// Diverged!
	assert.NotEqual(x, 1045.03)

	// Rectify
	Rectify(&x, 0.01)
	assert.Equal(x, 1045.03)
}
