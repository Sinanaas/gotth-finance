package managers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

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

func (m *BasicManager) FilterTransactions(userId, startDate, endDate, categoryID, accountID, search string, txType int, page, pageSize int) ([]models.Transaction, int64, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, 0, err
	}

	query := m.DB.Preload("Category").Preload("Account").
		Where("user_id = ? AND deleted_at IS NULL", userUUID)

	if startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("transaction_date <= ?", endDate)
	}
	if categoryID != "" {
		catUUID, err := uuid.Parse(categoryID)
		if err == nil {
			query = query.Where("category_id = ?", catUUID)
		}
	}
	if accountID != "" {
		accUUID, err := uuid.Parse(accountID)
		if err == nil {
			query = query.Where("account_id = ?", accUUID)
		}
	}
	if search != "" {
		query = query.Where("description ILIKE ?", "%"+search+"%")
	}
	// txType: -1 = all, 0 = expenses, 1 = income
	if txType >= 0 {
		query = query.Where("transaction_type = ?", txType)
	}

	var total int64
	query.Model(&models.Transaction{}).Count(&total)

	offset := (page - 1) * pageSize
	var transactions []models.Transaction
	if err := query.Order("transaction_date desc").Limit(pageSize).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
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

func (m *BasicManager) CreateTransfer(payload models.TransferRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}
	fromUUID, err := uuid.Parse(payload.FromAccountID)
	if err != nil {
		return err
	}
	toUUID, err := uuid.Parse(payload.ToAccountID)
	if err != nil {
		return err
	}

	transferDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return err
	}

	fromAccount, err := m.FindAccountById(payload.FromAccountID)
	if err != nil {
		return err
	}
	if fromAccount.Balance < payload.Amount {
		return fmt.Errorf("insufficient balance in source account")
	}

	// Fetch system categories for Transfer Out / Transfer In
	transferOutCat, err := m.FindCategoryByName("Transfer Out")
	if err != nil {
		return fmt.Errorf("system category 'Transfer Out' not found: seed it first")
	}
	transferInCat, err := m.FindCategoryByName("Transfer In")
	if err != nil {
		return fmt.Errorf("system category 'Transfer In' not found: seed it first")
	}

	groupID := uuid.New()
	desc := payload.Description
	if desc == "" {
		desc = "Transfer"
	}

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Debit from source account
	outTx := models.Transaction{
		UserID:          userUUID,
		Amount:          payload.Amount,
		TransactionType: constants.Expenses,
		Description:     desc + " (out)",
		CategoryID:      transferOutCat.ID,
		TransactionDate: transferDate,
		AccountID:       fromUUID,
		TransferGroupID: &groupID,
	}
	if err := tx.Create(&outTx).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Credit to destination account
	inTx := models.Transaction{
		UserID:          userUUID,
		Amount:          payload.Amount,
		TransactionType: constants.Income,
		Description:     desc + " (in)",
		CategoryID:      transferInCat.ID,
		TransactionDate: transferDate,
		AccountID:       toUUID,
		TransferGroupID: &groupID,
	}
	if err := tx.Create(&inTx).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Recalculate both accounts
	if err := m.RecalculateAccountBalance(payload.FromAccountID); err != nil {
		return err
	}
	if err := m.RecalculateAccountBalance(payload.ToAccountID); err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) FindTransactionById(transactionId string) (models.Transaction, error) {
	transactionUUID, err := uuid.Parse(transactionId)
	if err != nil {
		return models.Transaction{}, err
	}

	var transaction models.Transaction
	if err := m.DB.Preload("Category").Preload("Account").Where("id = ? AND deleted_at IS NULL", transactionUUID).First(&transaction).Error; err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func (m *BasicManager) UpdateTransaction(payload models.TransactionUpdateRequest) error {
	transactionUUID, err := uuid.Parse(payload.ID)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}

	newAccountUUID, err := uuid.Parse(payload.NewAccountID)
	if err != nil {
		return err
	}

	transactionDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return err
	}

	var transaction models.Transaction
	if err := m.DB.Where("id = ? AND deleted_at IS NULL", transactionUUID).First(&transaction).Error; err != nil {
		return err
	}

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Reverse old transaction effect on old account
	if err := m.RecalculateAccountBalance(payload.OldAccountID); err != nil {
		tx.Rollback()
		return err
	}

	transaction.Amount = payload.Amount
	transaction.TransactionType = constants.TransactionType(payload.Type)
	transaction.Description = payload.Description
	transaction.CategoryID = categoryUUID
	transaction.TransactionDate = transactionDate
	transaction.AccountID = newAccountUUID

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Recalculate new account balance
	if err := m.RecalculateAccountBalance(payload.NewAccountID); err != nil {
		return err
	}

	// Recalculate old account if it changed
	if payload.OldAccountID != payload.NewAccountID {
		if err := m.RecalculateAccountBalance(payload.OldAccountID); err != nil {
			return err
		}
	}

	return nil
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
