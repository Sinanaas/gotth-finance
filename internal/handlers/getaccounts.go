package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetAccountsHandler struct {
	BC *controllers.BasicController
}

func NewGetAccountsHandler(bc *controllers.BasicController) *GetAccountsHandler {
	return &GetAccountsHandler{BC: bc}
}

func (h *GetAccountsHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	accounts, err := h.BC.GetUserAccounts(userId)
	c := templates.Accounts(accounts)

	cookie, _ := ctx.Cookie("access_token")
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
