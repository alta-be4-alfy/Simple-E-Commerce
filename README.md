# Alta Store Kelompok 3

Program Alta Store ini dibuat untuk mengimplementasikan metode **SCRUM** untuk menyelesaikan tugas dalam tim. Program ini dibuat menggunakan bahasa **GO** dengan mengimplementasikan **STRUCTURING** , **ECHO** , **JWT** , **API**. 

Fitur yang tersedia pada program ini :
# 1. CREATE
Fitur **CREATE** ini mencakup tabel User, Shopping Cart, Order, dan Product.
# 2. READ
Fitur **READ** ini mencakup tabel User, Shopping Cart, Order, dan Product.
# 3.UPDATE
Fitur **UPDATE** ini mencakup tabel User, Shopping Cart, dan Product.
# 4. DELETE
Fitur **DELETE** ini mencakup tabel User, Shopping Cart, dan Product.


Berikut struktur folder dalam aplikasi ini:

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

## Requirements

* Visual Studio Code
* Postman
* Mysql Workbench

   
## Usage

1. Clone Repository **Alta Store** ke dalam folder yang diinginkan.
2. Buat Database dengan nama *project1_kelompok3*.
3. Buka Visual Studio Code, lalu ketikkan 
   > *go run main.go*.
4. Buka Postman , lalu jalankan sesuai routing yang tertera pada 
   > *./routes/route.go*

## Credits

1. https://github.com/alfynf
2. https://github.com/Nathannov24
3. https://github.com/armuh16