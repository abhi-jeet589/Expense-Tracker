package routes

import (
	"github.com/abhi-jeet589/Expense-Tracker/internal/web"
	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(router *gin.Engine) {
	router.GET("/app/transactions", web.ShowTransactionsPage)
}
