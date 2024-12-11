package handlers

import (
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
	h.BC.CreateAccount(c)
}
