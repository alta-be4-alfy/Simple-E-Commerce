package models

import "gorm.io/gorm"

type Shopping_Carts struct {
	gorm.Model
	Qty        int    `json:"qty" form:"qty"`
	Price      int    `json:"price" form:"price"`
	ProductsID int    `json:"product_id" form:"product_id"`
	Orders     Orders `gorm:"foreignKey:Shopping_CartsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersID    int    `json:"users_id" form:"users_id"`
	OrdersID   int    `json:"orders_id" form:"orders_id"`
}
