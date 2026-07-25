package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
)

func BuildJSONArray(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

func BuildJSONFloatArray(items []float64) string {
	parts := make([]string, len(items))
	for i, v := range items {
		parts[i] = fmt.Sprintf("%.2f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

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
