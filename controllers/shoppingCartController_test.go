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
	mock_data_shopping_cart = models.Shopping_Carts{
		Qty:        1,
		ProductsID: 1,
		UsersID:    1,
	}
	mock_dataQty0_shopping_cart = models.Shopping_Carts{
		Qty:        0,
		ProductsID: 1,
		UsersID:    1,
	}
	mock_update_shopping_cart = models.Shopping_Carts{
		Qty: 2,
	}
	mock_update1_shopping_cart = models.Shopping_Carts{
		Qty:     2,
		UsersID: 2,
	}
	mock_data_user_shoppingcart = models.Users{
		User_Name: "alfa",
		Email:     "alfa@gmail.com",
		Password:  "inipwd",
	}
	mock_update_user = models.Users{
		User_Name: "ajjo",
		Email:     "ajjo@gmail.com",
		Password:  "itupwd",
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

// Fungsi untuk memasukkan data shopping cart test ke dalam database
func InsertMockDataShoppingCartToDB() error {
	query := config.DB.Save(&mock_data_shopping_cart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
func InsertMockDataQty0ShoppingCartToDB() error {
	query := config.DB.Save(&mock_dataQty0_shopping_cart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func InsertMockDataUpdateShoppingCartToDB() error {
	query := config.DB.Save(&mock_update_shopping_cart)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func InsertMockDataUpdate1ShoppingCartToDB() error {
	query := config.DB.Save(&mock_update1_shopping_cart)
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

func InsertMockDataUsersUpdateShoppingCartsToDB() error {
	query := config.DB.Save(&mock_update_user)
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

// Fungsi testing CreateShoppingCartController
func CreateShoppingCartsControllerTesting() echo.HandlerFunc {
	return CreateShoppingCartsController
}

// Fungsi testing GetShoppingCartController
func GetShoppingCartsControllerTesting() echo.HandlerFunc {
	return GetShoppingCartsController

}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request success
func TestCreateShoppingCartControllerSuccess(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "success to create shopping cart",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestShoppingCartAPI()
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(mock_data_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}

	// Mengirim data menggunakan request body dengan HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/jwt/shopping_carts", bytes.NewBuffer(body))
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
		// assert.Equal(t, 1, shoppingCart.Data.Qty)
	})
}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request failed
func TestCreateProductControllerFailed(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "failed to create shopping cart",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataQty0ShoppingCartToDB()

	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(mock_dataQty0_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Menghapus tabel user untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Shopping_Carts{})

	req := httptest.NewRequest(http.MethodPost, "/jwt/shopping_carts", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(CreateShoppingCartsControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var shoppingCart ShoppingCartsResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &shoppingCart)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/shopping_carts", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", shoppingCart.Status)
	})
}

//Fungsi untuk melakukan testing fungsi GetProductsController
//kondisi request success
func TestGetShoppingCartControllerSuccess(t *testing.T) {
	var testCases = []struct {
		Name       string
		Path       string
		ExpectCode int
		ExpectSize int
	}{
		{
			Name:       "success to get all data shopping cart",
			Path:       "/shopping_carts",
			ExpectCode: http.StatusOK,
			ExpectSize: 1,
		},
	}

	e := InitEchoTestShoppingCartAPI()
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()

	req := httptest.NewRequest(http.MethodPost, "/jwt/shopping_carts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetShoppingCartsControllerTesting())(context)
	for index, testCase := range testCases {
		context.SetPath(testCase.Path)
		body := rec.Body.String()
		var responses ShoppingCartsResponseSuccess
		err := json.Unmarshal([]byte(body), &responses)
		if err != nil {
			assert.Error(t, err, "error")
		}
		t.Run("GET /jwt/shopping_carts", func(t *testing.T) {
			assert.Equal(t, testCases[index].ExpectSize, len(responses.Data))
			assert.Equal(t, "success", responses.Status)
		})

	}
}

// Fungsi untuk melakukan testing fungsi GetShoppingCartController menggunakan JWT
// kondisi request failed
func TestGetShoppingCartControllerFailed(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "failed to get shopping cart",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()

	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	// Menghapus tabel user untuk membuat request failed
	// config.DB.Migrator().DropTable(&models.Shopping_Carts{})

	req := httptest.NewRequest(http.MethodPost, "/jwt/shopping_carts", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(GetShoppingCartsControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var shoppingCart ShoppingCartsResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &shoppingCart)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/shopping_carts", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", shoppingCart.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetShoppingCartController menggunakan JWT dengan cara drop database
// kondisi request failed
func TestGetShoppingCartDropDatabaseControllerFailed(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "failed to get shopping cart",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()

	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	// Menghapus tabel user untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Shopping_Carts{})

	req := httptest.NewRequest(http.MethodPost, "/jwt/shopping_carts", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(GetShoppingCartsControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var shoppingCart ShoppingCartsResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &shoppingCart)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/shopping_carts", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", shoppingCart.Status)
	})
}

// Fungsi testing UpdateProductController
func UpdateShoppingCartControllerTesting() echo.HandlerFunc {
	return UpdateShoppingCartsController
}

// Fungsi untuk melakukan testing fungsi UpdateProductController menggunakan JWT
// kondisi request success
func TestUpdateShoppingCartControllerSuccess(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "success to update one data product",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan data update shopping cart
	body, err := json.Marshal(mock_update_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodPut, "/shopping_carts/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateShoppingCartControllerTesting())(context)

	var shopping_carts SingleShoppingCartsResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &shopping_carts)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("PUT /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, 2, shopping_carts.Data.Qty)
	})
}

// Fungsi untuk melakukan testing fungsi UpdateShoppingCartController menggunakan JWT
// kondisi request failed
func TestUpdateProductControllerFailedChar(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "false param",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan data update shopping cart
	body, err := json.Marshal(mock_update_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodPut, "/shopping_carts/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan tipe data id yang berbeda untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateShoppingCartControllerTesting())(context)

	var response ShoppingCartsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("PUT /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi UpdateShoppingCartController menggunakan JWT
// kondisi request failed
func TestUpdateProductControllerWrongId(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "wrong id",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan data update product
	body, err := json.Marshal(mock_update_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodPut, "/shopping_carts/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan Shopping_CartsID yang tidak tersimpan di database untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("3")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateShoppingCartControllerTesting())(context)

	var response ShoppingCartsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("PUT /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi UpdateProductController menggunakan JWT
// kondisi request failed
func TestUpdateShoppingCartControllerFailed(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "failed to update one data shopping cart",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan data update user
	body, err := json.Marshal(mock_update_shopping_cart)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	// Menghapus tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Shopping_Carts{})

	req := httptest.NewRequest(http.MethodPut, "/shopping_carts/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateShoppingCartControllerTesting())(context)

	var response ShoppingCartsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("PUT /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// // Fungsi untuk melakukan testing fungsi UpdateProductController menggunakan JWT
// // kondisi request failed
// func TestShoppingCartProductControllerNotAllowed(t *testing.T) {
// 	var testCases = ShoppingCartsTestCase{
// 		Name:       "not allowed update one data shopping cart",
// 		Path:       "/shopping_carts/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestShoppingCartAPI()
// 	// Mendapatkan data update user
// 	body, err := json.Marshal(mock_update_shopping_cart)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWTCart()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// InsertMockDataUpdate1ShoppingCartToDB()
// 	// InsertMockDataShoppingCartToDB()
// 	InsertMockDataUsersUpdateShoppingCartsToDB()

// 	req := httptest.NewRequest(http.MethodPut, "/shopping_carts/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	// Membuat userID pada shopping_cartsID berbeda dengan userID token untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("2")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateShoppingCartControllerTesting())(context)

// 	var response ShoppingCartsResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/shopping_carts/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// Fungsi testing DeleteShoppingCartController
func DeleteShoppingCartControllerTesting() echo.HandlerFunc {
	return DeleteShoppingCartController
}

// Fungsi untuk melakukan testing fungsi DeleteShoppingCartController menggunakan JWT
// kondisi request success
func TestDeleteShoppingCartControllerSuccess(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "success to delete one data shopping cart",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}
	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodDelete, "/jwt/shopping_carts/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteShoppingCartControllerTesting())(context)

	var response SingleShoppingCartsResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "success", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi DeleteShoppingCartController menggunakan JWT
// kondisi request failed
func TestDeleteShoppingCartControllerFailedChar(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "false param",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()
	// Mendapatkan token
	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodDelete, "/shopping_carts/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan tipe data id yang berbeda untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteShoppingCartControllerTesting())(context)

	var response ShoppingCartsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi DeleteShoppingCartController menggunakan JWT
// kondisi request failed
func TestDeleteShoppingCartControllerWrongId(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "wrong id",
		Path:       "/shopping_carts/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()

	token, err := UsingJWTCart()
	if err != nil {
		panic(err)
	}

	InsertMockDataUsersShoppingCartsToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	req := httptest.NewRequest(http.MethodDelete, "/shopping_carts/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan Shopping_CartsID yang tidak tersimpan di database untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("3")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteShoppingCartControllerTesting())(context)

	var response ShoppingCartsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/shopping_carts/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// // Fungsi untuk melakukan testing fungsi DeleteShoppingCartController menggunakan JWT
// // kondisi request failed
// func TestDeleteShoppingCartControllerNotAllowed(t *testing.T) {
// 	var testCases = ShoppingCartsTestCase{
// 		Name:       "not allowed delete one data shopping carts",
// 		Path:       "/shopping_carts/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestShoppingCartAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWTCart()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataUsersUpdateShoppingCartsToDB()
// 	InsertMockDataShoppingCartToDB()
// 	InsertMockDataUpdate1ShoppingCartToDB()

// 	req := httptest.NewRequest(http.MethodDelete, "/shopping_carts/:id", nil)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	// Membuat userID pada Shopping_CartsID berbeda dengan userID token untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("2")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteShoppingCartControllerTesting())(context)

// 	var response ShoppingCartsResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("DELETE /jwt/shopping_carts/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }
