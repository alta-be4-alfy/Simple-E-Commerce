package models

import "gorm.io/gorm"

type Shopping_Carts struct {
	gorm.Model
	Qty           int           `json:"qty" form:"qty"`
	Price         int           `json:"price" form:"price"`
	ProductsID    int           `json:"product_id" form:"product_id"`
	UsersID       int           `json:"users_id" form:"users_id"`
	Order_Details Order_Details `gorm:"foreignKey:Shopping_CartsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
