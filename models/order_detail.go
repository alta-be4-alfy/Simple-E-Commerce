package models

import "gorm.io/gorm"

type Order_Details struct {
	gorm.Model
	OrdersID         int `json:"order_id" form:"order_id"`
	Shopping_CartsID int `gorm:"unique" json:"shopping_cartsid" form:"shopping_cartsid"`
	Qty              int `json:"qty" form:"qty"`
	Total_Price      int `json:"total_price" form:"total_price"`
}

type Order_DetailsResponse struct {
	ID               uint
	OrdersID         int
	Shopping_CartsID int
	Qty              int
	Total_Price      int
	Product_Name     string
	User_Name        string
}
