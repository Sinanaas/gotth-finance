package handlers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostLoginHandler struct {
	AC *controllers.AuthController
}

func NewPostLoginHandler(ac *controllers.AuthController) *PostLoginHandler {
	return &PostLoginHandler{AC: ac}
}

func (h *PostLoginHandler) ServeHTTP(c *gin.Context) {
	err := h.AC.SignIn(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/login"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Header("HX-Redirect", "/")
	c.Status(200)
}
