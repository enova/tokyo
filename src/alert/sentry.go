package alert

import (
	"time"

	"github.com/getsentry/raven-go"
)

func sendToSentry(level int, msgs ...string) {

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

	// Message And Tags
	msg := msgs[0]
	tags := msgs[1:]
	TagLen := len(tags)

	// Packet
	packet := &raven.Packet{
		Message: msg,
		Level:   severity,
	}

	// Set Tags
	if TagLen > 0 {
		packet.Tags = make(raven.Tags, TagLen)

		for t, tag := range tags {

			// Skip-Sentry (Option)
			if tag == "skip_sentry" {
				return
			}

			packet.Tags[t] = raven.Tag{
				Key:   tag,
				Value: "true",
			}
		}
	}

	// Send To Sentry
	var err error
	_, ch := sentry.Capture(packet, nil)
	if err = <-ch; err != nil {
		Cerr("Failed to send packet to Sentry: " + err.Error())
	}
}
