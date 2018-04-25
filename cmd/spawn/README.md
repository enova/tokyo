# Spawn
A command-line utility to spawn process N at a time.
# Usage
Suppose the file `commands.txt` contains a list of 20 shell commands, one per line. The following will execute all 20 commands 5 at a time:
```
cat commands.txt | spawn 5
```
It will start by launching the first 5 commands in the file. Each time a command completes, a new command will be launched.
If a command returns with non-zero exit-code, a warning message will be displayed on the console. This will not however prevent other commands from being launched.
