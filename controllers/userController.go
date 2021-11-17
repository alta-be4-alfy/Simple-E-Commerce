package controllers

import (
	"net/http"
	"project1/lib/database"
	"project1/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

var user models.Users

type M map[string]interface{}

//login users
func LoginUsersController(c echo.Context) error {
	user := models.Users{}
	c.Bind(&user)

	users, err := database.LoginUsers(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "Login failed",
		})
	}
	return c.JSON(http.StatusOK, M{
		"message": "Login success",
		"users":   users,
	})
}

// get all users
func GetAllUsersController(c echo.Context) error {
	users, err := database.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "Bad Request",
		})
	}
	return c.JSON(http.StatusOK, M{
		"message": "Successful Operation",
		"users":   users,
	})
}

//create user by id
func RegisterUsersController(c echo.Context) error {
	c.Bind(&user)
	user, err := database.RegisterUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Failed to create user",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success Create User",
		"users":   user,
	})
}

//get user by id
func GetUsersController(c echo.Context) error {
	c.Bind(&user)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "False Param",
		})
	}
	user, err := database.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "Bad Request",
			"users":   user,
		})
	}

	return c.JSON(http.StatusOK, M{
		"message": "Success get user",
		"users":   user,
	})
}

//delete user by id
func DeleteUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "False Param",
		})
	}
	user, err := database.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "Failed to delete user",
			"users":   user,
		})
	}
	return c.JSON(http.StatusOK, M{
		"message": "Success Delete User",
	})
}

//update user by id
func UpdateUserController(c echo.Context) error {
	c.Bind(&user)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "False Param",
		})
	}
	user, err := database.UpdateUser(id, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, M{
			"message": "Bad Request",
			"users":   user,
		})
	}
	return c.JSON(http.StatusOK, M{
		"message": "Success get user",
		"users":   user,
	})
}
