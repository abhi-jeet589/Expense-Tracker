package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionsPageData struct {
	ListEndpoint   string
	CreateEndpoint string
}

func ShowTransactionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "transactions_page", TransactionsPageData{
		ListEndpoint:   "/api/transactions",
		CreateEndpoint: "/api/transactions",
	})
}
