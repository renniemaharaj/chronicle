package narrative

import (
	"encoding/json"
	"strings"
	"testing"

	"chronicle/internal/domain"
)

func TestStoriesQualifyEvidenceAndRetainCuratedText(t *testing.T) {
	project := domain.Project{
		Events: []domain.Event{{ID: "event", Title: "A documented event", ClaimIDs: []string{"claim"}},
			{ID: "missing", Title: "An unsupported event", Summary: "Never trust this unsupported summary"}},
		Claims: []domain.Claim{{ID: "claim", Text: "Two sources give different dates.", Status: "disputed"}},
	}
	EnsureStories(&project)
	if !strings.Contains(project.Events[0].Narration.Text, "Sources disagree") {
		t.Fatal("dispute was hidden")
	}
	if !strings.Contains(project.Events[1].Narration.Text, "unverified") || strings.Contains(project.Events[1].Narration.Text, "Never trust") {
		t.Fatal("unsupported summary became factual narration")
	}
	project.Events[0].Narration.Text = "A curated, evidence-linked reading."
	EnsureStories(&project)
	if project.Events[0].Narration.Text != "A curated, evidence-linked reading." {
		t.Fatal("curated text lost")
	}
	data, err := json.Marshal(project)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.Project
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Events[0].Narration.Text != project.Events[0].Narration.Text {
		t.Fatal("story did not round trip")
	}
	project.Claims[0].Status = "rejected_superseded"
	EnsureStories(&project)
	if strings.Contains(project.Events[0].Narration.Text, "curated") {
		t.Fatal("rejected claim still narrated")
	}
}
