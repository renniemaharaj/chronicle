package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"chronicle/internal/domain"
	"chronicle/internal/mcp"
	"chronicle/internal/storage"
)

func (s *Server) loadSettings(ctx context.Context) (domain.Settings, error) {
	settings := domain.Settings{Provider: "ollama", BaseURL: "http://localhost:11434", MaxRequests: 40, Concurrency: 2}
	data, err := s.store.Get(ctx, "settings", "providers")
	if errors.Is(err, storage.ErrNotFound) {
		return settings, nil
	}
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(data, &settings)
	return settings, err
}

func (s *Server) settings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.loadSettings(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	settings.HasAPIKey, settings.APIKey = settings.APIKey != "", ""
	write(w, 200, settings)
}

func (s *Server) saveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.Settings
	if err := decode(w, r, &settings); err != nil {
		fail(w, err, 400)
		return
	}
	if settings.Provider != "ollama" && settings.Provider != "openai" {
		fail(w, fmt.Errorf("unknown provider"), 422)
		return
	}
	if err := validateEndpoint(settings.BaseURL); err != nil {
		fail(w, err, 422)
		return
	}
	if settings.MCPURL != "" {
		if err := validateEndpoint(settings.MCPURL); err != nil {
			fail(w, err, 422)
			return
		}
	}
	if settings.MaxRequests < 4 || settings.MaxRequests > 100 || settings.Concurrency < 1 || settings.Concurrency > 4 {
		fail(w, fmt.Errorf("request budget must be 4–100 and concurrency 1–4"), 422)
		return
	}
	previous, err := s.loadSettings(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	if settings.APIKey == "" && settings.HasAPIKey {
		settings.APIKey = previous.APIKey
	}
	settings.HasAPIKey = false
	if err := s.store.Put(r.Context(), "settings", "providers", settings); err != nil {
		fail(w, err, 500)
		return
	}
	settings.HasAPIKey, settings.APIKey = settings.APIKey != "", ""
	write(w, 200, settings)
}

func (s *Server) tools(w http.ResponseWriter, r *http.Request) {
	settings, err := s.loadSettings(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	if settings.MCPURL == "" {
		fail(w, fmt.Errorf("configure an MCP server URL first"), 422)
		return
	}
	client, err := mcp.Connect(r.Context(), settings.MCPURL)
	if err != nil {
		fail(w, err, 502)
		return
	}
	defer client.Close()
	tools := []map[string]any{}
	for _, tool := range client.Tools {
		tools = append(tools, map[string]any{"name": tool.Name, "description": tool.Description, "inputSchema": tool.InputSchema})
	}
	write(w, 200, tools)
}
