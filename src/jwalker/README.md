# jwalker

A Go package to access JSON data bytes

# Description

The reflective properties of Go structs permit smooth marshalling and unmarshalling of JSON encoded bytes. But what if you don't want to unmarshal an entire JSON packet into a single struct? Suppose you want to extract entries like `name := data["players"][2]["name"]`. You certainly _can_ do this in Go, though it requires a few lines per level of depth. So the aforementioned example would require about nine lines of code. This package wraps that code to reduce such a query to a couple lines.

# Usage

To use this package you must first construct a `W` object using `jwalker.New(data)` where `data` is a slice of JSON-encoded bytes. If the underlying object is a map then use the method `Key(s)` to descend into the value for key `s`. If the underlying object is an array then use the method `At(i)` to descend into the value at index `i`. Each of these methods (`Key` and `At`) return a `W` object pointing to the correct _child_ data. If a lookup was invalid (e.g. calling `Key()` for a non-existent key), then the resulting object be set to invalid. You can call the method `Ok()` to check the validity of a `W` instance.

When `W` points to a terminal value you can use the methods `S(), I(), U32(), F64()` to extract that value. Each of these methods returns a value along with a `bool` indicating success.

Both the `Key()` and `At()` methods will always return an object of type `*W`. When an invalid lookup takes place, the value returned will still be a bonafide `W` object. This design allows the user to chain arbitrary lookups without worrying about `nil` pointer exceptions. For example, consider the following code:

```
data := []byte(`{ "name": "gopher" }`)
w, _ := jwalker.New(data)

child := w.Key("who").Key("what").Key("where").At[14].Key("when")
valid := child.Ok() // false

child = w.Key("name")
valid = child.Ok() // true
```

You can arbitrarily chain calls to `Key()` and `At()` on a `W` instance. If any of the calls is invalid, the final instance will return false when its `Ok()` method is invoked.

Here are the contents of `test/a.json`:

```
{
	"fruit": "apple",
	"width": 32,
	"ratio": 12.50,

	"owner": {
		"name": "gopher",
		"year": 2010,
		"ssid": 1234567890,
		"frac": 1.234
	},

	"teams": [
		"red",
		"green",
		"blue"
	]
}

```

Example Usage:

```
import (
  "io/ioutil"

  "github.com/enova/tokyo/src/jwalker"
)

func main() {
  bytes, _ := ioutil.ReadFile("test/a.json")
  w, _ := jwalker.New(bytes)
  
  // Look-Up Key/Value
  fruit, ok := w.Key("fruit").S()
  width, ok := w.Key("width").I()
  ratio, ok := w.Key("ratio").F64()

  // Look-Up Key/Value Short-Form
  fruit, ok := w.KeyS("fruit")
  width, ok := w.KeyI("width")
  ratio, ok := w.KeyF64("ratio")

  // Descend One Level
  name, ok := w.Key("owner").Key("name").S()
  year, ok := w.Key("owner").Key("year").I()
  ssid, ok := w.Key("owner").Key("ssid").U32()
  frac, ok := w.Key("owner").Key("frac").F64()

  // Descend One Level Short-Form
  name, ok := w.Key("owner").KeyS("name")
  year, ok := w.Key("owner").KeyI("year")
  ssid, ok := w.Key("owner").KeyU32("ssid")
  frac, ok := w.Key("owner").KeyF64("frac")

  // Descend Into Array
  team, ok := w.Key("teams").At(0).S()
  team, ok := w.Key("teams").At(1).S()
  team, ok := w.Key("teams").At(2).S()

  // Descend Into Array Short-Form
  team, ok := w.Key("teams").AtS(0)
  team, ok := w.Key("teams").AtS(1)
  team, ok := w.Key("teams").AtS(2)
}
```

Notice the short-form key-value methods `KeyS(), KeyI(), KeyU32(), KeyF64()`. Analogous to these are the array methods `AtS(), AtI(), AtU32() , AtF64()`.

#### Acknowledgment

A portion of the test-data was taken from [json.org](http://json.org/example.html).
