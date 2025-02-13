package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() error {
	var err error
	db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&User{}, &Car{})
	if err != nil {
		return err
	}

	log.Println("DB: up and running")
	return nil
}
