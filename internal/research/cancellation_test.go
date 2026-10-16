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

func TestResearchPauseAndResumePreservesSession(t *testing.T) {
	server := testMCP(t)
	defer server.Close()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	manager := NewManager(store)
	entered := make(chan struct{})
	manager.ModelFactory = func(domain.Settings) providers.Model {
		return modelFunc(func(ctx context.Context, system, input string) (json.RawMessage, error) {
			close(entered)
			<-ctx.Done()
			return nil, ctx.Err()
		})
	}
	settings := domain.Settings{Model: "mock", SearchTool: "search", RetrieveTool: "retrieve", MCPURL: server.URL, MaxRequests: 40, Concurrency: 2}
	run, err := manager.Start(context.Background(), "Meridian Observatory", settings)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("model request did not start")
	}
	paused, err := manager.Control(context.Background(), run.ID, "pause", settings)
	if err != nil {
		t.Fatal(err)
	}
	if paused.Status != "paused" || paused.Requests != 1 {
		t.Fatal("pause must preserve its request checkpoint")
	}
	manager.ModelFactory = func(domain.Settings) providers.Model { return modelFunc(syntheticModel) }
	deadline := time.Now().Add(5 * time.Second)
	for {
		manager.mu.Lock()
		active := manager.jobs[run.ID] != nil
		manager.mu.Unlock()
		if !active {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("cancelled worker did not stop")
		}
		time.Sleep(5 * time.Millisecond)
	}
	resumed, err := manager.Control(context.Background(), run.ID, "resume", settings)
	if err != nil {
		t.Fatal(err)
	}
	if resumed.ID != run.ID {
		t.Fatal("resumption must retain the reproducible session identity")
	}
	completed := awaitRun(t, store, run.ID, "completed")
	if completed.Requests <= paused.Requests {
		t.Fatal("resumption must preserve cumulative request accounting")
	}
}
