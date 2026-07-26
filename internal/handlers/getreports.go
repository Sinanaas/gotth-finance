package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetReportsHandler struct {
	BM *managers.BasicManager
}

func NewGetReportsHandler(bm *managers.BasicManager) *GetReportsHandler {
	return &GetReportsHandler{BM: bm}
}

func (h *GetReportsHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	breakdown, err := h.BM.CategoryBreakdown(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load category breakdown")
		return
	}
	catNames := make([]string, 0, len(breakdown))
	catTotals := make([]float64, 0, len(breakdown))
	for _, b := range breakdown {
		catNames = append(catNames, b.Name)
		catTotals = append(catTotals, b.Total)
	}

	nwLabels, nwValues, err := h.BM.NetWorthSeries(userId, 6)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load net worth series")
		return
	}

	c := templates.Reports(catNames, catTotals, nwLabels, nwValues)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		renderErrorPage(ctx, cookie, "Failed to render page")
	}
}
