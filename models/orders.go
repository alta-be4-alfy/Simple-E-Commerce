package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	Total_Qty         int             `json:"total_qty" form:"total_qty"`
	Total_Price       int             `json:"total_price" form:"total_price"`
	Order_Status      string          `json:"order_status" form:"order_status"`
	UsersID           int             `json:"user_id" form:"user_id"`
	Payment_MethodsID int             `json:"payment_methodid" form:"payment_methodid"`
	AddressID         int             `json:"address_id" form:"address_id"`
	Order_Details     []Order_Details `gorm:"foreignKey:OrdersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderResponse struct {
	Total_Qty         int
	Total_Price       int
	Order_Status      string
	Payment_MethodsID int
	AddressID         int
	Qty               int
	Price             int
	ProductsID        int
	UsersID           int
}
