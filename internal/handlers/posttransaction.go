package handlers

import (
	"fmt"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostTransactionHandler struct {
	BM *managers.BasicManager
}

func NewPostTransactionHandler(bm *managers.BasicManager) *PostTransactionHandler {
	return &PostTransactionHandler{BM: bm}
}

func (h *PostTransactionHandler) ServeHTTP(c *gin.Context) {
	var payload models.TransactionRequest
	var err error

	payload.Description = c.PostForm("Description")
	payload.CategoryID = c.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(c.PostForm("Amount"), 64)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/transaction"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.Date = c.PostForm("Date")
	payload.Type, err = strconv.Atoi(c.PostForm("Type"))
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/transaction"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.Account = c.PostForm("Account")
	session := sessions.Default(c)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = h.BM.CreateTransaction(payload)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/transaction"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Transaction Created!", "text": "Transaction has been successfully created.", "icon": "success", "redirect": "/transaction"}}`)
	c.Status(200)
}
