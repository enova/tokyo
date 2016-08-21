package jwalker

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReadFile ...
func ReadFile(assert *assert.Assertions, filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	assert.Nil(err)
	return bytes
}

// Contains ...
func Contains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func TestAccess(t *testing.T) {
	assert := assert.New(t)

	bytes := ReadFile(assert, "test/a.json")
	w, err := New(bytes)
	assert.Nil(err)

	// S
	fruit := w.Key("fruit")
	assert.True(fruit.Ok())

	s, ok := fruit.S()
	assert.True(ok)
	assert.Equal(s, "apple")

	// I
	width := w.Key("width")
	assert.True(width.Ok())

	i, ok := width.I()
	assert.True(ok)
	assert.Equal(i, 32)

	// U32
	u, ok := width.U32()
	assert.True(ok)
	assert.Equal(u, uint32(32))

	// F64
	ratio := w.Key("ratio")
	assert.True(ratio.Ok())

	f, ok := ratio.F64()
	assert.True(ok)
	assert.Equal(f, 12.50)

	//////////////////
	// Descend: Map //
	//////////////////
	
	owner := w.Key("owner")
	assert.True(owner.Ok())

	// KeyS
	name, ok := owner.KeyS("name")
	assert.True(ok)
	assert.Equal(name, "gopher")

	// KeyI
	year, ok := owner.KeyI("year")
	assert.True(ok)
	assert.Equal(year, 2010)

	// KeyU32
	ssid, ok := owner.KeyU32("ssid")
	assert.True(ok)
	assert.Equal(ssid, uint32(1234567890))

	// KeyF64
	frac, ok := owner.KeyF64("frac")
	assert.True(ok)
	assert.Equal(frac, 1.234)

	// Keys
	keys := owner.Keys()
	assert.Equal(len(keys), 4)
	assert.True(Contains(keys, "name"))
	assert.True(Contains(keys, "year"))
	assert.True(Contains(keys, "ssid"))
	assert.True(Contains(keys, "frac"))
	
	////////////////////
	// Descend: Array //
	////////////////////
	
	teams := w.Key("teams")
	assert.True(teams.Ok())
	assert.Equal(teams.Len(), 3)

	// Item 0
	at := teams.At(0)
	assert.True(at.Ok())

	team, ok := at.S()
	assert.True(ok)
	assert.Equal(team, "red")

	team, ok = teams.AtS(0)
	assert.True(ok)
	assert.Equal(team, "red")

	// Item 1
	at = teams.At(1)
	assert.True(at.Ok())

	team, ok = at.S()
	assert.True(ok)
	assert.Equal(team, "green")

	team, ok = teams.AtS(1)
	assert.True(ok)
	assert.Equal(team, "green")

	// Item 2
	at = teams.At(2)
	assert.True(at.Ok())

	team, ok = at.S()
	assert.True(ok)
	assert.Equal(team, "blue")

	team, ok = teams.AtS(2)
	assert.True(ok)
	assert.Equal(team, "blue")

	//////////////
	// Chaining //
	//////////////

	// Key + Key
	name, ok = w.Key("owner").KeyS("name")
	assert.True(ok)
	assert.Equal(name, "gopher")
	
	year, ok = w.Key("owner").KeyI("year")
	assert.True(ok)
	assert.Equal(year, 2010)

	ssid, ok = w.Key("owner").KeyU32("ssid")
	assert.True(ok)
	assert.Equal(ssid, uint32(1234567890))

	frac, ok = w.Key("owner").KeyF64("frac")
	assert.True(ok)
	assert.Equal(frac, 1.234)

	// Key + At
	team, ok = w.Key("teams").AtS(0)
	assert.True(ok)
	assert.Equal(team, "red")

	team, ok = w.Key("teams").AtS(1)
	assert.True(ok)
	assert.Equal(team, "green")

	team, ok = w.Key("teams").AtS(2)
	assert.True(ok)
	assert.Equal(team, "blue")

	team, ok = w.Key("teams").AtS(3)
	assert.False(ok)

	/////////////
	// Invalid //
	/////////////

	badKey := w.Key("fruitx")
	assert.False(badKey.Ok())

	badKey = w.Key("ownerx").Key("name")
	assert.False(badKey.Ok())

	badKey = w.Key("owner").Key("namex")
	assert.False(badKey.Ok())

	badAt := w.Key("owner").At(0)
	assert.False(badAt.Ok())

	badKey = w.Key("owner").Key("name").Key("x").Key("y").Key("z")
	assert.False(badKey.Ok())
}

func TestChain(t *testing.T) {
	assert := assert.New(t)

	bytes := ReadFile(assert, "test/b.json")
	w, err := New(bytes)
	assert.Nil(err)

	value, ok := w.Key("glossary").Key("title").S()
	assert.True(ok)
	assert.Equal(value, "example glossary")

	value, ok = w.Key("glossary").Key("Div").Key("title").S()
	assert.True(ok)
	assert.Equal(value, "S")

	value, ok = w.Key("glossary").Key("Div").Key("List").Key("Entry").Key("ID").S()
	assert.True(ok)
	assert.Equal(value, "SGML")

	value, ok = w.Key("glossary").Key("Div").Key("List").Key("Entry").Key("Def").Key("SeeAlso").At(0).S()
	assert.True(ok)
	assert.Equal(value, "GML")

	value, ok = w.Key("glossary").Key("Div").Key("List").Key("Entry").Key("Def").Key("SeeAlso").At(1).S()
	assert.True(ok)
	assert.Equal(value, "XML")
}
