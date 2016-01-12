args
====

A package to parse command-line arguments


Usage
-----
The `args` package allows the user to pass ordered command-line arguments along with unordered options. All options begin with a single dash, e.g. `-debug` or `-threads=4`. Unary options (also known as _boolean_ options or _flags_) require no values as in `-debug`. Binary options require a second value as in `-threads=4`. Here `threads` is the _option_ and `4` is the _option-value_.


To use `args` in your code you must import the package and create an `args` instance:

```go
import (
  "os"
  "github.com/tokyo/src/args"
)

func main() {
  args := args.New(os.Args) // <-- That's it!
}
```
The method `args.Size()` skips over options and returns the number of ordered arguments plus one (since position 0 is the executable name). Below are some example usages.

Use `args` to detect ordered arguments:
```go
$ ./myApp loans.txt payday

size := args.Size() // Returns 3
file := args.Get(1) // <-- loans.txt
type := args.Get(2) // <-- payday
```

Use `args` to detect unary options:

```go
$ ./myApp loans.txt payday -debug

size := args.Size() // Still returns 3 since -debug is an option, it is skipped over in the count
file := args.Get(1) // <-- loans.txt
type := args.Get(2) // <-- payday

debug := args.IsOn("debug") // <-- Returns boolean value (true)
```

Use `args` to detect binary options

```go
$ ./myApp loans.txt payday -debug -threads=7

size := args.Size() // Still returns 3, both options are skipped over in the count
file := args.Get(1) // <-- loans.txt
type := args.Get(2) // <-- payday

debug := args.IsOn("debug")       // <-- Returns boolean value (true)
threads := args.GetOpt("threads") // <-- Returns string ("7")

// Check if option is set
if args.HasOpt("threads") {
  threads := args.GetOptI("threads") // <-- Method GetOptI converts string to integer
}
```
Options can be interspersed within ordered arguments. The following two command lines are equivalent with respect to `args`:
```go
$ ./myApp loans.txt payday -debug -threads=7
$ ./myApp -debug loans.txt -threads=7 payday
```
In both examples `args.Size()` will return 3.
