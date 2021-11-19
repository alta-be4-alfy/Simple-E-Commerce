package database

import (
	"project1/config"
	"project1/models"
)

// Fungsi untuk membuat data products baru
func CreateProduct(product models.Products) (interface{}, error) {
	query := config.DB.Save(&product)
	// query := config.DB.Table("products").Select("shopping_carts.qty, shopping_carts.price, products.product_name, users.user_name, users.id").Joins("join products on products.id = shopping_carts.products_id").Joins("join users on users.id = shopping_carts.users_id").Where("shopping_carts.users_id = ?", id).Find(&shoppingCart)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return product, nil
	}
}

// Fungsi untuk mendapatkan seluruh data products
func GetProducts() (interface{}, error) {
	var products []models.Products
	query := config.DB.Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}
	return products, nil
}

// Fungsi untuk mendapatkan satu data product berdasarkan id product
func GetProduct(id int) (interface{}, error) {
	var product models.Products
	query := config.DB.Find(&product, id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return product, nil
}

// Fungsi untuk mendapatkan seluruh data product product tertentu berdasarkan id product
func GetUserProducts(id int) (interface{}, error) {
	var products []models.Products
	query := config.DB.Find(&products, "users_id = ?", id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return products, nil
}

//Fungsi untuk memperbaharui data product berdasarkan id product
func UpdateProduct(id int, updateProduct *models.Products) (interface{}, error) {
	var product models.Products
	query := config.DB.Find(&product, id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	updateQuery := config.DB.Model(&product).Updates(updateProduct)
	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return product, nil
}

// Fungsi untuk menghapus satu data product berdasarkan id product
func DeleteProduct(id int) (interface{}, error) {
	query := config.DB.Delete(&models.Products{}, id)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return "deleted", nil
}
