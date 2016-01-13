# tokyo
A collection of general-purpose Go packages and commands
```
alert     - package to send alert messages from your application
args      - package to make handling command-line arguments easier
cfg       - package to parse configuration files
dbl       - package to help avoid indeterminacy when dealing with floating-point numbers
details   - package to help control logging of details in Go apps
lax       - package containing a potpourri of utility functions commonly used by some of the authors
set       - package that implements a Set container (of strings)
spawn     - package that allows users to execute multiple shell commands in parallel
stopwatch - package that implements a simple stopwatch for inline benchmarking

cols      - command-line tool to help parse CSV and tabular data
spawn     - command-line tool to spawn multiple processes in parallel
```

# Building
The standard Go build-tools should work:
```
go test ./...
go build ./...
go install ./...
```
A few convenience scripts have been included as well:

- Run `./go_get.sh` to fetch all prerequisite packages including the `golint` command
- Run `./build.sh` to run the linter, formatter and test-suite
- Run `./install.sh` to install all commands to your go-path bin

# Licensing
Tokyo is released by [Enova](http://www.enova.com) under the
[MIT License](https://github.com/enova/tokyo/blob/master/LICENSE).
