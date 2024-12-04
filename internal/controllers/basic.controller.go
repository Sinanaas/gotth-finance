package controllers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

type BasicController struct {
	BM managers.BasicManager
}

func NewBasicController(bm managers.BasicManager) BasicController {
	return BasicController{bm}
}

func (bc *BasicController) CreateTransaction(ctx *gin.Context) {
	var payload models.TransactionRequest
	var err error

	payload.Description = ctx.PostForm("Description")
	payload.CategoryID = ctx.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return
	}
	payload.Date = ctx.PostForm("Date")
	payload.Type, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return
	}
	payload.Account = ctx.PostForm("Account")
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}
	payload.UserID = user_id

	err = bc.BM.CreateTransaction(payload)

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

func (bc *BasicController) GetTransactionWithCategoryName(user_id string) ([]models.TransactionCategoryAccounts, error) {
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

func (bc *BasicController) CreateAccount(ctx *gin.Context) {
	var payload models.AccountRequest
	var err error

	payload.Name = ctx.PostForm("Name")
	payload.Description = ctx.PostForm("Description")
	payload.Balance, err = strconv.ParseFloat(ctx.PostForm("Balance"), 64)
	if err != nil {
		return
	}
	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}
	payload.UserID = user_id

	err = bc.BM.CreateAccount(payload)
	if err != nil {
		return
	}
	ctx.Header("HX-Redirect", "/accounts")
}

func (bc *BasicController) GetAccountBalance(ctx *gin.Context) {
	accountID := ctx.DefaultQuery("Account", "")
	if accountID == "" {
		ctx.JSON(400, gin.H{"error": "Account ID is missing"})
		return
	}

	account, err := bc.BM.FindAccountById(accountID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Account not found"})
		return
	}

	balanceMessage := utils.FormatCurrency(account.Balance)
	renderedHTML := utils.GetMessageTemplate(balanceMessage)

	ctx.Data(200, "text/html", renderedHTML)
}
