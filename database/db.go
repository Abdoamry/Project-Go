package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("فشل الاتصال بقاعدة البيانات")
	}
	DB = db
}