package temporal

import "testing"

func TestUncertainDatesOverlapAsRanges(t *testing.T) {
	fluxions := TemporalValue{Year: 1665, Precision: PrecisionYear, Calendar: CalendarJulian, UncertaintyYears: 2}
	letter := TemporalValue{Year: 1667, Precision: PrecisionYear, Calendar: CalendarJulian}
	publication := TemporalValue{Year: 1684, Precision: PrecisionYear, Calendar: CalendarProlepticGregorian}

	if !fluxions.Overlaps(letter) {
		t.Fatal("uncertain development and correspondence dates should overlap")
	}
	if !publication.DefinitelyAfter(fluxions) {
		t.Fatal("publication should be definitely later than the uncertain development interval")
	}
}

func TestBCEValuesRemainSignedHistoricalYears(t *testing.T) {
	firstCenturyBCE := TemporalValue{Year: -80, Precision: PrecisionYear, Calendar: CalendarJulian}
	secondCenturyBCE := TemporalValue{Year: -140, Precision: PrecisionYear, Calendar: CalendarJulian}

	if !firstCenturyBCE.DefinitelyAfter(secondCenturyBCE) {
		t.Fatal("negative years should preserve chronological ordering")
	}
}

func TestRejectsInvalidPrecisionFields(t *testing.T) {
	date := TemporalValue{Year: 1687, Month: 7, Precision: PrecisionMonth, Calendar: CalendarJulian}
	if err := date.Validate(); err != nil {
		t.Fatal(err)
	}

	date.Day = 5
	if err := date.Validate(); err == nil {
		t.Fatal("month precision must not accept a day")
	}
}
