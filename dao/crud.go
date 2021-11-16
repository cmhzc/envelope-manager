package dao

import (
	"envelope_manager/entity"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initDB() *gorm.DB {
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DBNAME")

	// MYSQL dns格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	// 类似{username}使用花括号包着的名字都是需要替换的参数
	dns := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	db, err := gorm.Open("mysql", dns)
	if err != nil {
		panic("failed to connect mysql, error: " + err.Error())
	}

	return db
}

func FlushDB() {
	db := initDB()
	// delete all records
	db.Where("1 = 1").Delete(&entity.Envelope{})
	db.Close()
}
