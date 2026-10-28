package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"chronicle/internal/config"
	"chronicle/internal/research"
	"chronicle/internal/storage"
)

func TestDeploymentConfigurationDoesNotExposeOrPersistSecrets(t *testing.T) {
	t.Setenv("CHRONICLE_CONFIG", "")
	t.Setenv("CHRONICLE_LLM_API_KEY", "deployment-secret")
	t.Setenv("CHRONICLE_LLM_MODEL", "deployment-model")

	settings, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	settings.DatabaseURL = "postgres://private-database"
	settings.MapTileURL = "https://maps.example/{z}/{x}/{y}.png"
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	handler := New(store, research.NewManager(store), settings)

	for _, path := range []string{"/api/config", "/api/settings"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != 200 || strings.Contains(response.Body.String(), "deployment-secret") || strings.Contains(response.Body.String(), "private-database") {
			t.Fatal("browser configuration must not expose deployment secrets")
		}

		if path == "/api/config" && !strings.Contains(response.Body.String(), "https://maps.example/") {
			t.Fatal("runtime map settings must come from deployment configuration")
		}
	}

	payload, err := json.Marshal(settings.Research)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPut, "/api/settings", strings.NewReader(string(payload))))
	if response.Code != 200 || strings.Contains(response.Body.String(), "deployment-secret") {
		t.Fatal("settings writes must return a redacted effective configuration")
	}

	stored, err := store.Get(context.Background(), "settings", "providers")
	if err != nil || strings.Contains(string(stored), "deployment-secret") {
		t.Fatal("environment credentials must not be copied into workspace storage")
	}
}

func TestConfiguredOriginsAreExact(t *testing.T) {
	settings := config.Defaults()
	settings.AllowedOrigins = []string{"https://workspace.example"}
	server := &Server{config: settings}

	for _, test := range []struct {
		origin  string
		allowed bool
	}{
		{"", true},
		{"https://workspace.example", true},
		{"https://workspace.example.attacker.test", false},
		{"https://workspace.example/path", false},
		{"null", false},
	} {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Origin", test.origin)

		origin, err := url.Parse(test.origin)
		if err != nil {
			t.Fatal(err)
		}

		if server.allowsOrigin(origin, "local.test") != test.allowed {
			t.Fatalf("unexpected origin decision for %q", test.origin)
		}
	}
}
