package stopwatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyOneClickBeforeShow(t *testing.T) {
	assert := assert.New(t)

	s := New()
	s.Click()

	label, err := s.Show("Hello")
	assert.Error(err, "Expected to error because there has only been 1 click")
	assert.Equal(label, "", "Expected no label because there has only been 1 click")
}

func TestNormalShow(t *testing.T) {
	assert := assert.New(t)

	s := New()
	s.Click()
	s.Click()

	label, err := s.Show("Hello")
	assert.NoError(err, "Expected no error because there has been 2 clicks")
	assert.Contains(label, "Hello", "Expected label because there has been 2 clicks")
}

func TestMultipleShows(t *testing.T) {
	assert := assert.New(t)

	s := New()
	s.Click()
	s.Click()

	label, _ := s.Show("Hello")
	assert.Contains(label, "Hello", "Expected label because there has been 2 clicks")

	s.Click()
	newLabel, _ := s.Show("Hello")
	assert.NotEqual(label, newLabel, "Expected new label because there has been 3 clicks")
}
