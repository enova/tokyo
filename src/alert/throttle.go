package alert

import (
	"time"
)

// Throttle ...
type Throttle struct {
	last  time.Time
	count int
	limit int
}

// NewThrottle ...
func NewThrottle(limit int) *Throttle {
	t := &Throttle{limit: limit}
	return t
}

// Update ...
func (t *Throttle) Update(now time.Time) bool {

	// One-Hour Elapsed: Reset
	if now.Sub(t.last) >= time.Hour {
		t.count = 0
		t.last = now
	}

	// Too Many Updates
	if t.count >= t.limit {
		return false
	}

	// All Good
	t.count++
	return true
}
