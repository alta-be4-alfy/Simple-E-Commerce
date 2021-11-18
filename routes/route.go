package routes

import (
	"project1/constants"
	c "project1/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	
	e.POST("/register", controllers.RegisterUsersController)
	e.POST("/login", controllers.LoginUsersController) // jwt login
	e.GET("/products", controllers.GetProductsController)
	e.GET("/products/:id", controllers.GetProductController)
	e.GET("/products/users/:id", controllers.GetUserProductController)

	r := e.Group("/jwt")
	r.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	r.GET("/users/:id", controllers.GetAllUsersController)   // jwt
	r.DELETE("/users/:id", controllers.DeleteUserController) // jwt delete
	r.PUT("/users/:id", controllers.UpdateUserController)    // jwt put
	r.POST("/products", controllers.CreateProductController)
	r.PUT("/products/:id", controllers.UpdateProductController)
	r.DELETE("/products/:id", controllers.DeleteProductController)
  r.GET("/shopping_carts", c.GetShoppingCartsController)
	r.POST("/shopping_carts", c.CreateShoppingCartsController)
	r.PUT("/shopping_carts/:id", c.UpdateShoppingCartsController)
	r.DELETE("/shopping_carts/:id", c.DeleteShoppingCartController)

	return e
}
