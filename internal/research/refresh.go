package research

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"chronicle/internal/domain"
)

func projectHash(project domain.Project) string {
	data, _ := json.Marshal(project)
	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}

func (m *Manager) Refresh(ctx context.Context, id string, settings domain.Settings) (domain.ResearchRun, error) {
	if settings.Model == "" || settings.MCPURL == "" || settings.SearchTool == "" || settings.RetrieveTool == "" {
		return domain.ResearchRun{}, fmt.Errorf("configure a model, MCP endpoint, search tool, and retrieval tool in Settings")
	}

	data, err := m.Store.Get(ctx, "projects", id)
	if err != nil {
		return domain.ResearchRun{}, err
	}

	var project domain.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return domain.ResearchRun{}, err
	}

	if project.Fixture {
		return domain.ResearchRun{}, fmt.Errorf("refresh requires a saved research project")
	}

	date := time.Now().UTC().Format("2006-01-02")
	lastCapture := "the existing evidence"

	for _, source := range project.Sources {
		if source.RetrievedAt > lastCapture || lastCapture == "the existing evidence" {
			lastCapture = source.RetrievedAt
		}
	}

	plan := []domain.Question{
		{
			Question: fmt.Sprintf("For %s, find newly documented events and public developments since %s, up to %s. Search current primary and authoritative sources. Separate event dates from publication dates; future announcements remain plans. Do not invent news if none is found.", project.Title, lastCapture, date),
			Interval: "Recent developments",
			Status:   "pending",
		},
		{
			Question: fmt.Sprintf("Review new scholarship and corrections for %s as of %s. Preserve existing claims, quotations and competing chronologies; add attributed evidence rather than rewriting history to fit headlines.", project.Title, date),
			Interval: "Corrections and competing accounts",
			Status:   "pending",
		},
		{
			Question: fmt.Sprintf("Revisit unresolved coverage and missing places for %s. Seek explicit geographic references and source passages. Preserve unknowns. Do not create death or dissolution for an ongoing subject.", project.Title),
			Interval: "Coverage and geography",
			Status:   "pending",
		},
	}

	return m.startRun(ctx, project.Title, settings, &project, plan, projectHash(project))
}

func (j *job) saveResult(project domain.Project) error {
	run := j.snapshot()
	if run.BaseHash != "" {
		data, err := j.store.Get(j.ctx, "projects", project.ID)
		if err != nil {
			return err
		}

		var current domain.Project
		if err := json.Unmarshal(data, &current); err != nil {
			return err
		}

		if projectHash(current) != run.BaseHash {
			return fmt.Errorf("the saved story changed during refresh; findings remain in this session. Start a new refresh from the current story to preserve edits")
		}
	}

	if err := j.store.Put(j.ctx, "projects", project.ID, project); err != nil {
		return err
	}

	if run.BaseHash != "" {
		j.update(func(run *domain.ResearchRun) {
			run.BaseHash = projectHash(project)
		})
	}

	return nil
}
