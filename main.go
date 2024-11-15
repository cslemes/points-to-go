package main

import (
	"flag"
	"log"
	"points/db"
	"points/handlers"
	"points/models"

	"github.com/gin-gonic/gin"
)

func main() {

	migration := flag.Bool("migrations", false, "Realizar migrations do banco de dados")
	flag.Parse()

	db, err := db.OpenDBConnection()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	if *migration {
		log.Println("Executando migrations...")
		err := db.AutoMigrate(&models.Customer{}, &models.Transaction{}, &models.TransactionProduct{}, &models.Product{})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("ok")
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	r.GET("customers/:id", handlers.GetCustomerByID)
	r.GET("customers/", handlers.GetCustomers)
	r.POST("customers/", handlers.PostCustomer)
	r.PUT("customers/:id", handlers.PutCustomer)

	r.POST("/transactions", handlers.PostTransaction)

	if err := r.Run("0.0.0.0:8081"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

}
