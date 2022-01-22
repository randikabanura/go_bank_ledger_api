package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Account struct {
	ID             uuid.UUID     `gorm:"type:char(36);primary_key;"`
	NickName       string        `json:"nick_name"`
	Amount         float64       `gorm:"default:0.0;not_null" json:"amount"`
	CustomerID     uuid.UUID     `gorm:"type:char(36);not_null" json:"customer_id"`
	Customer       Customer      `json:"customer"`
	Transactions   []Transaction `json:"transactions"`
	ToTransactions []Transaction `gorm:"foreignkey:ToAccountID" json:"to_transactions"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (account *Account) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}
