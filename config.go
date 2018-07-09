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
	"io"
	"bufio"
	"strings"
	"unicode"
	"errors"
	"fmt"
	"strconv"
	"math"
	"time"
	"net/url"
	"path/filepath"
	"net"
)

var KeyNotFoundError = errors.New("key not found")
var ValueParseError = errors.New("failed to parse value into given type")

type Config struct {
	m map[string]string
}

func (c *Config) Set(key, val string) {
	c.m[key] = val
}

func (c *Config) String(key string) (val string, err error) {
	if val, ok := c.m[key]; !ok {
		return "", KeyNotFoundError
	} else {
		return val, nil
	}
}

func (c *Config) StringOrDefault(key string, def string) (val string, used bool) {
	if val, err := c.String(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Bool(key string) (val bool, err error) {
	if val, err := c.String(key); err != nil {
		return false, err
	} else {
		return strconv.ParseBool(val)
	}
}

func (c *Config) BoolOrDefault(key string, def bool) (val bool, used bool) {
	if val, err := c.Bool(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Float32(key string) (val float32, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else if i64, err := strconv.ParseFloat(str, 32); err != nil {
		return 0, err
	} else {
		return float32(i64), nil
	}
}

func (c *Config) Float32OrDefault(key string, def float32) (val float32, used bool) {
	if val, err := c.Float32(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Float64(key string) (val float64, err error) {
	if val, err := c.String(key); err != nil {
		return math.NaN(), err
	} else {
		return strconv.ParseFloat(val, 64)
	}
}

func (c *Config) Float64OrDefault(key string, def float64) (val float64, used bool) {
	if val, err := c.Float64(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Int(key string) (val int, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else if i64, err := strconv.ParseInt(str, 10, 0); err != nil {
		return 0, err
	} else {
		return int(i64), nil
	}
}

func (c *Config) IntOrDefault(key string, def int) (val int, used bool) {
	if val, err := c.Int(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Int32(key string) (val int32, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else if i64, err := strconv.ParseInt(str, 10, 32); err != nil {
		return 0, err
	} else {
		return int32(i64), nil
	}
}

func (c *Config) Int32OrDefault(key string, def int32) (val int32, used bool) {
	if val, err := c.Int32(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Int64(key string) (val int64, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(str, 10, 64)
	}
}

func (c *Config) Int64OrDefault(key string, def int64) (val int64, used bool) {
	if val, err := c.Int64(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Uint(key string) (val uint, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else if u64, err := strconv.ParseUint(str, 10, 0); err != nil {
		return 0, err
	} else {
		return uint(u64), nil
	}
}

func (c *Config) UintOrDefault(key string, def uint) (val uint, used bool) {
	if val, err := c.Uint(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Uint32(key string) (val uint32, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else if u64, err := strconv.ParseUint(str, 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(u64), nil
	}
}

func (c *Config) Uint32OrDefault(key string, def uint32) (val uint32, used bool) {
	if val, err := c.Uint32(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Uint64(key string) (val uint64, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else {
		return strconv.ParseUint(str, 10, 64)
	}
}

func (c *Config) Uint64OrDefault(key string, def uint64) (val uint64, used bool) {
	if val, err := c.Uint64(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) Duration(key string) (val time.Duration, err error) {
	if str, err := c.String(key); err != nil {
		return 0, err
	} else {
		return time.ParseDuration(str)
	}
}

func (c *Config) DurationOrDefault(key string, def time.Duration) (val time.Duration, used bool) {
	if val, err := c.Duration(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) URL(key string) (val *url.URL, err error) {
	if str, err := c.String(key); err != nil {
		return nil, err
	} else {
		return url.Parse(str)
	}
}

func (c *Config) URLOrDefault(key string, def *url.URL) (val *url.URL, used bool) {
	if val, err := c.URL(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) FilePath(key string) (val string, err error) {
	if str, err := c.String(key); err != nil {
		return "", err
	} else {
		return filepath.Clean(str), nil
	}
}

func (c *Config) FilePathOrDefault(key string, def string) (val string, used bool) {
	if val, err := c.FilePath(key); err != nil {
		return def, true
	} else {
		return val, false
	}
}

func (c *Config) IP(key string) (val net.IP, err error) {
	if str, err := c.String(key); err != nil {
		return nil, err
	} else if val = net.ParseIP(str); val == nil {
		return nil, ValueParseError
	} else {
		return val, nil
	}
}

func (c *Config) IPOrDefault(key string, def net.IP) (val net.IP, used bool) {
	if val, err := c.IP(key); err != nil {
		return def, true
	} else {
		return val, false
	}
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

func Read(r io.Reader) (map[string]*Config, error) {
	var m = make(map[string]*Config)
	var cfg = &Config{
		m: make(map[string]string),
	}
	m[""] = cfg

	var buf = bufio.NewReader(r)
	var lnum uint = 0
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
			return nil, errors.New(fmt.Sprintf("Unrecognized input at line %d: %s", lnum, line))
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return m, nil
}