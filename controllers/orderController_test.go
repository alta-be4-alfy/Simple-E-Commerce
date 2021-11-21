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
		Order_Status:      "done",
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
	config.DB.Save(&mock_data_user)
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
	InsertMockDataToDB()
	// Melakukan penghapusan tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Orders{})

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetAllOrderControllerTesting())(context)
	var responses OrdersResponseFailed

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
	var responses OrdersResponseFailed

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

// Fungsi testing CreateOrderController
func CreateOrderControllerTesting() echo.HandlerFunc {
	return CreateOrderController
}

// Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// kondisi request success
func TestCreateOrderControllerSuccess(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "success to create order",
		Path:       "/orders",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	config.DB.Save(&mock_data_address)
	config.DB.Save(&mock_data_payment)

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(mock_data_order)
	if err != nil {
		t.Error(t, err, "error")
	}

	// Mengirim data menggunakan request body dengan HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderControllerTesting())(context)

	bodyResponses := rec.Body.String()
	var order SingleOrderResponseSuccess

	er := json.Unmarshal([]byte(bodyResponses), &order)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/orders", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", order.Status)
	})
}

// Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// kondisi request failed
func TestCreateOrderControllerFailed(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "failed to create order",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(mock_data_order)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Menghapus tabel user untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Orders{})

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var order OrdersResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &order)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/orders", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", order.Status)
	})
}

func GetHistoryOrderControllerTesting() echo.HandlerFunc {
	return GetHistoryOrderController
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request success
func TestGetHistoryOrdersControllerSuccess(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "success to get all history data orders",
		Path:       "/orders/history",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	InsertMockDataToDB()

	req := httptest.NewRequest(http.MethodGet, "/orders/history", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetHistoryOrderControllerTesting())(context)
	var responses OrdersResponseSuccess
	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/history", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestGetHistoryOrdersControllerFailed(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "failed to get all history data orders",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	// Menghapus tabel order untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Orders{})

	req := httptest.NewRequest(http.MethodGet, "/orders/history", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetAllOrderControllerTesting())(context)
	var responses OrdersResponseFailed

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/history", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestHistoryGetOrdersControllerNoOrder(t *testing.T) {
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
	middleware.JWT([]byte(constants.SECRET_JWT))(GetHistoryOrderControllerTesting())(context)
	var responses OrdersResponseFailed

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/history", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

func GetCancelOrderControllerTesting() echo.HandlerFunc {
	return GetCancelOrderController
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request success
func TestGetCancelOrdersControllerSuccess(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "success to get all cancel data orders",
		Path:       "/orders/cancel",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	mock_data_order.Order_Status = "cancel"
	InsertMockDataToDB()

	req := httptest.NewRequest(http.MethodGet, "/orders/cancel", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetCancelOrderControllerTesting())(context)
	var responses OrdersResponseSuccess
	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/cancel", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestGetCancelOrdersControllerFailed(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "failed to get all cancel data orders",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	// Menghapus tabel order untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Orders{})

	req := httptest.NewRequest(http.MethodGet, "/orders/cancel", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetCancelOrderControllerTesting())(context)
	var responses OrdersResponseFailed

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/cancel", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetOrdersController
// kondisi request failed
func TestCancelGetOrdersControllerNoOrder(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "cancel orders not found",
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
	middleware.JWT([]byte(constants.SECRET_JWT))(GetCancelOrderControllerTesting())(context)
	var responses OrdersResponseFailed

	body := rec.Body.String()
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		assert.Error(t, err, "error")
	}
	t.Run("GET /jwt/orders/cancel", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", responses.Status)
	})
}

func CreateOrderDetailControllerTesting() echo.HandlerFunc {
	return CreateOrderDetailController
}

// Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// kondisi request success
func TestCreateOrderDetailControllerSuccess(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "success to create order detail",
		Path:       "/orders/detail",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	InsertMockDataToDB()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(mock_data_orderdetail)
	if err != nil {
		t.Error(t, err, "error")
	}

	// Mengirim data menggunakan request body dengan HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/orders/detail", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderDetailControllerTesting())(context)

	bodyResponses := rec.Body.String()
	var order SingleOrderResponseSuccess

	er := json.Unmarshal([]byte(bodyResponses), &order)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/orders/detail", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", order.Status)
	})
}

// Fungsi untuk melakukan testing fungsi CreateOrderController menggunakan JWT
// kondisi request failed
func TestCreateOrderDetailControllerFailed(t *testing.T) {
	var testCases = OrdersTestCase{
		Name:       "failed to create order detail",
		Path:       "/orders",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(mock_data_orderdetail)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Menghapus tabel user untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Order_Details{})

	req := httptest.NewRequest(http.MethodPost, "/orders/detail", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(CreateOrderDetailControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var order OrdersResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &order)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/orders/detail", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", order.Status)
	})
}
