package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestSourceOverride(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		valid       bool
	}{
		{"default", `{}`, true},
		{"override", `{"sourceOverride":{"provider":"reparser","baseUrl":"https://example.test"}}`, true},
		{"unknown", `{"sourceOverride":{"provider":"unknown","baseUrl":"https://example.test"}}`, false},
		{"missing URL", `{"sourceOverride":{"provider":"reparser"}}`, false},
		{"http", `{"sourceOverride":{"provider":"reparser","baseUrl":"http://example.test"}}`, false},
		{"credentials", `{"sourceOverride":{"provider":"reparser","baseUrl":"https://user:secret@example.test"}}`, false},
		{"query", `{"sourceOverride":{"provider":"reparser","baseUrl":"https://example.test?x=1"}}`, false},
		{"path", `{"sourceOverride":{"provider":"reparser","baseUrl":"https://example.test/path"}}`, false},
		{"short cache", `{"sourceOverride":{"provider":"reparser","baseUrl":"https://example.test","cacheMinutes":1}}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var u UniversityConfig
			if err := json.Unmarshal([]byte(tc.input), &u); err != nil {
				t.Fatal(err)
			}
			cfg := Config{Universities: []UniversityConfig{u}}
			if err := loadSourceSettings(&cfg); (err == nil) != tc.valid {
				t.Fatalf("unexpected validation result: %v", err)
			}
			if tc.valid {
				want := 12 * time.Hour
				if u.SourceOverride != nil {
					want = 30 * time.Minute
				}
				if got := cfg.Universities[0].CacheTTL(); got != want {
					t.Fatalf("TTL = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestRepositoryConfiguration(t *testing.T) {
	data, err := os.ReadFile("../../config.json")
	if err != nil {
		t.Fatal(err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if err := validateConfig(&cfg); err != nil {
		t.Fatal(err)
	}
	if err := loadSourceSettings(&cfg); err != nil {
		t.Fatal(err)
	}
}
