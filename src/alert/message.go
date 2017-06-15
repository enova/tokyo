package alert

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/mgutz/ansi"
)

// Meta ...
type Meta struct {
	User string
	App  string
	PID  string
}

// Message ...
type Message struct {
	Meta
	Now   time.Time
	Level Level
	Text  string
	Flags []Flag
}

// Initialize Meta-Data (Static)
func init() {
	meta.User = username()
	meta.App = os.Args[0]
	meta.PID = fmt.Sprintf("%d", os.Getpid())
}

// Copy returns a deep-copy of the message
func (m *Message) copy() Message {
	copy := Message{}
	copy.Now = m.Now
	copy.Level = m.Level
	copy.Text = m.Text
	copy.Flags = make([]Flag, len(m.Flags))
	for i, t := range m.Flags {
		copy.Flags[i] = t
	}

	return copy
}

// Whisper searches for the presence of the Whisper flag
func (m *Message) Whisper() bool {
	for _, t := range m.Flags {
		if t == Whisper {
			return true
		}
	}
	return false
}

// Pretty returns a colorful string for Console/Log-File
// It also appends a stack-trace for warnings and above
func (m *Message) Pretty() string {

	// Timestamp
	stamp := m.Now.Format("20060102-15:04:05")
	result := ansi.Color(stamp+" ", "cyan")

	// Prepend Level To Text
	text := m.Level.String() + ": " + m.Text

	// Message
	switch {
	case m.Level == LevelInfo:
		result += ansi.Color(text, "blue")
	case m.Level == LevelWarn:
		result += ansi.Color(text, "yellow") + "\n"
	case m.Level == LevelExit:
		result += ansi.Color(text, "red") + "\n"
	default:
		result += ansi.Color(text, "white") + "\n"
	}

	// Stack-Trace
	if m.Level > LevelInfo {
		result += stacktrace()
	}

	return result
}

// BuildMessage returns a newly instantiated message
func buildMessage(level Level, fields ...interface{}) *Message {

	// Message
	msg := &Message{
		Meta:  meta,
		Now:   time.Now(),
		Level: level,
	}

	// Add Fields
	for i, f := range fields {

		// Flags
		if flag, ok := f.(Flag); ok {

			// Whisper
			if flag == Whisper {
				msg.Flags = append(msg.Flags, flag)
			}
		}

		// Convert To String (Reflection)
		str := fmt.Sprintf("%+v", f)

		// Append To Message (Space Separator)
		if i > 0 {
			msg.Text += " "
		}

		// Append Message
		msg.Text += str
	}

	return msg
}

// Stacktrace ...
func stacktrace() string {
	var result string

	for i := 3; i < 8; i++ {
		if _, fn, line, ok := runtime.Caller(i); ok {
			result += ansi.Color(fmt.Sprintf("\t%s:%d\n", fn, line), "yellow")
		}
	}

	return result
}
