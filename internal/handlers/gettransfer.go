package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetTransferHandler struct {
	BC *controllers.BasicController
}

func NewGetTransferHandler(bc *controllers.BasicController) *GetTransferHandler {
	return &GetTransferHandler{BC: bc}
}

func (h *GetTransferHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	c := templates.TransferModal(accounts)
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render"})
	}
}
