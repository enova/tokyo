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
With configuration you must use the config package [cfg](github.com/tokyo/src/cfg). In your config file, include an `Alert` section:

```
test.cfg
--------

Alert.Sentry.Use   true
Alert.LogFile.Use  true
Alert.LogFile.Dir  log
```

The `Use` lines are optional. If they don't exist then `alert` will assume the value is false and that functionality will not be activated. If any of the `Use` lines are true, `alert` will complain if the config file does not contain required fields for that subsection. So for example, if `Alert.LogFile.Use` is set to `true`, then the config file must also have a line for `Alert.LogFile.Dir`.

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
(will finish later)
# Activating Log-File
(will finish later)
