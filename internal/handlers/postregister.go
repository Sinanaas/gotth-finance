package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostRegisterHandler struct {
	AC *controllers.AuthController
}

func NewPostRegisterHandler(ac *controllers.AuthController) *PostRegisterHandler {
	return &PostRegisterHandler{AC: ac}
}

func (h *PostRegisterHandler) ServeHTTP(c *gin.Context) {
	err := h.AC.SignUp(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/register"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Account Created!", "text": "Account has been successfully created.", "icon": "success", "redirect": "/login"}}`)
	c.Status(200)
}
