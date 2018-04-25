package args

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewArgs(s ...string) *Args {
	return New(s)
}

func TestArgs(t *testing.T) {
	assert := assert.New(t)

	a := NewArgs("-d", "b", "-a=5", "c")
	assert.Equal(len(a.Raws()), 4)

	// Non-Optional Arguments
	assert.Equal(a.Size(), 2)
	assert.Equal(a.Get(0), "b")
	assert.Equal(a.Get(1), "c")

	// Unary Options
	assert.True(a.IsOn("d"))
	assert.False(a.IsOff("d"))

	assert.False(a.IsOn("g"))
	assert.True(a.IsOff("g"))

	// Binary Options
	assert.True(a.HasOpt("a"))
	assert.Equal(a.GetOpt("a"), "5")
	assert.Equal(a.GetOptI("a"), 5)
}
