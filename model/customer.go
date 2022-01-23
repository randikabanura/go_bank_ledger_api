package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Customer struct {
	ID           uuid.UUID     `gorm:"type:char(36);primary_key;"`
	FirstName    string        `gorm:"not_null" json:"first_name"`
	LastName     string        `gorm:"not_null" json:"last_name"`
	Email        string        `gorm:"unique;not_null" json:"email"`
	Password     string        `gorm:"size:100;not null;" json:"password"`
	Accounts     []Account     `json:"accounts"`
	Transactions []Transaction `json:"transactions"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (customer *Customer) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (customer *Customer) BeforeSave() error {
	hashedPassword, err := Hash(customer.Password)
	if err != nil {
		return err
	}
	customer.Password = string(hashedPassword)
	return nil
}
