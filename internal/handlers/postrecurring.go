package handlers

import (
	"fmt"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostRecurringHandler struct {
	BM *managers.BasicManager
}

func NewPostRecurringHandler(bm *managers.BasicManager) *PostRecurringHandler {
	return &PostRecurringHandler{BM: bm}
}

func (h *PostRecurringHandler) ServeHTTP(c *gin.Context) {
	var payload models.RecurringRequest
	var err error

	payload.Name = c.PostForm("Name")
	payload.CategoryID = c.PostForm("Category")
	payload.Amount, err = strconv.ParseFloat(c.PostForm("Amount"), 64)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/recurring"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.Periodicity, err = strconv.Atoi(c.PostForm("Periodicity"))
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/recurring"}}`, err.Error()))
		c.Status(400)
		return
	}
	payload.TransactionType, err = strconv.Atoi(c.PostForm("Type"))
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/recurring"}}`, err.Error()))
		c.Status(400)
		return
	}

	payload.StartDate = c.PostForm("StartDate")
	payload.AccountID = c.PostForm("Account")
	session := sessions.Default(c)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = h.BM.CreateRecurring(payload)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/recurring"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Recurring Created!", "text": "Recurring has been successfully created.", "icon": "success", "redirect": "/recurring"}}`)
	c.Status(200)
}
