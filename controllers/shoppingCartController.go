package controllers

import (
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"

	"github.com/labstack/echo/v4"
)

// Controller untuk mendapatkan seluruh data shopping carts
func GetShoppingCartsController(c echo.Context) error {
	// mendapatkan seluruh data product menggunakan fungsi GetShoppingCarts
	shoppingCart, e := database.GetShoppingCarts()
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to load data"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to load data", shoppingCart))
}
