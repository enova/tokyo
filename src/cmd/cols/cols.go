package main

import (
	"bufio"
	"fmt"
	"github.com/enova/tokyo/src/alert"
	"github.com/enova/tokyo/src/args"
	"github.com/enova/tokyo/src/lax"
	"os"
	"strings"
)

func main() {
	args := args.New(os.Args)

	// Help
	if args.IsOn("h") {
		fmt.Fprintf(os.Stderr, usage())
		os.Exit(0)
	}

	// Read Columns
	var cols []uint32
	for a := 1; a < args.Size(); a++ {
		col := lax.ParseUint32(args.Get(a))
		if col <= 0 {
			alert.Cerr("Columns start at 1")
			os.Exit(1)
		}

		cols = append(cols, col-1)
	}

	// Skip Header Lines (Optional)
	var skip int

	// Option: -s (Skip One Header Line)
	if args.IsOn("s") {
		skip = 1
	}

	// Option: -s=N (Skip N Header Lines)
	if args.HasOpt("s") {
		skip = args.GetOptI("s")
	}

	// Read Stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip Header
		if skip > 0 {
			skip--
			continue
		}

		var tokens []string

		if args.IsOn("space") {

			// Split By Whitespace
			tokens = strings.Fields(line)

		} else {

			// Split By Comma
			tokens = strings.Split(line, ",")

		}

		// Display Columns
		for i, col := range cols {
			if col < uint32(len(tokens)) {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%s", tokens[col])
			}
		}

		fmt.Println()
	}
}
