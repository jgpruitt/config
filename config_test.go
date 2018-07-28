package config

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestIsComment(t *testing.T) {
	var tests = []struct {
		in  string
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
	var tests = []struct {
		in  string
		out bool
	}{
		{"", true},
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
	var tests = []struct {
		in  string
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
	var tests = []struct {
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
	var tests = []struct {
		in  string
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
	var tests = []struct {
		in  string
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

	var input = `
# these go into the "" config

bool=true
!bool = 1234

int64=1234567890
!int64 = alpha

int= 1
!int = beta

int32 =  123
!int32 = ^^^^^

	str=hello world

duration=16h12m
!duration=2018-06-15

# this is also a comment

foo:
	uint64=1234
!uint64 = aaaa
uint32 = 123
!uint32 = -1
	float64 = 1.234
	!float64 = gamma
float32 = 1.2233
		!float32 =delta
	
bar :
url = http://jgpruitt.com
!url = charlie
	ip=127.0.0.1
!ip=http://jgpruitt.com
filepath=/usr/bin/env
`
	// test read
	cfgs, err := Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// find all 3 configs?
	t.Run("len(cfg)=3", func(t *testing.T) {
		if len(cfgs) != 3 {
			t.Errorf("expected 3 Configs but got %d", len(cfgs))
		}
	})

	// is the "" default config present?
	t.Run(`cfgs[""]!=nil`, func(t *testing.T) {
		if _, prs := cfgs[""]; !prs {
			t.Error("missing the default config")
		}
	})

	// is the "foo" config present?
	t.Run(`cfgs["foo"]!=nil`, func(t *testing.T) {
		if _, prs := cfgs["foo"]; !prs {
			t.Error("missing the 'foo' config")
		}
	})

	// is the "bar" config present?
	t.Run(`cfgs["bar"]!=nil`, func(t *testing.T) {
		if _, prs := cfgs["bar"]; !prs {
			t.Error("missing the 'bar' config")
		}
	})
}

func TestConfig_String(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  string
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    string1 = hello world
	string2 = 1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.String("string0")
	t.Run("string0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.String("string1")
	t.Run("string1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != "hello world" {
			t.Errorf("expected 'hello world' but got '%s'", val)
		}
	})

	val, err = cfg.String("string2")
	t.Run("string2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != "1234" {
			t.Errorf("expected '1234' but got '%s'", val)
		}
	})
}

func TestConfig_StringOrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  string
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    string1 = fizzbuzz
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.StringOrDefault("string1", "alpha")
	t.Run("string1", func(t *testing.T) {
		if val != "fizzbuzz" {
			t.Errorf("expected val='fizzbuzz' but got '%s'", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.StringOrDefault("string2", "beta")
	t.Run("string2", func(t *testing.T) {
		if val != "beta" {
			t.Errorf("expected val='beta' but got '%s'", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Bool(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    bool1 = true
	bool2 = F
	bool3 = +++
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Bool("bool0")
	t.Run("bool0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Bool("bool1")
	t.Run("bool1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != true {
			t.Error("expected true but got false")
		}
	})

	val, err = cfg.Bool("bool2")
	t.Run("bool2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != false {
			t.Error("expected false but got true")
		}
	})

	val, err = cfg.Bool("bool3")
	t.Run("bool3", func(t *testing.T) {
		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestConfig_BoolOrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  bool
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    bool1 = true
	bool2 = notabool
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.BoolOrDefault("bool1", false)
	t.Run("bool1", func(t *testing.T) {
		if !val {
			t.Error("expected val=true but got false")
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.BoolOrDefault("bool2", true)
	t.Run("bool2", func(t *testing.T) {
		if !val {
			t.Error("expected val=true but got false")
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.BoolOrDefault("bool3", true)
	t.Run("bool3", func(t *testing.T) {
		if !val {
			t.Error("expected val=true but got false")
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Float32(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  float32
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    float32_1 = 123.4
	float32_2 = -12.34
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Float32("float32_0")
	t.Run("float32_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Float32("float32_1")
	t.Run("float32_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 123.4 {
			t.Errorf("expected 123.4 but got %f", val)
		}
	})

	val, err = cfg.Float32("float32_2")
	t.Run("float32_2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != -12.34 {
			t.Errorf("expected -12.34 but got %f", val)
		}
	})
}

func TestConfig_Float32OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  float32
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    float32_1 = 123.4
	float32_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Float32OrDefault("float32_0", 99.99)
	t.Run("float32_0", func(t *testing.T) {
		if val != 99.99 {
			t.Errorf("expected val=99.99 but got %f", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Float32OrDefault("float32_1", 99.99)
	t.Run("float32_1", func(t *testing.T) {
		if val != 123.4 {
			t.Errorf("expected val=123.4 but got %f", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Float32OrDefault("float32_2", 99.99)
	t.Run("float32_2", func(t *testing.T) {
		if val != 99.99 {
			t.Errorf("expected val=99.99 but got %f", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Float64(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  float64
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    float64_1 = 123.4
	float64_2 = -12.34
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Float64("float64_0")
	t.Run("float64_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Float64("float64_1")
	t.Run("float64_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 123.4 {
			t.Errorf("expected 123.4 but got %f", val)
		}
	})

	val, err = cfg.Float64("float64_2")
	t.Run("float64_2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != -12.34 {
			t.Errorf("expected -12.34 but got %f", val)
		}
	})
}

func TestConfig_Float64OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  float64
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    float64_1 = 123.4
	float64_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Float64OrDefault("float64_0", 99.99)
	t.Run("float64_0", func(t *testing.T) {
		if val != 99.99 {
			t.Errorf("expected val=99.99 but got %f", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Float64OrDefault("float64_1", 99.99)
	t.Run("float64_1", func(t *testing.T) {
		if val != 123.4 {
			t.Errorf("expected val=123.4 but got %f", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Float64OrDefault("float64_2", 99.99)
	t.Run("float64_2", func(t *testing.T) {
		if val != 99.99 {
			t.Errorf("expected val=99.99 but got %f", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Int(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int_1 = 1234
	int_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Int("int_0")
	t.Run("int_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Int("int_1")
	t.Run("int_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Int("int_2")
	t.Run("int_2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != -1234 {
			t.Errorf("expected -1234 but got %d", val)
		}
	})
}

func TestConfig_IntOrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int_1 = 1234
	int_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.IntOrDefault("int_0", 9999)
	t.Run("int_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.IntOrDefault("int_1", 9999)
	t.Run("int_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.IntOrDefault("int_2", 9999)
	t.Run("int_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Int32(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int32
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int32_1 = 1234
	int32_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Int32("int32_0")
	t.Run("int32_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Int32("int32_1")
	t.Run("int32_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Int32("int32_2")
	t.Run("int32_2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != -1234 {
			t.Errorf("expected -1234 but got %d", val)
		}
	})
}

func TestConfig_Int32OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int32
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int32_1 = 1234
	int32_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Int32OrDefault("int32_0", 9999)
	t.Run("int32_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Int32OrDefault("int32_1", 9999)
	t.Run("int32_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Int32OrDefault("int32_2", 9999)
	t.Run("int32_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Int64(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int64
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int64_1 = 1234
	int64_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Int64("int64_0")
	t.Run("int64_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Int64("int64_1")
	t.Run("int64_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Int64("int64_2")
	t.Run("int64_2", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != -1234 {
			t.Errorf("expected -1234 but got %d", val)
		}
	})
}

func TestConfig_Int64OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  int64
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    int64_1 = 1234
	int64_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Int64OrDefault("int64_0", 9999)
	t.Run("int64_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Int64OrDefault("int64_1", 9999)
	t.Run("int64_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Int64OrDefault("int64_2", 9999)
	t.Run("int64_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Uint(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint_1 = 1234
	uint_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Uint("uint_0")
	t.Run("uint_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Uint("uint_1")
	t.Run("uint_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Uint("uint_2")
	t.Run("uint_2", func(t *testing.T) {
		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestConfig_UintOrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint_1 = 1234
	uint_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.UintOrDefault("uint_0", 9999)
	t.Run("uint_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.UintOrDefault("uint_1", 9999)
	t.Run("uint_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.UintOrDefault("uint_2", 9999)
	t.Run("uint_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Uint32(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint32
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint32_1 = 1234
	uint32_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Uint32("uint32_0")
	t.Run("uint32_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Uint32("uint32_1")
	t.Run("uint32_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Uint32("uint32_2")
	t.Run("uint32_2", func(t *testing.T) {
		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestConfig_Uint32OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint32
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint32_1 = 1234
	uint32_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Uint32OrDefault("uint32_0", 9999)
	t.Run("uint32_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Uint32OrDefault("uint32_1", 9999)
	t.Run("uint32_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Uint32OrDefault("uint32_2", 9999)
	t.Run("uint32_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Uint64(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint64
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint64_1 = 1234
	uint64_2 = -1234
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Uint64("uint64_0")
	t.Run("uint64_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Uint64("uint64_1")
	t.Run("uint64_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		if val != 1234 {
			t.Errorf("expected 1234 but got %d", val)
		}
	})

	val, err = cfg.Uint64("uint64_2")
	t.Run("uint64_2", func(t *testing.T) {
		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestConfig_Uint64OrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  uint64
		used bool
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    uint64_1 = 1234
	uint64_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.Uint64OrDefault("uint64_0", 9999)
	t.Run("uint64_0", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.Uint64OrDefault("uint64_1", 9999)
	t.Run("uint64_1", func(t *testing.T) {
		if val != 1234 {
			t.Errorf("expected val=1234 but got %d", val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.Uint64OrDefault("uint64_2", 9999)
	t.Run("uint64_2", func(t *testing.T) {
		if val != 9999 {
			t.Errorf("expected val=9999 but got %d", val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}

func TestConfig_Duration(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  time.Duration
		err  error
	)

	cfgs, err = Read(strings.NewReader(`
    duration_1 = 3h15m22s
	duration_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	_, err = cfg.Duration("duration_0")
	t.Run("duration_0", func(t *testing.T) {
		if err != ErrKeyNotFound {
			t.Error("expected 'ErrKeyNotFound'")
		}
	})

	val, err = cfg.Duration("duration_1")
	t.Run("duration_1", func(t *testing.T) {
		if err != nil {
			t.Errorf("did not expect an error: %s", err)
		}
		exp, _ := time.ParseDuration("3h15m22s")
		if val != exp {
			t.Errorf("expected %s but got %s", exp, val)
		}
	})

	val, err = cfg.Duration("duration_2")
	t.Run("duration_2", func(t *testing.T) {
		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestConfig_DurationOrDefault(t *testing.T) {
	var (
		cfgs map[string]*Config
		cfg  *Config
		val  time.Duration
		used bool
		err  error
		exp  = 55 * time.Minute
		def  = 99 * time.Second
	)

	cfgs, err = Read(strings.NewReader(`
    duration_1 = 55m
	duration_2 = gamma
	`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cfg = cfgs[""]
	if cfg == nil {
		t.Fatal("default config missing")
	}

	val, used = cfg.DurationOrDefault("duration_0", def)
	t.Run("duration_0", func(t *testing.T) {
		if val != def {
			t.Errorf("expected val=%s but got %s", def, val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})

	val, used = cfg.DurationOrDefault("duration_1", def)
	t.Run("duration_1", func(t *testing.T) {
		if val != exp {
			t.Errorf("expected val=%s but got %s", exp, val)
		}
		if used {
			t.Error("did not expect to use default")
		}
	})

	val, used = cfg.DurationOrDefault("duration_2", def)
	t.Run("duration_2", func(t *testing.T) {
		if val != def {
			t.Errorf("expected val=%s but got %s", def, val)
		}
		if !used {
			t.Error("expected to use default")
		}
	})
}
