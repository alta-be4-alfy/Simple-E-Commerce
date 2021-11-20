package controllers

import (
	"net/http"
	"project1/lib/database"
	responses "project1/lib/response"
	"project1/middlewares"
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
	idUser := middlewares.ExtractTokenUserId(c)
	idProduct, e := database.CreateProduct(idUser)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	product, _ := database.UpdateProduct(idProduct, &input)
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success to create product", product))
}

// Controller untuk mendapatkan seluruh data products
func GetProductsController(c echo.Context) error {
	// mendapatkan seluruh data product menggunakan fungsi GetProducts
	products, e := database.GetProducts()
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch products"))
	}
	if products == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("product not found"))
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
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("product not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get product by id", product))
}

// Controller untuk mendapatkan seluruh data product user tertentu berdasarkan id user
func GetUserProductController(c echo.Context) error {
	// Mendapatkan id user dari token
	idToken := middlewares.ExtractTokenUserId(c)
	// Mendapatkan data seluruh product user tertentu menggunakan fungsi GetUserProduct
	product, e := database.GetUserProducts(idToken)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch product"))
	}
	if product == 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("product not found"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get product by user id", product))
}

// Fungsi untuk memperbaharui satu data product berdasarkan id product
func UpdateProductController(c echo.Context) error {
	// Mendapatkan id product yang diingikan client
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}

	// Pengecekan apakah id product memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	idOwner, _ := database.GetProductOwner(id)
	if idOwner != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to update product"))
	}
	// Mendapatkan data product yang akan diperbaharui dari client
	var updatedProduct models.Products
	c.Bind(&updatedProduct)
	// Memperbaharui data menggunakan fungsi UpdateProduct
	database.UpdateProduct(id, &updatedProduct)
	updateProduct, _ := database.GetProduct(id)
	return c.JSON(http.StatusOK, responses.StatusSuccessData("update success", updateProduct))
}

// Controller untuk menghapus satu data product berdasarkan id product
func DeleteProductController(c echo.Context) error {
	// Mendapatkan id product yang diingikan client
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	// Pengecekan apakah id product memiliki id user yang sama dengan id token
	idToken := middlewares.ExtractTokenUserId(c)
	idOwner, _ := database.GetProductOwner(id)
	if idOwner != idToken {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to delete product"))
	}
	// Mengapus data satu product menggunakan fungsi DeleteProduct
	database.DeleteProduct(id)
	return c.JSON(http.StatusOK, responses.StatusSuccess("success deleted one product"))
}
