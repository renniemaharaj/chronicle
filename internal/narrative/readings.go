package narrative

import (
	"sort"
	"strings"

	"chronicle/internal/domain"
)

type Reading struct {
	Event domain.Event
	Text  string
}

// Match the chronological story reader exactly so saved narration can be reused.
func Readings(project domain.Project) []Reading {
	events := append([]domain.Event{}, project.Events...)
	sort.SliceStable(events, func(i, j int) bool {
		if storyBefore(events[i], events[j]) {
			return true
		}

		if storyBefore(events[j], events[i]) {
			return false
		}

		return events[i].ID < events[j].ID
	})

	claims := map[string]domain.Claim{}

	for _, claim := range project.Claims {
		claims[claim.ID] = claim
	}

	result := []Reading{}
	previousChapter := ""

	for _, event := range events {
		chapter := ""
		if len(project.Chapters) > 0 {
			chapter = project.Chapters[0].ID
		}

		for _, item := range project.Chapters {
			if containsID(item.EventIDs, event.ID) {
				chapter = item.ID
				break
			}
		}

		introduction := ""
		if chapter != "" && chapter != previousChapter {
			introduction = readingContext(project, event, claims)
			previousChapter = chapter
		}

		result = append(result, Reading{
			Event: event,
			Text:  event.Title + ". " + introduction + " " + readingEvent(event, claims),
		})
	}

	return result
}

func containsID(ids []string, target string) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}

	return false
}

func readingEvent(event domain.Event, claims map[string]domain.Claim) string {
	ids := event.ClaimIDs
	if event.Narration != nil {
		ids = event.Narration.ClaimIDs
	}

	parts := []string{}
	disputed := false

	for _, id := range ids {
		claim, exists := claims[id]
		if !exists || claim.Status == "unverified" || claim.Status == "rejected_superseded" {
			continue
		}

		disputed = disputed || claim.Status == "disputed"
		parts = append(parts, claim.Text)
	}

	if len(parts) == 0 {
		return event.Title + ". Supporting evidence has not yet been verified."
	}

	text := strings.Join(parts, " ")
	if len(parts) == len(ids) && event.Narration != nil && event.Narration.Text != "" {
		text = event.Narration.Text
	}

	if disputed {
		text = "Sources disagree about part of this account. " + text
	}

	return text
}
