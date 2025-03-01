package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostLoanHandler struct {
	BC *controllers.BasicController
}

func NewPostLoanHandler(bc *controllers.BasicController) *PostLoanHandler {
	return &PostLoanHandler{BC: bc}
}

func (h *PostLoanHandler) ServeHTTP(c *gin.Context) {
	err := h.BC.CreateLoan(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Loan Created!", "text": "Loan has been successfully created.", "icon": "success", "redirect": "/loans"}}`)
	c.Status(200)
}
