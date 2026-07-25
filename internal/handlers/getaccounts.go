package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetAccountsHandler struct {
	BM *managers.BasicManager
}

func NewGetAccountsHandler(bm *managers.BasicManager) *GetAccountsHandler {
	return &GetAccountsHandler{BM: bm}
}

func (h *GetAccountsHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	cookie, _ := ctx.Cookie("access_token")

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load accounts")
		return
	}

	c := templates.Accounts(accounts)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}
}
