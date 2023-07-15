package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// USER := "root"
	// PASS := "Natoato19"
	// HOST := "localhost"
	// DBNAME := "article"
	// URL := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST)
	// root:sample-password@tcp(db:3306)/jwt_auth
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/?parseTime=true"))
	// db, err := gorm.Open(mysql.Open(URL))

	if err != nil {
		panic(err)
	}

	_ = db.Exec("CREATE DATABASE IF NOT EXISTS article;")

	dbc, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/article?parseTime=true"))

	if err != nil {
		panic(err)
	}

	dbc.AutoMigrate(&Posts{})

	DB = dbc
}
