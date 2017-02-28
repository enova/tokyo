package jwalker

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// W ...
type W struct {
	obj      interface{}
	location string
	failure  string
}

// New ...
func New(b []byte) (*W, error) {
	w := W{}
	err := json.Unmarshal(b, &w.obj)
	return &w, err
}

// String
func (w *W) String() string {
	return fmt.Sprintf("%v", w.obj)
}

// Ok return true if the instance is not invalid
func (w *W) Ok() bool {
	return w.failure == ""
}

// Location returns the current location
func (w *W) Location() string {
	return w.location
}

// Failure returns the current failure
func (w *W) Failure() string {
	return w.failure
}

// Trace returns a string with the current location and current failure (if any)
func (w *W) Trace() string {
	if w.failure == "" {
		return w.location + " (no failures)"
	}
	return "[location] => " + w.location + " [failure] => " + w.failure
}

// Key descends into the supplied key and returns a new W object
// constructed with the descended value.
func (w *W) Key(key string) *W {

	if !w.Ok() {
		return w
	}

	// Verify Map
	mapped, ok := w.obj.(map[string]interface{})
	if !ok {
		w.failure = "key: " + key + " (not a map)"
		return w
	}

	// Verify Key
	value, ok := mapped[key]
	if !ok {
		w.failure = "key: " + key + " (key does not exist)"
		return w
	}

	child := &W{obj: value}

	// Append Location
	child.location = w.location
	child.appendLocation("key: " + key)
	return child
}

// Keys returns a slice of strings containing all keys if the instance is a map
// If the instance is not a map, it returns an empty slice
func (w *W) Keys() []string {

	// Verify Map
	mapped, ok := w.obj.(map[string]interface{})
	if !ok {
		return []string{}
	}

	// Extract Keys
	result := make([]string, 0, len(mapped))
	for key := range mapped {
		result = append(result, key)
	}

	return result
}

// At descends into the supplied array index and returns a new W object
// constructed with the descended value.
func (w *W) At(i int) *W {

	if !w.Ok() {
		return w
	}

	// Verify Array
	array, ok := w.obj.([]interface{})
	if !ok {
		w.failure = fmt.Sprintf("at: %d (not an array)", i)
		return w
	}

	// Verify At
	if i < 0 || i >= len(array) {
		w.failure = fmt.Sprintf("at: %d (out of range, size=%d)", i, len(array))
		return w
	}

	value := array[i]
	child := &W{obj: value}

	// Append Location
	child.location = w.location
	child.appendLocation("at: " + strconv.Itoa(i))
	return child
}

// Len returns the length of the instance's object array
// If the object is not an array then it returns 0
func (w *W) Len() int {

	// Verify Array (Return Zero If Not Array)
	array, ok := w.obj.([]interface{})
	if !ok {
		return 0
	}

	return len(array)
}

// S returns a string if the object is a string
func (w *W) S() (string, bool) {
	if s, ok := w.obj.(string); ok {
		return s, true
	}
	return "", false
}

// I returns a string if the object is a string
func (w *W) I() (int, bool) {
	if f, ok := w.F64(); ok {
		return int(f), true
	}
	return 0, false
}

// U32 returns a uint32 if the object is a string
func (w *W) U32() (uint32, bool) {
	if f, ok := w.F64(); ok {
		return uint32(f), true
	}
	return 0, false
}

// F64 returns the float64 value of the instance's object data
func (w *W) F64() (float64, bool) {
	if f, ok := w.obj.(float64); ok {
		return f, true
	}
	return 0, false
}

// KeyS returns the string for the given key
func (w *W) KeyS(key string) (string, bool) {
	child := w.Key(key)
	if !child.Ok() {
		return "", false
	}

	v, ok := child.S()
	if !ok {
		return "", false
	}

	return v, true
}

// KeyF64 returns the float64 for the given key
func (w *W) KeyF64(key string) (float64, bool) {
	child := w.Key(key)
	if !child.Ok() {
		return 0, false
	}

	v, ok := child.F64()
	if !ok {
		return 0, false
	}

	return v, true
}

// KeyI returns the int value for the given key
func (w *W) KeyI(key string) (int, bool) {
	child := w.Key(key)
	if !child.Ok() {
		return 0, false
	}

	v, ok := child.I()
	if !ok {
		return 0, false
	}

	return v, true
}

// KeyU32 returns the uint32 value for the given key
func (w *W) KeyU32(key string) (uint32, bool) {
	child := w.Key(key)
	if !child.Ok() {
		return 0, false
	}

	v, ok := child.U32()
	if !ok {
		return 0, false
	}

	return v, true
}

// AtS returns the string at the given index
func (w *W) AtS(i int) (string, bool) {
	child := w.At(i)
	if !child.Ok() {
		return "", false
	}

	v, ok := child.S()
	if !ok {
		return "", false
	}

	return v, true
}

// AppendLocation ...
func (w *W) appendLocation(l string) {

	// Separator (If Non-Empty)
	if w.location != "" {
		w.location += " | "
	}

	w.location += l
}
