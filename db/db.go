package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDBConnection() (*gorm.DB, error) {

	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		return nil, fmt.Errorf("SECRET_NAME environment variable not set")
	}

	// Get database credentials from Secrets Manager
	dbSecret, err := getSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database credentials: %v", err)
	}

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// HOST_DB := os.Getenv("HOST_DB")
	// PORT_DB := os.Getenv("PORT_DB")
	// USER_DB := os.Getenv("USER_DB")
	// PASSWORD_DB := os.Getenv("PASSWORD_DB")
	HOST_DB := dbSecret.Host
	PORT_DB := dbSecret.Port
	USER_DB := dbSecret.Username
	PASSWORD_DB := dbSecret.Password

	log.Println(PORT_DB)

	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, USER_DB, PASSWORD_DB, HOST_DB, PORT_DB, "points")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
