package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostAccountHandler struct {
	BC *controllers.BasicController
}

func NewPostAccountHandler(bc *controllers.BasicController) *PostAccountHandler {
	return &PostAccountHandler{BC: bc}
}

func (h *PostAccountHandler) ServeHTTP(c *gin.Context) {
	err := h.BC.CreateAccount(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/accounts"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Account Created!", "text": "Account has been successfully created.", "icon": "success", "redirect": "/accounts"}}`)
	c.Status(200)
}
