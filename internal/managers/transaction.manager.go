package managers

import (
	"fmt"
	"strconv"
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

	// Guard overdraft on expense before recording it.
	if transaction.TransactionType == constants.Expenses {
		var account models.Account
		if err := m.DB.Where("id = ? AND deleted_at IS NULL", accountUUID).First(&account).Error; err != nil {
			return err
		}
		if account.Balance < transaction.Amount {
			return fmt.Errorf("insufficient balance")
		}
	}

	if err := m.DB.Create(&transaction).Error; err != nil {
		return err
	}

	// Balance is always derived from the transaction set, never mutated incrementally.
	return m.RecalculateAccountBalance(payload.Account)
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
	// Destination must exist too (prevents FK violations from stale/nil account refs).
	if _, err := m.FindAccountById(payload.ToAccountID); err != nil {
		return fmt.Errorf("destination account not found")
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

	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}

	transactionDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return err
	}

	var transaction models.Transaction
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", transactionUUID, userUUID).First(&transaction).Error; err != nil {
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

func (m *BasicManager) DeleteTransactionById(transactionId, userId string) error {
	transactionUUID, err := uuid.Parse(transactionId)
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	var transaction models.Transaction
	if err := m.DB.Where("id = ? AND user_id = ?", transactionUUID, userUUID).First(&transaction).Error; err != nil {
		return err
	}

	transaction.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&transaction).Error; err != nil {
		return err
	}

	return nil
}

// ImportResult summarizes a CSV import.
type ImportResult struct {
	Imported int
	Skipped  []string // human-readable "row N: reason"
}

// ImportTransactions bulk-creates transactions from parsed CSV rows (excluding
// the header). Category and account are resolved by name (owned by the user),
// matching the export format. Valid rows are inserted atomically; invalid rows
// are reported and skipped. Affected account balances are recomputed after.
func (m *BasicManager) ImportTransactions(userId string, records [][]string) (ImportResult, error) {
	var result ImportResult

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return result, err
	}

	// Build name -> id maps for this user's categories and accounts.
	var categories []models.Category
	if err := m.DB.Where("deleted_at IS NULL AND (user_id IS NULL OR user_id = ?)", userUUID).Find(&categories).Error; err != nil {
		return result, err
	}
	catByName := map[string]uuid.UUID{}
	for _, c := range categories {
		catByName[c.Name] = c.ID
	}
	var accounts []models.Account
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&accounts).Error; err != nil {
		return result, err
	}
	accByName := map[string]uuid.UUID{}
	for _, a := range accounts {
		accByName[a.Name] = a.ID
	}

	var toInsert []models.Transaction
	affected := map[uuid.UUID]bool{}

	for i, row := range records {
		if i == 0 {
			continue // header
		}
		lineNo := i + 1
		if len(row) < 6 {
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: expected 6 columns", lineNo))
			continue
		}
		dateStr, typeStr, desc, catName, accName, amountStr := row[0], row[1], row[2], row[3], row[4], row[5]

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: invalid date %q", lineNo, dateStr))
			continue
		}
		var txType constants.TransactionType
		switch typeStr {
		case "Income":
			txType = constants.Income
		case "Expenses":
			txType = constants.Expenses
		default:
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: unknown type %q", lineNo, typeStr))
			continue
		}
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: invalid amount %q", lineNo, amountStr))
			continue
		}
		catID, ok := catByName[catName]
		if !ok {
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: unknown category %q", lineNo, catName))
			continue
		}
		accID, ok := accByName[accName]
		if !ok {
			result.Skipped = append(result.Skipped, fmt.Sprintf("row %d: unknown account %q", lineNo, accName))
			continue
		}

		toInsert = append(toInsert, models.Transaction{
			UserID:          userUUID,
			Amount:          amount,
			TransactionType: txType,
			Description:     desc,
			CategoryID:      catID,
			TransactionDate: date,
			AccountID:       accID,
		})
		affected[accID] = true
	}

	if len(toInsert) > 0 {
		if err := m.DB.Transaction(func(tx *gorm.DB) error {
			for idx := range toInsert {
				if err := tx.Create(&toInsert[idx]).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return result, err
		}
		for accID := range affected {
			if err := m.RecalculateAccountBalance(accID.String()); err != nil {
				return result, err
			}
		}
	}

	result.Imported = len(toInsert)
	return result, nil
}
