package managers

import (
	"fmt"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *BasicManager) GetGoalStatus(userId string) ([]models.GoalStatus, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	var goals []models.Goal
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&goals).Error; err != nil {
		return nil, err
	}

	statuses := make([]models.GoalStatus, 0, len(goals))
	for _, g := range goals {
		percent := 0
		if g.TargetAmount > 0 {
			percent = int(g.SavedAmount / g.TargetAmount * 100)
			if percent > 100 {
				percent = 100
			}
		}
		statuses = append(statuses, models.GoalStatus{
			ID:        g.ID,
			Name:      g.Name,
			Target:    g.TargetAmount,
			Current:   g.SavedAmount,
			Remaining: g.TargetAmount - g.SavedAmount,
			Percent:   percent,
			Reached:   g.SavedAmount >= g.TargetAmount,
		})
	}
	return statuses, nil
}

func (m *BasicManager) CreateGoal(payload models.GoalRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}
	goal := models.Goal{UserID: userUUID, Name: payload.Name, TargetAmount: payload.TargetAmount}
	return m.DB.Create(&goal).Error
}

func (m *BasicManager) DeleteGoal(id, userId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	goalUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return m.DB.Where("id = ? AND user_id = ?", goalUUID, userUUID).Delete(&models.Goal{}).Error
}

// ContributeToGoal deducts amount from the chosen account as an expense transaction
// and increments the goal's SavedAmount. The goal itself is not bound to any account.
func (m *BasicManager) ContributeToGoal(goalId, userId, fromAccountId string, amount float64) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	goalUUID, err := uuid.Parse(goalId)
	if err != nil {
		return err
	}
	accountUUID, err := uuid.Parse(fromAccountId)
	if err != nil {
		return err
	}

	var goal models.Goal
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", goalUUID, userUUID).First(&goal).Error; err != nil {
		return err
	}

	// Find or create the global "Savings" category (atomic upsert avoids duplicates).
	savingsCat := models.Category{Name: "Savings", Description: "Goal contributions"}
	if err := m.DB.Where("name = ? AND user_id IS NULL AND deleted_at IS NULL", "Savings").
		FirstOrCreate(&savingsCat).Error; err != nil {
		return err
	}

	// Run the balance check and both writes inside a single transaction with a
	// row-level lock so concurrent contributions can't overdraft the account.
	err = m.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", accountUUID, userUUID).
			First(&account).Error; err != nil {
			return fmt.Errorf("account not found")
		}
		if account.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}

		transaction := models.Transaction{
			UserID:          userUUID,
			AccountID:       accountUUID,
			CategoryID:      savingsCat.ID,
			Amount:          amount,
			TransactionType: constants.Expenses,
			Description:     "Saving toward: " + goal.Name,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		goal.SavedAmount += amount
		return tx.Save(&goal).Error
	})
	if err != nil {
		return err
	}

	return m.RecalculateAccountBalance(fromAccountId)
}
