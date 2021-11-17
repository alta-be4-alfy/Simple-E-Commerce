package database

import (
	"project1/config"
	"project1/middlewares"
	"project1/models"
)

var users []models.Users
var user models.Users

func GetAllUsers() (interface{}, error) {
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(id int) (interface{}, error) {
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func RegisterUser(user models.Users) (interface{}, error) {
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id int) (interface{}, error) {
	if err := config.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(id int, User models.Users) (models.Users, error) {
	var user models.Users

	if err := config.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	user.User_Name = User.User_Name
	user.Email = User.Email
	user.Password = User.Password
	user.Phone_Number = User.Phone_Number
	user.Gender = User.Gender
	user.Address = User.Address

	if err := config.DB.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func LoginUsers(user *models.Users) (interface{}, error) {
	var err error
	if err = config.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(user).Error; err != nil {
		return nil, err
	}
	user.Token, err = middlewares.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	if err := config.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
