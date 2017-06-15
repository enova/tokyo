package jwalker

import (
	"encoding/json"
	"fmt"
	"reflect"
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

// String is for diagnostics only. Use S() to extract a string element.
func (w *W) String() string {
	return fmt.Sprintf("%v", w.obj)
}

// Pretty returns an indented, multi-line string representation of the underlying object
// (invokes json.MarshalIndent)
func (w *W) Pretty() string {
	bytes, _ := json.MarshalIndent(w.obj, "", "  ")
	return string(bytes)
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

	// Result
	child := &W{}
	child.location = w.location

	// Verify Map
	mapped, ok := w.obj.(map[string]interface{})
	if !ok {
		child.failure = "key: " + key
		child.failure += " (not a map,"
		child.failure += " rather it is of type " + reflect.TypeOf(w.obj).String() + ")"
		return child
	}

	// Verify Key
	value, ok := mapped[key]
	if !ok {
		child.failure = "key: " + key + " (key does not exist)"
		return child
	}

	// Append Location
	child.obj = value
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

	// Result
	child := &W{}
	child.location = w.location

	// Verify Array
	array, ok := w.obj.([]interface{})
	if !ok {
		child.failure = fmt.Sprintf("at: %d (not an array,", i)
		child.failure += " rather it is of type " + reflect.TypeOf(w.obj).String() + ")"
		return child
	}

	// Verify At
	if i < 0 || i >= len(array) {
		child.failure = fmt.Sprintf("at: %d (out of range, size=%d)", i, len(array))
		return child
	}

	value := array[i]

	// Append Location
	child.obj = value
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

// I returns an integer if the object is an integer
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

// B returns the bool value if the object is a bool
func (w *W) B() (bool, bool) {
	if b, ok := w.obj.(bool); ok {
		return b, true
	}
	return false, false
}

// KeyS returns the string value for the given key
func (w *W) KeyS(key string) (string, bool) {
	return w.Key(key).S()
}

// KeyF64 returns the float64 value for the given key
func (w *W) KeyF64(key string) (float64, bool) {
	return w.Key(key).F64()
}

// KeyI returns the int value for the given key
func (w *W) KeyI(key string) (int, bool) {
	return w.Key(key).I()
}

// KeyU32 returns the uint32 value for the given key
func (w *W) KeyU32(key string) (uint32, bool) {
	return w.Key(key).U32()
}

// KeyB returns the bool value for the given key
func (w *W) KeyB(key string) (bool, bool) {
	return w.Key(key).B()
}

// AtB returns the bool value at the given index
func (w *W) AtB(i int) (bool, bool) {
	return w.At(i).B()
}

// AtF64 returns the float64 value at the given index
func (w *W) AtF64(i int) (float64, bool) {
	return w.At(i).F64()
}

// AtI returns the int value at the given index
func (w *W) AtI(i int) (int, bool) {
	return w.At(i).I()
}

// AtU32 returns the uint32 value at the given index
func (w *W) AtU32(i int) (uint32, bool) {
	return w.At(i).U32()
}

// AtS returns the string at the given index
func (w *W) AtS(i int) (string, bool) {
	return w.At(i).S()
}

// AppendLocation ...
func (w *W) appendLocation(l string) {

	// Separator (If Non-Empty)
	if w.location != "" {
		w.location += " | "
	}

	w.location += l
}
