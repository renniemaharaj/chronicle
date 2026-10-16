package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chronicle/internal/domain"
)

// Old bundled examples stay on disk for recovery but never enter the research library.
func (s *Server) researchProjects(ctx context.Context) ([]domain.Project, error) {
	items, err := s.store.List(ctx, "projects")
	if err != nil {
		return nil, err
	}
	projects := []domain.Project{}
	for _, item := range items {
		var project domain.Project
		if err := json.Unmarshal(item, &project); err != nil {
			return nil, err
		}
		if !project.Fixture {
			projects = append(projects, project)
		}
	}
	return projects, nil
}

func (s *Server) project(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.Get(r.Context(), "projects", r.PathValue("id"))
	if err != nil {
		fail(w, err, 500)
		return
	}
	write(w, 200, json.RawMessage(data))
}

func (s *Server) workspaceEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		fail(w, fmt.Errorf("streaming is unavailable"), 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	var previous [32]byte
	for {
		projects, err := s.researchProjects(r.Context())
		if err != nil {
			return
		}
		data, err := json.Marshal(projects)
		if err != nil {
			return
		}
		hash := sha256.Sum256(data)
		if hash != previous {
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
			previous = hash
		} else if _, err := fmt.Fprint(w, ": heartbeat\n\n"); err != nil {
			return
		}
		flusher.Flush()
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}
