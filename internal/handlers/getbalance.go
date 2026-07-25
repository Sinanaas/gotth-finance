package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetBalanceHandler struct {
	BM *managers.BasicManager
}

func NewGetBalanceHandler(bm *managers.BasicManager) *GetBalanceHandler {
	return &GetBalanceHandler{BM: bm}
}

func (h *GetBalanceHandler) ServeHTTP(c *gin.Context) {
	accountID := c.DefaultQuery("Account", "")
	if accountID == "" {
		c.JSON(400, gin.H{"error": "Account ID is missing"})
		return
	}

	account, err := h.BM.FindAccountById(accountID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Account not found"})
		return
	}

	balanceMessage := utils.FormatCurrency(account.Balance)
	renderedHTML := utils.GetMessageTemplate(balanceMessage)

	c.Data(200, "text/html", renderedHTML)
}
