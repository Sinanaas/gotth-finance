package managers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
)

func (m *BasicManager) GetUserCategories(userId string) ([]models.Category, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	var categories []models.Category
	if err := m.DB.Where("deleted_at IS NULL AND (user_id IS NULL OR user_id = ?)", userUUID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *BasicManager) CreateUserCategory(userId, name, description string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	category := models.Category{
		Name:        name,
		Description: description,
		UserID:      &userUUID,
	}
	return m.DB.Create(&category).Error
}

func (m *BasicManager) DeleteUserCategory(categoryId, userId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	catUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return err
	}
	return m.DB.Where("id = ? AND user_id = ?", catUUID, userUUID).Delete(&models.Category{}).Error
}

func (m *BasicManager) FindCategoryByName(name string) (models.Category, error) {
	var category models.Category
	if err := m.DB.Where("name = ? AND deleted_at IS NULL AND user_id IS NULL", name).First(&category).Error; err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (m *BasicManager) GetUserTopCategories(id string) ([]models.CategoryWithTotal, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var categories []models.CategoryWithTotal
	if err := m.DB.Raw("SELECT c.name, SUM(t.amount) as total FROM transactions t JOIN categories c ON t.category_id = c.id WHERE t.user_id = ? AND t.deleted_at IS NULL GROUP BY c.name ORDER BY total DESC LIMIT 5", userUUID).Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
