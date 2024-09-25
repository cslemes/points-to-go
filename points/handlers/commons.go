package handlers

import (
	"points/db"
	"points/models"

	"gorm.io/gorm"
)

var dbConnection, _ = db.OpenDBConnection()

func TransactionProcess(
	c *models.Customer,
	t *models.Transaction,
	ps []models.TransactionProduct,
	db *gorm.DB) error {

	tx := dbConnection.Begin()
	defer tx.Commit()

	res := tx.Create(t)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	for _, p := range ps {
		res := tx.Create(p)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	c.NrPoints += t.VlPoints
	res = tx.Save(c)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}
