package database

import (
	"project1/config"
	"project1/models"
)

// Fungsi untuk membuat data products baru
func CreateProduct(product models.Products) (interface{}, error) {
	query := config.DB.Save(&product)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return product, nil
	}
}
