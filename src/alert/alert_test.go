package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type recorder struct {
	value string
}

func (r *recorder) Write(data []byte) (int, error) {
	r.value = string(data)
	return len(r.value), nil
}

func setup() *recorder {
	lastSentryMsg = ""
	r := &recorder{}
	SetErr(r)
	return r
}

func TestCerr(t *testing.T) {
	assert := assert.New(t)
	r := setup()

	Cerr("abc")
	assert.Contains(r.value, "abc")
}

func TestInfo(t *testing.T) {
	assert := assert.New(t)
	r := setup()

	Info("abc")
	assert.Contains(r.value, "abc")
	assert.NotContains(r.value, "alert_test.go")
}

func TestMultipleArguments(t *testing.T) {
	assert := assert.New(t)
	r := setup()

	// Test Multiple Values
	Info("abc", "def", 4, false, 8.9)
	assert.Contains(r.value, "abc def 4 false 8.9")

	// A Non-Trivial Data Type
	type NonTrivial struct {
		age  int
		name string
		list []int
	}

	n := NonTrivial{}
	n.age = 9
	n.name = "abc"
	n.list = append(n.list, 3, 4, 5)

	Info(n)
	assert.Contains(r.value, "age:9")
	assert.Contains(r.value, "name:abc")
	assert.Contains(r.value, "list:[3 4 5]")
}

func TestWarn(t *testing.T) {
	assert := assert.New(t)
	r := setup()

	Warn("abc")
	assert.Contains(r.value, "abc")
	assert.Contains(r.value, "alert_test.go")
}

func TestExit(t *testing.T) {
	assert := assert.New(t)
	r := setup()
	PanicOnExit()

	exit := func() {
		Exit("abc")
	}

	assert.Panics(exit)
	assert.Contains(r.value, "abc")
	assert.Contains(r.value, "alert_test.go")
}

func TestSentry(t *testing.T) {
	assert := assert.New(t)
	r := setup()

	Info("abc", SkipMail)
	assert.Empty(lastSentryMsg)
	assert.Contains(r.value, "abc")

	Info("abc")
	assert.Contains(lastSentryMsg, "abc")
	assert.Contains(r.value, "abc")
}
