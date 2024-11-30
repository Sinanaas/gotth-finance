package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GetTransactionHandler struct {
}

func NewGetTransaction() *GetTransactionHandler {
	return &GetTransactionHandler{}
}

func (h *GetTransactionHandler) ServeHTTP(ctx *gin.Context) {
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
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}

	transactions, err := basicController.GetTransactionWithCategoryName(user_id)
	if err != nil {
		return
	}

	c := templates.Transaction(categories, transactions)

	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
