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
type ProductsResponseSuccess struct {
	Status  string
	Message string
	Data    []models.Products
}

// Struct yang digunakan ketika test request success, hanya menampung satu data
type SingleProductResponseSuccess struct {
	Status  string
	Message string
	Data    models.Products
}

// Struct yang digunakan ketika test request failed
type ProductsResponseFailed struct {
	Status  string
	Message string
}

// Struct untuk menampung data test case
type ProductsTestCase struct {
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
	mock_data_product = models.Products{
		Product_Name:        "Android Mini",
		Product_Type:        "Elektronik",
		Product_Stock:       3,
		Product_Price:       5000000,
		Product_Description: "5 in, 64GB",
		UsersID:             1,
	}
	mock_update_product1 = models.Products{
		Product_Name:        "Android Jumbo",
		Product_Description: "5 in, 100GB",
	}
	mock_update_product = models.Products{
		Product_Name:        "Android Jumbo",
		Product_Description: "5 in, 100GB",
		UsersID:             2,
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
)

// Fungsi untuk memasukkan data product test ke dalam database
func InsertMockDataProductsToDB() error {
	query := config.DB.Save(&mock_data_product)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data update product test ke dalam database
func InsertMockDataUpdateProductsToDB() error {
	query := config.DB.Save(&mock_update_product)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data user test ke dalam database
func InsertMockDataUsersToDB() error {
	query := config.DB.Save(&mock_data_user)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk memasukkan data update user test ke dalam database
func InsertMockDataUpdateUsersToDB() error {
	query := config.DB.Save(&mock_update_user)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

// Fungsi untuk melakukan login dan ekstraksi token JWT
func UsingJWT() (string, error) {
	// Melakukan login data user test
	InsertMockDataUsersToDB()
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

// Fungsi untuk melakukan testing fungsi GetProductsController
// kondisi request success
func TestGetProductsControllerSuccess(t *testing.T) {
	var testCases = []struct {
		Name       string
		Path       string
		ExpectCode int
		ExpectSize int
	}{
		{
			Name:       "success to get all data products",
			Path:       "/products",
			ExpectCode: http.StatusOK,
			ExpectSize: 1,
		},
	}

	e := InitEchoTestAPI()
	InsertMockDataUsersToDB()
	InsertMockDataProductsToDB()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	for index, testCase := range testCases {
		context.SetPath(testCase.Path)

		if assert.NoError(t, GetProductsController(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			body := rec.Body.String()
			var responses ProductsResponseSuccess
			err := json.Unmarshal([]byte(body), &responses)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, testCases[index].ExpectSize, len(responses.Data))
			assert.Equal(t, "success", responses.Status)
		}
	}
}

// Fungsi untuk melakukan testing fungsi GetProductsController
// kondisi request failed
func TestGetProductsControllerFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "failed to get all data products",
		Path:       "/products",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	// Melakukan penghapusan tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	context.SetPath(testCases.Path)
	GetProductsController(context)

	body := rec.Body.String()
	var responses ProductsResponseFailed
	er := json.Unmarshal([]byte(body), &responses)
	assert.Equal(t, testCases.ExpectCode, rec.Code)
	if er != nil {
		assert.Error(t, er, "error")
	}
	assert.Equal(t, "failed", responses.Status)
}

// Fungsi untuk melakukan testing fungsi GetProductsController
// kondisi request failed
func TestGetProductsControllerNoProduct(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "products not found",
		Path:       "/products",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	context.SetPath(testCases.Path)
	GetProductsController(context)

	body := rec.Body.String()
	var responses ProductsResponseFailed
	er := json.Unmarshal([]byte(body), &responses)
	assert.Equal(t, testCases.ExpectCode, rec.Code)
	if er != nil {
		assert.Error(t, er, "error")
	}
	assert.Equal(t, "failed", responses.Status)
}

// Fungsi testing CreateProductController
func CreateProductControllerTesting() echo.HandlerFunc {
	return CreateProductController
}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request success
func TestCreateProductControllerSuccess(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "success to create product",
		Path:       "/products",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(mock_data_product)
	if err != nil {
		t.Error(t, err, "error")
	}

	// Mengirim data menggunakan request body dengan HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	middleware.JWT([]byte(constants.SECRET_JWT))(CreateProductControllerTesting())(context)

	bodyResponses := rec.Body.String()
	var product SingleProductResponseSuccess

	er := json.Unmarshal([]byte(bodyResponses), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/products", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "success", product.Status)
	})
}

// Fungsi untuk melakukan testing fungsi CreateProductController menggunakan JWT
// kondisi request failed
func TestCreateProductControllerFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "failed to create product",
		Path:       "/products",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(mock_data_product)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Menghapus tabel user untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(testCases.Path)

	// Call function on controller
	middleware.JWT([]byte(constants.SECRET_JWT))(CreateProductControllerTesting())(context)
	bodyResponses := rec.Body.String()
	var product ProductsResponseFailed

	er := json.Unmarshal([]byte(bodyResponses), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("POST /jwt/products", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, rec.Code)
		assert.Equal(t, "failed", product.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetProductController
// kondisi request success
func TestGetProductControllerSuccess(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "success to get one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()

	InsertMockDataUsersToDB()
	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodGet, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, GetProductController(context)) {
		res_body := res.Body.String()
		var response SingleProductResponseSuccess
		er := json.Unmarshal([]byte(res_body), &response)
		if er != nil {
			assert.Error(t, er, "error")
		}
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "success", response.Status)
	}
}

// Fungsi untuk melakukan testing fungsi GetProductController
// kondisi request failed
func TestGetProductControllerFailedChar(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "false param",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodGet, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan tipe data id yang berbeda untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("#")
	if assert.NoError(t, GetProductController(context)) {
		res_body := res.Body.String()
		var response ProductsResponseFailed
		er := json.Unmarshal([]byte(res_body), &response)
		if er != nil {
			assert.Error(t, er, "error")
		}
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	}
}

// Fungsi untuk melakukan testing fungsi GetProductController
// kondisi request failed
func TestGetProductControllerWrongId(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "wrong id",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	InsertMockDataProductsToDB()

	req := httptest.NewRequest(http.MethodGet, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan ProductID yang tidak tersimpan di database untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("3")
	if assert.NoError(t, GetProductController(context)) {
		var response ProductsResponseFailed
		res_body := res.Body.String()
		er := json.Unmarshal([]byte(res_body), &response)
		if er != nil {
			assert.Error(t, er, "error")
		}
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	}
}

// Fungsi untuk melakukan testing fungsi GetProductController menggunakan JWT
// kondisi request failed
func TestGetProductControllerFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "failed to get one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	// Menghapus tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodGet, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	if assert.NoError(t, GetProductController(context)) {
		var response ProductsResponseFailed
		res_body := res.Body.String()
		er := json.Unmarshal([]byte(res_body), &response)
		if er != nil {
			assert.Error(t, er, "error")
		}
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	}
}

// Fungsi testing CreateProductController
func DeleteProductControllerTesting() echo.HandlerFunc {
	return DeleteProductController
}

// Fungsi untuk melakukan testing fungsi DeleteProductController menggunakan JWT
// kondisi request success
func TestDeleteProductControllerSuccess(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "success to delete one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodDelete, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteProductControllerTesting())(context)

	var response SingleProductResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "success", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi DeleteProductController menggunakan JWT
// kondisi request failed
func TestDeleteProductControllerFailedChar(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "false param",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodDelete, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan tipe data id yang berbeda untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteProductControllerTesting())(context)

	var response ProductsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi DeleteProductController menggunakan JWT
// kondisi request failed
func TestDeleteProductControllerWrongId(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "product not found",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()

	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodDelete, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan ProductID yang tidak tersimpan di database untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("3")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteProductControllerTesting())(context)

	var response ProductsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi untuk melakukan testing fungsi DeleteProductController menggunakan JWT
// kondisi request failed
func TestDeleteProductControllerFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "failed to delete one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	// Menghapus tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodDelete, "/products/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteProductControllerTesting())(context)

	var response ProductsResponseFailed
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &response)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("DELETE /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", response.Status)
	})
}

// Fungsi testing UpdateProductController
func UpdateProductControllerTesting() echo.HandlerFunc {
	return UpdateProductController
}

// Fungsi untuk melakukan testing fungsi UpdateProductController menggunakan JWT
// kondisi request success
func TestUpdateProductControllerSuccess(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "success to update one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	// Mendapatkan data update product
	body, err := json.Marshal(mock_update_product1)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodPut, "/products/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateProductControllerTesting())(context)

	var product SingleProductResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("PUT /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "success", product.Status)
	})
}

// Fungsi untuk melakukan testing fungsi UpdateProductController menggunakan JWT
// kondisi request failed
func TestUpdateProductControllerFailedChar(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "false param",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan data update product
	body, err := json.Marshal(mock_update_product1)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodPut, "/products/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan tipe data id yang berbeda untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateProductControllerTesting())(context)

	var response ProductsResponseFailed
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
func TestUpdateProductControllerWrongId(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "wrong id",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan data update product
	body, err := json.Marshal(mock_update_product1)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodPut, "/products/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)

	// Memasukkan ProductID yang tidak tersimpan di database untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("3")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateProductControllerTesting())(context)

	var response ProductsResponseFailed
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
func TestUpdateProductControllerFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "failed to update one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan data update user
	body, err := json.Marshal(mock_update_product1)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	// Menghapus tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodPut, "/products/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateProductControllerTesting())(context)

	var response ProductsResponseFailed
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
func TestUpdateProductControllerNotAllowed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "not allowed update one data product",
		Path:       "/products/:id",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan data update user
	body, err := json.Marshal(mock_update_product)
	if err != nil {
		t.Error(t, err, "error")
	}
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataUpdateUsersToDB()
	InsertMockDataProductsToDB()
	InsertMockDataUpdateProductsToDB()

	req := httptest.NewRequest(http.MethodPut, "/products/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	// Membuat userID pada productID berbeda dengan userID token untuk membuat request failed
	context.SetParamNames("id")
	context.SetParamValues("2")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateProductControllerTesting())(context)

	var response ProductsResponseFailed
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

func GetUserProductControllerTesting() echo.HandlerFunc {
	return GetUserProductController
}

// Fungsi untuk melakukan testing fungsi GetUserProductController
// kondisi request success
func TestGetUserProductControllerSuccess(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "success to get one data product",
		Path:       "/products/users",
		ExpectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}

	InsertMockDataProductsToDB()
	req := httptest.NewRequest(http.MethodGet, "/products/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserProductControllerTesting())(context)

	var product ProductsResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("GET /jwt/products/users", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "success", product.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetUserProductController
// kondisi request failed
func TestGetUserProductControllerWrongId(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "id not found",
		Path:       "/products/users",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	// Memasukkan data produk dengan id user yang berbeda dengan token untuk membuat request failed
	InsertMockDataUpdateProductsToDB()

	req := httptest.NewRequest(http.MethodGet, "/products/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserProductControllerTesting())(context)

	var product ProductsResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("GET /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", product.Status)
	})
}

// Fungsi untuk melakukan testing fungsi GetUserProductController
// kondisi request failed
func TestGetUserProductControlleFailed(t *testing.T) {
	var testCases = ProductsTestCase{
		Name:       "wrong id",
		Path:       "/products/users",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Mendapatkan token
	token, err := UsingJWT()
	if err != nil {
		panic(err)
	}
	// Menghapus tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Products{})

	req := httptest.NewRequest(http.MethodGet, "/products/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.Path)
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserProductControllerTesting())(context)

	var product ProductsResponseSuccess
	res_body := res.Body.String()
	er := json.Unmarshal([]byte(res_body), &product)
	if er != nil {
		assert.Error(t, er, "error")
	}
	t.Run("GET /jwt/products/:id", func(t *testing.T) {
		assert.Equal(t, testCases.ExpectCode, res.Code)
		assert.Equal(t, "failed", product.Status)
	})
}
