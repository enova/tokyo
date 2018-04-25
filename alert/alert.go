package alert

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/enova/tokyo/src/cfg"
)

// Globals
var (
	cerrLock      sync.Mutex
	sendLock      sync.Mutex
	logFile       *os.File
	consoleStream io.Writer
	panicOnExit   bool
	handlers      []Handler
	meta          Meta
)

// PanicOnExit will cause the alert package to call panic() instead of os.exit()
func PanicOnExit() {
	panicOnExit = true
}

// SetErr sets the writer for console output. The default writer is os.Stderr.
func SetErr(w io.Writer) {
	consoleStream = w
}

// Set configures the alert settings
func Set(cfg *cfg.Config) {

	// Configure: Sentry-Client
	if alertCfg, ok := getCfg("Sentry", cfg); ok {
		setSentry(alertCfg)
	}

	// Configure: Multicast-Client
	if alertCfg, ok := getCfg("Multicast", cfg); ok {
		setMulticast(alertCfg)
	}

	// Configure: Log-File
	if alertCfg, ok := getCfg("LogFile", cfg); ok {
		setLogFile(alertCfg)
	}
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
		log.Fatal(err)
	}

	// Append Log-File
	logFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		Cerr("Alert: Can't append file: " + filename + ", " + err.Error())
		log.Fatal(err)
	}

	Cerr("LogFile: " + filename)
}

// AddHandler adds the supplied handler to the list of handlers
func AddHandler(h Handler) {
	handlers = append(handlers, h)
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

// AddPrefix prepends the supplied prefix to msgs slice
func addPrefix(prefix string, msgs ...interface{}) []interface{} {
	fields := make([]interface{}, 1, len(msgs)+1)
	fields[0] = prefix + ":"
	fields = append(fields, msgs...)
	return fields
}

// Returns the config for the given path, else nil and then also returns true if found
func getCfg(cfgPath string, cfg *cfg.Config) (*cfg.Config, bool) {

	// If path has a '.Use', ...
	if cfg.Has("Alert." + cfgPath + ".Use") {

		// ...and Use is True...
		use := cfg.Get("Alert." + cfgPath + ".Use")
		if use == "true" || use == "True" {

			// ...Get Alert configuration for path
			return cfg.Descend("Alert." + cfgPath), true
		}
	}
	return nil, false
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

// Configure Log-File
func setLogFile(cfg *cfg.Config) {

	// Get Log-Directory
	dir := cfg.Get("Dir")
	SetLogDir(dir)
}

// Send a message
func send(level Level, fields ...interface{}) {

	// Create New Message
	msg := buildMessage(level, fields...)

	// Sending Is Reentrant
	sendLock.Lock()
	defer sendLock.Unlock()

	// Write to Console
	Cerr(msg.Pretty())

	// Write to Log-File
	if logFile != nil {
		fmt.Fprintf(logFile, "%s\n", msg.Pretty())
	}

	// Send To External Services
	if !msg.Whisper() {

		// Write to Sentry
		sendToSentry(msg)

		// Write to multicast
		sendToMulticast(msg)
	}

	// Send To Custom-Handlers
	for _, h := range handlers {
		h.Handle(*msg)
	}
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
func Info(msgs ...interface{}) {
	send(LevelInfo, msgs...)
}

// Warn ...
func Warn(msgs ...interface{}) {
	send(LevelWarn, msgs...)
}

// WarnIf invokes a warning if there was a failure
func WarnIf(failure bool, msgs ...interface{}) {
	if failure {
		Warn(msgs)
	}
}

// WarnOn invokes a warning containing both the error message
// and the supplied message if there was an error
func WarnOn(err error, msgs ...interface{}) {
	if err != nil {
		fields := addPrefix("("+err.Error()+")", msgs...)
		Warn(fields)
	}
}

// Exit ...
func Exit(msgs ...interface{}) {
	send(LevelExit, msgs...)
	if panicOnExit {
		err := fmt.Errorf("%+v", msgs)
		panic(err)
	}
	os.Exit(1)
}

// ExitIf ...
func ExitIf(failure bool, msgs ...interface{}) {
	if failure {
		Exit(msgs)
	}
}

// ExitOn ...
func ExitOn(err error, msgs ...interface{}) {
	if err != nil {
		fields := addPrefix(err.Error(), msgs...)
		Exit(fields)
	}
}
