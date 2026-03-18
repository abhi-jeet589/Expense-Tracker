package routes

import (
	"github.com/abhi-jeet589/Expense-Tracker/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/transactions", services.GetTransactions)
	rg.POST("/transactions", services.CreateTransaction)
}
