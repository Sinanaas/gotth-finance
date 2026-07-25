package managers

import (
	"fmt"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
)

func (m *BasicManager) GetGoalStatus(userId string) ([]models.GoalStatus, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	var goals []models.Goal
	if err := m.DB.Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&goals).Error; err != nil {
		return nil, err
	}

	statuses := make([]models.GoalStatus, 0, len(goals))
	for _, g := range goals {
		current := g.Account.Balance // progress = the linked account's balance
		percent := 0
		if g.TargetAmount > 0 {
			percent = int(current / g.TargetAmount * 100)
		}
		statuses = append(statuses, models.GoalStatus{
			ID:          g.ID,
			Name:        g.Name,
			AccountID:   g.AccountID,
			AccountName: g.Account.Name,
			Target:      g.TargetAmount,
			Current:     current,
			Remaining:   g.TargetAmount - current,
			Percent:     percent,
			Reached:     current >= g.TargetAmount,
		})
	}
	return statuses, nil
}

func (m *BasicManager) CreateGoal(payload models.GoalRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}
	accountUUID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		return err
	}
	// Ensure the account belongs to the user.
	var count int64
	m.DB.Model(&models.Account{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", accountUUID, userUUID).
		Count(&count)
	if count == 0 {
		return fmt.Errorf("account not found")
	}

	goal := models.Goal{UserID: userUUID, Name: payload.Name, TargetAmount: payload.TargetAmount, AccountID: accountUUID}
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

// ContributeToGoal moves money from a source account into the goal's linked
// account as a real transfer — recorded as transactions, net worth preserved.
func (m *BasicManager) ContributeToGoal(goalId, userId, fromAccountId string, amount float64) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	goalUUID, err := uuid.Parse(goalId)
	if err != nil {
		return err
	}
	var goal models.Goal
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", goalUUID, userUUID).First(&goal).Error; err != nil {
		return err
	}
	if fromAccountId == goal.AccountID.String() {
		return fmt.Errorf("source account must differ from the goal's account")
	}

	return m.CreateTransfer(models.TransferRequest{
		FromAccountID: fromAccountId,
		ToAccountID:   goal.AccountID.String(),
		Amount:        amount,
		Date:          time.Now().Format("2006-01-02"),
		Description:   "Saving toward: " + goal.Name,
		UserID:        userId,
	})
}
