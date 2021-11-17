package routes

import (
	"project1/constants"
	"project1/controllers"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/register", controllers.RegisterUsersController)
	e.POST("/login", controllers.LoginUsersController) // jwt login

	r := e.Group("/jwt")
	r.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	r.GET("/users/:id", controllers.GetAllUsersController)   // jwt
	r.DELETE("/users/:id", controllers.DeleteUserController) // jwt delete
	r.PUT("/users/:id", controllers.UpdateUserController)    // jwt put
	return e
}
