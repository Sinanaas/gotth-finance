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
	accounts, err := h.BM.GetUserAccounts(userId)
	c := templates.Accounts(accounts)

	cookie, _ := ctx.Cookie("access_token")
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
