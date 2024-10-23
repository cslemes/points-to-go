package handlers

import (
	"net/http"
	"points/models"

	"github.com/gin-gonic/gin"
)

type PayloadTransaction struct {
	models.Transaction
	Products []models.TransactionProduct `json:"products" binding:"required"`
}

func PostTransaction(c *gin.Context) {

	payload := &PayloadTransaction{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	customer := &models.Customer{}
	res := dbConnection.First(&customer, "uuid = ?", payload.IdCustomer)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	t := models.NewTransaction(payload.IdCustomer, payload.VlPoints, payload.DescSysOrigin)
	ps := []models.TransactionProduct{}

	for _, p := range payload.Products {
		ps = append(ps, *models.NewTransactionProduct(t.UUID, p.CodProduct, p.QtdeProduct, p.VlProduct))
	}

	if err := TransactionProcess(customer, t, ps, dbConnection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "transacao criada e pontos do usuário atualizados"})
}
