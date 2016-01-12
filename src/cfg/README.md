# cfg [![Build Status](https://travis.enova.com/go/golios.svg?token=MxEEmacAvkpaxD6dqa3m&branch=master)](https://travis.enova.com/go/cfg)
A package to parse configuration files

Why
===

Yeah yeah, there's `YAML` and `JSON`. They're great if you want to read straight into a `struct`. But when you want optional parameters, Go code starts to look clunky. You need to rely on `[]Interface{}([]byte)` a little too much. This package allows you to _query_ for parameters and uses a format similar to the one found in Java's `Properties`.

Usage
=====

A `cfg` object is built from a file:
```
import (
  "github.com/tokyo/src/cfg"
)

cfg := cfg.New("file.txt")
```

Config files contain key-value pairs, each on a single line. Please refer to the following config contents when reading the examples below: 

```
LogFile log/prod.log
LogLevel Most

connection.dev.user techops
connection.dev.host 10.1.1.1
connection.dev.port 2345

connection.prod.user techops
connection.prod.host 10.1.1.1
connection.prod.port 2345

# Comments start with a hash!
email admin@firm.com
email staff@firm.com
email desks@firm.com
```


Keys
----

Keys are contiguous strings that can be subdivided into segments using dot (`"."`) as a delimiter. A key may NOT contain whitespace characters. If a key contains dots then `cfg` can interpret it as a concatenation of _sub-keys_.

Accessing complete keys:

```
cfg.Has("LogFile")             // True
cfg.Has("LogFileX")            // False
cfg.Has("connection.dev.user") // True
cfg.Has("connection.dev")      // False - Need complete key

cfg.Get("LogFile") // "log/prod.log"
cfg.Get("Junk")    // Dies with Exit(1)
cfg.Get("email")   // Dies with Exit(1) - duplicate keys, use GetN() - see below
```

Notice that the key `email` occurs multiple times. In this case calling `Get("email")` will error out. You must use the method `GetN()` to access duplicate keys (see section Duplicate Keys below).

Sub-Keys
--------

You can segment keys out into sub-keys by using the dot character (`"."`) as a delimiter. The `Has` and `Get` methods can take multiple strings as arguments. They will join the supplied strings with dots and treat the result as a complete key:

```
cfg.Has("connection", "dev", "user") // True - Equivalent to cfg.Has("connection.dev.user")
cfg.Has("connection", "dev")         // False - Need complete key

cfg.Has("connection", "dev", "user") // "techops" - Equivalent to cfg.Get("connection.dev.user")
```

The `cfg` package allows you to discover sub-keys for a given prefix:

```
cfg.SubKeys("connection")        // ["dev", "user"]          - All sub-keys for the prefix "connection"
cfg.SubKeys("connection", "dev") // ["user", "host", "port"] - All sub-keys for the prefix "connection.dev"
cfg.SubKeys("connection.dev")    // ["user", "host", "port"] - Same
```

To check if a particular sub-key exists for a given prefix, use `HasSubKey`. The method `HasSubKey` needs at least two string arguments. The final argument is the sub-key sought. The preceeding arguments comprise the prefix:

```
cfg.HasSubKey("connection", "dev")         // True  - The prefix "connection" has a sub-key named "dev"
cfg.HasSubKey("connection", "stg")         // False - The prefix "connection" has no sub-key named "stg"
cfg.HasSubKey("connection", "dev", "user") // True  - The prefix "connection.dev" has a sub-key named "user"
```

If you want to check if a particular prefix exists, use `HasPrefix()`:

```
cfg.HasPrefix("connection")                // True  - The prefix "connection" exists
cfg.HasPrefix("connection", "dev")         // True  - The prefix "connection.dev" exists
cfg.HasPrefix("connection", "stg")         // False - The prefix "connection.stg" does not exist
cfg.HasPrefix("connection", "dev", "user") // False - The prefix "connection.dev.user" is NOT a prefix, it is a complete key
```

Duplicate Keys
--------------
The `cfg` package supports duplicate keys. The methods `Size` and `GetN` can be used to iterate over values.

```
cfg.Size("email")    // 3 - The key "email" occurs three times

cfg.GetN(0, "email") // "admin@firm.com"
cfg.GetN(1, "email") // "staff@firm.com"
cfg.GetN(2, "email") // "desks@firm.com"
```

Sub-Configs
-----------
You can extract configurations grouped under a common prefix by using the `Descend` method:

```
cfg := cfg.New("file.txt")
dev := cfg.Descend("connection", "dev")
```

Then the contents of `dev` is:

```
user techops
host 10.1.1.1
port 2345
```

The `Descend` method returns a newly created `Config` instance. Its keys are all the keys of the original configuration that matched the supplied prefix. However, the prefix is stripped in the new configuration.
In the example above, the prefix `"connection.dev"` was matched by three entries. Upon removing the prefix from those entries, the resulting keys are `user`, `host` and `port`. If an unrecognized prefix is passed
to the `Descend` method, it will return a newly created `Config` instance with no entries. 
