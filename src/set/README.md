set
===

An implementation of a Set

What is Set?
------------

Go has built-in Arrays and Maps but no Sets. This package implements a Set by wrapping `map[string]bool` (i.e. map from `string` to `bool`). 

Importing This Package
----------------------
Add `github.com/tokyo/src/set` to your set of imports:

```
package main

import (
  "fmt"
  "github.com/tokyo/src/set"
)

func main() {
  s := NewS("A", "B", "C")
  fmt.Println("The set s has", s.Size(), "elements")
}
```

Usage
-----
A set can be created in a few different ways:
```
x := set.NewS("A", "B", "C")
y := set.Parse("A B C")

z := set.NewS()
z.Insert("A")
z.Insert("B")
z.Insert("C")
```

Set arithmetic:
```
x := set.NewS("A", "B")
y := set.NewS("B", "C")

set.Or(x, y)  // Union: {"A", "B", "C"}
set.And(x, y) // Intersection: {"B"}
set.Sub(x, y) // Subtraction: {"A"}
set.Sub(x, y) // Subtraction: {"C"}
```
Elements can be removed by using the `Delete` method:
```
x := set.NewS("A", "B")
x.Delete("A")
```
Other methods and functions:
```
x.Contains("A")    // Returns true if the set x contains the element "A"
x.Size()           // Returns the number of elements in the set x
x.Empty()          // Returns true is the set x is empty
set.IsSubset(x, y) // Returns true if the set x is a subset of the set y
set.EQ(x, y)       // Returns true if the set x equals the set y
```
Iterating over Elements
-----------------------
There are two immediate ways to iterate over a set. One way is to use the method `set.Elements()` which returns the elements of the set as an array `[]string`:
```
x := set.NewS("A", "B")
elements := x.Elements()
for _, e := range elements {
  fmt.Println(e)
}
```
The other way is to directly access the internal `map[string]bool` of the set using the method `set.Map()`:
```
x := set.NewS("A", "B")
for e, _ := range x.Map() {
  fmt.Println(e)
}
```
A Note on the Set Type
----------------------
The code uses a capital `S` to denote `string`. If in the future a set of integers is added to this project it would presumably be named `I` (pull-requests welcome!).
