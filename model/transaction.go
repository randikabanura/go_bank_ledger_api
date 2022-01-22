package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:char(36);primary_key;"`
	Amount          float64   `gorm:"default:0.0;not_null" json:"amount"`
	TransactionType string    `gorm:"not_null" json:"transaction_type"`
	Notes           string    `json:"notes"`
	AccountID       uuid.UUID `gorm:"not_null" json:"account_id"`
	Account         Account   `json:"account"`
	ToAccountID     uuid.UUID `json:"to_account_id"`
	ToAccount       Account   `gorm:"foreignkey:ToAccountID" json:"to_account"`
	CustomerID      uuid.UUID `json:"customer_id"`
	Customer        Customer  `json:"customer"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (transaction *Transaction) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}
