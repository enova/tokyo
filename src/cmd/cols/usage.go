package main

func usage() string {
	return `
cols
----

-h      Help
-space  Split using whitespace
-s      Skip first line
-s=N    Skip first N lines


Examples:

  # Extract columns 2 and 13
  cat junk.csv | cols 2 13

  # Skip the first line
  cat junk.csv | cols 2 13 -s

  # Skip the first two lines
  cat junk.csv | cols 2 13 -s=2

  # Split using whitespace
  cat junk.txt | cols 2 13 -s=2 -space

`
}
