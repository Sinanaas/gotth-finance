package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type PostImportTransactionsHandler struct {
	BM *managers.BasicManager
}

func NewPostImportTransactionsHandler(bm *managers.BasicManager) *PostImportTransactionsHandler {
	return &PostImportTransactionsHandler{BM: bm}
}

func (h *PostImportTransactionsHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	file, err := ctx.FormFile("file")
	if err != nil {
		swalError(ctx, "Please choose a CSV file")
		return
	}
	f, err := file.Open()
	if err != nil {
		swalError(ctx, "Could not open the file")
		return
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		swalError(ctx, "That doesn't look like a valid CSV")
		return
	}

	result, err := h.BM.ImportTransactions(userId, records)
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	text := fmt.Sprintf("Imported %d transactions", result.Imported)
	if len(result.Skipped) > 0 {
		text += fmt.Sprintf(", skipped %d (bad rows)", len(result.Skipped))
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"swal:alert": map[string]interface{}{
			"title":    "Import complete",
			"text":     text,
			"icon":     "success",
			"redirect": "/transaction",
		},
	})
	ctx.Header("HX-Trigger", string(payload))
	ctx.Status(200)
}
