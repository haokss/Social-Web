package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *gorm.DB

func DataBase(conn_string string) {
	db, err := gorm.Open("mysql", conn_string)
	if err != nil {
		panic("Mysql连接错误")
	}

	log.Info("Mysql DataBase Connect Success!")
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	db.SingularTable(true)       // 表明不加s
	db.DB().SetMaxIdleConns(20)  // 设置连接池大小
	db.DB().SetMaxOpenConns(100) // 设置最大连接数
	DB = db
	migration()
}
