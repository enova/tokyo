package alert

import (
	"fmt"
)

// Level ...
type Level int

// Levels are the levels at which messages are alerted.
const (
	LevelInfo Level = 0 // INFO
	LevelWarn Level = 1 // WARN
	LevelExit Level = 2 // EXIT
)

// Level-Text
var levelText = map[Level]string{
	LevelInfo: "INFO",
	LevelWarn: "WARN",
	LevelExit: "EXIT",
}

func (l Level) String() string {
	s, ok := levelText[l]
	if ok {
		return s
	}
	return fmt.Sprintf("Level-(%d)", int(l))
}
