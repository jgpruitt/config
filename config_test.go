package config

import (
	"testing"
	"strings"
)

func TestConfig_Read(t *testing.T) {
	input := `
# this is a comment

alpha=beta
zeta=gamma
bravo=charlie
	# this is also a comment

`
	cfg, err := Read(strings.NewReader(input))
	if err != nil {
		t.Fail()
		return
	}
	if len(cfg.m) != 3 {
		t.Fail()
		return
	}
	if val, prs := cfg.m["alpha"]; !prs || val != "beta" {
		t.Fail()
		return
	}
	if val, prs := cfg.m["zeta"]; !prs || val != "gamma" {
		t.Fail()
		return
	}
	if val, prs := cfg.m["bravo"]; !prs || val != "charlie" {
		t.Fail()
		return
	}
}
