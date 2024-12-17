package utils

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"time"
)

func GetRecurringDate(date time.Time, frequency constants.Periodicity) string {
	switch frequency {
	case constants.Daily:
		return date.AddDate(0, 0, 1).Format("2006-01-02")
	case constants.Weekly:
		return date.AddDate(0, 0, 7).Format("2006-01-02")
	case constants.Monthly:
		return date.AddDate(0, 1, 0).Format("2006-01-02")
	case constants.Yearly:
		return date.AddDate(1, 0, 0).Format("2006-01-02")
	default:
		return ""
	}
}

func GetRecurringDays(date time.Time, frequency constants.Periodicity) int {
	switch frequency {
	case constants.Daily:
		return 1
	case constants.Weekly:
		return 7 - int(date.Weekday())
	case constants.Monthly:
		return 30 - date.Day()
	case constants.Yearly:
		return 365 - date.YearDay()
	default:
		return 0
	}
}
