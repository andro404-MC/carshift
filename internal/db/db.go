package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbConf *gorm.Config = &gorm.Config{
		SkipDefaultTransaction: true,
	}
)

func Setup() error {
	var err error
	db, err = gorm.Open(sqlite.Open("app.db"), dbConf)
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
