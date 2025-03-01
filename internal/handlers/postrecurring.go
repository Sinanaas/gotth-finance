package handlers

import (
	"fmt"
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
	err := h.BC.CreateRecurring(c)
	if err != nil {
		c.Writer.Header().Set("HX-Trigger", fmt.Sprintf(`{"swal:alert": {"title": "Error!", "text": "%s", "icon": "error", "redirect": "/recurring"}}`, err.Error()))
		c.Status(400)
		return
	}
	c.Writer.Header().Set("HX-Trigger", `{"swal:alert": {"title": "Recurring Created!", "text": "Recurring has been successfully created.", "icon": "success", "redirect": "/recurring"}}`)
	c.Status(200)
}
