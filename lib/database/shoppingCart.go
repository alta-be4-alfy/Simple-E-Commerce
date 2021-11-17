package database

import (
	"project1/config"
	"project1/models"
)

// Fungsi untuk mendapatkan seluruh data shopping carts
func GetShoppingCarts(id int) (interface{}, error) {
	var shoppingCart models.Shopping_Carts
	query := config.DB.Find(&shoppingCart, id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return shoppingCart, nil
}

// Fungsi untuk membuat data shopping carts
func CreateShoppingCarts(shoppingCart models.Shopping_Carts) (interface{}, error) {
	query := config.DB.Save(&shoppingCart)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return shoppingCart, nil
	}
}

func UpdateShoppingCarts(id int, shoppingCarts models.Shopping_Carts) (interface{}, error) {
	var shoppingCart models.Shopping_Carts
	query := config.DB.Find(&shoppingCart, id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	updateQuery := config.DB.Model(&shoppingCart).Updates(UpdateShoppingCarts)
	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return shoppingCart, nil
}
