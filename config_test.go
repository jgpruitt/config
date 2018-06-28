package config

import (
	"testing"
	"strings"
	"fmt"
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
		t.Run(test.in, func(t *testing.T) {
			if isKeyValue(test.in) != test.out {
				t.Errorf(`Expected %v for input %#v but got %v`, test.out, test.in, !test.out)
			}
		})
	}
}

func TestParseKeyValue(t *testing.T) {
	var tests = []struct{
		in  string
		key string
		val string
	}{
		{"a=b", "a", "b"},
		{"foo=bar", "foo", "bar"},
		{"123=456", "123", "456"},
		{"@#$!%^=a b c d e", "@#$!%^", "a b c d e"},
		{"true = false", "true", "false"},
		{"x = y", "x", "y"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			var key, val = parseKeyValue(test.in)
			if key != test.key || val != test.val {
				t.Errorf(`Expected key=%#v and value=%#v but got key=%#v and value=%#v`, test.key, test.val, key, val)
			}
		})
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
		t.Run(fmt.Sprintf("%#v=%t", test.in, test.out), func(t *testing.T) {
			if isName(test.in) != test.out {
				t.Errorf(`Expected %v for input %q but got %v`, test.out, test.in, !test.out)
			}
		})
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
		t.Run(fmt.Sprintf("%#v=%#v", test.in, test.out), func(t *testing.T) {
			var out = parseName(test.in)
			if out != test.out {
				t.Errorf(`Expected %q for input %q but got %q`, test.out, test.in, out)
			}
		})
	}
}

func TestRead(t *testing.T) {
	input := `
# this is a comment

alpha=beta
zeta=gamma
bravo=charlie
	# this is also a comment

foo:
	123=456
	cat = dog
	bird = worm
tooth= nail

bar :
	a=b
	c= d
			   f=g
	baz:
1+2=3

`
	cfgs, err := Read(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	t.Run("len(cfg)=4", func(t *testing.T) {
		if len(cfgs) != 4 {
			t.Errorf("Expected 4 Configs but got %d", len(cfgs))
		}
	})
	cfg, prs := cfgs[""]
	t.Run(`cfgs[""]!=nil`, func(t *testing.T) {
		if !prs {
			t.Error("Missing the default config")
		}
	})
	if prs {
		t.Run(`alpha=beta`, func(t *testing.T) {
			if val, prs := cfg.m["alpha"]; !prs {
				t.Error(`Value for "alpha" was missing`)
			} else if val != "beta" {
				t.Errorf(`Expected "beta" but got %q`, val)
			}
		})
	}

}
