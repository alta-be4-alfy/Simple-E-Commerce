package models

import "gorm.io/gorm"

type Products struct {
	gorm.Model
	Product_Name        string           `json:"product_name" form:"product_name"`
	Product_Type        string           `json:"product_type" form:"product_type"`
	Product_Stock       int              `json:"product_stock" form:"product_stock"`
	Product_Price       int              `json:"product_price" form:"product_price"`
	Product_Description string           `json:"product_description" form:"product_description"`
	UsersID             int              `json:"user_id" form:"user_id"`
	Shopping_Carts      []Shopping_Carts `gorm:"foreignKey:ProductsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
