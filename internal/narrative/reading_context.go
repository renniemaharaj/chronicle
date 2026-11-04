package narrative

import (
	"strings"

	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

func readingContext(project domain.Project, event domain.Event, claims map[string]domain.Claim) string {
	date := event.Time.Start
	if date == nil || date.Precision == temporal.PrecisionUnknown {
		return ""
	}

	parts := []string{}

	for _, period := range project.ContextPeriods {
		start, end := period.Time.Start, period.Time.End
		if start != nil && start.Precision != temporal.PrecisionUnknown && date.Year < start.Year {
			continue
		}

		if end != nil && end.Precision != temporal.PrecisionUnknown && date.Year > end.Year {
			continue
		}

		text := period.Summary
		valid, disputed := len(period.ClaimIDs) > 0, false

		for _, id := range period.ClaimIDs {
			claim, exists := claims[id]
			valid = valid && exists && claim.Status != "unverified" && claim.Status != "rejected_superseded"
			disputed = disputed || claim.Status == "disputed"
		}

		if !valid {
			text = "This historical context still needs verified supporting evidence."
		} else {
			if period.Interpretation {
				text = "Historical interpretation. " + text
			}

			if disputed {
				text = "Sources disagree about this context. " + text
			}
		}

		parts = append(parts, period.Title+". "+text)
	}

	return strings.Join(parts, " ")
}
