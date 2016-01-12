package stopwatch

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	// EmptyTime is the zero value of time
	EmptyTime = time.Time{}
)

// Clock stores the last timestamp that was Click'd
type Clock struct {
	lastClick time.Time
	durations []time.Duration
}

// New returns a new Clock instance
func New() *Clock {
	return &Clock{durations: make([]time.Duration, 0)}
}

// Click stores the current time and appends a delta
// from the last time
func (c *Clock) Click() {
	now := time.Now()
	if c.lastClick != EmptyTime {
		c.durations = append(c.durations, now.Sub(c.lastClick))
	}
	c.lastClick = time.Now()
}

// Show takes the last duration and displays it along
// with the label
func (c *Clock) Show(label string) (string, error) {
	if len(c.durations) < 1 {
		return "", errors.New("Cannot calculate duration without at least 2 clicks")
	}

	lastDuration := c.durations[len(c.durations)-1]
	return fmt.Sprintf("%s: %s", label, lastDuration), nil
}

// Log writes the contents of Show() to stderr. If Show returns an error
// it does nothing.
func (c *Clock) Log(label string) {
	show, err := c.Show(label)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s\n", show)
}
