package research

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/providers"
	"chronicle/internal/storage"
)

func TestRefreshPreservesCanonicalProjectAndResearch(t *testing.T) {
	server := testMCP(t)
	defer server.Close()

	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	project := domain.EmptyProject("saved-story", "Meridian Observatory")
	project.SubjectID = "original-subject"
	project.Entities = []domain.Entity{{
		ID:      project.SubjectID,
		Name:    project.Title,
		Kind:    "organization",
		Version: 1,
	}}
	project.Notes = []domain.Note{{
		ID:       "note",
		TargetID: project.SubjectID,
		Text:     "Keep my editorial notes",
	}}
	project.Layout[project.SubjectID] = domain.Position{
		X: 40,
		Y: 80,
	}
	store.Put(context.Background(), "projects", project.ID, project)

	manager := NewManager(store)
	manager.ModelFactory = func(domain.Settings) providers.Model {
		return modelFunc(syntheticModel)
	}
	settings := domain.Settings{
		Model:        "mock",
		MCPURL:       server.URL,
		SearchTool:   "search",
		RetrieveTool: "retrieve",
		MaxRequests:  40,
		Concurrency:  1,
	}
	run, err := manager.Refresh(context.Background(), project.ID, settings)
	if err != nil {
		t.Fatal(err)
	}

	completed := awaitRun(t, store, run.ID, "completed")
	if completed.Project.ID != project.ID || completed.Project.SubjectID != project.SubjectID || len(completed.Plan) != 3 {
		t.Fatal("refresh must keep the canonical story and research a bounded update plan")
	}

	if completed.Project.Notes[0].Text != project.Notes[0].Text || completed.Project.Layout[project.SubjectID].X != 40 {
		t.Fatal("refresh must preserve notes and layout")
	}

	if len(completed.Project.Events) == 0 || !completed.Project.Evidence[0].Verified || len(completed.Project.ResearchRecords) != 1 {
		t.Fatal("refresh must save source-backed findings with an audit record")
	}

	data, _ := store.Get(context.Background(), "projects", project.ID)
	var saved domain.Project
	json.Unmarshal(data, &saved)
	if projectHash(saved) != completed.BaseHash {
		t.Fatal("checkpoint must track its own saved revision for safe resumption")
	}
}

func TestRefreshCannotOverwriteConcurrentStoryEdits(t *testing.T) {
	store, _ := storage.OpenFiles(t.TempDir())
	ctx := context.Background()
	project := domain.EmptyProject("saved", "Original")
	run := domain.ResearchRun{
		ID:       "refresh",
		BaseHash: projectHash(project),
	}
	project.Title = "User edited this while refresh ran"
	store.Put(ctx, "projects", project.ID, project)

	j := &job{
		run:   run,
		ctx:   ctx,
		store: store,
	}
	if err := j.saveResult(domain.EmptyProject("saved", "Would overwrite")); err == nil || !strings.Contains(err.Error(), "changed during refresh") {
		t.Fatal("refresh must leave concurrent edits intact")
	}
}
