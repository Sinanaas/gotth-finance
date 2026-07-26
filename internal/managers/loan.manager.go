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

	transaction := models.Transaction{
		UserID:          userUUID,
		Amount:          payload.Amount,
		TransactionType: constants.TransactionType(payload.TransactionType),
		CategoryID:      categoryUUID,
		TransactionDate: transactionDate,
		AccountID:       accountUUID,
	}
	if loan.TransactionType == constants.Expenses {
		transaction.Description = "Lent: " + payload.Description
	} else {
		transaction.Description = "Borrowed: " + payload.Description
	}

	// All writes share one transaction: a failure rolls the whole thing back.
	if err := m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&loan).Error; err != nil {
			return err
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return tx.Model(&loan).Update("initial_transaction_id", transaction.ID).Error
	}); err != nil {
		return err
	}

	return m.RecalculateAccountBalance(loan.AccountID.String())
}

func (m *BasicManager) FinishLoan(id, userId string) error {
	loanUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	var loan models.Loan
	if err := m.DB.Where("id = ? AND user_id = ?", loanUUID, userUUID).First(&loan).Error; err != nil {
		return err
	}

	repaid := models.Transaction{
		UserID:          loan.UserID,
		Amount:          loan.Amount,
		Description:     "Loan repaid: " + loan.Description,
		CategoryID:      loan.CategoryID,
		TransactionDate: time.Now(),
		AccountID:       loan.AccountID,
	}
	if loan.TransactionType == constants.Income {
		repaid.TransactionType = constants.Expenses
	} else {
		repaid.TransactionType = constants.Income
	}

	// Close the loan and record the repayment atomically.
	if err := m.DB.Transaction(func(tx *gorm.DB) error {
		loan.Status = true
		loan.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		if err := tx.Save(&loan).Error; err != nil {
			return err
		}
		return tx.Create(&repaid).Error
	}); err != nil {
		return err
	}

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

func (m *BasicManager) DeleteLoanById(loadId, userId string) error {
	loanUUID, err := uuid.Parse(loadId)
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	var loan models.Loan
	if err := m.DB.Where("id = ? AND user_id = ?", loanUUID, userUUID).First(&loan).Error; err != nil {
		return err
	}

	if err := m.DeleteTransactionById(loan.InitialTransactionID.String(), userId); err != nil {
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
