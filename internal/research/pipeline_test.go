package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"chronicle/internal/domain"
	"chronicle/internal/providers"
	"chronicle/internal/storage"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type modelFunc func(context.Context, string, string) (json.RawMessage, error)

func (f modelFunc) Complete(ctx context.Context, system, input string) (json.RawMessage, error) {
	return f(ctx, system, input)
}

const testDocument = "Meridian Archive. The synthetic observatory was established in 1800. Section: Foundation."
const testURL = "https://archive.example/meridian"

func testMCP(t *testing.T) *httptest.Server {
	t.Helper()
	server := sdk.NewServer(&sdk.Implementation{Name: "synthetic-archive", Version: "1"}, nil)
	for _, name := range []string{"search", "retrieve"} {
		server.AddTool(&sdk.Tool{Name: name, InputSchema: map[string]any{"type": "object"}}, func(ctx context.Context, request *sdk.CallToolRequest) (*sdk.CallToolResult, error) {
			text := testURL
			if request.Params.Name == "retrieve" {
				text = testDocument
			}
			return &sdk.CallToolResult{Content: []sdk.Content{&sdk.TextContent{Text: text}}}, nil
		})
	}
	handler := sdk.NewStreamableHTTPHandler(func(*http.Request) *sdk.Server { return server }, nil)
	return httptest.NewServer(handler)
}

func syntheticModel(ctx context.Context, system, input string) (json.RawMessage, error) {
	var value any
	switch {
	case strings.Contains(system, "Identify and provisionally"):
		value = planOutput{Name: "Meridian Observatory", Kind: "organization", Questions: []domain.Question{{Question: "Origins", Interval: "earliest records", Status: "pending"}, {Question: "Development", Interval: "early period", Status: "pending"}, {Question: "Relationships", Interval: "middle period", Status: "pending"}, {Question: "Legacy", Interval: "later period", Status: "pending"}}}
	case strings.Contains(system, "Create JSON arguments"):
		value = argumentsOutput{Arguments: map[string]any{"query": "Meridian Observatory"}}
	case strings.Contains(system, "Select one relevant"):
		value = selectionOutput{URL: testURL, Arguments: map[string]any{"url": testURL}}
	default:
		var extracted extraction
		json.Unmarshal([]byte(testBatch), &extracted)
		value = extracted
	}
	return json.Marshal(value)
}

func TestResearchFromMCPThroughVerifiedProject(t *testing.T) {
	server := testMCP(t)
	defer server.Close()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	manager := NewManager(store)
	manager.ModelFactory = func(domain.Settings) providers.Model { return modelFunc(syntheticModel) }
	run, err := manager.Start(context.Background(), "Meridian Observatory", domain.Settings{Model: "mock", SearchTool: "search", RetrieveTool: "retrieve", MCPURL: server.URL, MaxRequests: 40, Concurrency: 2})
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
