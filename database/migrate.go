package database

import (
	"bank_ledger_api/model"
	"github.com/jinzhu/gorm"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	//database.DropTableIfExists(&Account{}, &Customer{}, &Transaction{})
	//database.DropTableIfExists(&Transaction{})
	db.AutoMigrate(&model.Account{}, &model.Customer{}, &model.Transaction{})
	return db
}
