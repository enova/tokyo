# tokyo
A collection of general-purpose Go packages and commands

- alert
- args
- cfg
- cols
- dbl
- details
- lax
- alert
- set
- spawn
- stopwatch

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
