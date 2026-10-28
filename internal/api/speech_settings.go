package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"chronicle/internal/storage"
)

type speechCredentials struct {
	APIKey string `json:"apiKey"`
}

func (s *Server) speechKey(ctx context.Context) (string, error) {
	if s.config.Speech.ElevenLabsKey != "" {
		return s.config.Speech.ElevenLabsKey, nil
	}

	data, err := s.store.Get(ctx, "settings", "speech")
	if errors.Is(err, storage.ErrNotFound) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	var saved speechCredentials
	err = json.Unmarshal(data, &saved)

	return saved.APIKey, err
}

func (s *Server) speechSettings(w http.ResponseWriter, r *http.Request) {
	key, err := s.speechKey(r.Context())
	if err != nil {
		fail(w, err, 500)

		return
	}

	source := "none"
	if key != "" {
		source = "workspace"
		if s.config.Speech.ElevenLabsKey != "" {
			source = "environment/config"
		}
	}

	write(w, 200, map[string]any{
		"hasApiKey": key != "",
		"keySource": source,
	})
}

func (s *Server) saveSpeechSettings(w http.ResponseWriter, r *http.Request) {
	var input speechCredentials
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)

		return
	}

	if s.config.Speech.ElevenLabsKey != "" {
		fail(w, fmt.Errorf("the deployment manages this key; change ELEVENLABS_API_KEY or the server configuration"), 409)

		return
	}

	input.APIKey = strings.TrimSpace(input.APIKey)
	if len(input.APIKey) > 512 || strings.ContainsAny(input.APIKey, "\r\n\t ") {
		fail(w, fmt.Errorf("invalid ElevenLabs API key"), 422)

		return
	}

	if err := s.store.Put(r.Context(), "settings", "speech", input); err != nil {
		fail(w, err, 500)

		return
	}

	s.speechSettings(w, r)
}
