package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open("mysql", "root:@/contacts?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB = database

}
func GetDB() *gorm.DB {
	return DB
}
