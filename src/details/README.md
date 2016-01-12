details
=======

A package to help control logging of details in Go apps

Usage
=====

```
package main

import (
  "fmt"
  "github.com/enova/tokyo/src/details"
)

func main() {

  // Set Details-Level
  details.Set("Info")  // Can set to: None, Info, More, Most
  
  // Check Details-Level
  if details.None() { fmt.Println("Shhh!") }
  if details.Info() { fmt.Println("Info details") }
  if details.More() { fmt.Println("Info and More details") }
  if details.Most() { fmt.Println("Info and More and Most details") }
}
```

The `details.Set()` function also accepts numerical codes (as strings).
```
details.Set("0")  // Same as details.Set("None")
details.Set("1")  // Same as details.Set("Info")
details.Set("2")  // Same as details.Set("More")
details.Set("3")  // Same as details.Set("Most")
```

All functions is the `details` package are thread-safe.
