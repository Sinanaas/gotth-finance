package managers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func (m *BasicManager) CreateAccount(payload models.AccountRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}

	account := models.Account{
		UserID:      userUUID,
		Name:        payload.Name,
		Description: payload.Description,
		Balance:     payload.Balance,
	}

	if err := m.DB.Create(&account).Error; err != nil {
		return err
	}

	if payload.Balance > 0 {
		createdAccount, _ := m.GetLatestUserAccount(payload.UserID)
		var transaction models.TransactionRequest
		category, _ := m.FindCategoryByName("Initial")
		transaction.UserID = payload.UserID
		transaction.Amount = payload.Balance
		transaction.Type = 1
		transaction.Description = "Initial deposit"
		transaction.Account = createdAccount.ID.String()
		transaction.CategoryID = category.ID.String()
		transaction.Date = time.Now().Format("2006-01-02")
		err = m.CreateTransaction(transaction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BasicManager) GetUserAccounts(userId string) ([]models.Account, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var accounts []models.Account
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (m *BasicManager) GetLatestUserAccount(userId string) (models.Account, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Account{}, err
	}

	var account models.Account
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Last(&account).Error; err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (m *BasicManager) FindAccountById(accountId string) (models.Account, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return models.Account{}, err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountUUID).First(&account).Error; err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (m *BasicManager) DeleteAccountById(accountId string) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ?", accountUUID).First(&account).Error; err != nil {
		return err
	}

	account.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) RecalculateAccountBalance(accountId string) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountId).First(&account).Error; err != nil {
		return err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("account_id = ? AND deleted_at IS NULL", accountUUID).Find(&transactions).Error; err != nil {
		return err
	}

	var total float64
	for _, transaction := range transactions {
		if transaction.TransactionType == constants.Income {
			total += transaction.Amount
		} else if transaction.TransactionType == constants.Expenses {
			total -= transaction.Amount
		}
	}

	account.Balance = total
	if err := m.DB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}
