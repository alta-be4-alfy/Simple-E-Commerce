package database

import (
	"project1/config"
	"project1/models"
)

type Select struct {
	ID           int
	User_Name    string
	Product_Name string
	Price        int
	Qty          int
}

// Fungsi untuk menambahkan harga berdasarkan qty
func AddQtyPrice(products_id, id int) {
	config.DB.Table("shopping_carts").Joins("join products on products.id = shopping_carts.products_id")
	config.DB.Exec("UPDATE shopping_carts SET price = (SELECT product_price FROM products WHERE products.id =?) WHERE shopping_carts.id =?", products_id, id)
}

// Fungsi untuk mendapatkan seluruh data shopping carts
func GetShoppingCarts(id int) (interface{}, error) {
	var shoppingCart []Select

	query := config.DB.Table("shopping_carts").Select(
		"shopping_carts.qty, shopping_carts.price, products.product_name, users.user_name, shopping_carts.id").Joins(
		"join products on products.id = shopping_carts.products_id").Joins(
		"join users on users.id = shopping_carts.users_id").Where(
		"shopping_carts.users_id = ? AND shopping_carts.deleted_at is NULL", id).Find(&shoppingCart)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return shoppingCart, nil
}

// Fungsi untuk mendapatkan seluruh data shopping carts
func GetShoppingCartsTanpaJoin(id int) (interface{}, error) {
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
func CreateShoppingCarts(shoppingCart models.Shopping_Carts) (models.Shopping_Carts, error) {
	query := config.DB.Save(&shoppingCart)
	if query.Error != nil {
		return shoppingCart, query.Error
	} else {
		return shoppingCart, nil
	}
}

// Fungsi Update Data Shopping Cart
func UpdateShoppingCarts(id int, updateShoppingCart *models.Shopping_Carts) (interface{}, error) {
	var shoppingCart models.Shopping_Carts
	query := config.DB.Select("shopping_carts.id, users_id, products_id, qty, price").Find(&shoppingCart, id)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected == 0 {
		return 0, query.Error
	}

	updateQuery := config.DB.Model(&shoppingCart).Updates(updateShoppingCart)

	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return shoppingCart, nil
}

// Fungsi untuk menghapus satu data product berdasarkan id product
func DeleteShoppingCart(id int) (interface{}, error) {
	query := config.DB.Delete(&models.Shopping_Carts{}, id)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return "deleted", nil
}
