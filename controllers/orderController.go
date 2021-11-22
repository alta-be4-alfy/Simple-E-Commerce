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

// // Controller untuk memasukkan barang baru ke shopping cart
// func CreateOrderController(c echo.Context) error {
// 	// Mendapatkan data shopping carts baru dari client
// 	input := models.Orders{}
// 	c.Bind(&input)
// 	// Menyimpan data barang baru menggunakan fungsi CreateOrder
// 	order, e := database.CreateOrder(input)
// 	if e != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add order"))
// 	}
// 	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create order", order))
// }

func CreateOrderController(c echo.Context) error {
	// Mendapatkan data order id dan shopping id dari client
	input := models.OrderBody{}
	// orderBody := models.Orders{}
	// input := models.Order_Details{}
	orderDetailBody := models.Order_Details{}
	c.Bind(&input)
	// Membuat order baru
	if input.OrdersID == 0 {
		orderBody := models.Orders{
			UsersID:           input.UsersID,
			Payment_MethodsID: input.Payment_MethodsID,
			AddressID:         input.AddressID,
			Order_Status:      input.Order_Status,
		}
		order, e := database.CreateOrder(orderBody)
		if e != nil {
			return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add order"))
		}
		orderDetailBody = models.Order_Details{
			OrdersID:         int(order.ID),
			Shopping_CartsID: input.Shopping_CartsID,
		}
	} else {
		orderDetailBody = models.Order_Details{
			OrdersID:         input.OrdersID,
			Shopping_CartsID: input.Shopping_CartsID,
		}
	}

	// Membuat order detail baru
	orderDetail, er := database.CreateOrderDetail(orderDetailBody)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to create new order detail"))
	}
	// Input jumlah qty dan jumlah harga order id tertentu ke dalam tabel orders
	database.AddQtyPricetoOrderDetail(input.Shopping_CartsID)
	database.AddQtyPricetoOrder(orderDetail.OrdersID)
	// database.InsertOrderDetailtoOrder(orderDetail.OrdersID)

	showOrder, er := database.GetOrderDetail(int(orderDetail.ID))
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order detail"))
	}
	if showOrder == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch order detail"))
	}

	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create new order", showOrder))
}

// func ChangeOrderStatusController(c echo.Context) error {
// 	var input models.OrderStatusBody
// 	c.Bind(&input)

// 	idToken := middlewares.ExtractTokenUserId(c)
// 	idUser := database.GetOrderUserId(input.OrdersID)
// 	if idToken != database.GetOrderUserId(input.OrdersID) {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("not allowed to update order status"))
// 	}
// 	order, er := database.ChangeOrderStatus(input.OrdersID, input.Order_Status)
// 	if er != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to change order status"))
// 	}
// 	if order == 0 {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("order not found"))
// 	}
// 	return c.JSON(http.StatusOK, responses.StatusSuccess("success to update order status"))
// }
