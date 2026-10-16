package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"chronicle/internal/conflicts"
	"chronicle/internal/domain"
	"chronicle/internal/evidence"
	"chronicle/internal/narrative"
)

func (s *Server) prepareAssistantProject(ctx context.Context, input projectInput) (domain.Project, error) {
	project := input.Project
	if project.Fixture {
		return project, fmt.Errorf("assistant research must not be marked as a fixture")
	}
	if err := domain.ValidateProject(project); err != nil {
		return project, err
	}
	trusted := s.existingVerifiedEvidence(ctx, project)
	captures := map[string]evidence.Capture{}
	for _, source := range project.Sources {
		data, err := s.store.Get(ctx, "captures", source.ID)
		var capture evidence.Capture
		if err == nil && json.Unmarshal(data, &capture) == nil && capture.Source.ContentHash == source.ContentHash && capture.Source.URL == source.URL {
			captures[source.ID] = capture
		}
	}
	for index := range project.Evidence {
		item := &project.Evidence[index]
		capture, exists := captures[item.SourceID]
		passage := evidence.NormalizePassage(item.Excerpt)
		item.Verified = trusted[item.ID] || (exists && item.Mode == "quotation" && len(passage) >= 12 && strings.Contains(capture.Text, passage))
	}
	evidence.Assess(&project)
	conflicts.Analyze(&project)
	if input.RebuildNarrative || len(project.Chapters) == 0 {
		narrative.Build(&project)
	}
	narrative.EnsureStories(&project)
	s.assessMedia(ctx, &project)
	if err := domain.ValidateProject(project); err != nil {
		return project, err
	}
	return project, nil
}

func (s *Server) recordAssistantSave(project *domain.Project, researcher string) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if researcher == "" {
		researcher = "Assistant; model identifier not supplied"
	}
	project.ResearchRecords = append(project.ResearchRecords, domain.ResearchRecord{
		ID: fmt.Sprintf("assistant-%d", time.Now().UnixNano()), Subject: project.Title,
		Status: "completed", CreatedAt: now, Model: researcher,
		Plan: []domain.Question{}, Activities: []domain.Activity{{Time: now, Tool: "chronicle_save_project",
			Message: "Assistant research saved after deterministic provenance, claim and temporal checks. Coverage remains limited to retrieved evidence."}},
	})
}
