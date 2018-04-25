package alert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottle(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()

	l := NewThrottle(3)
	assert.True(l.Update(now))
	assert.True(l.Update(now))
	assert.True(l.Update(now))

	// At Limit!
	assert.False(l.Update(now))

	// Still Within One Hour
	now = now.Add(time.Hour / 2)
	assert.False(l.Update(now))

	// One Hour Passed
	now = now.Add(time.Hour / 2)
	assert.True(l.Update(now))
	assert.True(l.Update(now))
	assert.True(l.Update(now))
	assert.False(l.Update(now))
}
