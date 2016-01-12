package cfg

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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

// base.cfg
// -------------
// width 24
//
// #INCLUDE parent.cfg

// parent.cfg
// ----------
// height 36

func TestCfg(t *testing.T) {
	assert := assert.New(t)

	// Set Environment Variable To Test (Below)
	assert.Nil(os.Setenv("CFG_TEST", "all-good"))
	cfg := New("test.cfg")

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
