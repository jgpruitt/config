# config

[![GoDoc](https://godoc.org/github.com/jgpruitt/config?status.svg)](https://godoc.org/github.com/jgpruitt/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/jgpruitt/config)](https://goreportcard.com/report/github.com/jgpruitt/config)

A Go library for reading configuration files with simple, flexible syntax.

## Installation

```sh
go get -u github.com/jgpruitt/config
```

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
	* time of day (as `hour, minute int`)
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

## Example

```go
func ExampleRead() {
	var file = `

		number = 1234
		every = 3m20s

		database:
			username = admin
			port=5432

		log:
			path=../out/log.txt
			level=fatal
	`
	cfgs, _ := config.Read(strings.NewReader(file))

	// the "default" config contains key/values occurring before
	// the first named config appears
	def := cfgs[""]

	number, _ := def.IntOrDefault("number", 42)
	fmt.Println("number =", number)

	every, _ := def.DurationOrDefault("every", time.Minute * 9)
	fmt.Println("every =", every)

	db := cfgs["database"]

	username, _ := db.StringOrDefault("username", "not-admin")
	fmt.Println("username =", username)

	port, _ := db.IntOrDefault("port", 8086)
	fmt.Println("port =", port)

	// easily use a default in the case of a missing key/value pair
	ip, _ := db.IPOrDefault("ip", net.ParseIP("127.0.0.1"))
	fmt.Println("ip =", ip)

	log := cfgs["log"]

	path, _ := log.FilePathOrDefault("path", "./log.out")
	fmt.Println("path =", path)

	level, _ := log.StringOrDefault("level", "debug")
	fmt.Println("level =", level)
}
```
