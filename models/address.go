package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Street   string   `json:"address" form:"address"`
	City     string   `json:"city" form:"city"`
	Province string   `json:"province" form:"province"`
	Zip      int      `json:"zip" form:"zip"`
	Orders   []Orders `gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
