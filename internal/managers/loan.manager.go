package managers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

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
		UserID:          userUUID,
		Amount:          payload.Amount,
		ToWhom:          payload.ToWhom,
		CategoryID:      categoryUUID,
		AccountID:       accountUUID,
		LoanDate:        transactionDate,
		Status:          payload.Status,
		TransactionType: constants.TransactionType(payload.TransactionType),
		Description:     payload.Description,
	}

	if err := m.DB.Create(&loan).Error; err != nil {
		return err
	}

	// create transaction
	var transaction models.Transaction
	transaction.UserID = userUUID
	transaction.Amount = payload.Amount
	transaction.TransactionType = constants.TransactionType(payload.TransactionType)
	if loan.TransactionType == constants.Expenses {
		transaction.Description = "Lent: " + payload.Description
	} else {
		transaction.Description = "Borrowed: " + payload.Description
	}
	transaction.CategoryID = categoryUUID
	transaction.TransactionDate = transactionDate
	transaction.AccountID = accountUUID

	if err := m.DB.Create(&transaction).Error; err != nil {
		tx.Rollback()
	}

	// Update initial_transaction_id inside loan
	if err := tx.Model(&loan).Update("initial_transaction_id", transaction.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return m.RecalculateAccountBalance(loan.AccountID.String())
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
	transaction.Description = "Loan repaid: " + loan.Description
	transaction.CategoryID = loan.CategoryID
	transaction.TransactionDate = time.Now()
	transaction.AccountID = loan.AccountID

	if err := m.DB.Create(&transaction).Error; err != nil {
		tx.Rollback()
	}

	tx.Commit()

	return m.RecalculateAccountBalance(loan.AccountID.String())
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

	if err := m.DeleteTransactionById(loan.InitialTransactionID.String()); err != nil {
		return err
	}

	err = m.RecalculateAccountBalance(loan.AccountID.String())
	if err != nil {
		return err
	}

	loan.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&loan).Error; err != nil {
		return err
	}

	return nil
}
