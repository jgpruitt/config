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
)

type Config struct {
	m map[string]string
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