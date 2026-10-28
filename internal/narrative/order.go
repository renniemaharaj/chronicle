package narrative

import (
	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

// Editorial sequence follows recorded years without converting historical calendars.
// Temporal conflict checks continue to require supported absolute calendar bounds.
func storyBefore(left, right domain.Event) bool {
	a, b := left.Time.Start, right.Time.Start
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
}
