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

func CreateOrderController(c echo.Context) error {
	// Mendapatkan data order id dan shopping id dari client
	input := models.OrderBody{}
	orderDetailBody := models.Order_Details{}
	c.Bind(&input)

	// Membuat alamat baru
	idAddress, er := database.CreateAddress(input.Address)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to input address"))
	}
	// Membuat payment baru
	idPayment, er := database.CreatePayment(input.Payment_Methods)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to input payment methods"))
	}

	// Membuat order baru
	orderBody := models.Orders{
		UsersID:           input.UsersID,
		Payment_MethodsID: int(idPayment),
		AddressID:         int(idAddress),
	}
	order, e := database.CreateOrder(orderBody)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add order"))
	}

	for _, cart := range input.Shopping_CartsID {
		orderDetailBody = models.Order_Details{
			OrdersID:         int(order.ID),
			Shopping_CartsID: cart,
		}
		// Membuat order detail baru
		orderDetail, er := database.CreateOrderDetail(orderDetailBody)
		if er != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to create new order detail"))
		}
		// Input jumlah qty dan jumlah harga order id tertentu ke dalam tabel orders
		database.AddQtyPricetoOrderDetail(cart)
		database.AddQtyPricetoOrder(orderDetail.OrdersID)
	}

	return c.JSON(http.StatusOK, responses.StatusSuccess("success to create new order"))
}

func ChangeOrderStatusController(c echo.Context) error {
	var input models.OrderStatusBody
	c.Bind(&input)

	// Pengecekan apakah order yang ingin diubah statusnya merupakan order user yang sedang login
	idToken := middlewares.ExtractTokenUserId(c)
	idUser, _ := database.GetOrderUserId(input.OrdersID)

	if idToken != idUser {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("not allowed to update order status"))
	}

	// Mengubah status order
	order, er := database.ChangeOrderStatus(input.OrdersID, input.Order_Status)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to change order status"))
	}
	if order == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("order not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccess("success to update order status"))
}
