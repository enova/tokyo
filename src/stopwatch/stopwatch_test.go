package stopwatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyOneClickBeforeShow(t *testing.T) {
	s := New()

	s.Click()
	label, e := s.Show("Hello")
	assert.Error(t, e, "Expected to error because there has only been 1 click")
	assert.Equal(t, label, "", "Expected no label because there has only been 1 click")
}

func TestNormalShow(t *testing.T) {
	s := New()

	s.Click()
	s.Click()
	label, e := s.Show("Hello")
	assert.NoError(t, e, "Expected no error because there has been 2 clicks")
	assert.Contains(t, label, "Hello", "Expected label because there has been 2 clicks")
}

func TestMultipleShows(t *testing.T) {
	s := New()

	s.Click()
	s.Click()

	label, _ := s.Show("Hello")
	assert.Contains(t, label, "Hello", "Expected label because there has been 2 clicks")

	s.Click()
	newLabel, _ := s.Show("Hello")
	assert.NotEqual(t, label, newLabel, "Expected new label because there has been 3 clicks")
}
