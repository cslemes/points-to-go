package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	UUID          string    `json:"transaction_id" gorm:"primaryKey"`
	IdCustomer    string    `json:"customer_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	VlPoints      int64     `json:"points" binding:"required"`
	DescSysOrigin string    `json:"system_origin" binding:"required"`
}

func NewTransaction(idCustomer string, points int64, origin string) *Transaction {
	return &Transaction{
		UUID:          uuid.New().String(),
		IdCustomer:    idCustomer,
		VlPoints:      points,
		DescSysOrigin: origin,
	}
}
