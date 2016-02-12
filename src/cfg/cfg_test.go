package cfg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

// Code holds a block of code
type Code struct {
	lines []string
}

// Reset clears the lines
func (c *Code) Reset() {
	c.lines = make([]string, 0)
}

// Add adds a line
func (c *Code) Add(line string) {
	c.lines = append(c.lines, line)
}

// Text renders the code-text
func (c *Code) Text() string {
	var result string
	for _, line := range c.lines {
		result += line + "\n"
	}
	return result
}

// Run executes the code and returns the error
func (c *Code) Run() error {
	var result error

	// Header
	header := Code{}
	header.Add(`package main`)
	header.Add(`import "github.com/enova/tokyo/src/cfg"`)
	header.Add(``)
	header.Add(`func main() {`)

	// Footer
	footer := Code{}
	footer.Add(`}`)

	// Construct Code
	text := header.Text() + c.Text() + footer.Text()
	fmt.Println(text)

	// Create Test-Directory
	os.MkdirAll("test_code", 0755)

	// Create Test-App From Code
	ioutil.WriteFile("test_code/test_code.go", []byte(text), 0644)

	// Run Test-App
	cmd := "go run test_code/test_code.go"
	output, result := exec.Command("bash", "-c", cmd).CombinedOutput()

	// Logging For Failures
	fmt.Println("Output\n", string(output))
	fmt.Println("Result\n", result)

	// Delete Test-Directory
	cmd = "rm -rf test_code"
	exec.Command("bash", "-c", cmd).Output()

	return result
}

////////////
// test.cfg
// --------
// #DEFINE <path> /usr/share
// #ENV <status> CFG_TEST
// #INCLUDE base.cfg
//
// db.us.user bruce
// db.us.name inventory
// db.us.host 10.144.1.1
// db.us.port 1111
//
// db.uk.user leroy
// db.uk.name schedules
// db.uk.host 10.144.1.2
// db.uk.port 2222
//
// email support@firm.com
// email billing@firm.com
// email sales@firm.com
//
// slogan We love to code
//
// sentence   The cat ran after
// sentence+= the mouse

////////////
// base.cfg
// -------------
// width 24
//
// #INCLUDE parent.cfg

//////////////
// parent.cfg
// ----------
// height 36

func TestCfg(t *testing.T) {
	assert := assert.New(t)

	// Set Environment Variable To Test (Below)
	assert.Nil(os.Setenv("CFG_TEST", "all-good"))
	cfg := New("test/test.cfg")

	// Has
	assert.True(cfg.Has("db", "us", "user"))
	assert.True(cfg.Has("db", "us", "name"))
	assert.True(cfg.Has("db", "us", "host"))
	assert.True(cfg.Has("db", "us", "port"))
	assert.False(cfg.Has("db", "us", "namex"))

	// Get
	assert.Equal("bruce", cfg.Get("db", "us", "user"))
	assert.Equal("schedules", cfg.Get("db", "uk", "name"))
	assert.Equal("We love to code", cfg.Get("slogan"))
	assert.Equal("The cat ran after the mouse", cfg.Get("sentence"))

	// SubKeys
	assert.Equal(2, len(cfg.SubKeys("db")))
	assert.Equal("us", cfg.SubKeys("db")[0])
	assert.Equal("uk", cfg.SubKeys("db")[1])

	assert.Equal(4, len(cfg.SubKeys("db", "us")))
	assert.Equal("user", cfg.SubKeys("db", "us")[0])
	assert.Equal("name", cfg.SubKeys("db", "us")[1])
	assert.Equal("host", cfg.SubKeys("db", "us")[2])
	assert.Equal("port", cfg.SubKeys("db", "us")[3])

	// HasSubKey
	assert.True(cfg.HasSubKey("db", "us"))
	assert.False(cfg.HasSubKey("db", "ca"))
	assert.True(cfg.HasSubKey("db", "us", "name")) // "name" is a sub-key of "db.us"
	assert.False(cfg.HasSubKey("db"))              // Not enough arguments - Needs a prefix and a sub-key.

	// HasPrefix
	assert.True(cfg.HasPrefix("db"))
	assert.True(cfg.HasPrefix("db", "us"))
	assert.False(cfg.HasPrefix("db", "us", "name")) // Not a prefix, it's a complete key
	assert.False(cfg.HasPrefix("db", "ca"))
	assert.False(cfg.HasPrefix("d"))

	// Size
	assert.Equal(3, cfg.Size("email"))

	// GetN
	assert.Equal("support@firm.com", cfg.GetN(0, "email"))
	assert.Equal("billing@firm.com", cfg.GetN(1, "email"))
	assert.Equal("sales@firm.com", cfg.GetN(2, "email"))

	// Includes
	assert.Equal("24", cfg.Get("width"))
	assert.Equal("36", cfg.Get("height"))
	assert.Equal("48", cfg.Get("depth")) // Defined in base.cfg
	assert.Equal("60", cfg.Get("mass"))  // Defined in parent.cfg

	// Descend
	d := cfg.Descend("db", "us")
	assert.Equal(d.Get("user"), "bruce")
	assert.Equal(d.Get("name"), "inventory")
	assert.Equal(d.Get("host"), "10.144.1.1")
	assert.Equal(d.Get("port"), "1111")

	d = cfg.Descend("db")
	assert.Equal(d.Get("us", "user"), "bruce")
	assert.Equal(d.Get("us", "name"), "inventory")
	assert.Equal(d.Get("us", "host"), "10.144.1.1")
	assert.Equal(d.Get("us", "port"), "1111")

	assert.Equal(d.Get("uk", "user"), "leroy")
	assert.Equal(d.Get("uk", "name"), "schedules")
	assert.Equal(d.Get("uk", "host"), "10.144.1.2")
	assert.Equal(d.Get("uk", "port"), "2222")

	// Defines (#DEFINE)
	assert.Equal(cfg.Get("lib"), "/usr/share/lib")
	assert.Equal(cfg.Get("bin"), "/usr/share/bin")

	// Environment-Variable Substitution (#ENV)
	assert.Equal(cfg.Get("message"), "all-good")
}

// Test Exit-Points
func TestExit(t *testing.T) {
	assert := assert.New(t)

	code := Code{}

	code.Reset()
	code.Add(`cfg.New("test/parent.cfg")`)
	assert.Nil(code.Run(), "Read a valid config")

	code.Reset()
	code.Add(`cfg.New("test/bad/missing_env.cfg")`)
	assert.NotNil(code.Run(), "Missing environment variable")

	code.Reset()
	code.Add(`cfg.New("test/nonexistent.cfg")`)
	assert.NotNil(code.Run(), "Read a non-existent config file")

	code.Reset()
	code.Add(`cfg.New("test/bad/bad_define_key.cfg")`)
	assert.NotNil(code.Run(), "Bad DEFINE key in config")

	code.Reset()
	code.Add(`cfg.New("test/bad/bad_env_variable.cfg")`)
	assert.NotNil(code.Run(), "Bad ENV missing-variable name in config")

	code.Reset()
	code.Add(`cfg.New("test/bad_circular.cfg")`)
	assert.NotNil(code.Run(), "Circular file inclusion")

	code.Reset()
	code.Add(`cfg := cfg.New("test/parent.cfg")`)
	code.Add(`cfg.Get("nonexistent_key")`)
	assert.NotNil(code.Run(), "Get a non-existent key")

	code.Reset()
	code.Add(`cfg := cfg.New("test/bad/bad_suffix.cfg")`)
	code.Add(`cfg.Get("nonexistent_key")`)
	assert.NotNil(code.Run(), "Get a non-existent key")

	code.Reset()
	code.Add(`cfg := cfg.New("test/bad/duplicate_key.cfg")`)
	code.Add(`cfg.Get("fruits")`)
	assert.NotNil(code.Run(), "Can't call Get() when there are duplicate keys")

	code.Reset()
	code.Add(`cfg := cfg.New("test/bad/duplicate_key.cfg")`)
	code.Add(`cfg.GetN(0, "fruits")`)
	code.Add(`cfg.GetN(1, "fruits")`)
	assert.Nil(code.Run(), "Calling GetN() when there are duplicate keys")

	code.Reset()
	code.Add(`cfg := cfg.New("test/bad/duplicate_key.cfg")`)
	code.Add(`cfg.GetN(-1, "fruits")`)
	assert.NotNil(code.Run(), "Bad call to GetN(), out of range")

	code.Reset()
	code.Add(`cfg := cfg.New("test/bad/duplicate_key.cfg")`)
	code.Add(`cfg.GetN(2, "fruits")`)
	assert.NotNil(code.Run(), "Bad call to GetN(), out of range")
}
