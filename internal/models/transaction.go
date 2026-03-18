package models

import (
	"gorm.io/gorm"
)

type TransactionType string

const (
	DEBIT  TransactionType = "DEBIT"
	CREDIT TransactionType = "CREDIT"
)

type Transaction struct {
	gorm.Model

	Name          string
	Slug          string `gorm:"<-:create;uniqueIndex;size:8"`
	AmountInCents uint64
	Type          TransactionType `gorm:"size:6"`
}
