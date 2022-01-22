package model

import "github.com/jinzhu/gorm"

func DBMigrate(db *gorm.DB) *gorm.DB {
	//db.DropTableIfExists(&Account{}, &Customer{}, &Transaction{})
	db.DropTableIfExists(&Transaction{})
	db.AutoMigrate(&Account{}, &Customer{}, &Transaction{})
	return db
}
