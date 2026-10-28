package api

import (
	"context"
	"net/http"
	"time"

	"chronicle/internal/speech"
)

func (s *Server) speechRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/speech/settings", s.speechSettings)
	mux.HandleFunc("PUT /api/speech/settings", s.saveSpeechSettings)
	mux.HandleFunc("GET /api/speech/voices", s.speechVoices)
	mux.HandleFunc("GET /api/speech/models", s.speechModels)
	mux.HandleFunc("GET /api/speech/usage", s.speechUsage)
	mux.HandleFunc("POST /api/speech/synthesize", s.speechAudio)
	mux.HandleFunc("POST /api/speech/prepare", s.prepareSpeech)
	mux.HandleFunc("GET /api/speech/audio/{id}", s.savedSpeech)
}

func (s *Server) elevenLabs(ctx context.Context) (*speech.ElevenLabs, error) {
	key, err := s.speechKey(ctx)
	if err != nil {
		return nil, err
	}

	return speech.NewElevenLabs(s.config.Speech.ElevenLabsURL, key), nil
}

func (s *Server) speechVoices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if r.URL.Query().Get("provider") == "elevenlabs" {
		client, err := s.elevenLabs(ctx)
		if err != nil {
			fail(w, err, 500)

			return
		}

		voices, err := client.Voices(ctx)
		if err != nil {
			fail(w, err, 502)

			return
		}

		write(w, 200, voices)

		return
	}

	client, err := speech.New(ctx, s.config.Speech)
	if err != nil {
		fail(w, err, 503)

		return
	}

	voices, err := client.Voices(ctx, r.URL.Query().Get("language"))
	if err != nil {
		fail(w, err, 502)

		return
	}

	write(w, 200, voices)
}

func (s *Server) speechModels(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	client, err := s.elevenLabs(ctx)
	if err != nil {
		fail(w, err, 500)

		return
	}

	models, err := client.Models(ctx)
	if err != nil {
		fail(w, err, 502)

		return
	}

	write(w, 200, models)
}

func (s *Server) speechUsage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	client, err := s.elevenLabs(ctx)
	if err != nil {
		fail(w, err, 500)

		return
	}

	usage, err := client.Usage(ctx)
	if err != nil {
		fail(w, err, 502)

		return
	}

	write(w, 200, usage)
}
