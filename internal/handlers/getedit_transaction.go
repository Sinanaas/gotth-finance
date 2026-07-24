package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetEditTransactionHandler struct {
	BC *controllers.BasicController
}

func NewGetEditTransactionHandler(bc *controllers.BasicController) *GetEditTransactionHandler {
	return &GetEditTransactionHandler{BC: bc}
}

func (h *GetEditTransactionHandler) ServeHTTP(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := utils.GetSessionUserID(ctx)

	transaction, err := h.BC.FindTransactionById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	categories, err := h.BC.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	c := templates.EditTransactionModal(transaction, categories, accounts)
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render modal"})
	}
}
