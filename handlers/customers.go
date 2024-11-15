package handlers

import (
	"net/http"
	"points/models"
	"points/myerrors"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {

	uuid := c.Query("uuid")
	//cod_cpf := c.Query("cod_cpf")
	cod_cpf := "12394588900"
	desc_email := c.Query("desc_email")
	id_twitch := c.Query("id_twitch")
	id_youtube := c.Query("id_youtube")
	id_blue_sky := c.Query("id_blue_sky")
	id_instagram := c.Query("id_instagram")

	customers := models.GetCustomer(
		uuid,
		cod_cpf,
		desc_email,
		id_twitch,
		id_youtube,
		id_blue_sky,
		id_instagram,
		dbConnection,
	)

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "nenhum resultado encontrado"})
		return
	}

	c.JSON(http.StatusOK, customers)

}

func GetCustomerByID(c *gin.Context) {

	var customer models.Customer
	id := c.Param("id")

	result := dbConnection.First(&customer, "uuid = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func PostCustomer(c *gin.Context) {

	customer := models.NewCustomer()
	customerUUID := customer.UUID

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.UUID = customerUUID

	err := myerrors.GetCreateCustomerErrors(customer.Create(dbConnection))
	if err != nil {

		if err.Error() == myerrors.EmailCreateError.Error() {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "solicitação incorreta"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created", "customer": customer})
}

func PutCustomer(c *gin.Context) {

	id := c.Param("id")

	var customer *models.Customer

	res := dbConnection.First(&customer, "uuid = ?", id)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res = dbConnection.Save(customer)
	if res.Error != nil {
		res.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	res.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "cliente atualizado com sucesso"})
}
