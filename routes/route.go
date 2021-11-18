package routes

import (
	"project1/constants"
	c "project1/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	r := e.Group("/jwt")
	r.Use(middleware.JWT([]byte(constants.SECRET_JWT)))
	r.GET("/shopping_carts", c.GetShoppingCartsController)
	r.POST("/shopping_carts", c.CreateShoppingCartsController)
	r.PUT("/shopping_carts/:id", c.UpdateShoppingCartsController)
	r.DELETE("/shopping_carts/:id", c.DeleteShoppingCartController)
	e.POST("/login", c.LoginUsersController)
	return e
}
