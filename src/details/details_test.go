package details

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetails(t *testing.T) {
	assert := assert.New(t)

	Set("None")
	assert.True(None())
	assert.False(Info())
	assert.False(More())
	assert.False(Most())
	assert.Equal(LevelI(), 0)
	assert.Equal(LevelS(), "None")

	Set("0")
	assert.True(None())
	assert.False(Info())
	assert.False(More())
	assert.False(Most())
	assert.Equal(LevelI(), 0)
	assert.Equal(LevelS(), "None")

	Set("Info")
	assert.True(None())
	assert.True(Info())
	assert.False(More())
	assert.False(Most())
	assert.Equal(LevelI(), 1)
	assert.Equal(LevelS(), "Info")

	Set("1")
	assert.True(None())
	assert.True(Info())
	assert.False(More())
	assert.False(Most())
	assert.Equal(LevelI(), 1)
	assert.Equal(LevelS(), "Info")

	Set("More")
	assert.True(None())
	assert.True(Info())
	assert.True(More())
	assert.False(Most())
	assert.Equal(LevelI(), 2)
	assert.Equal(LevelS(), "More")

	Set("2")
	assert.True(None())
	assert.True(Info())
	assert.True(More())
	assert.False(Most())
	assert.Equal(LevelI(), 2)
	assert.Equal(LevelS(), "More")

	Set("Most")
	assert.True(None())
	assert.True(Info())
	assert.True(More())
	assert.True(Most())
	assert.Equal(LevelI(), 3)
	assert.Equal(LevelS(), "Most")

	Set("3")
	assert.True(None())
	assert.True(Info())
	assert.True(More())
	assert.True(Most())
	assert.Equal(LevelI(), 3)
	assert.Equal(LevelS(), "Most")
}
