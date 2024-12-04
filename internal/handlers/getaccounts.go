package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

type GetAccountsHandler struct{}

func NewGetAccountsHandler(ctx *gin.Context) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	basicManager := managers.NewBasicManager(initializers.DB)
	basicController := controllers.NewBasicController(basicManager)

	session := sessions.Default(ctx)
	var user_id string
	v := session.Get("user_id")
	if v != nil {
		user_id = v.(string)
	}

	accounts, err := basicController.GetUserAccounts(user_id)
	c := templates.Accounts(accounts)

	cookie, _ := ctx.Cookie("access_token")
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
