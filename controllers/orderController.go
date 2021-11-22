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
	// Mendapatkan data seluruh order user tertentu menggunakan fungsi GetAllOrder
	order, e := database.GetAllOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order"))
	}
	if order == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("order not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get order by user id", order))
}

func GetHistoryOrderController(c echo.Context) error {
	// Mendapatkan id user dari token
	idUser := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh history order user tertentu menggunakan fungsi GetHistoryOrder
	order, e := database.GetHistoryOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch history order"))
	}
	if order == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("order not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get order by user id", order))
}

func GetCancelOrderController(c echo.Context) error {
	// Mendapatkan id user dari token
	idUser := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh order yang di cancel user tertentu menggunakan fungsi GetCancelOrder
	order, e := database.GetCancelOrder(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch cancel order"))
	}
	if order == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("order not found"))
	}
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

func CreateOrderDetailController(c echo.Context) error {
	// Mendapatkan data order id dan shopping id dari client
	input := models.Order_Details{}
	c.Bind(&input)
	// Memasukkan data ke order detail
	orderDetail, er := database.CreateOrderDetail(input)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to create new order detail"))
	}
	// Input jumlah qty dan jumlah harga order id tertentu ke dalam tabel orders
	database.AddQtyPricetoOrderDetail(input.Shopping_CartsID)
	database.AddQtyPricetoOrder(input.OrdersID)

	order, er := database.GetOrderDetail(int(orderDetail.ID))
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order detail"))
	}
	if order == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order detail"))
	}

	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create new order", order))
}
