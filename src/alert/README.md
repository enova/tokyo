# alert
A package to send alert messages

# Usage
With no configuration:
```
package main

import (
  "github.com/tokyo/src/alert"
)

func main() {
  alert.Cerr("Write to stderr")
  alert.Info("Information message")
  alert.Warn("Warning message")
  alert.Exit("Major issue, exiting")
}
```
With configuration you must use the config package [cfg](github.com/tokyo/src/cfg).
In your config file, include an `Alert` section:

```
test.cfg
--------

Alert.Sentry.Use   true
Alert.LogFile.Use  true
Alert.LogFile.Dir  log
```

The `Use` lines are optional. If they don't exist then `alert` will assume the value is false
and that functionality will not be activated. If any of the `Use` lines are true, `alert` will
complain if the config file does not contain required fields for that subsection. So for example,
if `Alert.LogFile.Use` is set to `true`, then the config file must also have a line for `Alert.LogFile.Dir`.

In your application you must call `alert.Set()` to configure alerts:

```
package main

import (
  "github.com/tokyo/src/alert"
  "github.com/tokyo/src/cfg"
)

func main() {

  // Create Configuration
  cfg := cfg.New("test.cfg")

  // Configure Alerts
  alert.Set(cfg)
  
  alert.Cerr("Write to stderr")
  alert.Info("Information message")
  alert.Warn("Warning message")
  alert.Exit("Major issue, exiting")
}
```

# Activating Sentry

To activate Sentry, add the following lines to your config file:
```
Alert.Sentry.Use  True
Alert.Sentry.DSN  https://abc123...
```
Any calls to `alert.Info()`, `alert.Warn()` and `alert.Exit()` will then send your message to the specified Sentry DSN.
The following tags will be sent in the packet:
```
user - user executing the application
app  - application name
cmd  - the full command line being executed
pid  - the process PID
```
You can add additional tags within the config file:
```
Alert.Sentry.Tag  color blue
Alert.Sentry.Tag  blood-type 0
Alert.Sentry.Tag  cities Tokyo Osaka Kyoto
```
Each tag value must contain at least two tokens (key value). Additional values on the same
line will be joined into a single value (e.g. above "cities" => "Tokyo Osaka Kyoto"). These
tags will be added to every Sentry message that the application emits.

You can also add tags within individual `alert` calls:

```
alert.Info("Some alert message", "apple", "banana", "pear")
```
Here the three additional arguments will be interpreted as tags in the Sentry message:
```
"apple" => "true"
"banana" => "true"
"pear" => "true"
```
The special tag `"skip_sentry"` is used to inhibit emitting a Sentry packet:
```
for i := 0; i < 100000; i++ {

  // Message
  msg := fmt.Sprinf("The value of I is %d", i)

  // Send Alert (Skip Sentry)
  alert.Info(msg, "skip_sentry")
}
```

# Activating Log-File

To activate logging to a file, add the following lines to your config file:
```
Alert.LogFile.Use True
Alert.LogFile.Dir /var/log/
```
If the directory does not exist, it will be created.
Any calls to `alert.Info()`, `alert.Warn()` and `alert.Exit()` will then write messages
to a file in that directory. The filename will be generated using the application
name, the current date and time, and the process PID. For example:

```
/var/log/myapp_20151022_151843_000_46747.log
```
