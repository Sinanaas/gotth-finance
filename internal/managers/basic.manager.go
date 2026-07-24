package managers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BasicManager struct {
	DB     *gorm.DB
	GoCRON *gocron.Scheduler
}

func NewBasicManager(db *gorm.DB, goCRON *gocron.Scheduler) *BasicManager {
	return &BasicManager{
		DB:     db,
		GoCRON: goCRON,
	}
}

func (m *BasicManager) CalculateBalance(accountId string, amount float64, transactionType constants.TransactionType) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountUUID).First(&account).Error; err != nil {
		return err
	}
	_, err = m.FindAccountTransactions(accountId)
	if transactionType == constants.Income {
		account.Balance += amount
	} else if transactionType == constants.Expenses {
		if account.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}
		account.Balance -= amount
	}

	if err := m.DB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) GetUserMonthlyIncome(id string) (float64, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL AND transaction_date BETWEEN ? AND ?", userUUID, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
		return 0, err
	}

	var total float64
	for _, transaction := range transactions {
		if transaction.TransactionType == constants.Income {
			total += transaction.Amount
		}
	}

	return total, nil
}

func (m *BasicManager) GetUserMonthlyExpenses(id string) (float64, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL AND transaction_date BETWEEN ? AND ?", userUUID, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
		return 0, err
	}

	var total float64
	for _, transaction := range transactions {
		if transaction.TransactionType == constants.Expenses {
			total += transaction.Amount
		}
	}

	return total, nil
}

func (m *BasicManager) GetUserMonthlyTotals(userId string, year int, month time.Month) (income, expense float64, err error) {
	userUUID, parseErr := uuid.Parse(userId)
	if parseErr != nil {
		return 0, 0, parseErr
	}

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL AND transaction_date BETWEEN ? AND ?", userUUID, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
		return 0, 0, err
	}

	for _, t := range transactions {
		if t.TransactionType == constants.Income {
			income += t.Amount
		} else {
			expense += t.Amount
		}
	}
	return income, expense, nil
}

func (m *BasicManager) GetUserTotalBalance(userId string) float64 {
	userUUID, _ := uuid.Parse(userId)
	var accounts []models.Account
	m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&accounts)

	var total float64
	for _, account := range accounts {
		total += account.Balance
	}

	return total
}
