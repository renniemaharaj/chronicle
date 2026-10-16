package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) startResearch(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Subject string `json:"subject"`
	}
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)
		return
	}
	settings, err := s.loadSettings(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	run, err := s.research.Start(r.Context(), input.Subject, settings)
	if err != nil {
		fail(w, err, 422)
		return
	}
	write(w, 202, run)
}

func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.Get(r.Context(), "research", r.PathValue("id"))
	if err != nil {
		fail(w, err, 500)
		return
	}
	write(w, 200, json.RawMessage(data))
}

func (s *Server) control(w http.ResponseWriter, r *http.Request) {
	settings, err := s.loadSettings(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	run, err := s.research.Control(r.Context(), r.PathValue("id"), r.PathValue("action"), settings)
	if err != nil {
		fail(w, err, 409)
		return
	}
	write(w, 200, run)
}
