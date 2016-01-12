package cfg

import (
	"bufio"
	"fmt"
	"github.com/enova/tokyo/src/set"
	"github.com/mgutz/ansi"
	"os"
	"strings"
)

type entry struct {
	key string
	val string
}

// Config holds an ordered list of entries.
// Config also stores a "stem" to help make error messages
// more readable when dealing with descended configs.
type Config struct {
	entries  []entry
	defines  map[string]string
	stem     string
	includes *set.S
}

// New returns a new Config object constructed using the supplied filename
func New(filename string) *Config {
	c := &Config{
		defines:  make(map[string]string),
		includes: set.NewS(),
	}

	c.fromFile(filename)
	return c
}

func (c *Config) fromFile(filename string) {

	// Open File
	file, err := os.Open(filename)
	if err != nil {
		exit("Can't open config file: " + filename + ", " + err.Error())
	}
	defer file.Close()

	// Add File To Includes (To Prevent Circular inclusion)
	c.includes.Insert(filename)

	// Scan
	var prevKey string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Read-Line
		line := scanner.Text()

		// Apply Defines
		for word, definition := range c.defines {
			line = strings.Replace(line, word, definition, -1)
		}

		// Tokens
		tokens := strings.Fields(line)

		// Not Enough Tokens
		if len(tokens) < 2 {
			continue
		}

		///////////////////////////////
		// Substitution: User Define //
		///////////////////////////////

		if tokens[0] == "#DEFINE" {
			target := tokens[1]

			// Target Must Be Wrapped In Angle Brackets: <target>
			if !strings.HasPrefix(target, "<") || !strings.HasSuffix(target, ">") {
				exit("Bad Define - Target must be surrounded by <>: " + target + ", in line: " + line)
			}

			definition := strings.Join(tokens[2:], " ")
			c.defines[target] = definition
		}

		////////////////////////////////////////
		// Substitution: Environment Variable //
		////////////////////////////////////////

		if tokens[0] == "#ENV" {
			target := tokens[1]

			// Target Must Be Wrapped In Angle Brackets: <target>
			if !strings.HasPrefix(target, "<") || !strings.HasSuffix(target, ">") {
				exit("Bad Define - Target must be surrounded by <>: " + target + ", in line: " + line)
			}

			// Must Have Three Tokens: #ENV <target> variable
			if len(tokens) != 3 {
				exit("Bad Environment Substitution - Target must be followed with one token (representing environment-variable name): " + target + ", in line: " + line)
			}

			// Get Environment-Variable's Value
			variable := tokens[2]
			value := os.Getenv(variable)

			// Value Must Be Non-Empty
			if len(value) == 0 {
				exit("This config requires the environment variable " + variable + " to be defined according to line: " + line)
			}

			// Add To Definitions
			definition := value
			c.defines[target] = definition
		}

		//////////////////////////
		// Include Another File //
		//////////////////////////

		if tokens[0] == "#INCLUDE" {
			fmt.Println(filename, "Included", tokens[1])
			inclFile := tokens[1]

			// Check For Immediate Circular Inclusion
			if c.includes.Contains(inclFile) {
				exit("Circular or Duplicate file inclusion: " + inclFile + " found at " + filename)
			}

			// Build Config (Pass Current Includes Upward)
			i := &Config{
				defines:  make(map[string]string),
				includes: c.includes.Copy(),
			}

			// Add Defines To Include
			for k, v := range c.defines {
				i.defines[k] = v
			}

			// Construct Include Config
			i.fromFile(inclFile)

			// Add New Files To Includes
			c.includes = c.includes.Union(i.includes)

			// Add New Entries
			c.entries = append(c.entries, i.entries...)

			// Add New Defines
			for k, v := range i.defines {
				c.defines[k] = v
			}
		}

		// Comment (Skip)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Key
		key := tokens[0]

		// Join Remaining Tokens To Make Value (With Single-Space)
		val := strings.Join(tokens[1:], " ")

		// Key-Plus-Equals (Append Previous Value)
		if strings.HasSuffix(key, "+=") {

			// Confirm Key Matches Previous Key
			key = strings.TrimSuffix(key, "+=")
			if key != prevKey {
				exit("Config - Previous key does not match key with +=: " + line + ", " + prevKey + c.suffix())
			}

			c.entries[len(c.entries)-1].val += " " + val
			continue
		}

		// New-Key => Value
		e := entry{
			key: tokens[0],
			val: val,
		}
		c.entries = append(c.entries, e)

		// Set Previous-Key (for +=)
		prevKey = key
	}
}

// Has returns true if the key occurs.
func (c *Config) Has(key ...string) bool {
	joined := join(key...)

	for _, e := range c.entries {
		if e.key == joined {
			return true
		}
	}

	return false
}

// Is returns whether the value for the given key matches
// the argument passed does not exist it exits(1).
// If there are multiple occurrences of the key it exits(1)
func (c *Config) Is(key, val string) bool {
	V := c.Get(key)

	return val == V
}

// Get returns the value for the given key. If the key
// does not exist it exits(1). If there are multiple
// occurrences of the key it exits(1)
func (c *Config) Get(key ...string) string {
	joined := join(key...)
	vals := c.vals(key...)

	if len(vals) == 0 {
		exit("Config - Missing key: " + joined + c.suffix())
	}

	if len(vals) > 1 {
		exit("Config - Duplicate key: " + joined + c.suffix())
	}

	return vals[0]
}

// GetN returns the Nth value for the given key.
func (c *Config) GetN(i int, key ...string) string {
	joined := join(key...)
	vals := c.vals(key...)

	// Out-Of-Range
	if i < 0 || i >= len(vals) {
		msg := fmt.Sprintf("Config - Index out-of-range for key %s: %d (>= %d or negtive)", joined, i, len(vals))
		exit(msg + c.suffix())
	}

	return vals[i]
}

// Size returns the number of occurrences of the supplied key.
func (c *Config) Size(key ...string) int {
	var result int
	joined := join(key...)

	for _, e := range c.entries {
		if e.key == joined {
			result++
		}
	}

	return result
}

// Returns a list of strings containing all values that occur
// for the given key. The order is maintained and agrees with that
// of the file.
func (c *Config) vals(key ...string) []string {
	var result []string
	joined := join(key...)

	for _, e := range c.entries {
		if e.key == joined {
			result = append(result, e.val)
		}
	}

	return result
}

// SubKeys returns the sub-keys for the given prefix.
//
// Given:
// ------
// db.us.host www.host.com
// db.us.port 1234
// db.gb.host www.host.com
// db.gb.port 1234
//
// Then:
// -----
// SubKeys("db") => ["us", "gb"]
// SubKeys("db", "us") => ["host", "port"]
// SubKeys("db", "us", "port") => []
//
func (c *Config) SubKeys(stems ...string) []string {
	result := make([]string, 0, 1)
	prefix := join(stems...) + "."
	added := set.NewS()

	for _, e := range c.entries {

		// Found Prefix
		if strings.HasPrefix(e.key, prefix) {

			// Remove Prefix and Split Suffix
			suffix := strings.TrimPrefix(e.key, prefix)
			pieces := strings.Split(suffix, ".")

			// Add Sub-Key
			if len(pieces) > 0 {
				subKey := pieces[0]

				if !added.Contains(subKey) {
					result = append(result, subKey)
					added.Insert(subKey)
				}
			}
		}
	}

	return result
}

// HasSubKey ...
func (c *Config) HasSubKey(stems ...string) bool {

	// Sub-Key requires prefix and sub-key
	if len(stems) < 2 {
		return false
	}

	prefix := join(stems...)

	for _, e := range c.entries {
		if strings.HasPrefix(e.key, prefix) {
			return true
		}
	}

	return false
}

// HasPrefix ...
func (c *Config) HasPrefix(stems ...string) bool {
	prefix := join(stems...) + "."

	for _, e := range c.entries {
		if strings.HasPrefix(e.key, prefix) {
			return true
		}
	}

	return false
}

// Descend returns a newly created Config containing
// all key-value pairs of the original Config whose keys
// match the supplied prefix. The keys of the new Config
// are the original keys with the prefix removed. If there
// are no keys that match the supplied prefix, the returned
// Config will have no entries (but will still be a valid Config
// object).
//
// Config
// ------
// server.user Bruce
// server.host 10.1.1.1
// server.port 1234
// menu        lunch
// persons     12
//
//
// Descend("server")
// -----------------
// user Bruce
// host 10.1.1.1
// port 1234
//
func (c *Config) Descend(stems ...string) *Config {
	result := Config{}
	prefix := join(stems...) + "."

	for _, e := range c.entries {

		if strings.HasPrefix(e.key, prefix) {

			// Descended Entry (Remove Prefix)
			d := entry{
				key: strings.TrimPrefix(e.key, prefix),
				val: e.val,
			}

			// Add Descended Entry to Result
			result.entries = append(result.entries, d)
		}
	}

	// Expand Prefix
	result.stem += prefix
	return &result
}

// Joins together segments with a "."
func join(segments ...string) string {
	return strings.Join(segments, ".")
}

// Suffix returns a message containing the current stem (if it is non-empty)
// This is used to make error messages for useful
func (c *Config) suffix() string {
	if len(c.stem) == 0 {
		return ""
	}

	return " (stem=" + c.stem + ")"
}

// Exit on failure
func exit(msg string) {
	fmt.Fprintf(os.Stderr, ansi.Color(msg, "red")+"\n")
	os.Exit(1)
}
