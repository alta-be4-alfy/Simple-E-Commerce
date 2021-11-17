package models

import "gorm.io/gorm"

type Payment_Methods struct {
	gorm.Model
	Payment             string `json:"payment" form:"payment"`
	Payment_Description string `json:"payment_description" form:"payment_description"`
	Orders              Orders `gorm:"foreignKey:Payment_MethodsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
