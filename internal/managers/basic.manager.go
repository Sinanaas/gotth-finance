package managers

import (
	"fmt"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BasicManager struct {
	DB *gorm.DB
}

func NewBasicManager(DB *gorm.DB) BasicManager {
	return BasicManager{DB}
}

func (m *BasicManager) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := m.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *BasicManager) CreateTransaction(payload models.TransactionRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}

	transactionDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}

	accountUUID, err := uuid.Parse(payload.Account)
	if err != nil {
		return err
	}

	transaction := models.Transaction{
		UserID:          userUUID,
		Amount:          payload.Amount,
		TransactionType: constants.TransactionType(payload.Type),
		Description:     payload.Description,
		CategoryID:      categoryUUID,
		TransactionDate: transactionDate,
		AccountID:       accountUUID,
	}

	err = m.CalculateBalance(payload.Account, transaction.Amount, transaction.TransactionType)
	if err != nil {
		return err
	}
	if err := m.DB.Create(&transaction).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) GetUserTransactions(userId string) ([]models.Transaction, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ?", userUUID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *BasicManager) GetCategoryName(categoryId string) (string, error) {
	categoryUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return "", err
	}

	var category models.Category
	if err := m.DB.Where("id = ?", categoryUUID).First(&category).Error; err != nil {
		return "", err
	}

	return category.Name, nil
}

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
	}

	return nil
}

func (m *BasicManager) GetTransactionWithCategoryName(userId string) ([]models.TransactionCategoryAccounts, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ?", userUUID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	var transactionsWithCategory []models.TransactionCategoryAccounts
	for _, transaction := range transactions {
		categoryName, err := m.GetCategoryName(transaction.CategoryID.String())
		if err != nil {
			return nil, err
		}
		var accountName string
		if transaction.AccountID != uuid.Nil {
			accountName, err = m.GetAccountName(transaction.AccountID.String())
			if err != nil {
				return nil, err
			}
		} else {
			accountName = ""
		}

		transactionsWithCategory = append(transactionsWithCategory, models.TransactionCategoryAccounts{
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Description:     transaction.Description,
			CategoryName:    categoryName,
			AccountName:     accountName,
			TransactionDate: transaction.TransactionDate,
		})
	}

	return transactionsWithCategory, nil
}

func (m *BasicManager) GetAccountName(accountId string) (string, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return "", err
	}

	var account models.Account
	if err := m.DB.Where("id = ?", accountUUID).First(&account).Error; err != nil {
		return "", err
	}

	return account.Name, nil
}

func (m *BasicManager) GetUserAccounts(userId string) ([]models.Account, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var accounts []models.Account
	if err := m.DB.Where("user_id = ?", userUUID).Find(&accounts).Error; err != nil {
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
	if err := m.DB.Where("user_id = ?", userUUID).Last(&account).Error; err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (m *BasicManager) FindCategoryByName(name string) (models.Category, error) {
	var category models.Category
	if err := m.DB.Where("name = ?", name).First(&category).Error; err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (m *BasicManager) CalculateBalance(accountId string, amount float64, transactionType constants.TransactionType) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ?", accountUUID).First(&account).Error; err != nil {
		return err
	}
	transactions, err := m.FindAccountTransactions(accountId)
	if transactionType == constants.Income && len(transactions) != 0 {
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

func (m *BasicManager) FindAccountTransactions(accountId string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := m.DB.Where("account_id = ?", accountId).First(&transactions).Error; err != nil {
		return []models.Transaction{}, err
	}
	return transactions, nil
}

func (m *BasicManager) FindAccountById(accountId string) (models.Account, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return models.Account{}, err
	}

	var account models.Account
	if err := m.DB.Where("id = ?", accountUUID).First(&account).Error; err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (m *BasicManager) GetRecurringWithCategoryName(userId string) ([]models.RecurringWithCategoryName, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Where("user_id = ?", userUUID).Find(&recurrings).Error; err != nil {
		return nil, err
	}

	var recurringsWithCategory []models.RecurringWithCategoryName
	for _, recurring := range recurrings {
		categoryName, err := m.GetCategoryName(recurring.CategoryID.String())
		if err != nil {
			return nil, err
		}

		var accountName string
		if recurring.AccountID != uuid.Nil {
			accountName, err = m.GetAccountName(recurring.AccountID.String())
			if err != nil {
				return nil, err
			}
		} else {
			accountName = ""
		}

		recurringsWithCategory = append(recurringsWithCategory, models.RecurringWithCategoryName{
			Amount:          recurring.Amount,
			TransactionType: recurring.TransactionType,
			Name:            recurring.Name,
			CategoryName:    categoryName,
			Periodicity:     recurring.Periodicity,
			AccountName:     accountName,
		})
	}

	return recurringsWithCategory, nil

}

func (m *BasicManager) CreateRecurring(payload models.RecurringRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}

	accountUUID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		return err
	}

	recurring := models.Recurring{
		UserID:          userUUID,
		Name:            payload.Name,
		Amount:          payload.Amount,
		TransactionType: constants.TransactionType(payload.TransactionType),
		Periodicity:     constants.Periodicity(payload.Periodicity),
		CategoryID:      categoryUUID,
		AccountID:       accountUUID,
	}

	if err := m.DB.Create(&recurring).Error; err != nil {
		return err
	}

	return nil
}
