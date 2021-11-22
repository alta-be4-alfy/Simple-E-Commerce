package database

import (
	"project1/config"
	"project1/models"
)

var selectCart []models.OrderResponse

// Fungsi untuk mendapatkan seluruh order user tertentu
func GetOrderUserId(idOrder int) (int, error) {
	var orders models.Orders
	query := config.DB.Find(&orders, idOrder)
	if query.Error != nil {
		return 0, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return int(orders.ID), nil
}

// Fungsi untuk mendapatkan seluruh order user tertentu
func GetAllOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,order_details.qty, order_details.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
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
		return 0, nil
	}
	return selectCart, nil
}

// // Fungsi untuk mendapatkan seluruh order user tertentu
// func GetHistoryOrders(idUser int) (interface{}, error) {
// 	var orders []models.OrderResponse
// 	query := config.DB.Where("order_status = \"done\"").Find(&orders, idUser)
// 	if query.Error != nil {
// 		return nil, query.Error
// 	}
// 	if query.RowsAffected == 0 {
// 		return 0, nil
// 	}
// 	return orders, nil
// }

// // Fungsi untuk mendapatkan seluruh order user tertentu
// func GetCancelOrders(idUser int) (interface{}, error) {
// 	var orders []models.OrderResponse
// 	query := config.DB.Where("order_status = \"cancel\"").Find(&orders, idUser)
// 	if query.Error != nil {
// 		return nil, query.Error
// 	}
// 	if query.RowsAffected == 0 {
// 		return 0, nil
// 	}
// 	return orders, nil
// }

// Fungsi untuk mendapatkan seluruh order user tertentu
func GetOrderDetail(idOrderDetail int) (interface{}, error) {
	var orderDetail models.Order_DetailsResponse
	query := config.DB.Table("order_details").Select(
		"order_details.id,order_details.orders_id, order_details.shopping_carts_id, order_details.qty,order_details.total_price,products.product_name, users.user_name").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Joins(
		"join users on shopping_carts.users_id = users.id").Where("order_details.id = ?", idOrderDetail).Find(&orderDetail)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return orderDetail, nil
}

// Fungsi untuk mendapatkan seluruh history order user tertentu
func GetHistoryOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,order_details.qty, order_details.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
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
		return 0, nil
	}
	return selectCart, nil
}

// Fungsi untuk mendapatkan seluruh cancel order user tertentu
func GetCancelOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"orders.id,order_details.qty, order_details.total_price,orders.order_status, payment_methods.payment, addresses.street, users.user_name, products.product_name").Joins(
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
		return 0, nil
	}
	return selectCart, nil
}

// Fungsi untuk membuat order baru
func CreateOrder(order models.Orders) (models.Orders, error) {
	query := config.DB.Save(&order)
	if query.Error != nil {
		return order, query.Error
	} else {
		return order, nil
	}
}

// Fungsi untuk membuat order detail baru
func CreateOrderDetail(orderDetail models.Order_Details) (models.Order_Details, error) {
	query := config.DB.Save(&orderDetail)
	if query.Error != nil {
		return orderDetail, query.Error
	}
	return orderDetail, nil
}

func ChangeOrderStatus(idOrder int, orderStatus string) (interface{}, error) {
	var order models.Orders
	query := config.DB.Find(&order, idOrder)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	order.Order_Status = orderStatus
	updateQuery := config.DB.Save(&order)
	if updateQuery.Error != nil {
		return nil, query.Error
	}
	return order, nil
}

func AddQtyPricetoOrder(id int) {
	config.DB.Exec("UPDATE orders SET total_price = (SELECT SUM(order_details.total_price) FROM order_details WHERE order_details.orders_id =?) WHERE id =?", id, id)
	config.DB.Exec("UPDATE orders SET total_qty = (SELECT SUM(order_details.qty) FROM order_details WHERE order_details.orders_id =?) WHERE id =?", id, id)
}

func AddQtyPricetoOrderDetail(id int) {
	config.DB.Exec("UPDATE order_details SET qty = (SELECT qty FROM shopping_carts WHERE id = ?) WHERE order_details.shopping_carts_id = ?", id, id)
	config.DB.Exec("UPDATE order_details SET total_price = (SELECT qty*price FROM shopping_carts WHERE id = ?) WHERE order_details.shopping_carts_id = ?", id, id)
}

// func InsertOrderDetailtoOrder(idOrder int) {
// 	config.DB.Exec(`"UPDATE order SET order_details = (SELECT (order_details.id,order_details.orders_id, order_details.shopping_carts_id, order_details.qty,order_details.total_price,products.product_name, users.user_name)
// 	FROM order_details JOIN shopping_carts ON order_details.shopping_carts_id = shopping_carts.id JOIN products ON shopping_carts.products_id = products.id
// 	JOIN users on shopping_carts.users_id = users.id WHERE order_details.orders_id = ?) WHERE id = ?"`, idOrder, idOrder)
// }
