package temporal

import (
	"errors"
	"math"
)

func (v TemporalValue) Bounds() (int64, int64, bool) {
	if v.Precision == PrecisionUnknown || (v.Calendar != CalendarJulian && v.Calendar != CalendarProlepticGregorian) {
		return 0, 0, false
	}

	if v.Validate() != nil {
		return 0, 0, false
	}

	startYear, endYear := v.Year, v.Year
	if v.Precision == PrecisionDecade {
		startYear = int(floorDiv(int64(v.Year), 10)) * 10
		endYear = startYear + 9
	}

	if v.Precision == PrecisionCentury {
		startYear = int(floorDiv(int64(v.Year), 100)) * 100
		endYear = startYear + 99
	}

	startYear -= v.UncertaintyYears
	endYear += v.UncertaintyYears
	month, day := v.Month, v.Day
	if month == 0 {
		month = 1
	}

	if day == 0 {
		day = 1
	}

	lastMonth, lastDay := month, day
	if v.Precision == PrecisionMonth {
		lastDay = daysInMonth(endYear, month, v.Calendar)
	}

	if v.Precision == PrecisionYear || v.Precision == PrecisionDecade || v.Precision == PrecisionCentury {
		lastMonth, lastDay = 12, 31
	}

	start := julianDay(startYear, month, day, v.Calendar) * 86400
	end := julianDay(endYear, lastMonth, lastDay, v.Calendar)*86400 + 86399
	if v.Precision == PrecisionExact {
		seconds, _ := clockSeconds(v.TimeOfDay)
		start += seconds
		end = julianDay(endYear, lastMonth, lastDay, v.Calendar)*86400 + seconds
	}

	return start, end, true
}

func (v Interval) Bounds() (int64, int64, bool) {
	start, end := int64(math.MinInt64), int64(math.MaxInt64)
	if v.Start == nil && v.End == nil {
		return 0, 0, false
	}

	if v.Start != nil {
		first, last, known := v.Start.Bounds()
		if !known {
			return 0, 0, false
		}

		start, end = first, last
	}

	if v.End != nil {
		_, last, known := v.End.Bounds()
		if !known {
			return 0, 0, false
		}

		end = last
	}

	if v.OpenEnd {
		end = math.MaxInt64
	}

	return start, end, true
}

func (v Interval) Validate() error {
	for _, point := range []*TemporalValue{v.Start, v.End} {
		if point != nil {
			if err := point.Validate(); err != nil {
				return err
			}
		}
	}

	start, end, known := v.Bounds()
	if known && start > end {
		return errors.New("interval starts after it ends")
	}

	return nil
}
