package domain

import (
	"testing"

	"chronicle/internal/temporal"
)

func TestStoryReferencesAndContextPrecision(t *testing.T) {
	period := HistoricalContext{ID: "era", Title: "A regional era", Scope: "A region", Summary: "A sourced interpretation.", ClaimIDs: []string{"claim"},
		Time: temporal.Interval{Start: &temporal.TemporalValue{Year: -400, Precision: "year", Calendar: "historical"}}}
	p := Project{ContextPeriods: []HistoricalContext{period}, Events: []Event{{Narration: &EventNarration{Text: "A story.", ClaimIDs: []string{"claim"}}}}}
	if err := validateStories(p, map[string]bool{"claim": true}); err != nil {
		t.Fatal(err)
	}
	if err := validateStories(p, map[string]bool{}); err == nil {
		t.Fatal("missing claim accepted")
	}
	p.ContextPeriods = append(p.ContextPeriods, period)
	if err := validateStories(p, map[string]bool{"claim": true}); err == nil {
		t.Fatal("duplicate era accepted")
	}
}
