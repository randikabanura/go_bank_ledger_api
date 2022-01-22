package model

import "github.com/jinzhu/gorm"

func DBMigrate(db *gorm.DB) *gorm.DB {
	//db.DropTableIfExists(&Account{}, &Customer{})
	db.AutoMigrate(&Account{}, &Customer{})
	return db
}
