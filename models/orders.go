package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	Total_Qty         int              `json:"total_qty" form:"total_qty"`
	Total_Price       int              `json:"total_price" form:"total_price"`
	Order_Status      string           `json:"order_status" form:"order_status"`
	UsersID           int              `json:"user_id" form:"user_id"`
	Payment_MethodsID int              `json:"payment_methodid" form:"payment_methodid"`
	AddressID         int              `json:"address_id" form:"address_id"`
	Shopping_CartsID  []Shopping_Carts `gorm:"foreignKey:OrdersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
