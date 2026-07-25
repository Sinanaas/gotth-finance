package routers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/handlers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BasicRouter struct {
	BM *managers.BasicManager
	RG *gin.RouterGroup
}

func NewBasicRouter(bm *managers.BasicManager, rg *gin.RouterGroup) *BasicRouter {
	return &BasicRouter{BM: bm, RG: rg}
}

func (br *BasicRouter) BasicRoute(rg *gin.RouterGroup, bm *managers.BasicManager) {
	rg.GET("/", middleware.DeserializeUser(), handlers.NewGetHomeHandler(bm).ServeHTTP)
	rg.GET("/transaction", middleware.DeserializeUser(), handlers.NewGetTransaction(bm).ServeHTTP)
	rg.POST("/transaction", middleware.DeserializeUser(), handlers.NewPostTransactionHandler(bm).ServeHTTP)
	rg.GET("/accounts", middleware.DeserializeUser(), handlers.NewGetAccountsHandler(bm).ServeHTTP)
	rg.POST("/account", middleware.DeserializeUser(), handlers.NewPostAccountHandler(bm).ServeHTTP)
	rg.GET("/account/balance", middleware.DeserializeUser(), handlers.NewGetBalanceHandler(bm).ServeHTTP)
	rg.GET("/recurring", middleware.DeserializeUser(), handlers.NewGetRecurringHandler(bm).ServeHTTP)
	rg.POST("/recurring", middleware.DeserializeUser(), handlers.NewPostRecurringHandler(bm).ServeHTTP)
	rg.GET("/loans", middleware.DeserializeUser(), handlers.NewGetLoanHandler(bm).ServeHTTP)
	rg.POST("/loan", middleware.DeserializeUser(), handlers.NewPostLoanHandler(bm).ServeHTTP)
	rg.POST("/loan/finish", middleware.DeserializeUser(), handlers.NewPostFinishLoanHandler(bm).ServeHTTP)

	rg.GET("/transaction/list", middleware.DeserializeUser(), handlers.NewGetTransactionListHandler(bm).ServeHTTP)
	rg.GET("/transaction/:id/edit", middleware.DeserializeUser(), handlers.NewGetEditTransactionHandler(bm).ServeHTTP)
	rg.PATCH("/transaction/:id", middleware.DeserializeUser(), handlers.NewPatchTransactionHandler(bm).ServeHTTP)
	rg.PUT("/transaction", middleware.DeserializeUser(), handlers.NewDeleteTransactionHandler(bm).ServeHTTP)
	rg.GET("/accounts/transfer", middleware.DeserializeUser(), handlers.NewGetTransferHandler(bm).ServeHTTP)
	rg.POST("/accounts/transfer", middleware.DeserializeUser(), handlers.NewPostTransferHandler(bm).ServeHTTP)
	rg.PUT("/account", middleware.DeserializeUser(), handlers.NewDeleteAccountHandler(bm).ServeHTTP)
	rg.GET("/account/:id/edit", middleware.DeserializeUser(), handlers.NewGetEditAccountHandler(bm).ServeHTTP)
	rg.PATCH("/account/:id", middleware.DeserializeUser(), handlers.NewPatchAccountHandler(bm).ServeHTTP)
	rg.PUT("/recurring", middleware.DeserializeUser(), handlers.NewDeleteRecurringHandler(bm).ServeHTTP)
	rg.GET("/recurring/:id/edit", middleware.DeserializeUser(), handlers.NewGetEditRecurringHandler(bm).ServeHTTP)
	rg.PATCH("/recurring/:id", middleware.DeserializeUser(), handlers.NewPatchRecurringHandler(bm).ServeHTTP)
	rg.PUT("/loan", middleware.DeserializeUser(), handlers.NewDeleteLoanHandler(bm).ServeHTTP)

	rg.GET("/categories", middleware.DeserializeUser(), handlers.NewGetCategoriesHandler(bm).ServeHTTP)
	rg.POST("/categories", middleware.DeserializeUser(), handlers.NewPostCategoryHandler(bm).ServeHTTP)
	rg.DELETE("/categories/:id", middleware.DeserializeUser(), handlers.NewDeleteCategoryHandler(bm).ServeHTTP)

	rg.GET("/transactions/export", middleware.DeserializeUser(), handlers.NewExportTransactionsHandler(bm).ServeHTTP)
	rg.POST("/transactions/import", middleware.DeserializeUser(), handlers.NewPostImportTransactionsHandler(bm).ServeHTTP)

	rg.GET("/budgets", middleware.DeserializeUser(), handlers.NewGetBudgetsHandler(bm).ServeHTTP)
	rg.POST("/budgets", middleware.DeserializeUser(), handlers.NewPostBudgetHandler(bm).ServeHTTP)
	rg.DELETE("/budgets/:id", middleware.DeserializeUser(), handlers.NewDeleteBudgetHandler(bm).ServeHTTP)

	rg.GET("/goals", middleware.DeserializeUser(), handlers.NewGetGoalsHandler(bm).ServeHTTP)
	rg.POST("/goals", middleware.DeserializeUser(), handlers.NewPostGoalHandler(bm).ServeHTTP)
	rg.DELETE("/goals/:id", middleware.DeserializeUser(), handlers.NewDeleteGoalHandler(bm).ServeHTTP)
	rg.POST("/goals/:id/contribute", middleware.DeserializeUser(), handlers.NewPostGoalContributeHandler(bm).ServeHTTP)
}
