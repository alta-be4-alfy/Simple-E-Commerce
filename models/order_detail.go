package models

import "gorm.io/gorm"

type Order_Details struct {
	gorm.Model
	OrdersID         int `json:"order_id" form:"order_id"`
	Shopping_CartsID int `gorm:"unique" json:"shopping_cartsid" form:"shopping_cartsid"`
	Qty              int `json:"qty" form:"qty"`
	Total_Price      int `json:"total_price" form:"total_price"`
}

type GetOrderDetailResponse struct {
	ID           uint
	OrdersID     uint
	Order_Status string
	Qty          int
	Total_Price  int
	User_Name    string
	ProductsID   uint
	Product_Name string
}
