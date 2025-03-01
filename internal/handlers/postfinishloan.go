package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostFinishLoanHandler struct {
	BC *controllers.BasicController
}

func NewPostFinishLoanHandler(bc *controllers.BasicController) *PostFinishLoanHandler {
	return &PostFinishLoanHandler{BC: bc}
}

func (h *PostFinishLoanHandler) ServeHTTP(c *gin.Context) {
	err := h.BC.FinishLoan(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Loan Finished!", "text": "Loan has been successfully finished.", "icon": "success", "redirect": "/loans"}}`)
	c.Status(200)
}
