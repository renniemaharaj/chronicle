package api

import (
	"context"
	"net/http"
	"time"

	"chronicle/internal/speech"
)

func (s *Server) speechVoices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	client, err := speech.New(ctx)
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

func (s *Server) speechAudio(w http.ResponseWriter, r *http.Request) {
	var input speech.Request
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)
		return
	}
	if len(input.Text) == 0 || len(input.Text) > 5000 || input.VoiceName == "" || input.LanguageCode == "" {
		http.Error(w, "Speech needs text (at most 5000 bytes), voice and language", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	client, err := speech.New(ctx)
	if err != nil {
		fail(w, err, 503)
		return
	}
	audio, err := client.Synthesize(ctx, input)
	if err != nil {
		fail(w, err, 502)
		return
	}
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(audio)
}
