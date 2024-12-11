package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostRecurringHandler struct {
	BC *controllers.BasicController
}

func NewPostRecurringHandler(bc *controllers.BasicController) *PostRecurringHandler {
	return &PostRecurringHandler{BC: bc}
}

func (h *PostRecurringHandler) ServeHTTP(c *gin.Context) {
	h.BC.CreateRecurring(c)
}
