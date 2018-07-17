// MIT License
//
// Copyright (c) 2018 John Pruitt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ErrKeyNotFound returned when the Key queried does not exist in the configuration
var ErrKeyNotFound = errors.New("key not found")

// ErrParseValue returned when the value could not be parsed into the given type
var ErrParseValue = errors.New("failed to parse value into given type")

// Config is a set of key/value pairs
type Config struct {
	m map[string]string
}

// Set adds a key/value pair to the configuration.
// If the key already exists, the value will be replaced
func (c *Config) Set(key, val string) {
	c.m[key] = val
}

// String returns the value associated with the given key as a string.
// If the key does not exist, ErrKeyNotFound is returned.
func (c *Config) String(key string) (val string, err error) {
	var ok bool
	val, ok = c.m[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return
}

// StringOrDefault returns the value associated with the given key as a string.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) StringOrDefault(key string, def string) (val string, used bool) {
	var err error
	val, err = c.String(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Bool returns the value associated with the given key as a bool.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a bool.
func (c *Config) Bool(key string) (val bool, err error) {
	str, err := c.String(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(str)
}

// BoolOrDefault returns the value associated with the given key as a bool.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) BoolOrDefault(key string, def bool) (val bool, used bool) {
	var err error
	val, err = c.Bool(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Float32 returns the value associated with the given key as a float32.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a float32.
func (c *Config) Float32(key string) (val float32, err error) {
	var (
		str string
		f64 float64
	)
	str, err = c.String(key)
	if err != nil {
		return 0, err
	}
	f64, err = strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}
	return float32(f64), nil
}

// Float32OrDefault returns the value associated with the given key as a float32.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Float32OrDefault(key string, def float32) (val float32, used bool) {
	var err error
	val, err = c.Float32(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Float64 returns the value associated with the given key as a float64.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a float64.
func (c *Config) Float64(key string) (val float64, err error) {
	str, err := c.String(key)
	if err != nil {
		return math.NaN(), err
	}
	return strconv.ParseFloat(str, 64)
}

// Float64OrDefault returns the value associated with the given key as a float64.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Float64OrDefault(key string, def float64) (val float64, used bool) {
	var err error
	val, err = c.Float64(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Int returns the value associated with the given key as an int.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into an int.
func (c *Config) Int(key string) (val int, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	i64, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

// IntOrDefault returns the value associated with the given key as an int.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) IntOrDefault(key string, def int) (val int, used bool) {
	var err error
	val, err = c.Int(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Int32 returns the value associated with the given key as an int32.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into an int32.
func (c *Config) Int32(key string) (val int32, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	i64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i64), nil
}

// Int32OrDefault returns the value associated with the given key as an int32.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Int32OrDefault(key string, def int32) (val int32, used bool) {
	var err error
	val, err = c.Int32(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Int64 returns the value associated with the given key as an int64.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into an int64.
func (c *Config) Int64(key string) (val int64, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 10, 64)
}

// Int64OrDefault returns the value associated with the given key as an int64.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Int64OrDefault(key string, def int64) (val int64, used bool) {
	var err error
	val, err = c.Int64(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Uint returns the value associated with the given key as a uint.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a uint.
func (c *Config) Uint(key string) (val uint, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	u64, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}

// UintOrDefault returns the value associated with the given key as a uint.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) UintOrDefault(key string, def uint) (val uint, used bool) {
	var err error
	val, err = c.Uint(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Uint32 returns the value associated with the given key as a uint32.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a uint32.
func (c *Config) Uint32(key string) (val uint32, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(u64), nil
}

// Uint32OrDefault returns the value associated with the given key as a uint32.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Uint32OrDefault(key string, def uint32) (val uint32, used bool) {
	var err error
	val, err = c.Uint32(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Uint64 returns the value associated with the given key as a uint64.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a uint64.
func (c *Config) Uint64(key string) (val uint64, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(str, 10, 64)
}

// Uint64OrDefault returns the value associated with the given key as a uint64.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) Uint64OrDefault(key string, def uint64) (val uint64, used bool) {
	var err error
	val, err = c.Uint64(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// Duration returns the value associated with the given key as a time.Duration.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a time.Duration.
func (c *Config) Duration(key string) (val time.Duration, err error) {
	str, err := c.String(key)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(str)
}

// DurationOrDefault returns the value associated with the given key as a time.Duration.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) DurationOrDefault(key string, def time.Duration) (val time.Duration, used bool) {
	var err error
	val, err = c.Duration(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// URL returns the value associated with the given key as a *url.URL.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a *url.URL.
func (c *Config) URL(key string) (val *url.URL, err error) {
	str, err := c.String(key)
	if err != nil {
		return nil, err
	}
	return url.Parse(str)
}

// URLOrDefault returns the value associated with the given key as a *url.URL.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) URLOrDefault(key string, def *url.URL) (val *url.URL, used bool) {
	var err error
	val, err = c.URL(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// FilePath returns the value associated with the given key as a string that
// has been interpretted as a file path and cleaned.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a *url.URL.
func (c *Config) FilePath(key string) (val string, err error) {
	str, err := c.String(key)
	if err != nil {
		return "", err
	}
	return filepath.Clean(str), nil
}

// FilePathOrDefault returns the value associated with the given key as a string that
// has been interpretted as a file path and cleaned.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) FilePathOrDefault(key string, def string) (val string, used bool) {
	var err error
	val, err = c.FilePath(key)
	if err != nil {
		return def, true
	}
	return val, false
}

// IP returns the value associated with the given key as a net.IP.
// If the key does not exist, ErrKeyNotFound is returned.
// An error is returned if the value cannot be parsed into a net.IP.
func (c *Config) IP(key string) (val net.IP, err error) {
	str, err := c.String(key)
	if err != nil {
		return nil, err
	}
	val = net.ParseIP(str)
	if val == nil {
		return nil, ErrParseValue
	}
	return val, nil
}

// IPOrDefault returns the value associated with the given key as a net.IP.
// If the key does not exist or cannot be parsed appropriately, the default value "def" is returned.
// "used" will be true if the default value was used.
func (c *Config) IPOrDefault(key string, def net.IP) (val net.IP, used bool) {
	var err error
	val, err = c.IP(key)
	if err != nil {
		return def, true
	}
	return val, false
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func isEmpty(line string) bool {
	return line == ""
}

func isKeyValue(line string) bool {
	return strings.ContainsRune(line, '=') && len(line) >= 3
}

func parseKeyValue(line string) (key, value string) {
	strs := strings.SplitN(line, "=", 2)
	key, value = strings.TrimRightFunc(strs[0], unicode.IsSpace), strings.TrimLeftFunc(strs[1], unicode.IsSpace)
	return
}

func isName(line string) bool {
	return strings.HasSuffix(line, ":") && len(line) >= 2
}

func parseName(line string) string {
	return strings.TrimRightFunc(line, func(r rune) bool {
		return ':' == r || unicode.IsSpace(r)
	})
}

// Read parses one or more Configs out of the given io.Reader.
// An error is returned if there is a problem reading or
// unrecognized input.
func Read(r io.Reader) (map[string]*Config, error) {
	var m = make(map[string]*Config)
	var cfg = &Config{
		m: make(map[string]string),
	}
	m[""] = cfg

	var buf = bufio.NewReader(r)
	var lnum uint
	for {
		var line, err = buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		lnum = lnum + 1
		line = strings.TrimSpace(line)
		if isComment(line) || isEmpty(line) {
			// ignore
		} else if isKeyValue(line) {
			var key, value = parseKeyValue(line)
			cfg.m[key] = value
		} else if isName(line) {
			var name = parseName(line)
			if _, prs := m[name]; !prs {
				m[name] = &Config{
					m: make(map[string]string),
				}
			}
			cfg = m[name]
		} else {
			return nil, fmt.Errorf("unrecognized input at line %d: %s", lnum, line)
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return m, nil
}
