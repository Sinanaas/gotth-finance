package managers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
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
			percent = int(g.CurrentAmount / g.TargetAmount * 100)
		}
		statuses = append(statuses, models.GoalStatus{
			ID:        g.ID,
			Name:      g.Name,
			Target:    g.TargetAmount,
			Current:   g.CurrentAmount,
			Remaining: g.TargetAmount - g.CurrentAmount,
			Percent:   percent,
			Reached:   g.CurrentAmount >= g.TargetAmount,
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

// AddToGoal adds a contribution to the goal's current amount (scoped to user).
func (m *BasicManager) AddToGoal(id, userId string, amount float64) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	goalUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	var goal models.Goal
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", goalUUID, userUUID).First(&goal).Error; err != nil {
		return err
	}
	goal.CurrentAmount += amount
	return m.DB.Save(&goal).Error
}
