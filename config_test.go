package config

import (
	"testing"
	"strings"
	"fmt"
	"time"
	"net/url"
	"net"
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
		t.Run(test.in, func(t *testing.T) {
			if isComment(test.in) != test.out {
				t.Errorf(`Expected %v for input %#v but got %v`, test.out, test.in, !test.out)
			}
		})
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
		t.Run(test.in, func(t *testing.T) {
			if isEmpty(test.in) != test.out {
				t.Errorf(`Expected %v for input %#v but got %v`, test.out, test.in, !test.out)
			}
		})
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
				t.Errorf(`Expected %v for input %#v but got %v`, test.out, test.in, !test.out)
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
				t.Errorf(`Expected %#v for input %#v but got %#v`, test.out, test.in, out)
			}
		})
	}
}

func TestRead(t *testing.T) {
	input := `
# these go into the "" config

bool=true
int64=1234567890

	str=hello world

duration=16h12m

# this is also a comment

foo:
	uint64=1234
	float64 = 1.234

bar :
url = http://jgpruitt.com

	ip=127.0.0.1

`
	// test read
	cfgs, err := Read(strings.NewReader(input))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	t.Run("len(cfg)=3", func(t *testing.T) {
		if len(cfgs) != 3 {
			t.Errorf("expected 3 Configs but got %d", len(cfgs))
		}
	})

	// test the default "" config
	cfg, prs := cfgs[""]
	t.Run(`cfgs[""]!=nil`, func(t *testing.T) {
		if !prs {
			t.Error("missing the default config")
		}
	})

	if prs {

		// test bool methods
		t.Run(`bool=true`, func(t *testing.T) {
			if val, err := cfg.Bool("bool"); err != nil {
				t.Error(err)
			} else if !val {
				t.Error("expected true but got false")
			}
		})
		t.Run(`boolx`, func(t *testing.T) {
			if val, used := cfg.BoolOrDefault("boolx", true); !used {
				t.Error("expected to use default")
			} else if !val {
				t.Error("expected true but got false")
			}
		})

		// test int64 methods
		t.Run("int64=1234567890", func(t *testing.T) {
			if val, err := cfg.Int64("int64"); err != nil {
				t.Error(err)
			} else if val != 1234567890 {
				t.Errorf("expected 1234567890 but got %d", val)
			}
		})
		t.Run("int64x", func(t *testing.T) {
			if val, used := cfg.Int64OrDefault("int64x", 1234567890); !used {
				t.Error("expected to use default")
			} else if val != 1234567890 {
				t.Errorf("expected 1234567890 but got %d", val)
			}
		})

		// test string methods
		t.Run("str=hello world", func(t *testing.T) {
			if val, err := cfg.String("str"); err != nil {
				t.Error(err)
			} else if val != "hello world" {
				t.Errorf(`expected "hello world" but got %#v`, val)
			}
		})
		t.Run("strx", func(t *testing.T) {
			if val, used := cfg.StringOrDefault("strx", "hello world"); !used {
				t.Error("expected to use default")
			} else if val != "hello world" {
				t.Errorf(`expected "hello world" but got %#v`, val)
			}
		})

		// test duration methods
		t.Run("duration=16h12m", func(t *testing.T) {
			exp, _ := time.ParseDuration("16h12m")
			if val, err := cfg.Duration("duration"); err != nil {
				t.Error(err)
			} else if val != exp {
				t.Errorf(`expected %#v but got %#v`, exp, val)
			}
		})
		t.Run("durationx", func(t *testing.T) {
			exp, _ := time.ParseDuration("16h12m")
			if val, used := cfg.DurationOrDefault("durationx", exp); !used {
				t.Error("expected to use default")
			} else if val != exp {
				t.Errorf(`expected %#v but got %#v`, exp, val)
			}
		})
	}

	// test "foo" config
	cfg, prs = cfgs["foo"]
	t.Run(`cfgs["foo"]!=nil`, func(t *testing.T) {
		if !prs {
			t.Error("missing the foo config")
		}
	})
	if prs {

		// test uint64 methods
		t.Run(`uint64=1234`, func(t *testing.T) {
			if val, err := cfg.Uint64("uint64"); err != nil {
				t.Error(err)
			} else if val != 1234 {
				t.Errorf("expected 1234 but got %d", val)
			}
		})
		t.Run(`uint64x`, func(t *testing.T) {
			if val, used := cfg.Uint64OrDefault("uint64x", 1234); !used {
				t.Error("expected to use default")
			} else if val != 1234 {
				t.Errorf("expected 1234 but got %d", val)
			}
		})

		// test float64 methods
		t.Run(`float64=1.234`, func(t *testing.T) {
			if val, err := cfg.Float64("float64"); err != nil {
				t.Error(err)
			} else if val != 1.234 {
				t.Errorf("expected 1.234 but got %f", val)
			}
		})
		t.Run(`float64x`, func(t *testing.T) {
			if val, used := cfg.Float64OrDefault("float64x", 1.234); !used {
				t.Error("expected to use default")
			} else if val != 1.234 {
				t.Errorf("expected 1.234 but got %f", val)
			}
		})
	}

	// test "bar" config
	cfg, prs = cfgs["bar"]
	t.Run(`cfgs["bar"]!=nil`, func(t *testing.T) {
		if !prs {
			t.Error("missing the bar config")
		}
	})
	if prs {

		// test URL methods
		t.Run(`url=http://jgpruitt.com`, func(t *testing.T) {
			exp,_ := url.Parse(`http://jgpruitt.com`)
			if val, err := cfg.URL("url"); err != nil {
				t.Error(err)
			} else if val.String() != exp.String() {
				t.Errorf("expected %s but got %s", exp, val)
			}
		})
		t.Run(`urlx`, func(t *testing.T) {
			exp,_ := url.Parse(`http://jgpruitt.com`)
			if val, used := cfg.URLOrDefault("urlx", exp); !used {
				t.Error("expected to use default")
			} else if val.String() != exp.String() {
				t.Errorf("expected %s but got %s", exp, val)
			}
		})

		// test IP methods
		t.Run(`ip=127.0.0.1`, func(t *testing.T) {
			exp := net.ParseIP("127.0.0.1")
			if val, err := cfg.IP("ip"); err != nil {
				t.Error(err)
			} else if val.String() != exp.String() {
				t.Errorf("expected %s but got %s", exp, val)
			}
		})
		t.Run(`ipx`, func(t *testing.T) {
			exp := net.ParseIP("127.0.0.1")
			if val, used := cfg.IPOrDefault("ipx", exp); !used {
				t.Error("expected to use default")
			} else if val.String() != exp.String() {
				t.Errorf("expected %s but got %s", exp, val)
			}
		})
	}
}