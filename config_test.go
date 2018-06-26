package config

import (
	"testing"
	"strings"
)

func TestIsComment(t *testing.T) {
	var tests = []struct{
		in string
		out bool
	}{
		{"", false},
		{"#", true},
		{"#this is a comment", true},
		{"# this is a comment", true},
		{"nope", false},
	}

	for _, test := range tests {
		if isComment(test.in) != test.out {
			t.Errorf(`Expected %v for input %q but got %v`, test.out, test.in, !test.out)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	var tests = []struct{
		in string
		out bool
	}{
		{"", true },
		{"ab", false},
	}

	for _, test := range tests {
		if isEmpty(test.in) != test.out {
			t.Errorf(`Expected %v for input %q but got %v`, test.out, test.in, !test.out)
		}
	}
}

func TestIsKeyValue(t *testing.T) {
	var tests = []struct{
		in string
		out bool
	}{
		{"", false},
		{"a=b", true},
		{"=", false},
		{"a=", false},
		{"=b", false},
		{"key=val", true},
	}

	for _, test := range tests {
		if isKeyValue(test.in) != test.out {
			t.Errorf(`Expected %v for input %q but got %v`, test.out, test.in, !test.out)
		}
	}
}

func TestParseKeyValue(t *testing.T) {
	var tests = []struct{
		key string
		val string
	}{
		{"a", "b"},
		{"foo", "bar"},
		{"123", "456"},
		{"@#$!%^", "a b c d e"},
		{"true", "false"},
	}

	for _, test := range tests {
		var key, val = parseKeyValue(test.key + "=" + test.val)
		if key != test.key || val != test.val {
			t.Errorf(`Expected key=%q and value=%q but got key=%q and value=%q`, test.key, test.val, key, val)
		}
	}
	                                   
	if key, val := parseKeyValue("a = b"); key != "a" || val != "b" {
		t.Errorf(`Expected key="a" and value="b" but got key=%q and value=%q`, key, val)
	}
}

func TestIsName(t *testing.T) {
	var tests = []struct{
		in string
		out bool
	}{
		{"", false},
		{"a=b", false},
		{"foo:", true},
		{"bar :", true},
		{":baz", false},
		{": buz :", true},
	}

	for _, test := range tests {
		if isName(test.in) != test.out {
			t.Errorf(`Expected %v for input %q but got %v`, test.out, test.in, !test.out)
		}
	}
}

func TestParseName(t *testing.T) {
	var tests = []struct{
		in string
		out string
	}{
		{"foo:", "foo"},
		{"bar :", "bar"},
		{"baz buz:", "baz buz"},
		{"123 456 :", "123 456"},
	}

	for _, test := range tests {
		var out = parseName(test.in)
		if out != test.out {
			t.Errorf(`Expected %q for input %q but got %q`, test.out, test.in, out)
		}
	}
}

func TestRead(t *testing.T) {
	input := `
# this is a comment

alpha=beta
zeta=gamma
bravo=charlie
	# this is also a comment

`
	_, err := Read(strings.NewReader(input))
	if err != nil {
		t.Fail()
		return
	}
}
