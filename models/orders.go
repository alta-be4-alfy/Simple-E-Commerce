package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	Total_Qty         int             `json:"total_qty" form:"total_qty"`
	Total_Price       int             `json:"total_price" form:"total_price"`
	Order_Status      string          `json:"order_status" form:"order_status"`
	UsersID           int             `json:"users_id" form:"users_id"`
	Payment_MethodsID int             `json:"payment_method_id" form:"payment_method_id"`
	AddressID         int             `json:"address_id" form:"address_id"`
	Order_Details     []Order_Details `gorm:"foreignKey:OrdersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderResponse struct {
	ID           uint
	Qty          int
	Total_Price  int
	Order_Status string
	Payment      string
	Street       string
	User_Name    string
	Product_Name string
}
