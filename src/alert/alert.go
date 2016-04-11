package alert

import (
	"fmt"
	"github.com/enova/tokyo/src/cfg"
	"github.com/getsentry/raven-go"
	"github.com/mgutz/ansi"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// MaxSentryPerHour limits the number of messages sent to Sentry
const MaxSentryPerHour int = 100

// Globals
var lock sync.Mutex
var logFile *os.File
var sentry *raven.Client
var lastSentryHour time.Time
var sentryCnt int
var panicOnExit bool

func username() string {
	user, err := user.Current()
	if err != nil {
		return "unknown_user"
	}
	return user.Username
}

func stacktrace() string {
	var result string

	for i := 3; i < 8; i++ {
		if _, fn, line, ok := runtime.Caller(i); ok {
			result += ansi.Color(fmt.Sprintf("\t%s:%d\n", fn, line), "yellow")
		}
	}

	return result
}

// PanicOnExit will cause the alert package to call panic() instead of os.exit()
func PanicOnExit() {
	panicOnExit = true
}

// Set configures the alert settings
func Set(cfg *cfg.Config) {

	// Configure: Sentry-Client
	if cfg.Has("Alert.Sentry.Use") {

		// Use
		use := cfg.Get("Alert.Sentry.Use")
		if use == "true" || use == "True" {
			sentryCfg := cfg.Descend("Alert.Sentry")
			setSentry(sentryCfg)
		}
	}

	// Configure: Log-File
	if cfg.Has("Alert.LogFile.Use") {

		// Use
		use := cfg.Get("Alert.LogFile.Use")
		if use == "true" || use == "True" {
			logFileCfg := cfg.Descend("Alert.LogFile")
			setLogFile(logFileCfg)
		}
	}
}

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
	tags["user"] = username()
	tags["app"] = os.Args[0]
	tags["cmd"] = strings.Join(os.Args, " ")
	tags["pid"] = fmt.Sprintf("%d", os.Getpid())

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

	Cerr("Sentry: " + tagsToS(tags))
}

func tagsToS(tags map[string]string) string {
	result := "("
	first := true

	for _, v := range tags {
		if !first {
			result += ","
		}
		result += v
		first = false
	}
	result += ")"
	return result
}

// Configure Log-File
func setLogFile(cfg *cfg.Config) {

	// Construct Log-File Path
	dir := cfg.Get("Dir")
	_, app := filepath.Split(os.Args[0])
	now := time.Now().Format("20060102_150405_000")
	pid := os.Getpid()
	filename := fmt.Sprintf("%s/%s_%s_%d.log", dir, app, now, pid)

	// Create Directory
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		Cerr("Alert: Can't create log directory: " + dir)
		os.Exit(1)
	}

	// Append Log-File
	logFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		Cerr("Alert: Can't append file: " + filename + ", " + err.Error())
		os.Exit(1)
	}

	Cerr("LogFile: " + filename)
}

// Send a message
func send(level raven.Severity, msgs ...string) {

	// Rentrant
	lock.Lock()
	defer lock.Unlock()

	msg := msgs[0]
	tags := msgs[1:]
	TagLen := len(tags)

	// Local-Output for Stderr/Log-File (Colored-Msg + Stack-Trace)
	local := ansi.Color(timestamp()+" ", "cyan")

	switch {
	case level == raven.INFO:
		local += ansi.Color(msg, "blue") + "\n"
	case level == raven.WARNING:
		local += ansi.Color(msg, "yellow") + "\n"
		local += stacktrace()
	case level == raven.ERROR:
		local += ansi.Color(msg, "red") + "\n"
		local += stacktrace()
	}

	// Add Tags To Local-Output
	if TagLen > 0 {
		local += ansi.Color("tags: [", "cyan")
		for t, tag := range tags {
			if t > 0 {
				local += " "
			}
			local += ansi.Color(tag, "yellow")
		}
		local += ansi.Color("]", "cyan")
	}

	// Write to Stderr
	fmt.Fprintf(os.Stderr, "%s\n", local)

	// Write to Log-File
	if logFile != nil {
		fmt.Fprintf(logFile, "%s\n", local)
	}

	// Write to Sentry
	if sentry != nil {

		// Reset Last-Sentry-Hour
		if time.Since(lastSentryHour) >= time.Hour {
			sentryCnt = 0
			lastSentryHour = time.Now()
		}

		// Too Many Messages
		if sentryCnt >= MaxSentryPerHour {
			return
		}

		// Update Sentry Message-Count
		sentryCnt++

		// Packet
		packet := &raven.Packet{
			Message: msg,
			Level:   level,
		}

		// Set Tags
		if TagLen > 0 {
			packet.Tags = make(raven.Tags, TagLen)

			for t, tag := range tags {
				if tag == "skip_sentry" {
					// skip sending to sentry
					return
				}
				packet.Tags[t] = raven.Tag{
					Key:   tag,
					Value: "true",
				}
			}
		}

		var err error
		_, ch := sentry.Capture(packet, nil)
		if err = <-ch; err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send packet to Sentry: "+err.Error())
		}
	}
}

// Timestamp
func timestamp() string {
	return time.Now().Format("20060102-15:04:05")
}

// Cerr ...
func Cerr(msg string) {
	lock.Lock()
	defer lock.Unlock()
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

// Info ...
func Info(msgs ...string) {
	msgs[0] = "Info: " + msgs[0]
	send(raven.INFO, msgs...)
}

// Warn ...
func Warn(msgs ...string) {
	msgs[0] = "Warn: " + msgs[0]
	send(raven.WARNING, msgs...)
}

// WarnIf writes the message to stderr if there was a failure
func WarnIf(failure bool, msg string) {
	if failure {
		Warn(msg)
	}
}

// WarnOn writes the error message and the supplied message to stderr
// if there was an error
func WarnOn(err error, msg string) {
	if err != nil {
		Warn(msg + ": " + err.Error())
	}
}

// Exit ...
func Exit(msg string) {
	send(raven.ERROR, "Exit: "+msg)
	if panicOnExit {
		err := fmt.Errorf(msg)
		panic(err)
	}
	os.Exit(1)
}

// ExitIf ...
func ExitIf(failure bool, msg string) {
	if failure {
		Exit(msg)
	}
}

// ExitOn ...
func ExitOn(err error, msg string) {
	if err != nil {
		Exit(msg + ": " + err.Error())
	}
}
