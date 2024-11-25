package database

import (
	"fmt"
	"log"
	"microservice_grpc_auth/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() (*gorm.DB,error) {
	dsn := "root:kl18jda183079@tcp(127.0.0.1:3306)/users_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil,err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Printf("failed to migrate database: %v", err)
		return nil,err
	}
	return db,nil
}
