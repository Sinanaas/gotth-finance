package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostTransactionHandler struct {
	BC *controllers.BasicController
}

func NewPostTransactionHandler(bc *controllers.BasicController) *PostTransactionHandler {
	return &PostTransactionHandler{BC: bc}
}

func (h *PostTransactionHandler) ServeHTTP(c *gin.Context) {
	err := h.BC.CreateTransaction(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/transaction"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Transaction Created!", "text": "Transaction has been successfully created.", "icon": "success", "redirect": "/transaction"}}`)
	c.Status(200)
}
