package utils

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"time"
)

func GetNextOccurrence(startDate time.Time, frequency constants.Periodicity, today time.Time) time.Time {
	nextOccurrence := startDate
	for nextOccurrence.Before(today) {
		switch frequency {
		case constants.Daily:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 1)
		case constants.Weekly:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 7)
		case constants.Monthly:
			daysInMonth := daysIn(nextOccurrence.Month(), nextOccurrence.Year())
			nextOccurrence = nextOccurrence.AddDate(0, 1, 0)
			if nextOccurrence.Day() > daysInMonth {
				nextOccurrence = nextOccurrence.AddDate(0, 0, daysInMonth-nextOccurrence.Day())
			}
		}
	}
	return nextOccurrence
}

func GetRecurringDays(date time.Time, frequency constants.Periodicity) int {
	today := time.Now()
	nextOccurrence := GetNextOccurrence(date, frequency, today)
	return int(nextOccurrence.Sub(today).Hours() / 24)
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func CalculateInitialDelay(startDate time.Time, periodicity constants.Periodicity) time.Duration {
	now := time.Now()
	nextOccurrence := GetNextOccurrence(startDate, periodicity, now)
	return nextOccurrence.Sub(now)
}
