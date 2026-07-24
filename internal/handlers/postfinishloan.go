package handlers

import (
	"fmt"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type PostFinishLoanHandler struct {
	BM *managers.BasicManager
}

func NewPostFinishLoanHandler(bm *managers.BasicManager) *PostFinishLoanHandler {
	return &PostFinishLoanHandler{BM: bm}
}

func (h *PostFinishLoanHandler) ServeHTTP(c *gin.Context) {
	loanID := c.PostForm("LoanID")
	if loanID == "" {
		c.JSON(400, gin.H{"error": "Loan ID is missing"})
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, "Loan ID is missing"))
		c.Status(400)
		return
	}

	err := h.BM.FinishLoan(loanID)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Loan Finished!", "text": "Loan has been successfully finished.", "icon": "success", "redirect": "/loans"}}`)
	c.Status(200)
}
