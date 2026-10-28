package api

import (
	"context"
	"net/http"
	"time"

	"chronicle/internal/speech"
)

type audioPreparation struct {
	speech.Request
	Generate bool `json:"generate"`
}

func (s *Server) audioKey(input speech.Request) string {
	endpoint := s.config.Speech.Endpoint
	if input.Provider == "elevenlabs" {
		endpoint = s.config.Speech.ElevenLabsURL
	}

	return s.audio.Key(input, endpoint)
}

func (s *Server) ensureAudio(ctx context.Context, input speech.Request) (bool, error) {
	return s.audio.Ensure(ctx, s.audioKey(input), func() ([]byte, error) {
		if input.Provider == "elevenlabs" {
			client, err := s.elevenLabs(ctx)
			if err != nil {
				return nil, err
			}

			return client.Synthesize(ctx, input)
		}

		client, err := speech.New(ctx, s.config.Speech)
		if err != nil {
			return nil, err
		}

		return client.Synthesize(ctx, input)
	})
}

func (s *Server) prepareSpeech(w http.ResponseWriter, r *http.Request) {
	var input audioPreparation
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)

		return
	}

	input.Request = input.Request.Normalized()
	if err := input.Validate(); err != nil {
		fail(w, err, 422)

		return
	}

	id := s.audioKey(input.Request)
	exists, err := s.audio.Has(id)
	cached := exists
	if err == nil && !exists && input.Generate {
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()

		cached, err = s.ensureAudio(ctx, input.Request)
		exists = err == nil
	}

	if err != nil {
		failSpeech(w, err)

		return
	}

	write(w, 200, map[string]any{
		"available": exists,
		"cached":    cached,
		"url":       "/api/speech/audio/" + id,
	})
}

func (s *Server) speechAudio(w http.ResponseWriter, r *http.Request) {
	var input speech.Request
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)

		return
	}

	input = input.Normalized()
	if err := input.Validate(); err != nil {
		fail(w, err, 422)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	if _, err := s.ensureAudio(ctx, input); err != nil {
		failSpeech(w, err)

		return
	}

	r.SetPathValue("id", s.audioKey(input))
	s.savedSpeech(w, r)
}

func (s *Server) savedSpeech(w http.ResponseWriter, r *http.Request) {
	file, err := s.audio.Open(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)

		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fail(w, err, 500)

		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeContent(w, r, "narration.mp3", stat.ModTime(), file)
}
