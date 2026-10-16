package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"chronicle/internal/conflicts"
	"chronicle/internal/domain"
	"chronicle/internal/evidence"
)

func (s *Server) projects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.researchProjects(r.Context())
	if err != nil {
		fail(w, err, 500)
		return
	}
	write(w, 200, projects)
}

func (s *Server) saveProject(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	if err := decode(w, r, &project); err != nil {
		fail(w, err, 400)
		return
	}
	if project.ID != r.PathValue("id") {
		fail(w, fmt.Errorf("project identifier does not match URL"), 400)
		return
	}
	if err := domain.ValidateProject(project); err != nil {
		fail(w, err, 422)
		return
	}
	if !project.Fixture {
		evidence.Assess(&project)
		conflicts.Analyze(&project)
	}
	if err := s.store.Put(r.Context(), "projects", project.ID, project); err != nil {
		fail(w, err, 500)
		return
	}
	write(w, 200, project)
}

func (s *Server) runs(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.List(r.Context(), "research")
	if err != nil {
		fail(w, err, 500)
		return
	}
	result := []domain.ResearchRun{}
	for _, item := range items {
		var run domain.ResearchRun
		if err := json.Unmarshal(item, &run); err != nil {
			fail(w, err, 500)
			return
		}
		run.Config.APIKey = ""
		result = append(result, run)
	}
	write(w, 200, result)
}
