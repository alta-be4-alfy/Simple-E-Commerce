package database

import (
	"project1/config"
	"project1/models"
)

var orders []models.GetOrderDetailResponse

// Fungsi untuk mendapatkan seluruh id user order tertentu
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
		"order_details.id, order_details.orders_id, orders.order_status, order_details.qty, order_details.total_price, users.user_name, shopping_carts.products_id, products.product_name").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ?", id).Find(&orders)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return orders, nil
}

// Fungsi untuk mendapatkan seluruh history order user tertentu
func GetHistoryOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"order_details.id, order_details.orders_id, orders.order_status, order_details.qty, order_details.total_price, users.user_name, shopping_carts.products_id, products.product_name").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ? AND orders.order_status = \"done\"", id).Find(&orders)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return orders, nil
}

// Fungsi untuk mendapatkan seluruh cancel order user tertentu
func GetCancelOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select(
		"order_details.id, order_details.orders_id, orders.order_status, order_details.qty, order_details.total_price, users.user_name, shopping_carts.products_id, products.product_name").Joins(
		"join users on orders.users_id = users.id").Joins(
		"join order_details on orders.id = order_details.orders_id").Joins(
		"join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Joins(
		"join products on shopping_carts.products_id = products.id").Where("orders.users_id = ? AND orders.order_status = \"cancel\"", id).Find(&orders)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, nil
	}
	return orders, nil
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

// Fungsi untuk membuat alamat baru
func CreateAddress(address models.Address) (uint, error) {
	query := config.DB.Save(&address)
	if query.Error != nil {
		return 0, query.Error
	}
	return address.ID, nil
}

// Fungsi untuk membuat pembayaran baru
func CreatePayment(payment models.Payment_Methods) (uint, error) {
	query := config.DB.Save(&payment)
	if query.Error != nil {
		return 0, query.Error
	}
	return payment.ID, nil
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
