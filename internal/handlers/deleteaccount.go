package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type DeleteAccountHandler struct {
	BM *managers.BasicManager
}

func NewDeleteAccountHandler(bm *managers.BasicManager) *DeleteAccountHandler {
	return &DeleteAccountHandler{BM: bm}
}

func (h *DeleteAccountHandler) ServeHTTP(ctx *gin.Context) {
	accountId := ctx.PostForm("AccountID")
	if accountId == "" {
		ctx.JSON(400, gin.H{"error": "Account ID is missing"})
		return
	}

	err := h.BM.DeleteAccountById(accountId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete account"})
		return
	}

	ctx.Header("HX-Redirect", "/accounts")
}
