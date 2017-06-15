package alert

import (
	"fmt"
)

// Flag ...
type Flag int

// Flags
const (

	// Whisper inhibits a message from being sent externally (e.g. Sentry, Multicast)
	Whisper Flag = 0

	// SkipMail is the same as Whisper (Deprecated)
	SkipMail Flag = 0
)

// Flag-Text
var flagText = map[Flag]string{
	Whisper: "Whisper",
}

func (f Flag) String() string {
	s, ok := flagText[f]
	if ok {
		return "Flag-" + s
	}
	return fmt.Sprintf("Flag-(%d)", int(f))
}
