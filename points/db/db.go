package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("../data/database.db"), &gorm.Config{})
	return db, err
}
