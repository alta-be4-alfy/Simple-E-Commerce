package database

import (
	"project1/config"
	"project1/models"
)

// Fungsi untuk memasukkan id user pada product yang baru dibuat
func CreateProduct(id int) (int, error) {
	productUser := models.Products{
		UsersID: id,
	}
	query := config.DB.Select("UsersID").Create(&productUser)
	if query.Error != nil {
		return 0, query.Error
	}
	return int(productUser.ID), nil
}

// Fungsi untuk mendapatkan seluruh data products
func GetProducts() (interface{}, error) {
	var products []models.ProductResponse
	query := config.DB.Table("products").Select(
		"products.id, products.product_name,products.product_type, products.product_stock, products.product_price, products.product_description, users.user_name").Joins(
		"join users on users.id = products.users_id").Where("products.deleted_at IS NULL").Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return products, nil
}

// Fungsi untuk mendapatkan satu data product berdasarkan id product
func GetProduct(id int) (interface{}, error) {
	var product models.ProductResponse
	query := config.DB.Table("products").Select(
		"products.id, products.product_name,products.product_type, products.product_stock, products.product_price, products.product_description, users.user_name").Joins(
		"join users on users.id = products.users_id").Where("products.id = ? AND products.deleted_at IS NULL", id).Find(&product)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return product, nil
}

// Fungsi untuk mendapatkan user id pemilik product
func GetProductOwner(id int) (int, error) {
	var product models.Products
	query := config.DB.Find(&product, id)
	if query.Error != nil {
		return 0, nil
	}
	return product.UsersID, nil
}

// Fungsi untuk mendapatkan seluruh data product product tertentu berdasarkan id user
func GetUserProducts(id int) (interface{}, error) {
	var products []models.ProductResponse
	query := config.DB.Table("products").Select(
		"products.id, products.product_name,products.product_type, products.product_stock, products.product_price, products.product_description, users.user_name").Joins(
		"join users on users.id = products.users_id").Where("users.id = ? AND products.deleted_at IS NULL", id).Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
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
		return 0, nil
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
		return 0, nil
	}
	return "deleted", nil
}
