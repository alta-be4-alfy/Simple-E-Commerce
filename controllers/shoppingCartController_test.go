package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project1/config"
	"project1/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Struct yang digunakan ketika test request success, dapat menampung banyak data
type ShoppingCartsResponseSuccess struct {
	Status  string
	Message string
	Data    []models.Shopping_Carts
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
		UsersID:           1,
		Payment_MethodsID: 1,
		Shopping_CartsID:  []models.Shopping_Carts{},
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
		OrdersID:   1,
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

// Fungsi untuk melakukan testing fungsi GetShoppingCartsController
// kondisi request success
func TestGetShoppingCartsControllerSuccess(t *testing.T) {
	var testCases = []struct {
		Name       string
		Path       string
		ExpectCode int
		ExpectSize int
	}{
		{
			Name:       "success to get all data shopping carts",
			Path:       "/shopping_carts",
			ExpectCode: http.StatusOK,
			ExpectSize: 1,
		},
	}

	e := InitEchoTestShoppingCartAPI()
	InsertMockDataOrdersShoppingCartToDB()
	InsertMockDataProductsShoppingCartToDB()
	InsertMockDataShoppingCartToDB()
	InsertMockDataAddressShoppingCartToDB()
	InsertMockDataUsersShoppingCartsToDB()

	req := httptest.NewRequest(http.MethodGet, "/shopping_carts", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	for index, testCase := range testCases {
		context.SetPath(testCase.Path)

		if assert.NoError(t, GetShoppingCartsController(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			body := rec.Body.String()
			var responses ShoppingCartsResponseSuccess
			err := json.Unmarshal([]byte(body), &responses)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, testCases[index].ExpectSize, len(responses.Data))
			assert.Equal(t, "1", responses.Data[0].OrdersID)
		}
	}
}

// Fungsi untuk melakukan testing fungsi GetProductsController
// kondisi request failed
func TestGetShoppingCartsControllerFailed(t *testing.T) {
	var testCases = ShoppingCartsTestCase{
		Name:       "failed to get all data products",
		Path:       "/shopping_carts",
		ExpectCode: http.StatusBadRequest,
	}

	e := InitEchoTestShoppingCartAPI()

	// Melakukan penghapusan tabel untuk membuat request failed
	config.DB.Migrator().DropTable(&models.Shopping_Carts{})

	req := httptest.NewRequest(http.MethodGet, "/shopping_carts", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	context.SetPath(testCases.Path)
	GetShoppingCartsController(context)

	body := rec.Body.String()
	var responses ShoppingCartsResponseFailed
	er := json.Unmarshal([]byte(body), &responses)
	assert.Equal(t, testCases.ExpectCode, rec.Code)
	if er != nil {
		assert.Error(t, er, "error")
	}
	assert.Equal(t, "failed", responses.Status)
}
