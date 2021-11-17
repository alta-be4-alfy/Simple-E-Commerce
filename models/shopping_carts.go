package models

import "gorm.io/gorm"

type Shopping_Carts struct {
	gorm.Model
	Qty        int  `json:"qty" form:"qty"`
	Price      int  `json:"price" form:"price"`
	ProductsID uint `json:"product_id" form:"product_id"`
	OrdersID   uint `json:"order_id" form:"order_id"`
}
