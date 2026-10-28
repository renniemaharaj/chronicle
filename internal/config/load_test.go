package config

import (
	"os"
	"path/filepath"
	"testing"

	"chronicle/internal/domain"
)

func TestDeploymentSettingsOverrideSavedResearch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := []byte(`{"address":"127.0.0.1:9000","research":{"baseUrl":"http://model:11434","model":"file-model"}}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CHRONICLE_CONFIG", path)
	t.Setenv("CHRONICLE_LLM_MODEL", "environment-model")
	t.Setenv("CHRONICLE_LLM_API_KEY", "deployment-secret")
	t.Setenv("CHRONICLE_MAX_REQUESTS", "72")

	settings, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	saved := domain.Settings{
		BaseURL:    "http://old-model:11434",
		Model:      "old-model",
		SearchTool: "saved-search",
		APIKey:     "old-secret",
	}
	effective := settings.ApplyResearch(saved)
	if settings.Address != "127.0.0.1:9000" || settings.Research.MaxRequests != 72 {
		t.Fatal("file and environment settings must override defaults")
	}

	if effective.Model != "environment-model" || effective.BaseURL != "http://model:11434" {
		t.Fatal("deployment settings must override saved settings")
	}

	if effective.SearchTool != "saved-search" || effective.APIKey != "deployment-secret" || !settings.ManagesAPIKey() {
		t.Fatal("unconfigured saved fields must survive, while deployment credentials take precedence")
	}
}

func TestExplicitEmptyEnvironmentClearsConfiguredValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"databaseUrl":"postgres://db/workspace","research":{"apiKey":"file-secret"}}`), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CHRONICLE_CONFIG", path)
	t.Setenv("DATABASE_URL", "")
	t.Setenv("CHRONICLE_LLM_API_KEY", "")

	settings, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if settings.DatabaseURL != "" || settings.Research.APIKey != "" || !settings.ManagesAPIKey() {
		t.Fatal("an explicitly empty environment value must clear the configured value")
	}
}
