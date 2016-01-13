package stopwatch

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Clock stores the last timestamp that was Click'd
type Clock struct {
	twoClicks bool
	lastClick time.Time
	duration  time.Duration
}

// New returns a new Clock instance
func New() *Clock {
	return &Clock{}
}

// Click stores the current time and appends a delta
// from the last time
func (c *Clock) Click() {
	now := time.Now()
	if !c.lastClick.IsZero() {
		c.duration = now.Sub(c.lastClick)
		c.twoClicks = true
	}
	c.lastClick = now
}

// Show takes the last duration and displays it along
// with the label
func (c *Clock) Show(label string) (string, error) {
	if !c.twoClicks {
		return "", errors.New("Cannot calculate duration without at least 2 clicks")
	}

	return fmt.Sprintf("%s: %s", label, c.duration), nil
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
