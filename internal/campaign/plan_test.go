package campaign

import (
	"strings"
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

func syntheticProject() domain.Project {
	return domain.Project{
		ID:    "synthetic-campaign",
		Title: "Synthetic editorial test, not historical evidence",
		Claims: []domain.Claim{{
			ID:     "claim",
			Text:   "Synthetic assertion",
			Status: "disputed",
		}},
		Events: []domain.Event{{
			ID:       "event",
			Title:    "Synthetic event",
			ClaimIDs: []string{"claim"},
			Time: temporal.Interval{Start: &temporal.TemporalValue{
				Year:      1800,
				Calendar:  "julian",
				Precision: temporal.PrecisionYear,
			}},
			Narration: &domain.EventNarration{
				Text:     strings.Repeat("Café synthetic sentence.\n", 90),
				ClaimIDs: []string{"claim"},
			},
		}},
	}
}

func TestPlanPreservesEveryByteAndUncertainty(t *testing.T) {
	project := syntheticProject()
	items, _ := passages(project)
	value, err := Plan(project, 15)
	if err != nil || len(value.Parts) < 2 {
		t.Fatalf("expected a multi-part plan: %v", err)
	}

	text := ""
	next := 0

	for _, part := range value.Parts {
		count := 0

		for _, passage := range part.Passages {
			if passage.Start != next || passage.Text != items[0].Text[passage.Start:passage.End] {
				t.Fatal("passage ranges must be contiguous and exact UTF-8 byte slices")
			}

			if passage.Time.Start.Calendar != "julian" {
				t.Fatal("historical calendars must survive planning")
			}

			text += passage.Text
			next = passage.End
			count += len(words.FindAllString(passage.Text, -1))
		}

		if count > 15*140/60 {
			t.Fatal("part exceeded its editorial word budget")
		}
	}

	if text != items[0].Text || !strings.HasPrefix(text, "Sources disagree") {
		t.Fatal("the series must preserve every byte and the competing-claims qualification")
	}
}

func TestPlanDisclosesGapsAndRejectsFixtures(t *testing.T) {
	project := syntheticProject()
	project.Claims[0].Status = "unverified"
	value, err := Plan(project, 60)
	if err != nil || len(value.Gaps) == 0 {
		t.Fatalf("unverified passages need an explicit research gap: %v", err)
	}

	if strings.Contains(value.Parts[0].Passages[0].Text, "Café") {
		t.Fatal("unverified narration must not become historical evidence")
	}

	project.Fixture = true
	if _, err := Plan(project, 60); err == nil {
		t.Fatal("fixtures cannot become campaigns")
	}

	project.Fixture = false
	if _, err := Plan(project, 0); err == nil {
		t.Fatal("duration must be bounded")
	}
}

func TestLongSentenceStaysTogetherWithoutMutatingResearch(t *testing.T) {
	project := syntheticProject()
	project.Events[0].Narration.Text = strings.Repeat("synthetic ", 75) + "sentence."
	value, err := Plan(project, 15)
	if err != nil {
		t.Fatal(err)
	}

	if len(value.Parts) != 2 || len(value.Parts[1].Passages) != 4 {
		t.Fatal("a long sentence should stay in one reel with multiple readable captions")
	}

	project.Claims[0].Status = "unverified"
	original := project.Events[0].Narration.Text
	if _, err := Plan(project, 15); err != nil {
		t.Fatal(err)
	}

	if project.Events[0].Narration.Text != original {
		t.Fatal("editorial planning must not modify canonical research")
	}
}
