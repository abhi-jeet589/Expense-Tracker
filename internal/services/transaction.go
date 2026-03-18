package services

import (
	"crypto/rand"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/abhi-jeet589/Expense-Tracker/internal/models"
	"github.com/abhi-jeet589/Expense-Tracker/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	db := c.MustGet("db").(*gorm.DB)
	result := db.Find(&transactions)

	if result.Error != nil {
		transport.Fail(c, http.StatusInternalServerError, "ERR-1", "Something went wrong")
		return
	}

	transport.OK(c, transactions)
}

type TransactionCreateDTO struct {
	Name   string                 `json:"name" binding:"required"`
	Amount string                 `json:"amount" binding:"required"`
	Type   models.TransactionType `json:"type" binding:"required,oneof=DEBIT CREDIT"`
}

func generateNewSlug(length int) string {
	token := rand.Text()
	if len(token) < length {
		return token
	}
	return token[:length]
}

func validateCreateTransactionRequest(request *TransactionCreateDTO) error {
	request.Name = strings.TrimSpace(request.Name)
	request.Amount = strings.TrimSpace(request.Amount)

	if request.Name == "" {
		return errors.New("name is required")
	}

	if request.Amount == "" {
		return errors.New("amount is required")
	}

	switch request.Type {
	case models.DEBIT, models.CREDIT:
		return nil
	default:
		return errors.New("type must be either DEBIT or CREDIT")
	}
}

func parseAmountToCents(raw string) (uint64, error) {
	parts := strings.Split(raw, ".")
	switch len(parts) {
	case 1:
		dollars, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return 0, errors.New("amount must contain digits only")
		}
		if dollars == 0 {
			return 0, errors.New("amount must be greater than zero")
		}
		return dollars * 100, nil
	case 2:
		dollarsPart := parts[0]
		if dollarsPart == "" {
			dollarsPart = "0"
		}

		dollars, err := strconv.ParseUint(dollarsPart, 10, 64)
		if err != nil {
			return 0, errors.New("amount must contain digits only")
		}

		centsPart := parts[1]
		if len(centsPart) > 2 {
			return 0, errors.New("amount can have at most two decimal places")
		}
		if centsPart == "" {
			centsPart = "00"
		}
		if len(centsPart) == 1 {
			centsPart += "0"
		}

		cents, err := strconv.ParseUint(centsPart, 10, 64)
		if err != nil {
			return 0, errors.New("amount must contain digits only")
		}

		total := dollars*100 + cents
		if total == 0 {
			return 0, errors.New("amount must be greater than zero")
		}
		return total, nil
	default:
		return 0, errors.New("amount must be a valid number")
	}
}

func isUniqueSlugViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func CreateTransaction(c *gin.Context) {
	var request TransactionCreateDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		transport.Fail(c, http.StatusBadRequest, "ERR-VALIDATION", err.Error())
		return
	}

	if err := validateCreateTransactionRequest(&request); err != nil {
		transport.Fail(c, http.StatusBadRequest, "ERR-VALIDATION", err.Error())
		return
	}

	amountInCents, err := parseAmountToCents(request.Amount)
	if err != nil {
		transport.Fail(c, http.StatusBadRequest, "ERR-VALIDATION", err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	for range 5 {
		newTransaction := models.Transaction{
			Name:          request.Name,
			AmountInCents: amountInCents,
			Type:          request.Type,
			Slug:          generateNewSlug(8),
		}

		result := db.Create(&newTransaction)

		if result.Error == nil {
			transport.Created(c, newTransaction)
			return
		}

		if isUniqueSlugViolation(result.Error) {
			continue
		}

		transport.Fail(c, http.StatusInternalServerError, "ERR-1", "Something went wrong")
		return
	}

	transport.Fail(c, http.StatusInternalServerError, "ERR-SLUG", "failed to generate a unique transaction reference")
}
