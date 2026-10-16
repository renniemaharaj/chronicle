package temporal

import (
	"errors"
	"fmt"
)

type Precision string

type Calendar string

const (
	PrecisionExact             Precision = "exact"
	PrecisionDay               Precision = "day"
	PrecisionMonth             Precision = "month"
	PrecisionYear              Precision = "year"
	PrecisionDecade            Precision = "decade"
	PrecisionCentury           Precision = "century"
	PrecisionUnknown           Precision = "unknown"
	CalendarProlepticGregorian Calendar  = "proleptic_gregorian"
	CalendarJulian             Calendar  = "julian"
	CalendarHistorical         Calendar  = "historical"
)

type TemporalValue struct {
	Year             int       `json:"year"`
	Month            int       `json:"month,omitempty"`
	Day              int       `json:"day,omitempty"`
	Precision        Precision `json:"precision"`
	Calendar         Calendar  `json:"calendar"`
	Approximate      bool      `json:"approximate,omitempty"`
	UncertaintyYears int       `json:"uncertaintyYears,omitempty"`
	TimeOfDay        string    `json:"timeOfDay,omitempty"`
	Label            string    `json:"label,omitempty"`
}

type Interval struct {
	OpenEnd bool           `json:"openEnd,omitempty"`
	Start   *TemporalValue `json:"start,omitempty"`
	End     *TemporalValue `json:"end,omitempty"`
}

func (v TemporalValue) Validate() error {
	valid := map[Precision]bool{PrecisionExact: true, PrecisionDay: true, PrecisionMonth: true, PrecisionYear: true, PrecisionDecade: true, PrecisionCentury: true, PrecisionUnknown: true}
	if !valid[v.Precision] || v.Calendar == "" {
		return errors.New("valid precision and explicit calendar are required")
	}
	if v.Year < -1000000 || v.Year > 1000000 || v.UncertaintyYears < 0 || v.UncertaintyYears > 1000000 {
		return errors.New("year or uncertainty is outside supported bounds")
	}
	if v.Precision == PrecisionUnknown {
		return nil
	}
	if v.Month < 0 || v.Month > 12 || v.Day < 0 || v.Day > 31 {
		return errors.New("invalid calendar month or day")
	}
	if v.Precision == PrecisionMonth && (v.Month < 1 || v.Day != 0) {
		return errors.New("month precision requires a month and no day")
	}
	if v.Precision == PrecisionDay || v.Precision == PrecisionExact {
		if v.Month < 1 || v.Day < 1 || v.Day > daysInMonth(v.Year, v.Month, v.Calendar) {
			return errors.New("invalid day for the supplied month and calendar")
		}
	}
	if v.Precision == PrecisionExact {
		if _, err := clockSeconds(v.TimeOfDay); err != nil {
			return err
		}
	}
	if v.Precision == PrecisionYear || v.Precision == PrecisionDecade || v.Precision == PrecisionCentury {
		if v.Month != 0 || v.Day != 0 {
			return errors.New("coarse precision must not include a month or day")
		}
	}

	return nil
}

func (v TemporalValue) DefinitelyAfter(other TemporalValue) bool {
	start, _, known := v.Bounds()
	_, end, otherKnown := other.Bounds()
	return known && otherKnown && start > end
}

func (v TemporalValue) Overlaps(other TemporalValue) bool {
	start, end, known := v.Bounds()
	otherStart, otherEnd, otherKnown := other.Bounds()
	return known && otherKnown && start <= otherEnd && otherStart <= end
}

func (v TemporalValue) String() string {
	if v.Precision == PrecisionUnknown {
		return "date unknown"
	}
	if v.Label != "" {
		return v.Label
	}
	label := fmt.Sprint(v.Year)
	if v.Year <= 0 {
		label = fmt.Sprintf("%d BCE", 1-v.Year)
	}
	if v.Approximate {
		label = "c. " + label
	}
	return label
}
