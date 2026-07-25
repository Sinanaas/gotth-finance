package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
)

// swalError responds to a failed htmx mutation: a toast fires via HX-Trigger and
// the 422 status stops htmx from swapping, so the page keeps its current state.
func swalError(ctx *gin.Context, text string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"swal:alert": map[string]interface{}{
			"title": "Error",
			"text":  text,
			"icon":  "error",
		},
	})
	ctx.Header("HX-Trigger", string(payload))
	ctx.Status(http.StatusUnprocessableEntity)
}

// renderErrorPage responds to a failed full-page GET with a friendly error page.
func renderErrorPage(ctx *gin.Context, cookie, message string) {
	ctx.Status(http.StatusInternalServerError)
	_ = templates.Layout(templates.ErrorPage(message), cookie).Render(ctx.Request.Context(), ctx.Writer)
}
