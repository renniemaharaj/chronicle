package temporal

import (
	"errors"
	"time"
)

func julianDay(year, month, day int, calendar Calendar) int64 {
	a := floorDiv(int64(14-month), 12)
	y := int64(year) + 4800 - a
	m := int64(month) + 12*a - 3
	result := int64(day) + floorDiv(153*m+2, 5) + 365*y + floorDiv(y, 4)
	if calendar == CalendarJulian {
		return result - 32083
	}

	return result - floorDiv(y, 100) + floorDiv(y, 400) - 32045
}

func floorDiv(value, divisor int64) int64 {
	quotient := value / divisor
	if value < 0 && value%divisor != 0 {
		quotient--
	}

	return quotient
}

func daysInMonth(year, month int, calendar Calendar) int {
	if month < 1 || month > 12 {
		return 0
	}

	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	leap := year%4 == 0
	if calendar == CalendarProlepticGregorian {
		leap = leap && (year%100 != 0 || year%400 == 0)
	}

	if month == 2 && leap {
		return 29
	}

	return days[month-1]
}

func clockSeconds(value string) (int64, error) {
	clock, err := time.Parse("15:04:05Z07:00", value)
	if err != nil {
		return 0, errors.New("exact time requires HH:MM:SS and an explicit UTC offset")
	}

	_, offset := clock.Zone()

	return int64(clock.Hour()*3600 + clock.Minute()*60 + clock.Second() - offset), nil
}
