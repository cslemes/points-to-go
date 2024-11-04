package main

import (
	"flag"
	"log"
	"points/db"
	"points/handlers"
	"points/models"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		db.AutoMigrate(&models.Customer{}, &models.Transaction{}, &models.TransactionProduct{}, &models.Product{})
		log.Println("ok")
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(handlers.MetricsMiddleware) // prometheus

	r.GET("customers/:id", handlers.GetCustomerByID)
	r.GET("customers/", handlers.GetCustomers)
	r.POST("customers/", handlers.PostCustomer)
	r.PUT("customers/:id", handlers.PutCustomer)

	r.POST("/transactions", handlers.PostTransaction)

	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // prometheus

	r.Run("0.0.0.0:8081")



}
