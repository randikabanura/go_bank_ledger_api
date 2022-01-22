package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Account struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key;"`
	NickName   string    `json:"nick_name"`
	Amount     float64   `gorm:"not null" json:"amount"`
	CustomerID uuid.UUID `gorm:"type:char(36);not_null" json:"customer_id"`
	Customer   Customer  `json:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (account *Account) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}
