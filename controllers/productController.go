package controllers

import (
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/models"
	"strconv"

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

// Controller untuk mendapatkan seluruh data products
func GetProductsController(c echo.Context) error {
	// mendapatkan seluruh data product menggunakan fungsi GetProducts
	products, e := database.GetProducts()
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to load products"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to load products", products))
}

// Controller untuk mendapatkan satu data product berdasarkan id product
func GetProductController(c echo.Context) error {
	// Mendapatkan id buku yang diingikan client
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	// Mendapatkan data satu buku menggunakan fungsi GetProduct
	product, e := database.GetProduct(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	if product == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("product id not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get product by id", product))
}
