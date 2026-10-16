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

func TestFabricatedPassageNeverBecomesVerifiedEvidence(t *testing.T) {
	server := testMCP(t)
	defer server.Close()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	manager := NewManager(store)
	manager.ModelFactory = func(domain.Settings) providers.Model {
		return modelFunc(func(ctx context.Context, system, input string) (json.RawMessage, error) {
			raw, err := syntheticModel(ctx, system, input)
			if strings.Contains(system, extractionInstructions) {
				var batch extraction
				json.Unmarshal(raw, &batch)
				batch.Evidence[0].Excerpt = "An invented passage not present in the retrieved source."
				raw, err = json.Marshal(batch)
			}
			return raw, err
		})
	}
	run, err := manager.Start(context.Background(), "Meridian Observatory", domain.Settings{Model: "mock", SearchTool: "search", RetrieveTool: "retrieve", MCPURL: server.URL, MaxRequests: 40, Concurrency: 2})
	if err != nil {
		t.Fatal(err)
	}
	completed := awaitRun(t, store, run.ID, "completed")
	if len(completed.Project.Evidence) != 0 {
		t.Fatal("a fabricated excerpt must not enter verified provenance")
	}
	if completed.Project.Claims[0].Status != "unverified" {
		t.Fatal("claims without retrieved passages must stay unverified")
	}
	if strings.Contains(completed.Project.Chapters[0].Summary, "established in 1800") {
		t.Fatal("unsupported claims must not become factual narration")
	}
}
