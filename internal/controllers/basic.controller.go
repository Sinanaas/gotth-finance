package controllers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

type BasicController struct {
	BM managers.BasicManager
}

func NewBasicController(bm managers.BasicManager) BasicController {
	return BasicController{bm}
}

func (bc *BasicController) CreateTransaction(ctx *gin.Context) {
	err := bc.BM.CraeteTransaction(ctx)
	if err != nil {
		return
	}
	ctx.Header("HX-Redirect", "/transaction")
}

func (bc *BasicController) GetAllCategories() ([]models.Category, error) {
	categories, err := bc.BM.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (bc *BasicController) GetUserTransactions(user_id string) ([]models.Transaction, error) {
	transactions, err := bc.BM.GetUserTransactions(user_id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) GetTransactionWithCategoryName(user_id string) ([]models.TransactionWithCategory, error) {
	transactions, err := bc.BM.GetTransactionWithCategoryName(user_id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) GetCategoryName(category_id string) (string, error) {
	categoryName, err := bc.BM.GetCategoryName(category_id)
	if err != nil {
		return "", err
	}
	return categoryName, nil
}

func (bc *BasicController) GetAccountName(account_id string) (string, error) {
	accountName, err := bc.BM.GetAccountName(account_id)
	if err != nil {
		return "", err
	}
	return accountName, nil
}

func (bc *BasicController) GetUserAccounts(user_id string) ([]models.Account, error) {
	accounts, err := bc.BM.GetUserAccounts(user_id)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
