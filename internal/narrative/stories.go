package narrative

import (
	"strings"

	"chronicle/internal/domain"
)

// EnsureStories provides an evidence-qualified reading for every event.
// Curated stories remain portable; rejected claims never become narration.
func EnsureStories(project *domain.Project) {
	claims := map[string]domain.Claim{}
	for _, claim := range project.Claims {
		claims[claim.ID] = claim
	}
	for index := range project.Events {
		event := &project.Events[index]
		if usableStory(event.Narration, claims) {
			if strings.TrimSpace(event.Summary) == "" {
				event.Summary = event.Narration.Text
			}
			continue
		}
		parts, ids := []string{}, []string{}
		for _, id := range event.ClaimIDs {
			claim, exists := claims[id]
			if !exists || claim.Status == "unverified" || claim.Status == "rejected_superseded" {
				continue
			}
			text := claim.Text
			if claim.Status == "disputed" {
				text = "Sources disagree about this assertion: " + text
			}
			if len(strings.Join(parts, " "))+len(text) > 2800 {
				break
			}
			parts, ids = append(parts, text), append(ids, id)
		}
		text := strings.Join(parts, " ")
		if text == "" {
			text = "This event has no assessed supporting evidence yet. Its details remain unverified."
		} else {
			if strings.TrimSpace(event.Summary) == "" {
				event.Summary = text
			}
			text += " The linked evidence records the basis and limitations of this account."
		}
		event.Narration = &domain.EventNarration{Text: text, ClaimIDs: ids}
	}
}

func usableStory(story *domain.EventNarration, claims map[string]domain.Claim) bool {
	if story == nil || strings.TrimSpace(story.Text) == "" || len(story.ClaimIDs) == 0 {
		return false
	}
	for _, id := range story.ClaimIDs {
		claim, exists := claims[id]
		if !exists || claim.Status == "unverified" || claim.Status == "rejected_superseded" {
			return false
		}
	}
	return true
}
