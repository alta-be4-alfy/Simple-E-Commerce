package models

import "gorm.io/gorm"

type Order_Details struct {
	gorm.Model
	OrdersID         int `json:"order_id" form:"order_id"`
	Shopping_CartsID int `json:"shopping_cartsid" form:"shopping_cartsid"`
	Qty              int `json:"qty" form:"qty"`
	Price            int `json:"price" form:"price"`
	ProductsID       int `json:"product_id" form:"product_id"`
	UsersID          int `json:"users_id" form:"users_id"`
}

type CartItem struct {
	Qty        int
	Price      int
	ProductsID int
	UsersID    int
}

type Order_DetailsResponse struct {
	ID               int
	OrdersID         int
	Shopping_CartsID int
	Qty              int
	Price            int
	ProductsID       int
	UsersID          int
}
