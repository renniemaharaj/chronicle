package domain

import (
	"fmt"
	"strings"

	"chronicle/internal/temporal"
)

type EventNarration struct {
	Text           string   `json:"text"`
	ClaimIDs       []string `json:"claimIds"`
	Interpretation bool     `json:"interpretation"`
}

// Context windows are scoped interpretations, not universal historical eras.
type HistoricalContext struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Scope          string            `json:"scope"`
	Time           temporal.Interval `json:"time"`
	Summary        string            `json:"summary"`
	ClaimIDs       []string          `json:"claimIds"`
	Interpretation bool              `json:"interpretation"`
}

func validateStories(p Project, claims map[string]bool) error {
	check := func(text string, ids []string) error {
		if strings.TrimSpace(text) == "" || len(text) > 3500 {
			return fmt.Errorf("story text must contain 1–3500 bytes")
		}
		for _, id := range ids {
			if !claims[id] {
				return fmt.Errorf("story references missing claim %s", id)
			}
		}
		return nil
	}
	for _, event := range p.Events {
		if event.Narration != nil {
			if err := check(event.Narration.Text, event.Narration.ClaimIDs); err != nil {
				return err
			}
		}
	}
	seen := map[string]bool{}
	for _, period := range p.ContextPeriods {
		if period.ID == "" || seen[period.ID] || period.Title == "" || period.Scope == "" || len(period.ClaimIDs) == 0 {
			return fmt.Errorf("historical context needs unique id, title, scope, and claim references")
		}
		seen[period.ID] = true
		if err := period.Time.Validate(); err != nil {
			return err
		}
		if err := check(period.Summary, period.ClaimIDs); err != nil {
			return err
		}
	}
	return nil
}
