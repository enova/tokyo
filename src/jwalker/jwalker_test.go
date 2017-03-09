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

// Validity ...
func validity(v interface{}, ok bool) bool {
	return ok
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

	// B
	tasty := w.Key("tasty")
	assert.True(tasty.Ok())

	b, ok := tasty.B()
	assert.True(ok)
	assert.True(b)

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

	// KeyB
	rich, ok := owner.KeyB("rich")
	assert.True(ok)
	assert.False(rich)

	// Keys
	keys := owner.Keys()
	assert.Equal(len(keys), 5)
	assert.True(Contains(keys, "name"))
	assert.True(Contains(keys, "year"))
	assert.True(Contains(keys, "ssid"))
	assert.True(Contains(keys, "frac"))
	assert.True(Contains(keys, "rich"))

	// Miss Then Hit
	owner = w.Key("owner_x")
	assert.False(owner.Ok())

	owner = w.Key("owner")
	assert.True(owner.Ok())

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

	// Miss Then Hit
	at = teams.At(3)
	assert.False(at.Ok())

	at = teams.At(2)
	assert.True(at.Ok())

	// Short-Forms At
	mixed := w.Key("mixed")
	assert.True(mixed.Ok())

	s, ok = mixed.AtS(0)
	assert.True(ok)
	assert.Equal(s, "banana")

	u, ok = mixed.AtU32(1)
	assert.True(ok)
	assert.Equal(u, uint32(64))

	i, ok = mixed.AtI(1)
	assert.True(ok)
	assert.Equal(i, 64)

	f, ok = mixed.AtF64(2)
	assert.True(ok)
	assert.Equal(f, 29.98)

	b, ok = mixed.AtB(3)
	assert.True(ok)
	assert.True(b)

	b, ok = mixed.AtB(4)
	assert.True(ok)
	assert.False(b)

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
}

func TestFailedAccess(t *testing.T) {
	assert := assert.New(t)

	bytes := ReadFile(assert, "test/a.json")
	w, err := New(bytes)
	assert.Nil(err)

	badKey := w.Key("fruitx")
	assert.False(badKey.Ok())

	badKey = w.Key("ownerx").Key("name")
	assert.False(badKey.Ok())

	badKey = w.Key("owner").Key("namex")
	assert.False(badKey.Ok())

	badKey = w.Key("owner").Key("name").Key("x").Key("y").Key("z")
	assert.False(badKey.Ok())

	badAt := w.Key("owner").At(0)
	assert.False(badAt.Ok())

	mixed := w.Key("mixed")
	badAt = mixed.At(-1)
	assert.False(badAt.Ok())

	badAt = mixed.At(5)
	assert.False(badAt.Ok())

	s := mixed.At(0)
	assert.False(validity(s.B()))
	assert.False(validity(s.F64()))
	assert.False(validity(s.I()))
	assert.False(validity(s.U32()))
	assert.True(validity(s.S()))

	i := mixed.At(1)
	assert.False(validity(i.B()))
	assert.True(validity(i.F64()))
	assert.True(validity(i.I()))
	assert.True(validity(i.U32()))
	assert.False(validity(i.S()))

	f := mixed.At(2)
	assert.False(validity(f.B()))
	assert.True(validity(f.F64()))
	assert.True(validity(f.I()))
	assert.True(validity(f.U32()))
	assert.False(validity(f.S()))

	b := mixed.At(3)
	assert.True(validity(b.B()))
	assert.False(validity(b.F64()))
	assert.False(validity(b.I()))
	assert.False(validity(b.U32()))
	assert.False(validity(b.S()))

	b = mixed.At(4)
	assert.True(validity(b.B()))
	assert.False(validity(b.F64()))
	assert.False(validity(b.I()))
	assert.False(validity(b.U32()))
	assert.False(validity(b.S()))
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

func TestTrace(t *testing.T) {
	assert := assert.New(t)

	file := ReadFile(assert, "test/b.json")
	w, err := New(file)
	assert.Nil(err)

	value := w.Key("glossary").At(100)
	assert.Equal(value.Trace(), "[location] => key: glossary [failure] => at: 100 (not an array)")
}

func TestDiagnostics(t *testing.T) {
	assert := assert.New(t)

	file := ReadFile(assert, "test/b.json")
	w, err := New(file)
	assert.Nil(err)

	value := w.Key("glossary").Key("title2")
	assert.Equal(value.Location(), "key: glossary")
	assert.Equal(value.Failure(), "key: title2 (key does not exist)")

	value = w.Key("glossary").Key("title2").Key("thiswontexist")
	assert.Equal(value.Location(), "key: glossary")
	assert.Equal(value.Failure(), "key: title2 (key does not exist)")

	value = w.Key("glossary").At(100)
	assert.Equal(value.Location(), "key: glossary")
	assert.Equal(value.Failure(), "at: 100 (not an array)")

	value = w.Key("glossary").Key("Div").Key("List").Key("Entry").Key("Def").Key("SeeAlso").At(10)
	assert.Equal(value.Location(), "key: glossary | key: Div | key: List | key: Entry | key: Def | key: SeeAlso")
	assert.Equal(value.Failure(), "at: 10 (out of range, size=2)")

	value = w.Key("glossary").Key("Div").Key("List").Key("Entry").Key("Def").Key("SeeAlso").At(-5)
	assert.Equal(value.Location(), "key: glossary | key: Div | key: List | key: Entry | key: Def | key: SeeAlso")
	assert.Equal(value.Failure(), "at: -5 (out of range, size=2)")

	value = w.Key("glossary").Key("title").Key("another_title")
	assert.Equal(value.Failure(), "key: another_title (not a map)")
}
