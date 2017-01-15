package alert

import (
	"time"

	"github.com/getsentry/raven-go"
)

// For Testing
var lastSentryMsg string

func sendToSentry(level int, msg string) {

	// Testing
	lastSentryMsg = msg

	// Sentry-Enabled
	if sentry == nil {
		return
	}

	// Severity
	var severity raven.Severity

	switch {
	case level == LevelInfo:
		severity = raven.INFO
	case level == LevelWarn:
		severity = raven.WARNING
	case level == LevelError:
		severity = raven.ERROR
	default:
		severity = raven.INFO
	}

	// Check Sentry-Throttle
	if ok := sentryThrottle.Update(time.Now()); !ok {
		return
	}

	// Packet
	packet := &raven.Packet{
		Message: msg,
		Level:   severity,
	}

	// Send To Sentry
	var err error
	_, ch := sentry.Capture(packet, nil)
	if err = <-ch; err != nil {
		Cerr("Failed to send packet to Sentry: " + err.Error())
	}
}
