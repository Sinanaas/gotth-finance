package handlers

import (
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
	h.AC.SignIn(c)
}
