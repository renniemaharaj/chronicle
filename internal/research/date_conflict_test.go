package research

import (
	"encoding/json"
	"testing"

	"chronicle/internal/conflicts"
	"chronicle/internal/domain"
)

func TestSameOccurrenceWithCompetingDatesBranchesWithoutDiscardingClaims(t *testing.T) {
	var first extraction
	if err := json.Unmarshal([]byte(testBatch), &first); err != nil {
		t.Fatal(err)
	}
	project := domain.EmptyProject("dates", "Synthetic date disagreement")
	project.SubjectID = "subject"
	project.Entities = []domain.Entity{{ID: "subject", Name: "Meridian Observatory", Kind: "organization", Version: 1, Attributes: map[string]any{}}}
	project, err := mergeBatch(project, first, domain.Source{ID: "source-one", Title: "Synthetic archive one", Classification: "primary"})
	if err != nil {
		t.Fatal(err)
	}
	var second extraction
	json.Unmarshal([]byte(testBatch), &second)
	second.Events[0].Time.Start.Year = 1810
	second.Claims[0].Text = "The synthetic observatory was established in 1810."
	second.Evidence[0].Excerpt = second.Claims[0].Text
	project, err = mergeBatch(project, second, domain.Source{ID: "source-two", Title: "Synthetic archive two", Classification: "primary"})
	if err != nil {
		t.Fatal(err)
	}
	conflicts.Analyze(&project)
	if len(project.Events) != 1 || len(project.Claims) != 2 {
		t.Fatal("event identity must not be split solely because its sources disagree about a date")
	}
	if len(project.Conflicts) != 1 {
		t.Fatal("disjoint source dates should produce an explicit conflict")
	}
}
