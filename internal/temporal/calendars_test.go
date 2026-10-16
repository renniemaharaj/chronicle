package temporal

import "testing"

func TestJulianGregorianEquivalentDays(t *testing.T) {
	julian := TemporalValue{Year: 1642, Month: 12, Day: 25, Precision: PrecisionDay, Calendar: CalendarJulian}
	gregorian := TemporalValue{Year: 1643, Month: 1, Day: 4, Precision: PrecisionDay, Calendar: CalendarProlepticGregorian}
	if !julian.Overlaps(gregorian) {
		t.Fatal("calendar-equivalent birth dates should overlap")
	}
	start, end, _ := julian.Bounds()
	otherStart, otherEnd, _ := gregorian.Bounds()
	if start != otherStart || end != otherEnd {
		t.Fatal("equivalent calendars should produce identical normalized intervals")
	}
}

func TestLeapYearsFollowExplicitCalendar(t *testing.T) {
	date := TemporalValue{Year: 1700, Month: 2, Day: 29, Precision: PrecisionDay, Calendar: CalendarJulian}
	if err := date.Validate(); err != nil {
		t.Fatal(err)
	}
	date.Calendar = CalendarProlepticGregorian
	if err := date.Validate(); err == nil {
		t.Fatal("1700 is not a Gregorian leap year")
	}
}

func TestUnknownCalendarDoesNotProduceFalseCertainty(t *testing.T) {
	unknown := TemporalValue{Year: -100, Precision: PrecisionYear, Calendar: "ancient_egyptian"}
	if _, _, known := unknown.Bounds(); known {
		t.Fatal("an unsupported historical calendar must remain incomparable")
	}
}

func TestAncientIntervalsAndUncertainty(t *testing.T) {
	first := TemporalValue{Year: -43, Precision: PrecisionYear, Calendar: CalendarJulian, UncertaintyYears: 3}
	second := TemporalValue{Year: -45, Precision: PrecisionYear, Calendar: CalendarJulian}
	if !first.Overlaps(second) {
		t.Fatal("BCE uncertainty must expand the historical interval")
	}
}
