package managers

import (
	"fmt"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
)

func (m *BasicManager) GetUserBudgets(userId string) ([]models.Budget, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	var budgets []models.Budget
	if err := m.DB.Preload("Category").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (m *BasicManager) CreateBudget(payload models.BudgetRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}
	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}

	// One budget per (user, category).
	var count int64
	m.DB.Model(&models.Budget{}).
		Where("user_id = ? AND category_id = ? AND deleted_at IS NULL", userUUID, categoryUUID).
		Count(&count)
	if count > 0 {
		return fmt.Errorf("a budget for this category already exists")
	}

	budget := models.Budget{UserID: userUUID, CategoryID: categoryUUID, Amount: payload.Amount}
	return m.DB.Create(&budget).Error
}

func (m *BasicManager) DeleteBudget(id, userId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	budgetUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return m.DB.Where("id = ? AND user_id = ?", budgetUUID, userUUID).Delete(&models.Budget{}).Error
}

// GetBudgetStatus derives current-month spend per budget from the transaction set.
func (m *BasicManager) GetBudgetStatus(userId string) ([]models.BudgetStatus, error) {
	budgets, err := m.GetUserBudgets(userId)
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, -1)

	statuses := make([]models.BudgetStatus, 0, len(budgets))
	for _, b := range budgets {
		var spent float64
		m.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND category_id = ? AND transaction_type = ? AND deleted_at IS NULL AND transaction_date BETWEEN ? AND ?",
				userUUID, b.CategoryID, constants.Expenses, start, end).
			Select("COALESCE(SUM(amount), 0)").Scan(&spent)

		percent := 0
		if b.Amount > 0 {
			percent = int(spent / b.Amount * 100)
		}
		statuses = append(statuses, models.BudgetStatus{
			ID:           b.ID,
			CategoryName: b.Category.Name,
			Limit:        b.Amount,
			Spent:        spent,
			Remaining:    b.Amount - spent,
			Over:         spent > b.Amount,
			Percent:      percent,
		})
	}
	return statuses, nil
}
