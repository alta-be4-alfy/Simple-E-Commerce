package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	Total_Qty         int             `json:"total_qty" form:"total_qty"`
	Total_Price       int             `json:"total_price" form:"total_price"`
	Order_Status      string          `gorm:"default:pending;type:enum('done','cancel','pending');" json:"order_status" form:"order_status"`
	UsersID           int             `json:"users_id" form:"users_id"`
	Payment_MethodsID int             `json:"payment_method_id" form:"payment_method_id"`
	AddressID         int             `json:"address_id" form:"address_id"`
	Order_Details     []Order_Details `gorm:"foreignKey:OrdersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderBody struct {
	Shopping_CartsID []int           `gorm:"unique" json:"shopping_cartsid" form:"shopping_cartsid"`
	Address          Address         `json:"address_id" form:"shopping_cartsid"`
	Payment_Methods  Payment_Methods `json:"payment_method_id" form:"payment_method_id"`
	UsersID          int             `json:"users_id" form:"users_id"`
}

type OrderStatusBody struct {
	OrdersID     int    `json:"orders_id" form:"orders_id"`
	Order_Status string `gorm:"default:pending;type:enum('done','cancel','pending');" json:"order_status" form:"order_status"`
}
