package details

import (
	"fmt"
	"os"
	"sync/atomic"
)

// Level defaults to None (0)
var level int32

// Levels
const (
	NoneT = iota
	InfoT
	MoreT
	MostT
)

// Set sets the details level
func Set(s string) {
	switch s {
	case "None":
		set(NoneT)
	case "Info":
		set(InfoT)
	case "More":
		set(MoreT)
	case "Most":
		set(MostT)

	case "0":
		set(NoneT)
	case "1":
		set(InfoT)
	case "2":
		set(MoreT)
	case "3":
		set(MostT)

	default:
		fmt.Fprintf(os.Stderr, "Unsupported level for details: [%s]", s)
	}
}

// None returns true if the level exceeds None
func None() bool {
	return get() >= NoneT
}

// Info returns true if the level exceeds Info
func Info() bool {
	return get() >= InfoT
}

// More returns true if the level exceeds More
func More() bool {
	return get() >= MoreT
}

// Most returns true if the level exceeds Most
func Most() bool {
	return get() >= MostT
}

// LevelI returns the details level as an integer
func LevelI() int {
	return get()
}

// LevelS returns the details level as a string
func LevelS() string {
	switch get() {
	case 0:
		return "None"
	case 1:
		return "Info"
	case 2:
		return "More"
	case 3:
		return "Most"
	}

	return "FixMePlease"
}

// Set the level (atomically)
func set(v int) {
	atomic.StoreInt32(&level, int32(v))
}

// Get the level (atomically)
func get() int {
	return int(atomic.LoadInt32(&level))
}
