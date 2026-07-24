package handlers

import (
	"fmt"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostAccountHandler struct {
	BM *managers.BasicManager
}

func NewPostAccountHandler(bm *managers.BasicManager) *PostAccountHandler {
	return &PostAccountHandler{BM: bm}
}

func (h *PostAccountHandler) ServeHTTP(c *gin.Context) {
	var payload models.AccountRequest
	var err error

	payload.Name = c.PostForm("Name")
	payload.Description = c.PostForm("Description")
	payload.Balance, err = strconv.ParseFloat(c.PostForm("Balance"), 64)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/accounts"}}`, err.Error()))
		c.Status(400)
		return
	}
	session := sessions.Default(c)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	err = h.BM.CreateAccount(payload)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/accounts"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Account Created!", "text": "Account has been successfully created.", "icon": "success", "redirect": "/accounts"}}`)
	c.Status(200)
}
