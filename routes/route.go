package routes

import (
	c "project1/controllers"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	e.GET("/shopping_carts", c.GetShoppingCartsController)
	return e
}
