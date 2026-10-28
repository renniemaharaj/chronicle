package research

import (
	"context"
	"fmt"
	"time"

	"chronicle/internal/domain"
)

func (m *Manager) Start(ctx context.Context, subject string, settings domain.Settings) (domain.ResearchRun, error) {
	if subject == "" || len(subject) > 400 {
		return domain.ResearchRun{}, fmt.Errorf("subject must contain between 1 and 400 characters")
	}

	if settings.Model == "" || settings.MCPURL == "" || settings.SearchTool == "" || settings.RetrieveTool == "" {
		return domain.ResearchRun{}, fmt.Errorf("configure a model, MCP endpoint, search tool, and retrieval tool in Settings")
	}

	return m.startRun(ctx, subject, settings, nil, nil, "")
}

func (m *Manager) startRun(ctx context.Context, subject string, settings domain.Settings, project *domain.Project, plan []domain.Question, baseHash string) (domain.ResearchRun, error) {
	if plan == nil {
		plan = []domain.Question{}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	id := fmt.Sprintf("run-%d", time.Now().UnixNano())
	config := settings
	config.APIKey = ""
	run := domain.ResearchRun{
		ID:         id,
		Subject:    subject,
		Status:     "pending",
		Phase:      "Planning",
		CreatedAt:  now,
		UpdatedAt:  now,
		Model:      settings.Model,
		Config:     config,
		Plan:       plan,
		Project:    project,
		BaseHash:   baseHash,
		Activities: []domain.Activity{},
	}
	if project != nil {
		run.SubjectID = project.SubjectID
		run.Sources, run.Claims = len(project.Sources), len(project.Claims)
	}

	if err := m.Store.Put(ctx, "research", id, run); err != nil {
		return run, err
	}

	started, err := m.launch(run, settings)
	if err != nil {
		run.Status, run.Phase, run.Error = "failed", "Not started", err.Error()
		m.Store.Put(ctx, "research", id, run)
	}

	return started, err
}
