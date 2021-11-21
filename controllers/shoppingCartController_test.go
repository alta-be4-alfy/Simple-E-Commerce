package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"project1/config"
	"project1/constants"
	"project1/middlewares"
	"project1/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

// Struct yang digunakan ketika test request success, dapat menampung banyak data
type ShoppingCartsResponseSuccess struct {
	Status  string
	Message string
	Data    []models.Shopping_Carts
}

// Struct yang digunakan ketika test request success, hanya menampung satu data
type SingleShoppingCartsResponseSuccess struct {
	Status  string
	Message string
	Data    models.Shopping_Carts
}

// Struct yang digunakan ketika test request failed
type ShoppingCartsResponseFailed struct {
	Status  string
	Message string
}

// Struct untuk menampung data test case
type ShoppingCartsTestCase struct {
	Name       string
	Path       string
	ExpectCode int
}

// Fungsi untuk menginisiasi koneksi ke database test
func InitEchoTestShoppingCartAPI() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

var (
	mock_data_product_shoppingcart = models.Products{
		Product_Name:        "Android Mini",
		Product_Type:        "Elektronik",
		Product_Stock:       3,
		Product_Price:       5000000,
		Product_Description: "5 in, 64GB",
		UsersID:             1,
	}
	mock_data_order_shoppingcart = models.Orders{
		Total_Qty:         1,
		Total_Price:       1000,
		Order_Status:      "pending",
		AddressID:         1,
		Payment_MethodsID: 1,
	}
	mock_data_address_shoppingcart = models.Address{
		Street:   "Dahlia",
		City:     "Bekasi",
		Province: "West Java",
		Zip:      17520,
	}
	mock_data_shopping_cart = models.Shopping_Carts{
		Qty:        1,
		Price:      1,
		ProductsID: 1,
		UsersID:    1,
	}
	mock_data_user_shoppingcart = models.Users{
		User_Name: "alfa",
		Email:     "alfa@gmail.com",
		Password:  "inipwd",
	}
)

// Fungsi untuk memasukkan data product test ke dalam database
func InsertMockDataProductsShoppingCartToDB() error {
	query := config.DB.Save(&mock_data_product_shoppingcart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data address test ke dalam database
func InsertMockDataAddressShoppingCartToDB() error {
	query := config.DB.Save(&mock_data_address_shoppingcart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data shopping cart test ke dalam database
func InsertMockDataShoppingCartToDB() error {
	query := config.DB.Save(&mock_data_shopping_cart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data order test ke dalam database
func InsertMockDataOrdersShoppingCartToDB() error {
	query := config.DB.Save(&mock_data_order_shoppingcart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data user test ke dalam database
func InsertMockDataUsersShoppingCartsToDB() error {
	query := config.DB.Save(&mock_data_user_shoppingcart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk melakukan login dan ekstraksi token JWT
func UsingJWTCart() (string, error) {
	// Melakukan login data user test
	InsertMockDataUsersShoppingCartsToDB()
	var user models.Users
	tx := config.DB.Where("email = ? AND password = ?", mock_data_user_shoppingcart.Email, mock_data_user_shoppingcart.Password).First(&user)
	if tx.Error != nil {
		return "", tx.Error
	}
	// Mengektraksi token data user test
	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Fungsi testing CreateProductController
func CreateShoppingCartsControllerTesting() echo.HandlerFunc {
	return CreateShoppingCartsController
}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request success
func TestCreateProductControllerSuccess(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "success to create shopping cart",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestShoppingCartAPI()

	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(mock_data_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}

	// Mengirim data menggunakan request body dengan HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/shopping_carts", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	middleware.JWT([]byte(constants.SECRET_JWT))(CreateShoppingCartsControllerTesting())(context)

	bodyResponses := rec.Body.String()
	var shoppingCart SingleShoppingCartsResponseSuccess

	er := json.Unmarshal([]byte(bodyResponses), &shoppingCart)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/shopping_carts", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, 1, shoppingCart.Data.Qty)
	})
}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request failed
// func TestCreateProductControllerFailed(t *testing.T) {
// 	var testCases = ShoppingCartsTestCase{
// 		Name:       "failed to create shopping cart",
// 		Path:       "/shopping_carts",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestShoppingCartAPI()

// 	token, err := UsingJWTCart()
// 	if err != nil {
// 		panic(err)
// 	}
// 	body, err := json.Marshal(mock_data_shopping_cart)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Menghapus tabel user untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Shopping_Carts{})

// 	req := httptest.NewRequest(http.MethodPost, "/shopping_carts", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	rec := httptest.NewRecorder()
// 	context := e.NewContext(req, rec)
// 	context.SetPath(testCases.Path)

// 	// Call function on controller
// 	middleware.JWT([]byte(constants.SECRET_JWT))(CreateShoppingCartsControllerTesting())(context)
// 	bodyResponses := rec.Body.String()
// 	var shoppingCart ShoppingCartsResponseFailed

// 	er := json.Unmarshal([]byte(bodyResponses), &shoppingCart)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("POST /jwt/shopping_carts", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, rec.Code)
// 		assert.Equal(t, "failed", shoppingCart.Status)
// 	})
// }
