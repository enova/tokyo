package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	assert := assert.New(t)

	msg := buildMessage(LevelInfo, "abc")
	assert.False(msg.Whisper())

	msg = buildMessage(LevelInfo, "abc", Whisper)
	assert.True(msg.Whisper())
}
