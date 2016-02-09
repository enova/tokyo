package preprdsql

import (
	"bufio"
	"regexp"
	"strings"
)

// Scanner ...
type Scanner struct {
	currentLine string
	statements  map[string]string
	currentTag  string
}

type scanner func(*Scanner) scanner

// getTag checks if the current line contains a tag or non sql strings
func getTag(line string) string {
	re := regexp.MustCompile("^\\s*--\\s*name:\\s*(\\S+)")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return ""
	}
	return matches[1]
}

// readLine processes the current line by determining if a tag is located
// or not. If a tag is matched, it resets the currentTag to the newly matched
// tag and returns a queryLine.
func readLine(s *Scanner) scanner {
	if tag := getTag(s.currentLine); len(tag) > 0 {
		s.currentTag = tag
		return processLine
	}
	return readLine
}

// processLine processes the current line and if the current line matches a
// tag, it sets the *Scanner.currentTag. If the current line is not a tag,
// it appends the line.
func processLine(s *Scanner) scanner {
	if tag := getTag(s.currentLine); len(tag) > 0 {
		s.currentTag = tag
	} else {
		s.appendLine()
	}
	return processLine
}

// appendLine appends sql statements to the current tag being processed.
// If the line is blank or empty it does not append that line
func (s *Scanner) appendLine() {
	tag := s.statements[s.currentTag]
	line := strings.Trim(s.currentLine, " \t")
	if len(line) == 0 {
		return
	}

	if len(tag) > 0 {
		tag = tag + "\n"
	}

	tag = tag + line
	s.statements[s.currentTag] = tag
}

// Scan scans the given file and creates a map using the tag as the key and the
// values below the tag till the next tag as the map value.
func (s *Scanner) Scan(io *bufio.Scanner) map[string]string {
	s.statements = make(map[string]string)

	for line := readLine; io.Scan(); {
		s.currentLine = io.Text()
		line = line(s)
	}

	return s.statements
}
