# Restful API Project1 Group 3 E-Commerces

[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=blue)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=blue)](https://github.com/labstack/echo)


# Table of Content

- [Description](#description)
- [Requirements](#Requirements)
- [How to use](#how-to-use)
- [Our Feature](#Our-Feature)
- [Structuring](#Structuring)
- [Endpoints](#endpoints)
- [Credits](#credits)


# Description

Project Base task Alterra Academy ini dibuat untuk mengimplementasikan SCRUM method using GO Language STRUCTURING , ECHO , JWT , API. 

# Requirements

* Visual Studio Code
* Postman
* Mysql Workbench

# How to use
- Install Go, Postman, MySQL Workbench
- Clone this repository in your $PATH:
```
git clone https://github.com/alta-be4-alfy/project1_kelompok3.git
```
* CREATE DATABASE IF NOT EXISTS `project1_kelompok3`;
* USE `project1_kelompok3`;
* Run `main.go`
```
$ go run main.go
```
* Open Postman run with your local host, follow the routes in Visual Studio Code folder.


# Our Feature
* CREATE : User, Shopping Cart, Order, dan Product.
* READ : User, Shopping Cart, Order, dan Product.
* UPDATE User, Shopping Cart, dan Product.
* DELETE : User, Shopping Cart, dan Product.


# Structuring
```
📦project1_kelompok3
 ┣ 📂.vscode
 ┃   ┗ 📜settings.json
 ┣ 📂config
 ┃   ┗ 📜config.go
 ┣ 📂constants
 ┃   ┗ 📜constant.go
 ┣ 📂controllers
 ┃   ┗ 📜orderController.go
 ┃   ┗ 📜productController.go
 ┃   ┗ 📜productController_test.go
 ┃   ┗ 📜shoppingCartController.go
 ┃   ┗ 📜shoppingCartController_test.go
 ┃   ┗ 📜userController.go
 ┣ 📂lib
 ┃   ┗ 📂database
 ┃     ┗ 📜order.go
 ┃     ┗ 📜product.go
 ┃     ┗ 📜shoppingCart.go
 ┃     ┗ 📜user.go
 ┃   ┗ 📂response
 ┃     ┗ 📜response.go
 ┣ 📂middlewares
 ┃   ┗ 📜logMiddleware.go
 ┃   ┗ 📜middlewares.go
 ┣ 📂models
 ┃   ┗ 📜address.go
 ┃   ┗ 📜order_detail.go
 ┃   ┗ 📜orders.go
 ┃   ┗ 📜payment_methods.go
 ┃   ┗ 📜products.go
 ┃   ┗ 📜shopping_carts.go
 ┃   ┗ 📜users.go
 ┣ 📂routes
 ┃   ┗ 📜route.go
 ┣ 📜.env
 ┣ 📜.gitignore
 ┣ 📜cover.out
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┣ 📜main.go
 ┣ 📜profile.cov
 ┗ 📜README.MD
```

# Endpoints

| Method | Endpoint | Description| Authentication | Authorization
|:-----|:--------|:----------| :----------:| :----------:|
| POST  | /register | Register a new user | No | No
| POST | /login | Login existing user| No | No
|---|---|---|---|---|
| GET    | /users/:id | Get list of all user | Yes | Yes
| PUT | /users/:id | Update user profile | Yes | Yes
| DELETE | /users/:id | Delete user profile | Yes | Yes
|---|---|---|---|---|
| POST   | /products | Add products | Yes | Yes
| GET   | /products | Get products list | Yes | Yes
| DELETE   | /categories/:id | Delete category by id | Yes | Yes
|---|---|---|---|---|
| GET | /products | Get list of all products | No | No
| GET | /products/:id | Get product by product id | No | No
| POST | /products | Add products by admin | Yes | Yes
| PUT | /products/:id | Update products by admin | Yes | Yes
| DELETE | /products/:id | Delete products by admin | Yes | Yes
|---|---|---|---|---|
| POST | /carts | Add products to cart | Yes | Yes
| GET | /carts | Get list of all cart item | Yes | Yes
| PUT | /cartitems/:id | Update cart item by id | Yes | Yes
| DELETE | /cartitems/:id | Delete cart item by id | Yes | Yes
|---|---|---|---|---|
| POST | /checkout | List of products checkout | Yes | Yes
|---|---|---|---|---|
| GET | /payments/:id | Get transaction status | Yes | Yes
|---|---|---|---|---|
| GET | /transactions | Get list of all transaction | Yes | Yes
| GET | /transactionreport?range={range} | Get transactions with range date | Yes | Yes
|---|---|---|---|---|

<br>


# Credits

1. https://github.com/alfynf
2. https://github.com/Nathannov24
3. https://github.com/armuh16
