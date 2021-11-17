package main

import (
	"project1/config"
	m "project1/middlewares"
	"project1/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	m.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(":8080"))
}
