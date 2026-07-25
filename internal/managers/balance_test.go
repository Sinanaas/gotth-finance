package managers

import (
	"testing"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
)

func tx(amount float64, t constants.TransactionType) models.Transaction {
	return models.Transaction{Amount: amount, TransactionType: t}
}

func TestSumBalance(t *testing.T) {
	cases := []struct {
		name string
		in   []models.Transaction
		want float64
	}{
		{"empty", nil, 0},
		{"income only", []models.Transaction{tx(1000, constants.Income)}, 1000},
		{"expense only", []models.Transaction{tx(400, constants.Expenses)}, -400},
		{
			"mixed",
			[]models.Transaction{
				tx(1000, constants.Income),
				tx(250, constants.Expenses),
				tx(500, constants.Income),
				tx(100, constants.Expenses),
			},
			1150,
		},
	}
	for _, c := range cases {
		if got := SumBalance(c.in); got != c.want {
			t.Errorf("%s: SumBalance = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestGetNextOccurrence(t *testing.T) {
	today := time.Date(2026, 7, 25, 0, 0, 0, 0, time.UTC)

	// Daily starting well in the past → next occurrence is today.
	start := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	if got := GetNextOccurrence(start, constants.Daily, today); got.Before(today) {
		t.Errorf("daily next occurrence %v is before today %v", got, today)
	}

	// Monthly on the 10th → next is 2026-08-10 (this month's 10th already passed).
	start = time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	got := GetNextOccurrence(start, constants.Monthly, today)
	want := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("monthly next occurrence = %v, want %v", got, want)
	}

	// A future start date returns itself.
	future := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	if got := GetNextOccurrence(future, constants.Weekly, today); !got.Equal(future) {
		t.Errorf("future start = %v, want %v", got, future)
	}
}
