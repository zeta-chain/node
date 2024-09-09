package main

import (
	"strings"
	"testing"
)

func TestIPAlteration(t *testing.T) {
	// ARRANGE
	const data = `{
	  "@type": "config.global",
	  "dht": { ... },
	  "liteservers":[ {"id":{ ... }, "port": 4443, "ip": 2130706433} ],
	  "validator": { ... }
	}`

	// https://www.browserling.com/tools/ip-to-dec
	const (
		ip       = "142.250.186.78"
		expected = "2398796366"
	)

	// ACT
	out, err := alterConfigIP([]byte(data), ip)

	// ASSERT
	if err != nil {
		t.Errorf("Failed: %s", err)
	}

	t.Logf("Output: %s", out)

	if !strings.Contains(string(out), expected) {
		t.Errorf("Unexpected outcome")
	}
}
