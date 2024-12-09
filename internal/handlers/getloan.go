package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GetLoanHandler struct{}

func NewGetLoanHandler() *GetLoanHandler {
	return &GetLoanHandler{}
}

func (h *GetLoanHandler) ServeHTTP(ctx *gin.Context) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	basicManager := managers.NewBasicManager(initializers.DB)
	basicController := controllers.NewBasicController(basicManager)

	cookie, _ := ctx.Cookie("access_token")
	categories, err := basicController.GetAllCategories()
	if err != nil {
		return
	}

	session := sessions.Default(ctx)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}

	loans, err := basicController.GetLoanWithCategoryName(userId)
	if err != nil {
		return
	}

	accounts, err := basicController.GetUserAccounts(userId)
	if err != nil {
		return
	}

	var transactionType constants.TransactionType
	transactionTypeArray := transactionType.ToArrayString()

	c := templates.Loans(categories, loans, accounts, transactionTypeArray)

	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
