package campaign

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"time"

	"chronicle/internal/domain"
)

var words = regexp.MustCompile(`\S+`)

func Plan(project domain.Project, seconds int) (Campaign, error) {
	if project.Fixture || len(project.Events) == 0 {
		return Campaign{}, fmt.Errorf("campaigns require a saved research project with events, not a fixture")
	}

	if seconds < 15 || seconds > 180 {
		return Campaign{}, fmt.Errorf("duration target must be 15 to 180 seconds")
	}

	fingerprint := Fingerprint(project)
	items, gaps := passages(project)
	result := Campaign{
		ID:          "campaign-" + rand.Text(),
		ProjectID:   project.ID,
		ProjectHash: fingerprint,
		Title:       project.Title,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Revision:    1,
		Seconds:     seconds,
		Parts:       []Part{},
		Gaps:        gaps,
		Audit:       []Activity{},
	}
	budget := seconds * 140 / 60

	for _, source := range project.Sources {
		result.Sources = append(result.Sources, Reference{
			ID:    source.ID,
			Title: source.Title,
			URL:   source.URL,
			Hash:  source.ContentHash,
		})
	}

	part := Part{
		Number:   1,
		Status:   "planned",
		Passages: []Passage{},
		Metrics:  []Metrics{},
	}
	used := 0

	for _, unit := range readingUnits(items) {
		count := 0

		for _, passage := range unit {
			count += len(words.FindAllString(passage.Text, -1))
		}

		if used > 0 && used+count > budget {
			result.Parts = append(result.Parts, part)
			part = Part{
				Number:   part.Number + 1,
				Status:   "planned",
				Passages: []Passage{},
				Metrics:  []Metrics{},
			}
			used = 0
		}

		part.Passages = append(part.Passages, unit...)
		used += count
	}

	if len(part.Passages) > 0 {
		result.Parts = append(result.Parts, part)
	}

	return result, nil
}
