package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	User_Name      string           `gorm:"type:varchar(25);unique;not null" json:"user_name" form:"user_name"`
	Name           string           `gorm:"type:varchar(255)" json:"name" form:"name"`
	Email          string           `gorm:"type:varchar(100);unique;not null" json:"email" form:"email"`
	Password       string           `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	Gender         string           `gorm:"type:enum('male','female','other');" json:"gender" form:"gender"`
	Birth          string           `gorm:"type:date" json:"birth" form:"birth"`
	Phone_Number   string           `gorm:"type:varchar(15);unique;not null" json:"phone_number" form:"phone_number"`
	Token          string           `json:"token" form:"token"`
	Orders         []Orders         `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Products       []Products       `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Shopping_Carts []Shopping_Carts `gorm:"foreignKey:UsersID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
