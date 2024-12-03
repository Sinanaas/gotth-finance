package managers

import (
	"log"
	"strconv"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func (m *BasicManager) CraeteTransaction(ctx *gin.Context) error {
	var payload models.TransactionRequest
	var err error
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		log.Fatal(err)
		return err
	}
	payload.Date = ctx.PostForm("Date")
	payload.Type, err = strconv.Atoi(ctx.PostForm("Type"))
	payload.Description = ctx.PostForm("Description")
	payload.Category = ctx.PostForm("Category")

	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}

	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		return err
	}

	transactionDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(payload.Category)
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
	}

	if err := m.DB.Create(&transaction).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) GetUserTransactions(user_id string) ([]models.Transaction, error) {
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ?", userUUID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *BasicManager) GetCategoryName(category_id string) (string, error) {
	categoryUUID, err := uuid.Parse(category_id)
	if err != nil {
		return "", err
	}

	var category models.Category
	if err := m.DB.Where("id = ?", categoryUUID).First(&category).Error; err != nil {
		return "", err
	}

	return category.Name, nil
}

func (m *BasicManager) GetTransactionWithCategoryName(user_id string) ([]models.TransactionWithCategory, error) {
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := m.DB.Where("user_id = ?", userUUID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	var transactionsWithCategory []models.TransactionWithCategory
	for _, transaction := range transactions {
		categoryName, err := m.GetCategoryName(transaction.CategoryID.String())
		if err != nil {
			return nil, err
		}

		transactionsWithCategory = append(transactionsWithCategory, models.TransactionWithCategory{
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Description:     transaction.Description,
			CategoryName:    categoryName,
			TransactionDate: transaction.TransactionDate,
		})
	}

	return transactionsWithCategory, nil
}

func (m *BasicManager) GetAccountName(account_id string) (string, error) {
	accountUUID, err := uuid.Parse(account_id)
	if err != nil {
		return "", err
	}

	var account models.Account
	if err := m.DB.Where("id = ?", accountUUID).First(&account).Error; err != nil {
		return "", err
	}

	return account.Name, nil
}

func (m *BasicManager) GetUserAccounts(user_id string) ([]models.Account, error) {
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		return nil, err
	}

	var accounts []models.Account
	if err := m.DB.Where("user_id = ?", userUUID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}
