package models

import (
	"log"
	// the mysql import

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

// DB database global variable
var DB *gorm.DB

// ConnectDatabase connect to db
func ConnectDatabase() {
	database, err := gorm.Open("mysql", "root:12345678@(localhost)/authorize?charset=utf8&parseTime=True&loc=Local") // add pass when connect to db

	if err != nil {
		log.Println(err)
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{})

	DB = database
}

//User struct
type User struct {
	ID       int    ` json:"-" gorm:"primary_key;auto_increment:true"`
	UserName string `json:"username" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required"`
}
