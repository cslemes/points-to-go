package models

type Product struct {
	UUID                   string `gorm:"primaryKey"`
	DescProduct            string
	DescProductDescription string
	DescProductCategory    string
}
