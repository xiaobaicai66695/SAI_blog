package repository

import (
	"fmt"
	"sync"
)

type User struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	IcoUrl   string `gorm:"column:ico_url" json:"ico_url"`
	Account  string `gorm:"column:account" json:"account"`
	Password string `gorm:"column:password" json:"password"`
	Follower string `gorm:"column:follower" json:"follower"`
	Followed string `gorm:"column:followed" json:"followed"`
}

type UserDao struct{}

var userDao *UserDao
var userOnce sync.Once

func NewUserDao() *UserDao {
	userOnce.Do(func() {
		userDao = &UserDao{}
	})
	return userDao
}

func (User) TableName() string {
	return "user"
}
func (*UserDao) QueryByAccount(account string) (*User, error) {
	var user User
	if err := db.Where("account = ?", account).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (*UserDao) Save(user *User) error {
	fmt.Print("要保存用户了")
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
