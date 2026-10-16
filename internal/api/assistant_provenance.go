package api

import (
	"context"
	"encoding/json"

	"chronicle/internal/domain"
)

// Preserve unchanged quotations already verified by a previous backend research run.
// Assistant edits must not erase existing provenance merely because that run used an MCP retriever.
func (s *Server) existingVerifiedEvidence(ctx context.Context, project domain.Project) map[string]bool {
	trusted := map[string]bool{}
	data, err := s.store.Get(ctx, "projects", project.ID)
	if err != nil {
		return trusted
	}
	var previous domain.Project
	if json.Unmarshal(data, &previous) != nil || previous.Fixture {
		return trusted
	}
	sources := map[string]domain.Source{}
	for _, source := range previous.Sources {
		sources[source.ID] = source
	}
	for _, item := range project.Evidence {
		for _, stored := range previous.Evidence {
			if !stored.Verified || item.ID != stored.ID || item.SourceID != stored.SourceID || item.ClaimID != stored.ClaimID || item.Excerpt != stored.Excerpt || item.Mode != stored.Mode {
				continue
			}
			for _, source := range project.Sources {
				original, exists := sources[source.ID]
				if exists && source.ID == item.SourceID && source.ContentHash == original.ContentHash && source.URL == original.URL {
					trusted[item.ID] = true
				}
			}
		}
	}
	return trusted
}
