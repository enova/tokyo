# cols
Console app to process textual data stored in columns.

# Installing
First get the package:
```
go get github.com/tokyo/src/cols
```
Then enter the project directory and install:
```
cd $GOPATH/src/github.com/tokyo/src/cols
go install
```
# Usage
Output of `cols -h`:
```

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
```
