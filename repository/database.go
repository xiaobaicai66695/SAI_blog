package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	dsn := "root:210618@tcp(127.0.0.1:3306)/sai_blog?charset=utf8&parseTime=True&loc=Local"
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}
	err = db.AutoMigrate(&User{}, &Blog{}, &BlogComment{}, &BlogLikes{}, &Follow{})
	if err != nil {
		panic(err)
	}
	return err
}
