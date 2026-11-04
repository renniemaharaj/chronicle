package campaign

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"chronicle/internal/domain"
	"chronicle/internal/narrative"
	"chronicle/internal/temporal"
)

func Fingerprint(project domain.Project) string {
	events := append([]domain.Event{}, project.Events...)

	for index := range events {
		events[index].Visual = nil
	}

	data, _ := json.Marshal([]any{events, project.ContextPeriods, project.Chapters, project.Claims, project.Evidence, project.Sources})
	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}

func passages(project domain.Project) ([]Passage, []string) {
	project.Events = append([]domain.Event{}, project.Events...)
	narrative.EnsureStories(&project)

	claims := map[string]domain.Claim{}

	for _, claim := range project.Claims {
		claims[claim.ID] = claim
	}

	events := append([]domain.Event{}, project.Events...)
	sort.SliceStable(events, func(i, j int) bool {
		a, b := events[i].Time.Start, events[j].Time.Start
		if a == nil || a.Precision == temporal.PrecisionUnknown {
			return false
		}

		if b == nil || b.Precision == temporal.PrecisionUnknown {
			return true
		}

		if a.Year != b.Year {
			return a.Year < b.Year
		}

		if max(1, a.Month) != max(1, b.Month) {
			return max(1, a.Month) < max(1, b.Month)
		}

		return max(1, a.Day) < max(1, b.Day)
	})

	result, gaps := []Passage{}, []string{}

	for _, period := range project.ContextPeriods {
		text, usable := qualified(period.Summary, period.ClaimIDs, claims, period.Interpretation)
		if !usable {
			gaps = append(gaps, "Context needs evidence: "+period.Title)
		}

		result = append(result, Passage{
			ContextID: period.ID,
			Title:     period.Title,
			Text:      text,
			ClaimIDs:  period.ClaimIDs,
			Time:      period.Time,
		})
	}

	for _, event := range events {
		text, usable := qualified(event.Narration.Text, event.Narration.ClaimIDs, claims, event.Narration.Interpretation)
		if !usable {
			gaps = append(gaps, "Event needs evidence: "+event.Title)
		}

		result = append(result, Passage{
			EventID:  event.ID,
			Title:    event.Title,
			Text:     text,
			ClaimIDs: event.Narration.ClaimIDs,
			Time:     event.Time,
		})
	}

	for _, interval := range project.Coverage {
		if interval.Status != "researched" {
			gaps = append(gaps, interval.Status+": "+interval.Question+" "+interval.Notes)
		}
	}

	return result, gaps
}

func qualified(text string, ids []string, claims map[string]domain.Claim, interpretation bool) (string, bool) {
	if len(ids) == 0 {
		return "Supporting evidence has not yet been verified for this passage.", false
	}

	disputed := false

	for _, id := range ids {
		claim, exists := claims[id]
		if !exists || claim.Status == "unverified" || claim.Status == "rejected_superseded" || claim.Status == "" {
			return "Supporting evidence has not yet been verified for this passage.", false
		}

		disputed = disputed || claim.Status == "disputed"
	}

	if interpretation {
		text = "Historical interpretation. " + text
	}

	if disputed {
		text = "Sources disagree about part of this account. " + text
	}

	return text, true
}
