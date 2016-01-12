stopwatch
=========

A simple stopwatch to use for quick inline benchmarking

Usage
=====
Example:
```
package main

import (
	"github.com/enova/tokyo/src/stopwatch"
)

func main() {

	// Create a stopwatch
	clock := stopwatch.New()

	// Click the stopwatch to start
	clock.Click()

	// Do some work
	for i := 0; i < 1000000; i++ {
		// ...
	}

	// Click the stopwatch to stop
	clock.Click()

	// Display Duration between clicks
	clock.Log("First-Part")
}
```
