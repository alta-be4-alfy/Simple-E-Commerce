package database

import (
	"project1/config"
	"project1/models"
)

var selectCart []models.OrderResponse

func GetAllOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select("*").Joins("join order_details on orders.id = order_details.orders_id").Joins("join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Where("shopping_carts.users_id = ?", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	// if query.RowsAffected == 0 {
	// 	return 0, query.Error
	// }
	return selectCart, nil
}

func GetHistoryOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select("*").Joins("join order_details on orders.id = order_details.orders_id").Joins("join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Where("orders.order_status = \"done\" AND shopping_carts.users_id = ? ", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	// if query.RowsAffected == 0 {
	// 	return 0, query.Error
	// }
	return selectCart, nil
}

func GetCancelOrder(id int) (interface{}, error) {
	query := config.DB.Table("orders").Select("*").Joins("join order_details on orders.id = order_details.orders_id").Joins("join shopping_carts on order_details.shopping_carts_id = shopping_carts.id").Where("orders.order_status = \"cancel\" AND shopping_carts.users_id = ? ", id).Find(&selectCart)
	if query.Error != nil {
		return nil, query.Error
	}
	// if query.RowsAffected == 0 {
	// 	return 0, query.Error
	// }
	return selectCart, nil
}

func CreateOrder(order models.Orders) (interface{}, error) {
	query := config.DB.Save(&order)
	if query.Error != nil {
		return nil, query.Error
	} else {
		return order, nil
	}
}
