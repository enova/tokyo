package alert

import (
	"fmt"
)

// Tag ...
type Tag int

// Tags
const (

	// Whisper inhibits a message from being sent externally (e.g. Sentry, Multicast)
	Whisper Tag = 0

	// SkipMail is the same as Whisper (Deprecated)
	SkipMail Tag = 0
)

// Tag-Text
var tagText = map[Tag]string{
	Whisper: "Whisper",
}

func (t Tag) String() string {
	s, ok := tagText[t]
	if ok {
		return "Tag-" + s
	}
	return fmt.Sprintf("Tag-(%d)", int(t))
}
