package lax

// This package contains random snippets of code that
// we seem to use a lot

import (
	"bytes"
	"fmt"
	"github.com/enova/tokyo/src/alert"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// StackTrace returns a string containing stack-trace information up to the supplied
// depth. Supplying a depth of zero results in a full stack-trace.
func StackTrace(depth int) string {
	var result string

	i := int(1)
	for depth == 0 || i <= depth {
		_, fn, line, ok := runtime.Caller(i)

		// Reached Limit
		if !ok {
			break
		}

		result += fmt.Sprintf("Level-%d => %s:%d\n", i, fn, line)
		i++
	}
	return result
}

// Open a file. alert.Exit on failure.
func Open(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		alert.Exit("Can't open file: " + filename + ", " + err.Error())
	}
	return file
}

// Create a file. alert.Exit on failure.
func Create(filename string, label string) *os.File {
	file, err := os.Create(filename)
	if err != nil {
		alert.Exit("Can't create " + label + " file: " + filename + ", " + err.Error())
	}
	return file
}

// Append a file. alert.Exit on failure
func Append(filename string, label string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		alert.Exit("Can't append " + label + " file: " + filename + ", " + err.Error())
	}
	return file
}

// FileExists return true of the supplied filename exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// ReadFile reads a file and returns the contents as a string. alert.Exit on failure.
func ReadFile(filename string, label string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		alert.Exit("Can't read " + label + " file: " + filename + ", " + err.Error())
	}
	return string(bytes)
}

// WriteFile writes a string to file. If a file with that name exists,
// it will overwrite it. alert.Exit on failure.
func WriteFile(filename, body, label string) {
	err := ioutil.WriteFile(filename, []byte(body), 0644)
	if err != nil {
		alert.Exit("Can't create " + label + " file: " + filename + ", " + err.Error())
	}
}

// CopyFile copies a given file. alert.Exit on failure (via WriteFile)
func CopyFile(fromPath, toPath, label string) {
	body := ReadFile(fromPath, label+" - read")
	WriteFile(toPath, body, label+" - write")
}

// MkDir creates a directory (including all parents i.e. mkdir -p)
// alert.Exit on failure)
func MkDir(dir string, label string) {
	err := os.MkdirAll(dir, 0755)
	alert.ExitOn(err, "Could not create "+label+" directory: "+dir)
}

// FileMode returns the file-mode of the supplied filename.
func FileMode(filename string) (os.FileMode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}

// FileLen returns the length of a file
func FileLen(filename string) int64 {
	if !IsFile(filename) {
		return 0
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0
	}

	return stat.Size()
}

// IsFile returns true if the supplied path is a file (as opposed to a directory)
func IsFile(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	if stat.Mode().IsDir() {
		return false
	}

	if stat.Mode().IsRegular() {
		return true
	}

	return false
}

// IsDir returns true if the supplied path is a directory (as opposed to a file)
func IsDir(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	if stat.Mode().IsDir() {
		return true
	}

	return false
}

// AllFilesUnderDir returns a list of all filenames that exist under the
// supplied directory. If the supplied directory does not exist then it
// returns an empty list.
func AllFilesUnderDir(dir string) []string {
	var result []string

	// List Directory
	items, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		return result
	}

	for _, item := range items {

		switch {

		// File: Append Item
		case IsFile(item):
			result = append(result, item)

		// Dir: Append Item's Children
		case IsDir(item):
			children := AllFilesUnderDir(item)
			for _, c := range children {
				result = append(result, c)
			}
		}
	}

	return result
}

// ReadYAML reads a YAML file into an arbitrary structure.
func ReadYAML(filename string, config interface{}) {
	data := ReadFile(filename, "YAML")
	err := yaml.Unmarshal([]byte(data), config)
	if err != nil {
		alert.Exit("Can't unmarshal " + filename + " into config: " + err.Error())
	}
}

// Separator returns a string of dashes (for displaying).
func Separator(length int) string {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString("-")
	}
	return buffer.String()
}

// Heading returns a pretty separator with the given label
func Heading(label string, dashes int) string {
	return Separator(dashes) + " " + label + " " + Separator(dashes)
}

// ParseUint32 parses a string and returns a uint32. On failure it displays
// an error message and calls exit.
func ParseUint32(str string) uint32 {
	result, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		alert.Exit("Failed to parse uint32: " + str + ", " + err.Error())
	}
	return uint32(result)
}

// ParseFloat64 parses a string and returns a float64. On failure it displays
// an error message and calls exit.
func ParseFloat64(str string) float64 {
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		alert.Exit("Failed to parse float64: " + str + ", " + err.Error())
	}
	return result
}

// F32toa converts a float32 into a string. This is simply a convenience wrapper
// around `strconv.FormatFloat`.
func F32toa(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// F64toa converts a float64 into a string. This is simply a convenience wrapper
// around `strconv.FormatFloat`.
func F64toa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// U32toa converts a uint32 into a string.
func U32toa(u uint32) string {
	return fmt.Sprintf("%d", u)
}

// StripColors returns a string by removing color-coding characters from the supplied string
func StripColors(s string) string {
	var buf bytes.Buffer
	buf.Grow(len(s))

	for i := 0; i < len(s); i++ {

		// Skip Color Code
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}

			if i == len(s) {
				break
			}
		}

		// Add Byte
		if s[i] != 'm' {
			buf.WriteByte(s[i])
		}
	}

	return buf.String()
}

// MakeDate converts a date in literal-integer format (i.e. 2015-06-01 is
// represented by the integer 20150601) to a time.Time value in the local
// timezone.
func MakeDate(yyyymmdd int) time.Time {
	yyyy := (yyyymmdd / 10000)
	mm := (yyyymmdd % 10000) / 100
	dd := (yyyymmdd % 100)
	return time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
}

// MakeDateI converts a date in literal-integer format (i.e. 2015-06-01 is
// represented by the integer 20150601) to a time.Time value in the local
// timezone.
func MakeDateI(yyyymmdd int) time.Time {
	return MakeDate(yyyymmdd)
}

// MakeDateS converts a date from string formats "2015-06-01" and "20150601"
// to a time.Time value. The time will be midnight UTC.
func MakeDateS(s string) time.Time {
	switch len(s) {
	case 10:
		// YYYY-MM-DD
		yyyy := int(ParseUint32((s)[0:4]))
		mm := int(ParseUint32((s)[5:7]))
		dd := int(ParseUint32((s)[8:10]))
		return time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	case 8:
		// YYYYMMDD
		yyyy := int(ParseUint32((s)[0:4]))
		mm := int(ParseUint32((s)[4:6]))
		dd := int(ParseUint32((s)[6:8]))
		return time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	default:
		alert.Warn("Bad date format: " + s)
		return MakeDateI(0)
	}
}

// AddDays adds the given number of days to the given date. The number of days
// can be negative.
func AddDays(d time.Time, n int) time.Time {
	oneDay := 24 * time.Hour
	return d.Add(time.Duration(n) * oneDay)
}

// SubDays subtracts the given number of days to the given date. The number of days
// can be negative.
func SubDays(d time.Time, n int) time.Time {
	oneDay := 24 * time.Hour
	return d.Add(time.Duration(-n) * oneDay)
}

// InTimeRange returns true if the first time is between the second and
// third times inclusive
func InTimeRange(t time.Time, bgn time.Time, end time.Time) bool {
	return (bgn.UnixNano() <= t.UnixNano()) && (t.UnixNano() <= end.UnixNano())
}

// DaysBetween returns the number of days between the two dates.
// Remember that the difference between a date and itself is zero.
func DaysBetween(bgn time.Time, end time.Time) int {
	oneDay := int64(24) * time.Hour.Nanoseconds()
	return int((end.UnixNano() - bgn.UnixNano()) / oneDay)
}

// MakeDateList returns a list of consecutive dates starting and ending on
// the given dates (inclusive)
func MakeDateList(bgn time.Time, end time.Time) []time.Time {
	var result []time.Time
	oneDay := 24 * time.Hour

	for d := bgn; !d.After(end); d = d.Add(oneDay) {
		result = append(result, d)
	}

	return result
}

// NowDate returns the current date
func NowDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// HasS returns true if the supplied slice contains the supplied string
func HasS(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

////////////////////
// Sorting Slices //
////////////////////

type byString []string

func (s byString) Len() int {
	return len(s)
}

func (s byString) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byString) Less(i int, j int) bool {
	return s[i] < s[j]
}

// SortS sorts a slice of strings
func SortS(s []string) {
	sort.Sort(byString(s))
}

type byInt []int

func (s byInt) Len() int {
	return len(s)
}

func (s byInt) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byInt) Less(i int, j int) bool {
	return s[i] < s[j]
}

// SortI sorts a slice of integers
func SortI(s []int) {
	sort.Sort(byInt(s))
}

// SideBySide takes two (multi-line) strings and returns a string
// that displays them side by side:
//
// Apple         Pear
// Banana        Grapefruit
// Canteloupe
//
// They do not have to have the same number of lines
//
func SideBySide(a string, b string, divider string) string {
	var result string

	a = strings.TrimSuffix(a, "\n")
	b = strings.TrimSuffix(b, "\n")

	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")

	lenA := len(linesA)
	lenB := len(linesB)

	rows := lenA
	if lenB > rows {
		rows = lenB
	}

	var width int
	for _, line := range linesA {
		if len(line) > width {
			width = len(line)
		}
	}

	for i := 0; i < rows; i++ {
		var left string
		var right string

		if i < lenA {
			left = linesA[i]
		}

		if i < lenB {
			right = linesB[i]
		}

		size := len(StripColors(left))          // Colors mess with the length
		left += strings.Repeat(" ", width-size) // Pad Right

		result += fmt.Sprintf("%s%s%s\n", left, divider, right)
	}

	return result
}

/////////////////
// Environment //
/////////////////

// Username returns a string containing the current user. Returns the value
// "unknown_user" on failure
func Username() string {
	user, err := user.Current()
	if err != nil {
		return "unknown_user"
	}
	return user.Username
}

// Hostname returns a string containing the hostname. Returns the value
// "unknown_host" on failure
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown_host"
	}
	return hostname
}

// ExecPath returns a string containing the full executable path.
func ExecPath() string {
	filename, err := osext.Executable()
	if err != nil {
		return "unknown_exec"
	}
	return filename
}
