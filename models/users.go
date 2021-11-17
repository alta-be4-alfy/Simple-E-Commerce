package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	User_Name    string     `json:"user_name" form:"user_name"`
	Email        string     `json:"email" form:"email"`
	Password     string     `json:"password" form:"password"`
	Gender       string     `json:"gender" form:"gender"`
	Address      string     `json:"address" form:"address"`
	Phone_Number int        `json:"phone_number" form:"phone_number"`
	Token        string     `json:"token" form:"token"`
	Orders       []Orders   `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Products     []Products `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
