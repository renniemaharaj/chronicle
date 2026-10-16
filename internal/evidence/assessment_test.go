package evidence

import (
	"testing"

	"chronicle/internal/domain"
)

func TestCopiedSourcesDoNotBecomeIndependentCorroboration(t *testing.T) {
	p := domain.EmptyProject("test", "Synthetic source assessment")
	p.Claims = []domain.Claim{{ID: "claim", Text: "Synthetic assertion", Status: "unverified"}}
	p.Sources = []domain.Source{{ID: "one", Classification: "primary", IndependenceGroup: "original", IndependenceReviewed: true}, {ID: "copy", Classification: "secondary", IndependenceGroup: "original", IndependenceReviewed: true}}
	p.Evidence = []domain.Evidence{{ID: "a", ClaimID: "claim", SourceID: "one", Relation: "supporting", Verified: true}, {ID: "b", ClaimID: "claim", SourceID: "copy", Relation: "supporting", Verified: true}}
	Assess(&p)
	if p.Claims[0].Status != "supported_limited" {
		t.Fatal("copied reports must remain one source group")
	}
	p.Sources[1].IndependenceGroup = "independent"
	Assess(&p)
	if p.Claims[0].Status != "corroborated" {
		t.Fatal("reviewed independent evidence should be explainably corroborated")
	}
	p.Evidence[1].Relation = "contradicting"
	Assess(&p)
	if p.Claims[0].Status != "disputed" {
		t.Fatal("opposing evidence must remain visible")
	}
}
