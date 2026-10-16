package conflicts

import (
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

func year(value int) *temporal.TemporalValue {
	return &temporal.TemporalValue{Year: value, Precision: temporal.PrecisionYear, Calendar: temporal.CalendarJulian}
}

func TestDisjointAlternativeDatesArePreserved(t *testing.T) {
	p := domain.EmptyProject("conflict", "Synthetic contradiction")
	p.Claims = []domain.Claim{{ID: "a", EventID: "event", Time: &temporal.Interval{Start: year(1890)}}, {ID: "b", EventID: "event", Time: &temporal.Interval{Start: year(1895)}}}
	Analyze(&p)
	if len(p.Conflicts) != 1 || p.Conflicts[0].Severity != "definite" {
		t.Fatal("disjoint dates for the same event should create a definite conflict")
	}
	if len(p.Claims) != 2 {
		t.Fatal("alternative claims must not be discarded")
	}
	p.Claims[0].Time.Start.UncertaintyYears = 8
	p.Conflicts = nil
	Analyze(&p)
	if len(p.Conflicts) != 0 {
		t.Fatal("overlapping uncertain dates are not definite contradictions")
	}
}

func TestPosthumousPublicationIsNotPhysicalParticipation(t *testing.T) {
	p := domain.EmptyProject("life", "Synthetic lifespan")
	p.Entities = []domain.Entity{{ID: "person", Kind: "person", Time: temporal.Interval{Start: year(1800), End: year(1850)}}}
	p.Claims = []domain.Claim{{ID: "death", EntityID: "person"}, {ID: "publication", EventID: "posthumous"}}
	p.Events = []domain.Event{{ID: "posthumous", Time: temporal.Interval{Start: year(1860)}, Participants: []string{"person"}, ClaimIDs: []string{"publication"}}}
	Analyze(&p)
	if len(p.Conflicts) != 0 {
		t.Fatal("a posthumous publication is possible")
	}
	p.Events[0].RequiresPresence = true
	Analyze(&p)
	if len(p.Conflicts) != 1 {
		t.Fatal("required physical participation after documented death should be flagged")
	}
}
