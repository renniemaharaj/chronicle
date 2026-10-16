package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/fixtures"
	"chronicle/internal/research"
	"chronicle/internal/storage"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestAssistantResearchTools(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	projects, err := fixtures.Projects()
	if err != nil {
		t.Fatal(err)
	}
	project := projects[0]
	if err := store.Put(ctx, "projects", project.ID, project); err != nil {
		t.Fatal(err)
	}
	backend := httptest.NewServer(New(store, research.NewManager(store)))
	defer backend.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: backend.URL + "/api/mcp", DisableStandaloneSSE: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	tools, err := session.ListTools(ctx, nil)
	if err != nil || len(tools.Tools) != 9 {
		t.Fatalf("expected nine discoverable research tools: %v", err)
	}
	call := func(name string, args any, output any) {
		t.Helper()
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil || result.IsError {
			t.Fatalf("%s: %v, %+v", name, err, result)
		}
		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, output); err != nil {
			t.Fatal(err)
		}
	}
	var list projectList
	call("chronicle_list_projects", struct{}{}, &list)
	if len(list.Projects) != 0 {
		t.Fatal("bundled fixtures must be hidden from research libraries")
	}
	archive := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if r.URL.Path == "/robots.txt" {
			w.Write([]byte("User-agent: *\nDisallow:"))
			return
		}
		w.Write([]byte("The synthetic institution was established in 1800."))
	}))
	defer archive.Close()
	var captured captureOutput
	call("chronicle_capture_source", captureInput{URL: archive.URL, Title: "Synthetic archive", Classification: "primary"}, &captured)
	project.Fixture = false
	project.ID = "assistant-test"
	project.Sources = []domain.Source{captured.Source}
	claim := &project.Claims[0]
	claim.EvidenceIDs = []string{"passage-test"}
	project.Evidence = []domain.Evidence{{ID: "passage-test", SourceID: captured.Source.ID, ClaimID: claim.ID, Excerpt: "The synthetic institution was established in 1800.", Location: "Opening paragraph", Relation: "supporting", Mode: "quotation"}}
	var validated domain.Project
	call("chronicle_validate_project", projectInput{Project: project, RebuildNarrative: true}, &validated)
	if !validated.Evidence[0].Verified || validated.Claims[0].Status != "supported_limited" {
		t.Fatal("only matching backend-captured text should support a claim")
	}
	var saved domain.Project
	call("chronicle_save_project", projectInput{Project: validated, RebuildNarrative: true}, &saved)
	if len(saved.Chapters) == 0 {
		t.Fatal("saved research requires a deterministic narrative")
	}
	call("chronicle_list_projects", struct{}{}, &list)
	if len(list.Projects) != 1 || list.Projects[0].ID != project.ID {
		t.Fatal("saved assistant research must be visible in the app library")
	}
	if err := store.Put(ctx, "captures", captured.Source.ID, struct{}{}); err != nil {
		t.Fatal(err)
	}
	call("chronicle_validate_project", projectInput{Project: saved}, &validated)
	if !validated.Evidence[0].Verified {
		t.Fatal("unchanged, previously verified backend evidence must survive assistant edits")
	}
	project.Evidence[0].Excerpt = "An invented assertion that was never retrieved."
	project.Evidence[0].Verified = true
	call("chronicle_validate_project", projectInput{Project: project}, &validated)
	if validated.Evidence[0].Verified || validated.Claims[0].Status != "unverified" {
		t.Fatal("a caller cannot promote fabricated passages to verified evidence")
	}
}
