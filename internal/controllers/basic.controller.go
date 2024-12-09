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
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

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

func (bc *BasicController) GetUserTransactions(userId string) ([]models.Transaction, error) {
	transactions, err := bc.BM.GetUserTransactions(userId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) GetTransactionWithCategoryName(userId string) ([]models.TransactionCategoryAccounts, error) {
	transactions, err := bc.BM.GetTransactionWithCategoryName(userId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) GetLoanWithCategoryName(userId string) ([]models.LoanCategoryAccount, error) {
	loans, err := bc.BM.GetLoanWithCategoryName(userId)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (bc *BasicController) GetRecurringWithCategoryName(userId string) ([]models.RecurringWithCategoryName, error) {
	recurrings, err := bc.BM.GetRecurringWithCategoryName(userId)
	if err != nil {
		return nil, err
	}
	return recurrings, nil
}

func (bc *BasicController) GetCategoryName(categoryId string) (string, error) {
	categoryName, err := bc.BM.GetCategoryName(categoryId)
	if err != nil {
		return "", err
	}
	return categoryName, nil
}

func (bc *BasicController) GetAccountName(accountId string) (string, error) {
	accountName, err := bc.BM.GetAccountName(accountId)
	if err != nil {
		return "", err
	}
	return accountName, nil
}

func (bc *BasicController) GetUserAccounts(userId string) ([]models.Account, error) {
	accounts, err := bc.BM.GetUserAccounts(userId)
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
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

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

func (bc *BasicController) CreateRecurring(ctx *gin.Context) {
	var payload models.RecurringRequest
	var err error

	payload.Name = ctx.PostForm("Name")
	payload.CategoryID = ctx.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return
	}
	payload.Periodicity, err = strconv.Atoi(ctx.PostForm("Periodicity"))
	if err != nil {
		return
	}
	payload.TransactionType, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return
	}
	payload.AccountID = ctx.PostForm("Account")
	session := sessions.Default(ctx)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = bc.BM.CreateRecurring(payload)

	if err != nil {
		return
	}

	ctx.Header("HX-Redirect", "/recurring")
}

func (bc *BasicController) CreateLoan(ctx *gin.Context) {
	var payload models.LoanRequest
	var err error

	payload.Description = ctx.PostForm("Description")
	payload.CategoryID = ctx.PostForm("Category")
	payload.ToWhom = ctx.PostForm("Towhom")
	payload.Status = false
	payload.LoanDate = ctx.PostForm("Date")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return
	}
	payload.TransactionType, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return
	}
	payload.AccountID = ctx.PostForm("Account")
	session := sessions.Default(ctx)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = bc.BM.CreateLoan(payload)

	if err != nil {
		return
	}

	ctx.Header("HX-Redirect", "/loans")
}

func (bc *BasicController) FinishLoan(ctx *gin.Context) {
	loanID := ctx.PostForm("LoanID")
	if loanID == "" {
		ctx.JSON(400, gin.H{"error": "Loan ID is missing"})
		return
	}

	err := bc.BM.FinishLoan(loanID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Loan not found"})
		return
	}

	ctx.Header("HX-Redirect", "/loans")
}
