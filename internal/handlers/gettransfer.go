package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetTransferHandler struct {
	BM *managers.BasicManager
}

func NewGetTransferHandler(bm *managers.BasicManager) *GetTransferHandler {
	return &GetTransferHandler{BM: bm}
}

func (h *GetTransferHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	c := templates.TransferModal(accounts)
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render"})
	}
}
