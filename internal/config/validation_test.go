package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInvalidConfigurationFailsAtStartup(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"unknown key", `{"unknownService":"value"}`},
		{"unknown research key", `{"research":{"unknown":"value"}}`},
		{"second object", `{} {}`},
		{"null document", `null`},
		{"invalid port", `{"address":"127.0.0.1:0"}`},
		{"bad endpoint", `{"geocoderUrl":"file:///tmp/map"}`},
		{"bad provider", `{"research":{"provider":"invalid"}}`},
		{"bad budget", `{"research":{"maxRequests":1}}`},
		{"origin path", `{"allowedOrigins":["https://example.org/path"]}`},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(test.input), 0600); err != nil {
				t.Fatal(err)
			}

			t.Setenv("CHRONICLE_CONFIG", path)
			if _, err := Load(); err == nil {
				t.Fatal("invalid deployment configuration must be rejected")
			}
		})
	}
}

func TestInvalidNumericEnvironmentFails(t *testing.T) {
	t.Setenv("CHRONICLE_CONFIG", "")
	t.Setenv("CHRONICLE_CONCURRENCY", "many")
	if _, err := Load(); err == nil {
		t.Fatal("invalid numeric environment values must be rejected")
	}
}
