package controllers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

type BasicController struct {
	BM *managers.BasicManager
}

func NewBasicController(bm *managers.BasicManager) *BasicController {
	return &BasicController{BM: bm}
}

// ACCOUNT methods

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

func (bc *BasicController) DeleteAccountById(ctx *gin.Context) {
	accountId := ctx.PostForm("AccountID")
	if accountId == "" {
		return
	}

	err := bc.BM.DeleteAccountById(accountId)
	if err != nil {
		return
	}

	ctx.Header("HX-Redirect", "/accounts")
}

func (bc *BasicController) GetUserAccounts(userId string) ([]models.Account, error) {
	accounts, err := bc.BM.GetUserAccounts(userId)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (bc *BasicController) CreateAccount(ctx *gin.Context) error {
	var payload models.AccountRequest
	var err error

	payload.Name = ctx.PostForm("Name")
	payload.Description = ctx.PostForm("Description")
	payload.Balance, err = strconv.ParseFloat(ctx.PostForm("Balance"), 64)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// CATEGORY methods

func (bc *BasicController) GetAllCategories() ([]models.Category, error) {
	categories, err := bc.BM.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (bc *BasicController) GetUserTopCategories(id string) ([]models.CategoryWithTotal, error) {
	categories, err := bc.BM.GetUserTopCategories(id)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// LOAN methods

func (bc *BasicController) GetLoans(userId string) ([]models.Loan, error) {
	loans, err := bc.BM.GetLoans(userId)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (bc *BasicController) CreateLoan(ctx *gin.Context) error {
	var payload models.LoanRequest
	var err error

	payload.Description = ctx.PostForm("Description")
	payload.CategoryID = ctx.PostForm("Category")
	payload.ToWhom = ctx.PostForm("Towhom")
	payload.Status = false
	payload.LoanDate = ctx.PostForm("Date")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return err
	}
	payload.TransactionType, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (bc *BasicController) FinishLoan(ctx *gin.Context) error {
	loanID := ctx.PostForm("LoanID")
	if loanID == "" {
		ctx.JSON(400, gin.H{"error": "Loan ID is missing"})
		return fmt.Errorf("loan ID is missing")
	}

	err := bc.BM.FinishLoan(loanID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Loan not found"})
		return fmt.Errorf("loan not found")
	}

	return nil
}

func (bc *BasicController) GetUserActiveLoans(userId string) ([]models.Loan, error) {
	loans, err := bc.BM.GetUserActiveLoans(userId)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (bc *BasicController) DeleteLoanById(ctx *gin.Context) {
	loanId := ctx.PostForm("LoanID")
	if loanId == "" {
		return
	}

	err := bc.BM.DeleteLoanById(loanId)
	if err != nil {
		return
	}
	ctx.Header("HX-Redirect", "/loans")
}

// RECURRING methods

func (bc *BasicController) GetRecurrings(userId string) ([]models.Recurring, error) {
	recurrings, err := bc.BM.GetRecurrings(userId)
	if err != nil {
		return nil, err
	}
	return recurrings, nil
}

func (bc *BasicController) DeleteRecurringById(ctx *gin.Context) {
	recurringId := ctx.PostForm("RecurringID")
	if recurringId == "" {
		return
	}

	err := bc.BM.DeleteRecurringById(recurringId)
	if err != nil {
		return
	}
	ctx.Header("HX-Redirect", "/recurring")
}

func (bc *BasicController) GetUserUpcomingRecurring(userId string) (models.Recurring, error) {
	recurring, err := bc.BM.GetUserUpcomingRecurring(userId)
	if err != nil {
		return models.Recurring{}, err
	}
	return recurring, nil
}

func (bc *BasicController) CreateRecurring(ctx *gin.Context) error {
	var payload models.RecurringRequest
	var err error

	payload.Name = ctx.PostForm("Name")
	payload.CategoryID = ctx.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return err
	}
	payload.Periodicity, err = strconv.Atoi(ctx.PostForm("Periodicity"))
	if err != nil {
		return err
	}
	payload.TransactionType, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return err
	}

	payload.StartDate = ctx.PostForm("StartDate")
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
		return err
	}

	return nil
}

// TRANSACTION methods

func (bc *BasicController) CreateTransaction(ctx *gin.Context) error {
	var payload models.TransactionRequest
	var err error

	payload.Description = ctx.PostForm("Description")
	payload.CategoryID = ctx.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		return err
	}
	payload.Date = ctx.PostForm("Date")
	payload.Type, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (bc *BasicController) GetUserTransactions(userId string) ([]models.Transaction, error) {
	transactions, err := bc.BM.GetUserTransactions(userId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) GetUserLatestSixTransactions(userId string) ([]models.Transaction, error) {
	transactions, err := bc.BM.GetUserLatestSixTransactions(userId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (bc *BasicController) DeleteTransactionById(ctx *gin.Context) {
	transactionId := ctx.PostForm("TransactionID")
	if transactionId == "" {
		return
	}

	err := bc.BM.DeleteTransactionById(transactionId)
	if err != nil {
		return
	}

	accountId := ctx.PostForm("AccountID")

	err = bc.BM.RecalculateAccountBalance(accountId)
	if err != nil {
		return
	}

	ctx.Header("HX-Redirect", "/transaction")
}

// RANDOM methods

func (bc *BasicController) GetUserMonthlyExpenses(userId string) (float64, error) {
	expenses, err := bc.BM.GetUserMonthlyExpenses(userId)
	if err != nil {
		return -1, err
	}
	return expenses, nil
}

func (bc *BasicController) GetUserMonthlyIncome(userId string) (float64, error) {
	income, err := bc.BM.GetUserMonthlyIncome(userId)
	if err != nil {
		return -1, err
	}
	return income, nil
}

func (bc *BasicController) GetUserTotalBalance(userId string) float64 {
	balance := bc.BM.GetUserTotalBalance(userId)
	return balance
}
