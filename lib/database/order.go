package database

import (
	"project1/config"
	"project1/models"
)

var selectCart []models.OrderResponse

// Fungsi untuk mendapatkan seluruh order user tertentu
func GetAllOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,orders.total_qty, orders.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join addresses on orders.address_id = addresses.id").Joins(
		"join payment_methods on orders.payment_methods_id = payment_methods.id").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ?", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return selectCart, nil
}

// Fungsi untuk mendapatkan seluruh history order user tertentu
func GetHistoryOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,orders.total_qty, orders.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join addresses on orders.address_id = addresses.id").Joins(
		"join payment_methods on orders.payment_methods_id = payment_methods.id").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ? AND orders.order_status = \"done\"", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return selectCart, nil
}

// Fungsi untuk mendapatkan seluruh cancel order user tertentu
func GetCancelOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,orders.total_qty, orders.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join addresses on orders.address_id = addresses.id").Joins(
		"join payment_methods on orders.payment_methods_id = payment_methods.id").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ? AND orders.order_status = \"cancel\"", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return selectCart, nil
}

// Fungsi untuk membuat order baru
func CreateOrder(order models.Orders) (interface{}, error) {
	query := config.DB.Save(&order)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return order, nil
	}
}

// Fungsi untuk membuat order detail
func CreateOrderDetail(orderDetail models.Order_Details) (models.Order_Details, error) {
	query := config.DB.Save(&orderDetail)
	if query.Error != nil {
		return orderDetail, query.Error
	}
	return orderDetail, nil
}

func InsertOrderDetail(updateOrder models.Order_Details, id int) (interface{}, error) {
	var orderDetail models.Order_Details
	query := config.DB.Find(&orderDetail, id)
	if query.Error != nil {
		return nil, query.Error
	}
	updateQuery := config.DB.Model(&orderDetail).Updates(updateOrder)
	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return orderDetail, nil
}

func GetCartItem(id int) (models.CartItem, error) {
	var shoppingCart models.Shopping_Carts
	var cartItem models.CartItem
	query := config.DB.Model(&shoppingCart).Find(&cartItem)
	if query.Error != nil {
		return cartItem, query.Error
	}
	return cartItem, nil
}

func AddQtyPrice(id int) {
	config.DB.Exec("UPDATE orders SET total_price = (SELECT SUM(price) FROM order_details WHERE orders_id =?) WHERE id =?", id, id)
	config.DB.Exec("UPDATE orders SET total_qty = (SELECT SUM(qty) FROM order_details WHERE orders_id =?) WHERE id =?", id, id)
}
