package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chronicle/internal/domain"
)

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		fail(w, fmt.Errorf("streaming is unavailable"), 500)
		return
	}
	data, err := s.store.Get(r.Context(), "research", r.PathValue("id"))
	if err != nil {
		fail(w, err, 404)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	ticker := time.NewTicker(600 * time.Millisecond)
	defer ticker.Stop()
	previous := ""
	for {
		if string(data) != previous {
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			previous = string(data)
			var run domain.ResearchRun
			json.Unmarshal(data, &run)
			if run.Status == "completed" || run.Status == "failed" || run.Status == "paused" || run.Status == "cancelled" {
				return
			}
		}
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			data, err = s.store.Get(r.Context(), "research", r.PathValue("id"))
			if err != nil {
				return
			}
		}
	}
}
