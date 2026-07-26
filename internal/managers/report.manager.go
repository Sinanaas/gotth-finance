package managers

import (
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
)

// CategoryBreakdown returns current-month expense totals grouped by category.
func (m *BasicManager) CategoryBreakdown(userId string) ([]models.CategoryWithTotal, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, -1)

	var rows []models.CategoryWithTotal
	err = m.DB.Raw(`
		SELECT c.name, SUM(t.amount) AS total
		FROM transactions t JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.transaction_type = ? AND t.deleted_at IS NULL
		  AND t.transaction_date BETWEEN ? AND ?
		GROUP BY c.name ORDER BY total DESC`,
		userUUID, constants.Expenses, start, end).Scan(&rows).Error
	return rows, err
}

// NetWorthSeries returns net worth (cumulative income - expense across all
// accounts) at the end of each of the last `months` months.
func (m *BasicManager) NetWorthSeries(userId string, months int) (labels []string, values []float64, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		monthEnd := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, -1)

		var transactions []models.Transaction
		if err := m.DB.Where("user_id = ? AND deleted_at IS NULL AND transaction_date <= ?", userUUID, monthEnd).
			Find(&transactions).Error; err != nil {
			return nil, nil, err
		}

		// Goal contributions are recorded as "Savings" expenses, which reduces
		// SumBalance. Add them back because the money is still the user's.
		var savingsTotal float64
		m.DB.Raw(`
			SELECT COALESCE(SUM(t.amount), 0)
			FROM transactions t
			JOIN categories c ON t.category_id = c.id
			WHERE t.user_id = ? AND t.deleted_at IS NULL
			  AND t.transaction_date <= ?
			  AND c.name = 'Savings' AND c.user_id IS NULL`,
			userUUID, monthEnd).Scan(&savingsTotal)

		labels = append(labels, t.Format("Jan '06"))
		values = append(values, SumBalance(transactions)+savingsTotal)
	}
	return labels, values, nil
}
