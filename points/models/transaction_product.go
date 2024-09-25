package models

import "github.com/google/uuid"

type TransactionProduct struct {
	UUID          string `json:"transaction_product_id" gorm:"primaryKey"`
	IdTransaction string `json:"transaction_id"`
	CodProduct    string `json:"product_id" binding:"required"`
	QtdeProduct   int64  `json:"product_qtd" binding:"required"`
	VlProduct     int64  `json:"points" binding:"required"`
}

func NewTransactionProduct(
	idTransaction, idProduct string,
	qtdProduct, points int64) *TransactionProduct {

	return &TransactionProduct{
		UUID:          uuid.New().String(),
		IdTransaction: idTransaction,
		CodProduct:    idProduct,
		QtdeProduct:   qtdProduct,
		VlProduct:     points,
	}

}
