dbl
===
A collection of functions to help avoid indeterminacy when working with floating-point numbers in Go.

What is it?
-----------
Consider the following snippet:

```
import "fmt"

x := 0.1
y := 0.2
z := x + y

switch {
case z == 0.3: fmt.Println("Z equals 0.3.")
case z <  0.3: fmt.Println("Z is less than 0.3.")
case z >  0.3: fmt.Println("Z is greater than 0.3.")
}
```

One would hope that the condition `z == 0.3` be true, however, due to floating-point
error, the actual value of z may very well be `0.300000000007` (or `0.2999999999998`). So
running this code will produce seemingly random results each time it is compiled.

One way to help fix this is to set a threshold within which two decimals are considered
equal. Determining such a threshold will depend on the domain in use. For example, in financial
contexts a number like `1e-8` works well ... okay maybe `1e-12` if you're working in Turkish
Lira exchange rates. The `dbl` package defaults to `1e-10`.

Usage
-----
```
import (
   "dbl"
   "fmt"
)

var x float64
var y float64

dbl.SetEpsilon(1e-8) // Default value is 1e-10

if dbl.LT(x, y)  { fmt.Println("x is less than y")                }
if dbl.GT(x, y)  { fmt.Println("x is greater than y")             }
if dbl.LE(x, y)  { fmt.Println("x is less than or equal to y")    }
if dbl.GE(x, y)  { fmt.Println("x is greater than or equal to y") }
if dbl.EQ(x, y)  { fmt.Println("x is equal to y")                 }
if dbl.NE(x, y)  { fmt.Println("x is not equal to y")             }
if dbl.IsZero(x) { fmt.Println("x is zero")                       }
if dbl.IsPos(x)  { fmt.Println("x is positive")                   }
if dbl.IsNeg(x)  { fmt.Println("x is negative")                   }
```

Smoothened Math functions:
```
// Divides x and y, returns zero if IsZero(y)
z := dbl.SafeDiv(x, y)
```
