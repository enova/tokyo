package args

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// KeyVal stores a key and a value
type keyVal struct {
	key string
	val string
}

// Args contains arguments and their derivatives
type Args struct {
	raws    []string // Raw Arguments
	vals    []string // Non-Optional Arguments
	uniOpts []string // Unary Flags: -Flag
	binOpts []keyVal // Binary Flags: -Key=Value
}

// Parse returns a newly created Args using os.Args
func Parse() *Args {
	return New(os.Args)
}

// New returns a newly created Args
func New(raws []string) *Args {
	a := Args{raws: make([]string, len(raws))}

	// Raws
	copy(a.raws, raws)

	// Values
	for i := 0; i < len(raws); i++ {
		raw := raws[i]

		if strings.HasPrefix(raw, "-") {

			// Option
			if strings.Contains(raw, "=") {

				// Binary Option
				if tokens := strings.SplitN(raw, "=", 2); len(tokens) == 2 {
					key := strings.TrimPrefix(tokens[0], "-")
					binOpt := keyVal{key, tokens[1]}
					a.binOpts = append(a.binOpts, binOpt)
				}
			} else {

				// Unary Option
				uniOpt := strings.TrimPrefix(raw, "-")
				a.uniOpts = append(a.uniOpts, uniOpt)
			}
		} else {

			// Non-Option
			a.vals = append(a.vals, raw)
		}
	}

	return &a
}

// Raws returns the list of raw arguments
func (a *Args) Raws() []string {
	return a.raws
}

// Size returns the number of non-optional arguments
func (a *Args) Size() int {
	return len(a.vals)
}

// Get returns the non-optional argument in the supplied position. If
// the supplied index is out of range it Exits(1). If you don't want
// your application to die in case of an out-of-range index, use the
// method Size().
func (a *Args) Get(i int) string {
	if i < 0 || i >= a.Size() {
		fmt.Fprintf(os.Stderr, "Args: Index out of range: %d  [0, %d)\n", i, a.Size())
		os.Exit(1)
	}

	return a.vals[i]
}

// IsOn returns true if the supplied unary argument is present
func (a *Args) IsOn(s string) bool {
	for _, u := range a.uniOpts {
		if u == s {
			return true
		}
	}
	return false
}

// IsOff returns true if the supplied unary argument is not present
func (a *Args) IsOff(s string) bool {
	return !a.IsOn(s)
}

// HasOpt return true if the supplied binary option is present
func (a *Args) HasOpt(key string) bool {
	for _, b := range a.binOpts {
		if b.key == key {
			return true
		}
	}
	return false
}

// GetOpt return the value corresponding to the supplied key. It
// exits if the option is not found (Exit(1)). If you don't want
// your application to die in case of a missing option, then use
// the method HasOpt() to check whether an option has been supplied.
func (a *Args) GetOpt(key string) string {
	for _, b := range a.binOpts {
		if b.key == key {
			return b.val
		}
	}

	// Invalid Option
	fmt.Fprintf(os.Stderr, "Args - Missing option: %s\n", key)
	os.Exit(1)

	// Never Gets Here
	return ""
}

// GetOptI return the value corresponding to the supplied key as an
// integer (int). It exits if the option is not found (Exit(1))
func (a *Args) GetOptI(key string) int {
	s := a.GetOpt(key)

	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Args: Invalid option value - %s for key %s - Must be an integer", s, key)
		os.Exit(1)
	}

	return i
}
