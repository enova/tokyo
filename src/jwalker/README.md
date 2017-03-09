# jwalker

A Go package to walk through JSON objects without having to unmarshal their bytes into a `struct`.

# Description

The reflective properties of Go structs permit smooth marshalling and unmarshalling of JSON encoded bytes. But what if you don't want to unmarshal an entire JSON packet into a single struct? Suppose you want to extract entries like `name := data["players"][2]["name"]`. You certainly _can_ do this in Go, though it requires a few lines per level of depth. So the aforementioned example would require about nine lines of code. This package wraps that code to reduce such a query to a couple lines.

# Usage

To use this package you must first construct a `W` object using `jwalker.New(data)` where `data` is a slice of JSON-encoded bytes. If the underlying object is a map then use the method `Key(s)` to descend into the value for key `s`. If the underlying object is an array then use the method `At(i)` to descend into the value at index `i`. Each of these methods (`Key` and `At`) returns a `W` object pointing to the correct _child_ data. If the lookup was invalid (e.g. calling `Key()` for a non-existent key), then the resulting object be set to invalid. You can call the method `Ok()` to check the validity of a `W` instance.

When `W` points to a terminal value you can use the methods `S(), I(), U32(), F64(), B()` to extract that value. Each of these methods returns a value along with a `bool` indicating success.

Both the `Key()` and `At()` methods will always return an object of type `*W`. When an invalid lookup takes place, the value returned will still be a bonafide `W` object. This design allows the user to chain arbitrary lookups without worrying about `nil` pointer exceptions. For example, consider the following code:

```go
data := []byte(`{ "name": "gopher" }`)
w, _ := jwalker.New(data)

child := w.Key("who").Key("what").Key("where").At[14].Key("when")
valid := child.Ok() // false

child = w.Key("name")
valid = child.Ok() // true
```

You can arbitrarily chain calls to `Key()` and `At()` on a `W` instance. If any of the calls is invalid, the final instance will return false when its `Ok()` method is invoked.

## Diagnostics

The diagnostic methods `Location()` and `Failure()` reveal where a failure occurred. `Location()` indicates the sequence of successful accesses, and `Failure()` indicates the failed operation. The convenience method `Trace()` includes both location and failure in a single string. In the following code assume a failure occurs at the call `.Key("sign")`:

```go
// Assume a failure happens at Key("sign")
value := w.Key("fruits").At(2).Key("sign").Key("primary").At(3)

// Diagnostics
location := value.Location() // "key: fruits | at: 2"
failure := value.Failure()   // "key: sign (key does not exist)"
trace := value.Trace()       // "[location] => key: fruits | at: 2 [failure] => key: sign (key does not exist)"
```

The location indicates that the successful retrievals were `.Key("fruits")` and `.At(2)`. The failure shows there was an unsuccessful call to `Key("sign")` and it also provides a reason `(key does not exist)`. The other possible reason is `(not a map)`. For unsuccessful calls to `At(i)` the two reasons are `(out of range)` and `(not an array)`.

## Example

Here are the contents of `test/a.json`:

```json
{
  "fruit": "apple",
  "width": 32,
  "ratio": 12.50,
  "tasty": true,

  "mixed": [
	"banana",
	64,
	29.98,
	true,
	false
  ],

  "owner": {
	"name": "gopher",
	"year": 2010,
	"ssid": 1234567890,
	"frac": 1.234,
	"rich": false
  },

  "teams": [
	"red",
	"green",
	"blue"
  ]
}
```

Example code to interact with the above:

```go
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
  tasty, ok := w.Key("tasty").B()

  // Look-Up Key/Value Short-Form
  fruit, ok = w.KeyS("fruit")
  width, ok = w.KeyI("width")
  ratio, ok = w.KeyF64("ratio")
  tasty, ok = w.KeyB("tasty")

  // Descend One Level
  name, ok := w.Key("owner").Key("name").S()
  year, ok := w.Key("owner").Key("year").I()
  ssid, ok := w.Key("owner").Key("ssid").U32()
  frac, ok := w.Key("owner").Key("frac").F64()
  rich, ok := w.Key("owner").Key("rich").B()

  // Descend One Level Short-Form
  name, ok = w.Key("owner").KeyS("name")
  year, ok = w.Key("owner").KeyI("year")
  ssid, ok = w.Key("owner").KeyU32("ssid")
  frac, ok = w.Key("owner").KeyF64("frac")
  rich, ok = w.Key("owner").KeyB("rich")

  // Descend Into Array
  s, ok := w.Key("mixed").At(0).S()
  i, ok := w.Key("mixed").At(1).I()
  f, ok := w.Key("mixed").At(2).F64()
  b, ok := w.Key("mixed").At(3).B()

  // Descend Into Array Short-Form
  s, ok = w.Key("mixed").AtS(0)
  i, ok = w.Key("mixed").AtI(1)
  f, ok = w.Key("mixed").AtF64(2)
  b, ok = w.Key("mixed").AtB(3)

  // Iteration (Map)
  keys := w.Key("owner").Keys()
  for _, k := range keys {
    value := w.Key("owner").Key(k)
  }

  // Iteration (Array)
  for i = 0; i < w.Key("mixed").Len(); i++ {
    value := w.Key("mixed").At(i)
  }

  // Diagnostics (Key Not Found)
  fruit, ok = w.Key("owner").Key("favorite_fruit")
  if !ok {
    panic(fruit.Trace()) // Displays: [location] => key: owner [failure] => key: favorite_fruit (key not found)
  }

  // Diagnostics (Not A Map)
  team, ok := w.Key("teams").Key("first")
  if !ok {
    panic(team.Trace()) // Displays: [location] => key: teams [failure] => key: first (not a map)
  }

  // Diagnostics (Out Of Range)
  fruit, ok = w.Key("teams").At(5)
  if !ok {
    panic(fruit.Trace()) // Displays: [location] => key: teams [failure] => at: 5 (out of range, size 3)
  }

  // Diagnostics (Not An Array)
  fruit, ok = w.Key("owner").At(0)
  if !ok {
    panic(fruit.Trace()) // Displays: [location] => key: owner [failure] => at: 0 (key not found)
  }
}
```

Notice the short-form key-value methods `KeyS(), KeyI(), KeyU32(), KeyF64(), KeyB()`. Analogous to these are the array methods `AtS(), AtI(), AtU32() , AtF64(), AtB()`.
