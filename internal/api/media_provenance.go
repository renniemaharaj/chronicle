package api

import (
	"context"
	"encoding/json"

	"chronicle/internal/domain"
	"chronicle/internal/evidence"
	"chronicle/internal/media"
)

func (s *Server) assessMedia(ctx context.Context, project *domain.Project) {
	var canonical domain.Project
	if data, err := s.store.Get(ctx, "projects", project.ID); err == nil {
		_ = json.Unmarshal(data, &canonical)
	}
	for index := range project.Sources {
		source := &project.Sources[index]
		source.Resources = nil
		var capture evidence.Capture
		data, err := s.store.Get(ctx, "captures", source.ID)
		if err == nil && json.Unmarshal(data, &capture) == nil && capture.Source.URL == source.URL && capture.Source.ContentHash == source.ContentHash {
			source.Resources = capture.Source.Resources
			continue
		}
		for _, existing := range canonical.Sources {
			if existing.ID == source.ID && existing.URL == source.URL && existing.ContentHash == source.ContentHash {
				source.Resources = existing.Resources
			}
		}
	}
	media.Assess(project)
}
