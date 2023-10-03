package connection

import (
	"auth/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

func InitMySQL() {

	dsn := "root:root@tcp(:3306)/user_docker?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to Connect to database")
		os.Exit(2)
	}

	log.Println("Database has Connected")
	db.AutoMigrate(&models.User{})
	DBConn = db
}
