package controllers

import (
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/models"

	"github.com/labstack/echo/v4"
)

// Controller untuk membuat data product baru
func CreateProductController(c echo.Context) error {
	// Mendapatkan data product baru dari client
	input := models.Products{}
	c.Bind(&input)

	// Menyimpan data buku baru menggunakan fungsi CreateProduct
	product, e := database.CreateProduct(input)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to create product"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create product", product))
}
