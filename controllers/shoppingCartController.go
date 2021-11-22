package controllers

import (
	"encoding/json"
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/middlewares"
	"project1/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Controller untuk mendapatkan seluruh data shopping carts
func GetShoppingCartsController(c echo.Context) error {
	id := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data satu buku menggunakan fungsi GetProduct
	shoppingCart, e := database.GetShoppingCarts(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	if shoppingCart == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("shopping cart is null"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get shopping cart by id", shoppingCart))
}

// Controller untuk memasukkan barang baru ke shopping cart
func CreateShoppingCartsController(c echo.Context) error {
	// Mendapatkan data shopping carts baru dari client
	input := models.Shopping_Carts{}
	idToken := middlewares.ExtractTokenUserId(c)
	c.Bind(&input)
	if input.UsersID != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Wrong Users ID"))
	}
	if input.Qty <= 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Must Add product"))
	}
	// Menyimpan data barang baru menggunakan fungsi CreateShoppingCarts
	shoppingCart, e := database.CreateShoppingCarts(input)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to add cart"))
	}
	database.AddQtyPrice(input.ProductsID, int(shoppingCart.ID))
	return c.JSON(http.StatusOK, responses.StatusSuccess("success to create shopping cart"))
}

// Fungsi untuk memperbaharui satu data product berdasarkan id product
func UpdateShoppingCartsController(c echo.Context) error {
	// Memabuat parameter
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}

	// Pengecekan apakah id product memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	getShoppingCart, err := database.GetShoppingCartsTanpaJoin(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}

	if getShoppingCart == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("shopping cart id not found"))
	}

	getShoppingCartJSON, err := json.Marshal(getShoppingCart)
	if err != nil {
		panic(err)
	}

	var shoppingCarts models.Shopping_Carts
	json.Unmarshal([]byte(getShoppingCartJSON), &shoppingCarts)

	if shoppingCarts.UsersID != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("not allowed to update"))
	}

	// Mendapatkan data product yang akan diperbaharui dari client
	var updatedShoppingCart models.Shopping_Carts
	c.Bind(&updatedShoppingCart)

	// Memperbaharui data menggunakan fungsi UpdateProduct
	shoppingCart, er := database.UpdateShoppingCarts(id, &updatedShoppingCart)
	if er != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to update shopping cart"))
	}
	database.AddQtyPrice(shoppingCarts.ProductsID, id)
	return c.JSON(http.StatusOK, responses.StatusSuccessData("update success", shoppingCart))
}

func DeleteShoppingCartController(c echo.Context) error {
	// Mendapatkan id cart yang diingikan client
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}

	// Pengecekan apakah id cart memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	getShoppingCart, err := database.GetShoppingCartsTanpaJoin(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	getShoppingCartJSON, err := json.Marshal(getShoppingCart)
	if err != nil {
		panic(err)
	}

	var responseShoppingCart models.Shopping_Carts
	json.Unmarshal([]byte(getShoppingCartJSON), &responseShoppingCart)

	if responseShoppingCart.UsersID != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("not allowed to delete"))
	}

	// Mengapus data satu product menggunakan fungsi DeleteShoppingCart
	shoppingCart, e := database.DeleteShoppingCart(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to delete product"))
	}
	if shoppingCart == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("product id not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccess("success deleted one product"))
}
