package config

import (
	"os"
	"project1/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	connectionString := os.Getenv("CONNECTION_STRING")

	var e error
	DB, e = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if e != nil {
		panic(e)
	}
	InitMigrate()
}

func InitMigrate() {
	DB.AutoMigrate(&models.Payment_Methods{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.Orders{})
	DB.AutoMigrate(&models.Products{})
	DB.AutoMigrate(&models.Shopping_Carts{})
	DB.AutoMigrate(&models.Order_Details{})
}

// Inisiasi koneksi ke database untuk melakukan unit testing
func InitDBTest() {
	connectionStringTest := os.Getenv("CONNECTION_STRING_TEST")

	// Koneksi ke DB
	var err error
	DB, err = gorm.Open(mysql.Open(connectionStringTest), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	InitiateMigrateTest()
}

// Migrasi tabel untuk unit testing
// Drop tabel dilakukan agar tabel selalu terinisiasi kembali setiap terbuhung ke database
func InitiateMigrateTest() {
	DB.Migrator().DropTable(&models.Order_Details{})
	DB.Migrator().DropTable(&models.Shopping_Carts{})
	DB.Migrator().DropTable(&models.Products{})
	DB.Migrator().DropTable(&models.Orders{})
	DB.Migrator().DropTable(&models.Users{})
	DB.Migrator().DropTable(&models.Address{})
	DB.Migrator().DropTable(&models.Payment_Methods{})
	DB.AutoMigrate(&models.Payment_Methods{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.Orders{})
	DB.AutoMigrate(&models.Products{})
	DB.AutoMigrate(&models.Shopping_Carts{})
	DB.AutoMigrate(&models.Order_Details{})
}
