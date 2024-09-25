package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	UUID             string    `json:"uuid" gorm:"primaryKey"`
	DescCustomerName string    `json:"customer_name"`
	CodCPF           *string   `json:"cpf" gorm:"unique"`
	DescEmail        string    `json:"email" binding:"required,email" gorm:"unique"`
	IdTwitch         *string   `json:"twitch" gorm:"unique"`
	IdYouTube        *string   `json:"youtube" gorm:"unique"`
	IdBlueSky        *string   `json:"bluesky" gorm:"unique"`
	IdInstagram      *string   `json:"instagram" gorm:"unique"`
	NrPoints         int64     `json:"points"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime:true"`
}

func NewCustomer() *Customer {

	newUUID := uuid.New()
	strUUID := newUUID.String()

	return &Customer{
		UUID: strUUID,
	}

}

func (c *Customer) Create(db *gorm.DB) error {

	res := db.Create(c)
	if res.Error != nil {
		res.Rollback()
		return res.Error
	}

	defer res.Commit()
	return nil

}

func GetCustomer(id, cpf, email, twitch, youtube, bsky, instagram string, db *gorm.DB) []Customer {

	var customers []Customer

	query := db.Model(&Customer{})

	if id != "" {
		query = query.Where("uuid = ?", id)
	}

	if cpf != "" {
		query = query.Where("cod_cpf = ?", cpf)
	}

	if email != "" {
		query = query.Where("desc_email = ?", email)
	}

	if twitch != "" {
		query = query.Where("id_twitch = ?", twitch)
	}

	if youtube != "" {
		query = query.Where("id_youtube = ?", youtube)
	}

	if bsky != "" {
		query = query.Where("id_blue_sky = ?", bsky)
	}

	if instagram != "" {
		query = query.Where("id_instagram = ?", instagram)
	}

	query.Find(&customers)
	return customers

}
