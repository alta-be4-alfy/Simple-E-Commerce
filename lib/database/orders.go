package database

import (
	"project1/config"
	"project1/models"
)

func CreateOrder(order models.Orders) (interface{}, error) {
	if err := config.DB.Create(&order).Error; err != nil {
		return order, err
	}

	if err := config.DB.First(&order, order).Error; err != nil {
		return order, nil
	}
	query := config.DB.Save(&order)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return order, nil
	}
}
