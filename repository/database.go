package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	config, err = LoadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Charset,
		config.Database.ParseTime,
		config.Database.Loc)
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}
	err = db.AutoMigrate(&User{}, &Blog{}, &BlogComment{}, &BlogLikes{}, &Relationship{}, &FollowBlog{}, &UserGroup{})
	if err != nil {
		panic(err)
	}
	return err
}

func GetDB() *gorm.DB {
	return db
}
