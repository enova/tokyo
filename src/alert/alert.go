package alert

import (
	"fmt"
	"github.com/enova/tokyo/src/cfg"
	"github.com/getsentry/raven-go"
	"github.com/mgutz/ansi"
	"io"
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
var (
	cerrLock      sync.Mutex
	sendLock      sync.Mutex
	logFile       *os.File
	consoleStream io.Writer
	panicOnExit   bool
)

// Globals: Sentry
var (
	sentry         *raven.Client
	lastSentryHour time.Time
	sentryCnt      int
	sentryThrottle *Throttle
)

// Levels
const (
	LevelInfo  = iota // INFO
	LevelWarn         // WARN
	LevelError        // ERROR
)

// Init ...
func init() {
	sentryThrottle = NewThrottle(MaxSentryPerHour)
}

func console() io.Writer {
	if consoleStream != nil {
		return consoleStream
	}
	return os.Stderr
}

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

// SetErr sets the writer for console output. The default
// writer is os.Stderr.
func SetErr(w io.Writer) {
	consoleStream = w
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

// TagsToS ...
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

// PrettyTagsToS ...
func prettyTagsToS(tags []string) string {

	if len(tags) == 0 {
		return ""
	}

	result := ansi.Color("tags: [", "cyan")
	for t, tag := range tags {
		if t > 0 {
			result += " "
		}
		result += ansi.Color(tag, "yellow")
	}
	result += ansi.Color("]", "cyan")

	return result
}

// Configure Log-File
func setLogFile(cfg *cfg.Config) {

	// Get Log-Directory
	dir := cfg.Get("Dir")
	SetLogDir(dir)
}

// SetLogDir ...
func SetLogDir(dir string) {

	// Construct Log-File Path
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
func send(level int, msgs ...string) {

	// Sending Is Reentrant
	sendLock.Lock()
	defer sendLock.Unlock()

	msg := msgs[0]
	tags := msgs[1:]

	// Local-Output for Console/Log-File (Colored-Msg + Stack-Trace)
	local := ansi.Color(timestamp()+" ", "cyan")

	switch {
	case level == LevelInfo:
		local += ansi.Color(msg, "blue")
	case level == LevelWarn:
		local += ansi.Color(msg, "yellow") + "\n"
		local += stacktrace()
	case level == LevelError:
		local += ansi.Color(msg, "red") + "\n"
		local += stacktrace()
	}

	// Add Tags To Local-Output
	local += prettyTagsToS(tags)

	// Write to Console
	Cerr(local)

	// Write to Log-File
	if logFile != nil {
		fmt.Fprintf(logFile, "%s\n", local)
	}

	// Write to Sentry
	sendToSentry(level, msgs...)
}

// Timestamp
func timestamp() string {
	return time.Now().Format("20060102-15:04:05")
}

// Cerr ...
func Cerr(msg string) {
	cerrLock.Lock()
	defer cerrLock.Unlock()
	fmt.Fprintf(console(), "%s\n", msg)
}

// Info ...
func Info(msgs ...string) {
	msgs[0] = "Info: " + msgs[0]
	send(LevelInfo, msgs...)
}

// Warn ...
func Warn(msgs ...string) {
	msgs[0] = "Warn: " + msgs[0]
	send(LevelWarn, msgs...)
}

// WarnIf invokes a warning if there was a failure
func WarnIf(failure bool, msg string) {
	if failure {
		Warn(msg)
	}
}

// WarnOn invokes a warning containing both the error message
// and the supplied message if there was an error
func WarnOn(err error, msg string) {
	if err != nil {
		Warn(msg + ": " + err.Error())
	}
}

// Exit ...
func Exit(msg string) {
	send(LevelError, "Exit: "+msg)
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
