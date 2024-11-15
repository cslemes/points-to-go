package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDBConnection() (*gorm.DB, error) {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	HOST_DB := os.Getenv("HOST_DB")
	PORT_DB := os.Getenv("PORT_DB")
	USER_DB := os.Getenv("USER_DB")
	PASSWORD_DB := os.Getenv("PASSWORD_DB")

	log.Println(PORT_DB)
	// UNIX dsn
	//dsn := "%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, USER_DB, PASSWORD_DB, HOST_DB, PORT_DB, "points")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
