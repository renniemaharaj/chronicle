package api

import (
	"chronicle/internal/storage"
	"context"
	"strings"
	"testing"
)

func TestSpeechKeySourcesAreRedactedAndEnvironmentWins(t *testing.T) {
	handler, store := speechTestServer(t, "environment-secret", "https://example.org")
	response := speechCall(handler, "GET", "/api/speech/settings", "")
	if response.Code != 200 || !strings.Contains(response.Body.String(), "environment/config") || strings.Contains(response.Body.String(), "environment-secret") {
		t.Fatal("the UI must detect deployment keys without receiving their value")
	}

	response = speechCall(handler, "PUT", "/api/speech/settings", `{"apiKey":"replacement"}`)
	if response.Code != 409 {
		t.Fatal("UI must not silently override a deployment-managed key")
	}

	if _, err := store.Get(context.Background(), "settings", "speech"); err != storage.ErrNotFound {
		t.Fatal("deployment keys must not be copied to workspace storage")
	}

	handler, store = speechTestServer(t, "", "https://example.org")
	response = speechCall(handler, "PUT", "/api/speech/settings", `{"apiKey":"workspace-secret"}`)
	if response.Code != 200 || strings.Contains(response.Body.String(), "workspace-secret") {
		t.Fatal("saved key responses must be redacted")
	}

	data, err := store.Get(context.Background(), "settings", "speech")
	if err != nil || !strings.Contains(string(data), "workspace-secret") {
		t.Fatal("a UI-configured key must be stored on the server")
	}

	response = speechCall(handler, "PUT", "/api/speech/settings", `{"apiKey":""}`)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"hasApiKey":false`) {
		t.Fatal("removing a workspace key must be supported")
	}
}
