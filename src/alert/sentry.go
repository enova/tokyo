package alert

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/enova/tokyo/src/cfg"
	"github.com/getsentry/raven-go"
)

// MaxSentryPerHour limits the number of messages sent to Sentry
const MaxSentryPerHour int = 100

// Globals: Sentry
var (
	sentry         *raven.Client
	sentryThrottle *Throttle
)

// For Testing
var lastSentryMsg string

// Configure Sentry
func setSentry(cfg *cfg.Config) {
	var err error

	// URL (From Env)
	url := os.Getenv("SENTRY_DSN")

	// Overwrite With Config URL
	if cfg.Has("DSN") {
		url = cfg.Get("DSN")
	}

	// Tags
	tags := make(map[string]string)
	tags["user"] = meta.User
	tags["app"] = meta.App
	tags["cmd"] = strings.Join(os.Args, " ")
	tags["pid"] = meta.PID

	// Custom-Tags
	for t := 0; t < cfg.Size("Tag"); t++ {
		line := cfg.GetN(t, "Tag")
		tokens := strings.Fields(line)

		// Invalid Tag
		if len(tokens) < 2 {
			Cerr("Tag should have at least two tokens - key value..., " + line)
			continue
		}

		// Add Tag
		key := tokens[0]
		val := strings.Join(tokens[1:], " ")
		tags[key] = val
	}

	// Create Sentry-Client
	sentry, err = raven.NewWithTags(url, tags)
	if err != nil {
		log.Fatal(err)
	}

	// Create Sentry-Throttle
	sentryThrottle = NewThrottle(MaxSentryPerHour)
	Cerr("Sentry: " + tagsToS(tags))
}

// Send Message To Sentry
func sendToSentry(msg *Message) {

	// Testing
	lastSentryMsg = msg.Text

	// Sentry-Enabled
	if sentry == nil {
		return
	}

	// Severity
	var severity raven.Severity

	switch {
	case msg.Level == LevelInfo:
		severity = raven.INFO
	case msg.Level == LevelWarn:
		severity = raven.WARNING
	case msg.Level == LevelExit:
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
		Message: msg.Text,
		Level:   severity,
	}

	// Send To Sentry
	var err error
	_, ch := sentry.Capture(packet, nil)
	if err = <-ch; err != nil {
		Cerr("Failed to send packet to Sentry: " + err.Error())
	}
}
