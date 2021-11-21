package models

import "gorm.io/gorm"

type Order_Details struct {
	gorm.Model
	OrdersID         int `json:"order_id" form:"order_id"`
	Shopping_CartsID int `json:"shopping_cartsid" form:"shopping_cartsid"`
}

type Order_DetailsResponse struct {
	ID               uint
	OrdersID         int
	Shopping_CartsID int
	Qty              int
	Price            int
	Product_Name     string
	User_Name        string
}
