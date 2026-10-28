package research

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"chronicle/internal/domain"
	"chronicle/internal/providers"
	"chronicle/internal/storage"
)

func TestResearchFromMCPThroughVerifiedProject(t *testing.T) {
	server := testMCP(t)
	defer server.Close()

	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	manager := NewManager(store)
	manager.ModelFactory = func(domain.Settings) providers.Model {
		return modelFunc(syntheticModel)
	}
	run, err := manager.Start(context.Background(), "Meridian Observatory", domain.Settings{
		Model:        "mock",
		SearchTool:   "search",
		RetrieveTool: "retrieve",
		MCPURL:       server.URL,
		MaxRequests:  40,
		Concurrency:  2,
	})
	if err != nil {
		t.Fatal(err)
	}

	completed := awaitRun(t, store, run.ID, "completed")
	if len(completed.Project.Events) != 1 || len(completed.Project.Sources) != 1 {
		t.Fatal("repeated research must deduplicate source and event identities")
	}

	if len(completed.Project.Evidence) != 1 || !completed.Project.Evidence[0].Verified {
		t.Fatal("quotation must match actual retrieved content")
	}

	if completed.Project.Claims[0].Status != "supported_limited" {
		t.Fatal("a repeated source must not imply independent corroboration")
	}

	if len(completed.Project.Chapters) == 0 {
		t.Fatal("evidence-backed narrative should be assembled")
	}

	if err := domain.ValidateProject(*completed.Project); err != nil {
		t.Fatal(err)
	}
}

func awaitRun(t *testing.T, store storage.Store, id, status string) domain.ResearchRun {
	t.Helper()

	deadline := time.After(10 * time.Second)

	for {
		data, err := store.Get(context.Background(), "research", id)
		if err != nil {
			t.Fatal(err)
		}

		var run domain.ResearchRun
		json.Unmarshal(data, &run)
		if run.Status == status {
			return run
		}

		if run.Status == "failed" {
			t.Fatalf("research failed: %s", run.Error)
		}

		select {
		case <-deadline:
			t.Fatalf("research did not reach %s", status)
		case <-time.After(10 * time.Millisecond):
		}
	}
}
