package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `gorm:"not null" json:"email"`
	Accounts  []Account `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (customer *Customer) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}
