package routes

import (
	"project1/constants"
	controllers "project1/controllers"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/register", controllers.RegisterUsersController)
	e.POST("/login", controllers.LoginUsersController) // jwt login
	e.GET("/products", controllers.GetProductsController)
	e.GET("/products/:id", controllers.GetProductController)

	r := e.Group("/jwt")
	r.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	r.GET("/users/:id", controllers.GetUsersController)      //jwt
	r.DELETE("/users/:id", controllers.DeleteUserController) // jwt delete
	r.PUT("/users/:id", controllers.UpdateUserController)    // jwt put
	r.GET("/products/users", controllers.GetUserProductController)
	r.POST("/products", controllers.CreateProductController)
	r.PUT("/products/:id", controllers.UpdateProductController)
	r.DELETE("/products/:id", controllers.DeleteProductController)
	r.GET("/orders", controllers.GetAllOrderController)
	r.GET("/orders/history", controllers.GetHistoryOrderController)
	r.GET("/orders/cancel", controllers.GetCancelOrderController)
	r.POST("/orders", controllers.CreateOrderController)
	r.POST("/orders/status", controllers.ChangeOrderStatusController)
	r.GET("/shopping_carts", controllers.GetShoppingCartsController)
	r.POST("/shopping_carts", controllers.CreateShoppingCartsController)
	r.PUT("/shopping_carts/:id", controllers.UpdateShoppingCartsController)
	r.DELETE("/shopping_carts/:id", controllers.DeleteShoppingCartController)
	return e
}
