package handlers

import (
	"fmt"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostLoanHandler struct {
	BM *managers.BasicManager
}

func NewPostLoanHandler(bm *managers.BasicManager) *PostLoanHandler {
	return &PostLoanHandler{BM: bm}
}

func (h *PostLoanHandler) ServeHTTP(c *gin.Context) {
	var payload models.LoanRequest
	var err error

	payload.Description = c.PostForm("Description")
	payload.CategoryID = c.PostForm("Category")
	payload.ToWhom = c.PostForm("Towhom")
	payload.Status = false
	payload.LoanDate = c.PostForm("Date")
	payload.Amount, err = strconv.ParseFloat(c.PostForm("Amount"), 64)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.TransactionType, err = strconv.Atoi(c.PostForm("Type"))
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.AccountID = c.PostForm("Account")
	session := sessions.Default(c)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = h.BM.CreateLoan(payload)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/loans"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Loan Created!", "text": "Loan has been successfully created.", "icon": "success", "redirect": "/loans"}}`)
	c.Status(200)
}
