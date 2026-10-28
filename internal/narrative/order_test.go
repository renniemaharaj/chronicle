package narrative

import (
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

func TestStoryRebuildKeepsHistoricalOriginsBeforeImperialEndings(t *testing.T) {
	project := domain.EmptyProject("synthetic-story", "Synthetic chronological ordering")
	project.Events = []domain.Event{
		{
			ID:    "ending",
			Phase: "Ending",
			Time: temporal.Interval{Start: &temporal.TemporalValue{
				Year:      476,
				Calendar:  temporal.CalendarJulian,
				Precision: temporal.PrecisionYear,
			}},
		},
		{
			ID:    "unknown",
			Phase: "Unknown",
			Time: temporal.Interval{Start: &temporal.TemporalValue{
				Year:      0,
				Calendar:  temporal.CalendarHistorical,
				Precision: temporal.PrecisionUnknown,
			}},
		},
		{
			ID:    "origins",
			Phase: "Origins",
			Time: temporal.Interval{Start: &temporal.TemporalValue{
				Year:        -899,
				Calendar:    temporal.CalendarHistorical,
				Precision:   temporal.PrecisionYear,
				Approximate: true,
			}},
		},
	}
	Build(&project)
	if project.Events[0].ID != "origins" || project.Events[1].ID != "ending" || project.Events[2].ID != "unknown" {
		t.Fatal("story presentation must follow recorded historical years, leaving unknown dates last")
	}

	if project.Chapters[0].Title != "Origins" || project.Events[0].Time.Start.Calendar != temporal.CalendarHistorical {
		t.Fatal("chapter ordering must retain the original historical calendar")
	}

	if _, _, known := project.Events[0].Time.Bounds(); known {
		t.Fatal("presentation order must not invent absolute historical calendar bounds")
	}
}
