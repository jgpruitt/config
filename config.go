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
)

type Config struct {
	m map[string]string
}

func Read(r io.Reader) (*Config, error) {
	m := make(map[string]string)
	buf := bufio.NewReader(r)
	scope := ""
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if strings.HasPrefix(strings.TrimLeftFunc(line, unicode.IsSpace),"#") {
			// line is a comment
		} else if strings.ContainsRune(line, '=') {
			// key+value
			strs := strings.SplitN(line, "=", 2)
			key, val := strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1])
			m[key] = val
		} else if strings.TrimSpace(line) == "" {
			// empty line
		} else {
			// scope
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return &Config{m: m}, nil
}