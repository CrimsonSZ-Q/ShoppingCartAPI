package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db

}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=pohon4785 dbname=restapi port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error...")
		return nil
	}
	return Db
}
