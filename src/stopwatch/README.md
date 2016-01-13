stopwatch
=========

A simple stopwatch to use for quick inline benchmarking

Usage
=====
First create a `Clock` by `clock = stopwatch.New()`. Next call `clock.Click()` across any two points in the code. Finally, call `duration, error := clock.Show("your-label")` to retrieve the duration between the last two clicks. You can make as many calls to `clock.Click()` as you want. The call to `clock.Show()` will always return the time duration between the last two clicks. If you simply want to display the duration to `stderr` then call `clock.Log("your-label")`.

Example:
```
package main

import (
	"git.enova.com/go/stopwatch"
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
Output:
```
First-Part: 270.957µs
```
