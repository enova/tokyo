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
