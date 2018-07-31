# config
[![Go Report Card](https://goreportcard.com/badge/github.com/jgpruitt/config)](https://goreportcard.com/report/github.com/jgpruitt/config)

A Go library for reading configuration files.

## Features

* Lightweight syntax
* Supports multiple named configurations in a single file
* Configurations consist of key/value pairs
* Many data types supported for values including:
    * string
    * bool
    * float32
    * float64
    * int
    * int32
    * int64
    * uint
    * uint32
    * uint64
    * time.Duration
    * *url.URL
    * file path
    * net.IP
* Easily substitute defaults for missing keys or incorrectly specified values
* Heavily unit tested

## Syntax

A line where the first non-whitespace character is a ``#`` is a comment and is ignored. 

```# this is a comment```

A line consisting of characters followed by an ``=`` followed by more characters is a key/value pair.

```
key=value
  this_is_a_key = this_is_a_value
```

A line consisting of characters followed by a ``:`` starts a new named configuration.

```
# the first two key/value pairs go in a default configuration named "" (empty string)

key1 = 1234
key2 = 5678

# the next line starts a new configuration named "database"
database:
    username=admin
    password=default

```

Whitespace around configuration names, keys, and values is ignored. So, uses spaces and tabs to your heart's content
to make your configuration more readable.
