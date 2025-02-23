package managers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"google.golang.org/appengine/log"
	"gorm.io/gorm"
	"sort"
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

// Accounts
func (m *BasicManager) GetAccountName(accountId string) (string, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return "", err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountUUID).First(&account).Error; err != nil {
		return "", err
	}

	return account.Name, nil
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

// Categories
func (m *BasicManager) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	// add soft delete
	if err := m.DB.Where("deleted_at IS NULL").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *BasicManager) GetCategoryName(categoryId string) (string, error) {
	categoryUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return "", err
	}

	var category models.Category
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", categoryUUID).First(&category).Error; err != nil {
		return "", err
	}

	return category.Name, nil
}

func (m *BasicManager) FindCategoryByName(name string) (models.Category, error) {
	var category models.Category
	if err := m.DB.Where("name = ? AND deleted_at IS NULL", name).First(&category).Error; err != nil {
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
	if err := m.DB.Raw("SELECT c.name, SUM(t.amount) as total FROM transactions t JOIN categories c ON t.category_id = c.id WHERE t.user_id = ? AND t.deleted_at IS NULL GROUP BY c.name ORDER BY total DESC LIMIT 1", userUUID).Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Loans
func (m *BasicManager) GetLoans(userId string) ([]models.Loan, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var loans []models.Loan
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
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

func (m *BasicManager) GetUserActiveLoans(id string) ([]models.Loan, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var loans []models.Loan
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}

func (m *BasicManager) DeleteLoanById(loadId string) error {
	loanUUID, err := uuid.Parse(loadId)
	if err != nil {
		return err
	}

	var loan models.Loan
	if err := m.DB.Where("id = ?", loanUUID).First(&loan).Error; err != nil {
		return err
	}

	loan.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&loan).Error; err != nil {
		return err
	}

	return nil
}

// Recurring
func (m *BasicManager) GetRecurrings(userId string) ([]models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
		return nil, err
	}

	return recurrings, nil
}

func (m *BasicManager) GetRecurringWithCategoryName(userId string) ([]models.RecurringWithCategoryName, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
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

	fn := constants.Fn(func() {
		transactionRequest := models.TransactionRequest{
			UserID:      payload.UserID,
			Amount:      payload.Amount,
			Type:        payload.TransactionType,
			Description: payload.Name,
			CategoryID:  payload.CategoryID,
			Account:     payload.AccountID,
			Date:        time.Now().Format("2006-01-02"),
		}

		if err := m.CreateTransaction(transactionRequest); err != nil {
			log.Errorf(nil, "❌ failed to create transaction: %v", err)
		}
	})

	if err := m.DB.Create(&recurring).Error; err != nil {
		return err
	}

	if err := m.SetUserCRONJob(recurring, *m.GoCRON, fn); err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) GetUserUpcomingRecurring(userId string) (models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Recurring{}, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
		return models.Recurring{}, err
	}

	if len(recurrings) == 0 {
		return models.Recurring{}, nil
	}

	today := time.Now()
	var closestRecurring models.Recurring
	minDiff := time.Duration(1<<63 - 1) // Max duration

	for _, recurring := range recurrings {
		nextOccurrence := GetNextOccurrence(recurring.StartDate, recurring.Periodicity, today)
		diff := nextOccurrence.Sub(today)
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

func (m *BasicManager) DeleteRecurringById(recurringId string) error {
	recurringUUID, err := uuid.Parse(recurringId)
	if err != nil {
		return err
	}

	var recurring models.Recurring
	if err := m.DB.Where("id = ?", recurringUUID).First(&recurring).Error; err != nil {
		return err
	}

	if err := m.RemoveCRONJob(recurring, *m.GoCRON); err != nil {
		return err
	}

	recurring.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&recurring).Error; err != nil {
		return err
	}

	return nil
}

// Transactions
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
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *BasicManager) GetTransactionWithCategoryName(userId string) ([]models.TransactionCategoryAccounts, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&transactions).Error; err != nil {
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

	sort.SliceStable(transactionsWithCategory, func(i, j int) bool {
		return transactionsWithCategory[i].TransactionDate.After(transactionsWithCategory[j].TransactionDate)
	})

	return transactionsWithCategory, nil
}

func (m *BasicManager) FindAccountTransactions(accountId string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := m.DB.Where("account_id = ? AND deleted_at IS NULL", accountId).First(&transactions).Error; err != nil {
		return []models.Transaction{}, err
	}
	return transactions, nil
}

func (m *BasicManager) GetUserLatestSixTransactions(userId string) ([]models.Transaction, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Order("transaction_date desc").Limit(7).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *BasicManager) DeleteTransactionById(transactionId string) error {
	transactionUUID, err := uuid.Parse(transactionId)
	if err != nil {
		return err
	}

	var transaction models.Transaction
	if err := m.DB.Where("id = ?", transactionUUID).First(&transaction).Error; err != nil {
		return err
	}

	transaction.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&transaction).Error; err != nil {
		return err
	}

	return nil
}

// Random
func (m *BasicManager) CalculateBalance(accountId string, amount float64, transactionType constants.TransactionType) error {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}

	var account models.Account
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountUUID).First(&account).Error; err != nil {
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

func GetNextOccurrence(startDate time.Time, frequency constants.Periodicity, today time.Time) time.Time {
	nextOccurrence := startDate
	for nextOccurrence.Before(today) {
		switch frequency {
		case constants.Daily:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 1)
		case constants.Weekly:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 7)
		case constants.Monthly:
			nextOccurrence = nextOccurrence.AddDate(0, 1, 0)
		}
	}
	return nextOccurrence
}

func (m *BasicManager) SetUserCRONJob(recurring models.Recurring, scheduler gocron.Scheduler, taskFn constants.Fn) error {
	var err error
	var j gocron.Job

	switch recurring.Periodicity {
	case constants.Daily:
		j, err = scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	case constants.Weekly:
		j, err = scheduler.NewJob(
			gocron.WeeklyJob(
				1,
				gocron.NewWeekdays(recurring.StartDate.Weekday()),
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	case constants.Monthly:
		j, err = scheduler.NewJob(
			gocron.MonthlyJob(
				1,
				gocron.NewDaysOfTheMonth(recurring.StartDate.Day()),
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	default:
		return fmt.Errorf("❌ invalid periodicity")
	}

	if err != nil {
		return fmt.Errorf("❌ failed to create new job: %w", err)
	}

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	recurring.JobID = j.ID()
	recurring.JobName = recurring.Name

	if err := m.DB.Save(&recurring).Error; err != nil {
		tx.Rollback()
	}

	tx.Commit()

	return nil
}

func (m *BasicManager) RemoveCRONJob(recurring models.Recurring, scheduler gocron.Scheduler) error {
	if err := scheduler.RemoveJob(recurring.JobID); err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) LoadAndScheduleJobs() error {
	var recurrings []models.Recurring
	if err := m.DB.Find(&recurrings).Error; err != nil {
		return err
	}

	for _, recurring := range recurrings {
		fn := constants.Fn(func() {
			transactionRequest := models.TransactionRequest{
				UserID:      recurring.UserID.String(),
				Amount:      recurring.Amount,
				Type:        int(recurring.TransactionType),
				Description: recurring.Name,
				CategoryID:  recurring.CategoryID.String(),
				Account:     recurring.AccountID.String(),
				Date:        time.Now().Format("2006-01-02"),
			}

			if err := m.CreateTransaction(transactionRequest); err != nil {
				log.Errorf(nil, "❌ failed to create transaction: %v", err)
			}
		})

		if err := m.SetUserCRONJob(recurring, *m.GoCRON, fn); err != nil {
			return err
		}
	}
	return nil
}
