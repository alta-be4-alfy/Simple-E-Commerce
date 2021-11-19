package models

import "gorm.io/gorm"

type Order_Details struct {
	gorm.Model
	OrdersID         int `json:"order_id" form:"order_id"`
	Shopping_CartsID int `json:"shopping_cartsid" form:"shopping_cartsid"`
}
