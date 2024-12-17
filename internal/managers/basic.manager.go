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

func NewBasicManager(db *gorm.DB) *BasicManager {
	return &BasicManager{DB: db}
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
			StartDate:       recurring.StartDate,
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

	transactionDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		return err
	}

	recurring := models.Recurring{
		UserID:          userUUID,
		Name:            payload.Name,
		StartDate:       transactionDate,
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

func (m *BasicManager) GetLoanWithCategoryName(id string) ([]models.LoanCategoryAccount, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var loans []models.Loan
	var accountName string
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&loans).Error; err != nil {
		return nil, err
	}

	var loansWithCategory []models.LoanCategoryAccount
	for _, loan := range loans {
		categoryName, err := m.GetCategoryName(loan.CategoryID.String())
		if err != nil {
			return nil, err
		}

		if loan.AccountID != uuid.Nil {
			accountName, err = m.GetAccountName(loan.AccountID.String())
			if err != nil {
				return nil, err
			}
		} else {
			accountName = ""
		}

		loansWithCategory = append(loansWithCategory, models.LoanCategoryAccount{
			Amount:          loan.Amount,
			ToWhom:          loan.ToWhom,
			CategoryName:    categoryName,
			AccountName:     accountName,
			LoanDate:        loan.LoanDate,
			Status:          loan.Status,
			TransactionType: loan.TransactionType,
			Description:     loan.Description,
			ID:              loan.ID.String(),
		})
	}

	return loansWithCategory, nil

}

func (m *BasicManager) CreateLoan(payload models.LoanRequest) error {
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

	transactionDate, err := time.Parse("2006-01-02", payload.LoanDate)
	if err != nil {
		return err
	}

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	loan := models.Loan{
		UserID:     userUUID,
		Amount:     payload.Amount,
		ToWhom:     payload.ToWhom,
		CategoryID: categoryUUID,
		AccountID:  accountUUID,
		LoanDate:   transactionDate,
		Status:     payload.Status,

		TransactionType: constants.TransactionType(payload.TransactionType),
		Description:     payload.Description,
	}

	if err := m.DB.Create(&loan).Error; err != nil {
		return err
	}

	// create trnsaction
	var transaction models.Transaction
	transaction.UserID = userUUID
	transaction.Amount = payload.Amount
	transaction.TransactionType = constants.TransactionType(payload.TransactionType)
	if loan.TransactionType == constants.Expenses {
		transaction.Description = "Memberi Hutang: " + payload.Description
	} else {
		transaction.Description = "Hutang: " + payload.Description
	}
	transaction.CategoryID = categoryUUID
	transaction.TransactionDate = transactionDate
	transaction.AccountID = accountUUID

	if err := m.CalculateBalance(loan.AccountID.String(), loan.Amount, loan.TransactionType); err != nil {
		tx.Rollback()
	}

	if err := m.DB.Create(&transaction).Error; err != nil {
		tx.Rollback()
	}

	tx.Commit()
	return nil
}

func (m *BasicManager) FinishLoan(id string) error {
	loanUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var loan models.Loan
	if err := m.DB.Where("id = ?", loanUUID).First(&loan).Error; err != nil {
		return err
	}

	loan.Status = true
	loan.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&loan).Error; err != nil {
		return err
	}

	// create transaction
	var transaction models.Transaction
	transaction.UserID = loan.UserID
	transaction.Amount = loan.Amount
	if loan.TransactionType == constants.Income {
		transaction.TransactionType = constants.Expenses
	} else if loan.TransactionType == constants.Expenses {
		transaction.TransactionType = constants.Income
	}
	transaction.Description = "Bayar Hutang: " + loan.Description
	transaction.CategoryID = loan.CategoryID
	transaction.TransactionDate = time.Now()
	transaction.AccountID = loan.AccountID

	if err := m.CalculateBalance(loan.AccountID.String(), loan.Amount, transaction.TransactionType); err != nil {
		tx.Rollback()
	}

	if err := m.DB.Create(&transaction).Error; err != nil {
		tx.Rollback()
	}

	tx.Commit()
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
	if err := m.DB.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userUUID, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
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
	if err := m.DB.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userUUID, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
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

func (m *BasicManager) GetUserActiveLoans(id string) ([]models.Loan, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var loans []models.Loan
	if err := m.DB.Where("user_id = ?", userUUID).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}

func (m *BasicManager) GetUserTotalBalance(userId string) float64 {
	userUUID, _ := uuid.Parse(userId)
	var accounts []models.Account
	m.DB.Where("user_id = ?", userUUID).Find(&accounts)

	var total float64
	for _, account := range accounts {
		total += account.Balance
	}

	return total
}

func (m *BasicManager) GetUserLatestSixTransactions(userId string) ([]models.Transaction, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ?", userUUID).Order("transaction_date desc").Limit(6).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *BasicManager) GetUserUpcomingRecurring(userId string) (models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Recurring{}, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Where("user_id = ?", userUUID).Find(&recurrings).Error; err != nil {
		return models.Recurring{}, err
	}

	if len(recurrings) == 0 {
		return models.Recurring{}, gorm.ErrRecordNotFound
	}

	today := time.Now()
	var closestRecurring models.Recurring
	minDiff := time.Duration(1<<63 - 1) // Max duration

	for _, recurring := range recurrings {
		diff := recurring.StartDate.Sub(today)
		if diff >= 0 && diff < minDiff {
			minDiff = diff
			closestRecurring = recurring
		}
	}

	if minDiff == time.Duration(1<<63-1) {
		return models.Recurring{}, gorm.ErrRecordNotFound
	}

	return closestRecurring, nil
}

func (m *BasicManager) GetUserTopCategories(id string) ([]models.CategoryWithTotal, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var categories []models.CategoryWithTotal
	if err := m.DB.Raw("SELECT c.name, SUM(t.amount) as total FROM transactions t JOIN categories c ON t.category_id = c.id WHERE t.user_id = ? GROUP BY c.name ORDER BY total DESC LIMIT 1", userUUID).Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
