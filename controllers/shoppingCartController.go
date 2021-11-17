package controllers

import (
	"encoding/json"
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/middlewares"
	"project1/models"

	"github.com/labstack/echo/v4"
)

// Controller untuk mendapatkan seluruh data shopping carts
func GetShoppingCartsController(c echo.Context) error {
	id := middlewares.CurrentLoginUser(c)
	// Mendapatkan data satu buku menggunakan fungsi GetProduct
	shoppingCart, e := database.GetShoppingCarts(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	if shoppingCart == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("shopping cart is null"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get product by id", shoppingCart))
}

// Controller untuk memasukkan barang baru ke shopping cart
func CreateShoppingCartsController(c echo.Context) error {
	// Mendapatkan data shopping carts baru dari client
	input := models.Shopping_Carts{}
	c.Bind(&input)
	// Menyimpan data barang baru menggunakan fungsi CreateShoppingCarts
	shoppingCart, e := database.CreateShoppingCarts(input)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add cart"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create shopping cart", shoppingCart))
}

// Fungsi untuk memperbaharui satu data product berdasarkan id product
func UpdateShoppingCartsController(c echo.Context) error {
	// Memanggil id user
	id := middlewares.CurrentLoginUser(c)

	// Pengecekan apakah id product memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	getShoppingCart, err := database.GetShoppingCarts(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	getShoppingCartJSON, err := json.Marshal(getShoppingCart)
	if err != nil {
		panic(err)
	}
	var responseGetShoppingCarts models.Shopping_Carts
	json.Unmarshal([]byte(getShoppingCartJSON), &responseGetShoppingCarts)

	if responseGetShoppingCarts.UsersID != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("not allowed to update"))
	}

	// Mendapatkan data product yang akan diperbaharui dari client
	var updatedShoppingCart models.Shopping_Carts
	c.Bind(&updatedShoppingCart)
	// Memperbaharui data menggunakan fungsi UpdateProduct "&"
	shoppingCart, er := database.UpdateShoppingCarts(id, updatedShoppingCart)
	if shoppingCart == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("shopping cart id not found"))
	}
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to update shopping cart"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("update success", shoppingCart))
}
