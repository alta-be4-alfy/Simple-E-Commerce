package controllers

import (
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
type OrdersResponseSuccess struct {
	Status  string
	Message string
	Data    []models.Orders
}

// Struct yang digunakan ketika test request success, hanya menampung satu data
type SingleOrderResponseSuccess struct {
	Status  string
	Message string
	Data    models.Orders
}

// Struct yang digunakan ketika test request failed
type OrdersResponseFailed struct {
	Status  string
	Message string
}

// Struct untuk menampung data test case
type OrdersTestCase struct {
	Name       string
	Path       string
	ExpectCode int
}

// Fungsi untuk menginisiasi koneksi ke database test
func InitEchoTestAPI() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

var (
	mock_data_address = models.Address{
		Street: "Jl.Suaka",
	}
	mock_data_payment = models.Payment_Methods{
		Payment: "OVO",
	}
	mock_data_user = models.Users{
		User_Name: "alfa",
		Email:     "alfa@gmail.com",
		Password:  "inipwd",
	}
	mock_update_user = models.Users{
		User_Name: "ajjo",
		Email:     "ajjo@gmail.com",
		Password:  "itupwd",
	}
	mock_data_product = models.Products{
		Product_Name:        "Android Mini",
		Product_Type:        "Elektronik",
		Product_Stock:       3,
		Product_Price:       100000,
		Product_Description: "5 in, 64GB",
		UsersID:             1,
	}
	mock_data_shoppingcart = models.Shopping_Carts{
		Qty:        1,
		Price:      100000,
		ProductsID: 1,
		UsersID:    1,
	}
	mock_data_order = models.Orders{
		Payment_MethodsID: 1,
		AddressID:         1,
		UsersID:           1,
	}
	mock_data_orderdetail = models.Order_Details{
		OrdersID:         1,
		Shopping_CartsID: 1,
	}
)

// Fungsi untuk memasukkan data test ke dalam database
func InsertMockDataToDB() {
	config.DB.Save(&mock_data_address)
	config.DB.Save(&mock_data_payment)
	config.DB.Save(&mock_data_user)
	config.DB.Save(&mock_data_order)
	config.DB.Save(&mock_data_product)
	config.DB.Save(&mock_data_shoppingcart)
	config.DB.Save(&mock_data_orderdetail)
}

// // // Fungsi untuk memasukkan data update order test ke dalam database
// func InsertMockDataUpdateOrdersToDB() error {
// 	query := config.DB.Save(&mock_update_order)
// 	if query.Error != nil {
// 		return query.Error
// 	}
// 	return nil
// }

// Fungsi untuk melakukan login dan ekstraksi token JWT
func UsingJWT() (string, error) {
	// Melakukan login data user test
	InsertMockDataToDB()
	var user models.Users
	tx := config.DB.Where("email = ? AND password = ?", mock_data_user.Email, mock_data_user.Password).First(&user)
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

func GetAllOrderControllerTesting() echo.HandlerFunc {
	return GetAllOrderController
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request success
func TestGetOrdersControllerSuccess(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "success to get all data orders",
		Path:       "/orders",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	InsertMockDataToDB()
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetAllOrderControllerTesting())(context)
	var responses SingleOrderResponseSuccess
	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestGetOrdersControllerFailed(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "failed to get all data orders",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	// Melakukan penghapusan tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Orders{})

	InsertMockDataToDB()
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetAllOrderControllerTesting())(context)
	var responses SingleOrderResponseSuccess

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestGetOrdersControllerNoOrder(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "orders not found",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	context.SetPath(testCases.Path)
	middleware.JWT([]byte(constants.SECRET_JWT))(GetAllOrderControllerTesting())(context)
	var responses SingleOrderResponseSuccess

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

// // Fungsi testing CreateOrderController
// func CreateOrderControllerTesting() echo.HandlerFunc {
// 	return CreateOrderController
// }

// // Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// // kondisi request success
// func TestCreateOrderControllerSuccess(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "success to create order",
// 		Path:       "/orders",
// 		ExpectCode: http.StatusOK,
// 	}

// 	e := InitEchoTestAPI()

// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	body, err := json.Marshal(mock_data_order)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}

// 	// Mengirim data menggunakan request body dengan HTTP Method POST
// 	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	rec := httptest.NewRecorder()
// 	context := e.NewContext(req, rec)

// 	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderControllerTesting())(context)

// 	bodyResponses := rec.Body.String()
// 	var order SingleOrderResponseSuccess

// 	er := json.Unmarshal([]byte(bodyResponses), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("POST /jwt/orders", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, rec.Code)
// 		assert.Equal(t, "success", order.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// // kondisi request failed
// func TestCreateOrderControllerFailed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "failed to create order",
// 		Path:       "/orders",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()

// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}
// 	body, err := json.Marshal(mock_data_order)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Menghapus tabel user untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Orders{})

// 	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	rec := httptest.NewRecorder()
// 	context := e.NewContext(req, rec)
// 	context.SetPath(testCases.Path)

// 	// Call function on controller
// 	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderControllerTesting())(context)
// 	bodyResponses := rec.Body.String()
// 	var order OrdersResponseFailed

// 	er := json.Unmarshal([]byte(bodyResponses), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("POST /jwt/orders", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, rec.Code)
// 		assert.Equal(t, "failed", order.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi GetOrderController
// // kondisi request success
// func TestGetOrderControllerSuccess(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "success to get one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusOK,
// 	}

// 	e := InitEchoTestAPI()

// 	InsertMockDataUsersToDB()
// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodGet, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")

// 	if assert.NoError(t, GetOrderController(context)) {
// 		res_body := res.Body.String()
// 		var response SingleOrderResponseSuccess
// 		er := json.Unmarshal([]byte(res_body), &response)
// 		if er != nil {
// 			assert.Error(t, er, "error")
// 		}
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "success", response.Status)
// 	}
// }

// // Fungsi untuk melakukan testing fungsi GetOrderController
// // kondisi request failed
// func TestGetOrderControllerFailedChar(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "false param",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodGet, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan tipe data id yang berbeda untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("#")
// 	if assert.NoError(t, GetOrderController(context)) {
// 		res_body := res.Body.String()
// 		var response OrdersResponseFailed
// 		er := json.Unmarshal([]byte(res_body), &response)
// 		if er != nil {
// 			assert.Error(t, er, "error")
// 		}
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	}
// }

// // Fungsi untuk melakukan testing fungsi GetOrderController
// // kondisi request failed
// func TestGetOrderControllerWrongId(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "wrong id",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	InsertMockDataOrdersToDB()

// 	req := httptest.NewRequest(http.MethodGet, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan OrderID yang tidak tersimpan di database untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("3")
// 	if assert.NoError(t, GetOrderController(context)) {
// 		var response OrdersResponseFailed
// 		res_body := res.Body.String()
// 		er := json.Unmarshal([]byte(res_body), &response)
// 		if er != nil {
// 			assert.Error(t, er, "error")
// 		}
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	}
// }

// // Fungsi untuk melakukan testing fungsi GetOrderController menggunakan JWT
// // kondisi request failed
// func TestGetOrderControllerFailed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "failed to get one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()

// 	// Menghapus tabel untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Orders{})

// 	req := httptest.NewRequest(http.MethodGet, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")
// 	if assert.NoError(t, GetOrderController(context)) {
// 		var response OrdersResponseFailed
// 		res_body := res.Body.String()
// 		er := json.Unmarshal([]byte(res_body), &response)
// 		if er != nil {
// 			assert.Error(t, er, "error")
// 		}
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	}
// }

// // Fungsi testing CreateOrderController
// func DeleteOrderControllerTesting() echo.HandlerFunc {
// 	return DeleteOrderController
// }

// // Fungsi untuk melakukan testing fungsi DeleteOrderController menggunakan JWT
// // kondisi request success
// func TestDeleteOrderControllerSuccess(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "success to delete one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusOK,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodDelete, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteOrderControllerTesting())(context)

// 	var response SingleOrderResponseSuccess
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("DELETE /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "success", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi DeleteOrderController menggunakan JWT
// // kondisi request failed
// func TestDeleteOrderControllerFailedChar(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "false param",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodDelete, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan tipe data id yang berbeda untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("#")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("DELETE /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi DeleteOrderController menggunakan JWT
// // kondisi request failed
// func TestDeleteOrderControllerWrongId(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "order not found",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()

// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodDelete, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan OrderID yang tidak tersimpan di database untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("3")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("DELETE /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi DeleteOrderController menggunakan JWT
// // kondisi request failed
// func TestDeleteOrderControllerFailed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "failed to delete one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Menghapus tabel untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Orders{})

// 	req := httptest.NewRequest(http.MethodDelete, "/orders/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("DELETE /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi testing UpdateOrderController
// func UpdateOrderControllerTesting() echo.HandlerFunc {
// 	return UpdateOrderController
// }

// // Fungsi untuk melakukan testing fungsi UpdateOrderController menggunakan JWT
// // kondisi request success
// func TestUpdateOrderControllerSuccess(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "success to update one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusOK,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan data update order
// 	body, err := json.Marshal(mock_update_order1)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}
// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodPut, "/orders/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateOrderControllerTesting())(context)

// 	var order SingleOrderResponseSuccess
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "success", order.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi UpdateOrderController menggunakan JWT
// // kondisi request failed
// func TestUpdateOrderControllerFailedChar(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "false param",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan data update order
// 	body, err := json.Marshal(mock_update_order1)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodPut, "/orders/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan tipe data id yang berbeda untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("#")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi UpdateOrderController menggunakan JWT
// // kondisi request failed
// func TestUpdateOrderControllerWrongId(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "wrong id",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan data update order
// 	body, err := json.Marshal(mock_update_order1)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodPut, "/orders/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)

// 	// Memasukkan OrderID yang tidak tersimpan di database untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("3")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi UpdateOrderController menggunakan JWT
// // kondisi request failed
// func TestUpdateOrderControllerFailed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "failed to update one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan data update user
// 	body, err := json.Marshal(mock_update_order1)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Menghapus tabel untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Orders{})

// 	req := httptest.NewRequest(http.MethodPut, "/orders/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	context.SetParamNames("id")
// 	context.SetParamValues("1")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi UpdateOrderController menggunakan JWT
// // kondisi request failed
// func TestUpdateOrderControllerNotAllowed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "not allowed update one data order",
// 		Path:       "/orders/:id",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan data update user
// 	body, err := json.Marshal(mock_update_order)
// 	if err != nil {
// 		t.Error(t, err, "error")
// 	}
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataUpdateUsersToDB()
// 	InsertMockDataOrdersToDB()
// 	InsertMockDataUpdateOrdersToDB()

// 	req := httptest.NewRequest(http.MethodPut, "/orders/:id", bytes.NewBuffer(body))
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	// Membuat userID pada orderID berbeda dengan userID token untuk membuat request failed
// 	context.SetParamNames("id")
// 	context.SetParamValues("2")
// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateOrderControllerTesting())(context)

// 	var response OrdersResponseFailed
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &response)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("PUT /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", response.Status)
// 	})
// }

// func GetUserOrderControllerTesting() echo.HandlerFunc {
// 	return GetUserOrderController
// }

// // Fungsi untuk melakukan testing fungsi GetUserOrderController
// // kondisi request success
// func TestGetUserOrderControllerSuccess(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "success to get one data order",
// 		Path:       "/orders/users",
// 		ExpectCode: http.StatusOK,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}

// 	InsertMockDataOrdersToDB()
// 	req := httptest.NewRequest(http.MethodGet, "/orders/users", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserOrderControllerTesting())(context)

// 	var order OrdersResponseSuccess
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("GET /jwt/orders/users", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "success", order.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi GetUserOrderController
// // kondisi request failed
// func TestGetUserOrderControllerWrongId(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "id not found",
// 		Path:       "/orders/users",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Memasukkan data produk dengan id user yang berbeda dengan token untuk membuat request failed
// 	InsertMockDataUpdateOrdersToDB()

// 	req := httptest.NewRequest(http.MethodGet, "/orders/users", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserOrderControllerTesting())(context)

// 	var order OrdersResponseSuccess
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("GET /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", order.Status)
// 	})
// }

// // Fungsi untuk melakukan testing fungsi GetUserOrderController
// // kondisi request failed
// func TestGetUserOrderControlleFailed(t *testing.T) {
// 	var testCases = OrdersTestCase{
// 		Name:       "wrong id",
// 		Path:       "/orders/users",
// 		ExpectCode: http.StatusBadRequest,
// 	}

// 	e := InitEchoTestAPI()
// 	// Mendapatkan token
// 	token, err := UsingJWT()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Menghapus tabel untuk membuat request failed
// 	config.DB.Migrator().DropTable(&models.Orders{})

// 	req := httptest.NewRequest(http.MethodGet, "/orders/users", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath(testCases.Path)
// 	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserOrderControllerTesting())(context)

// 	var order OrdersResponseSuccess
// 	res_body := res.Body.String()
// 	er := json.Unmarshal([]byte(res_body), &order)
// 	if er != nil {
// 		assert.Error(t, er, "error")
// 	}
// 	t.Run("GET /jwt/orders/:id", func(t *testing.T) {
// 		assert.Equal(t, testCases.ExpectCode, res.Code)
// 		assert.Equal(t, "failed", order.Status)
// 	})
// }
