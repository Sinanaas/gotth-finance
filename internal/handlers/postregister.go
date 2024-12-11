package handlers

import (
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
	h.AC.SignUp(c)
}
