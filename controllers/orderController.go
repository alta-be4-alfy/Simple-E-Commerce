package controllers

import (
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/middlewares"
	"project1/models"

	"github.com/labstack/echo/v4"
)

func GetAllOrderController(c echo.Context) error {
	// Mendapatkan id user dari token
	idUser := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh order user tertentu menggunakan fungsi GetUserorder
	order, e := database.GetAllOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order"))
	}
	// if order == 0 {
	// 	return c.JSON(http.StatusBadRequest, responses.StatusFailed("user id not found"))
	// }
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get order by user id", order))
}

func GetHistoryOrderController(c echo.Context) error {
	// Mendapatkan id user dari token
	idUser := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh order user tertentu menggunakan fungsi GetUserorder
	order, e := database.GetHistoryOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch history order"))
	}
	// if order == 0 {
	// 	return c.JSON(http.StatusBadRequest, responses.StatusFailed("user id not found"))
	// }
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get order by user id", order))
}

func GetCancelOrderController(c echo.Context) error {
	// Mendapatkan id user dari token
	idUser := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh order user tertentu menggunakan fungsi GetUserorder
	order, e := database.GetCancelOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch cancel order"))
	}
	// if order == 0 {
	// 	return c.JSON(http.StatusBadRequest, responses.StatusFailed("user id not found"))
	// }
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get cancel order by user id", order))
}

// Controller untuk memasukkan barang baru ke shopping cart
func CreateOrderController(c echo.Context) error {
	// Mendapatkan data shopping carts baru dari client
	input := models.Orders{}
	c.Bind(&input)
	// Menyimpan data barang baru menggunakan fungsi CreateOrder
	order, e := database.CreateOrder(input)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add order"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create order", order))
}
