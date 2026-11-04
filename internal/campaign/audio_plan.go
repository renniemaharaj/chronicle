package campaign

import (
	"fmt"

	"chronicle/internal/domain"
)

func PlanAudio(project domain.Project, ready []ReadyReading, input AudioInput) (Campaign, error) {
	value, err := Plan(project, input.Seconds)
	if err != nil {
		return value, err
	}

	if len(ready) == 0 {
		return value, fmt.Errorf("no saved narration is available at the beginning of this story for the selected voice")
	}

	if err := ValidatePlatforms(input.Platforms); err != nil {
		return value, err
	}

	value.Platforms, value.Engine, value.Parts = input.Platforms, "ffmpeg", []Part{}
	value.ReadyThrough = ready[len(ready)-1].Event.ID
	value.RemainingEvents = len(project.Events) - len(ready)
	if value.RemainingEvents > 0 {
		value.Gaps = append(value.Gaps, fmt.Sprintf("Saved narration stops at %s; %d later events await audio.", ready[len(ready)-1].Event.Title, value.RemainingEvents))
	}

	part := Part{
		Number:   1,
		Status:   "planned",
		Passages: []Passage{},
		Metrics:  []Metrics{},
	}
	used := 0.0

	for _, reading := range ready {
		if reading.Duration <= 0 || reading.Duration > 180 {
			return value, fmt.Errorf("saved narration for %s needs splitting before short-form export", reading.Event.Title)
		}

		if used > 0 && used+reading.Duration > float64(input.Seconds) {
			value.Parts = append(value.Parts, part)
			part = Part{
				Number:   part.Number + 1,
				Status:   "planned",
				Passages: []Passage{},
				Metrics:  []Metrics{},
			}
			used = 0
		}

		ids := reading.Event.ClaimIDs
		if reading.Event.Narration != nil {
			ids = reading.Event.Narration.ClaimIDs
		}

		passage := Passage{
			EventID:  reading.Event.ID,
			Title:    reading.Event.Title,
			Text:     reading.Text,
			ClaimIDs: ids,
			Time:     reading.Event.Time,
		}
		units := readingUnits([]Passage{passage})
		count, spoken := len(words.FindAllString(reading.Text, -1)), 0

		for _, unit := range units {
			for _, chunk := range unit {
				chunk.AudioID, chunk.AudioStart = reading.AudioID, float64(spoken)/float64(count)*reading.Duration
				spoken += len(words.FindAllString(chunk.Text, -1))
				chunk.AudioEnd = float64(spoken) / float64(count) * reading.Duration
				if chunk.AudioStart == 0 {
					chunk.Visual = freezeVisual(project, reading.Event)
				}

				part.Passages = append(part.Passages, chunk)
			}
		}

		used += reading.Duration
	}

	value.Parts = append(value.Parts, part)

	return value, nil
}
