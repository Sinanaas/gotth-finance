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

	// Balance starts at 0: the initial-deposit transaction below brings it to
	// payload.Balance, keeping the stored balance equal to the transaction sum.
	account := models.Account{
		UserID:      userUUID,
		Name:        payload.Name,
		Description: payload.Description,
		Balance:     0,
	}

	if err := m.DB.Create(&account).Error; err != nil {
		return err
	}

	if payload.Balance > 0 {
		category, err := m.FindCategoryByName("Initial")
		if err != nil {
			return err
		}
		transaction := models.TransactionRequest{
			UserID:      payload.UserID,
			Amount:      payload.Balance,
			Type:        1,
			Description: "Initial deposit",
			Account:     account.ID.String(),
			CategoryID:  category.ID.String(),
			Date:        time.Now().Format("2006-01-02"),
		}
		if err := m.CreateTransaction(transaction); err != nil {
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

func (m *BasicManager) DeleteAccountById(accountId, userId string) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND user_id = ?", accountUUID, userUUID).First(&account).Error; err != nil {
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

	account.Balance = SumBalance(transactions)
	if err := m.DB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) FindUserAccountById(userId, id string) (models.Account, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Account{}, err
	}
	accUUID, err := uuid.Parse(id)
	if err != nil {
		return models.Account{}, err
	}
	var account models.Account
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", accUUID, userUUID).First(&account).Error; err != nil {
		return models.Account{}, err
	}
	return account, nil
}

// UpdateAccount edits name/description only; balance stays derived from transactions.
func (m *BasicManager) UpdateAccount(userId, id, name, description string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	accUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return m.DB.Model(&models.Account{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", accUUID, userUUID).
		Updates(map[string]interface{}{"name": name, "description": description}).Error
}

// SumBalance derives an account balance from its transactions: income adds,
// expense subtracts. This is the single source of truth for balances.
func SumBalance(transactions []models.Transaction) float64 {
	var total float64
	for _, t := range transactions {
		if t.TransactionType == constants.Income {
			total += t.Amount
		} else if t.TransactionType == constants.Expenses {
			total -= t.Amount
		}
	}
	return total
}
