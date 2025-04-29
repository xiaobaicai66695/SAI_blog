package repository

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type User struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"column:name;type:varchar(255)" json:"name"`
	IcoUrl   string `gorm:"column:ico_url;type:varchar(255)" json:"ico_url"`
	Account  string `gorm:"column:account;type:varchar(255)" json:"account"`
	Password string `gorm:"column:password;type:varchar(255)" json:"password"`
	Follower int64  `gorm:"column:follower" json:"follower"`
	Followed int64  `gorm:"column:followed" json:"followed"`
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

func UpdateUser(tx *gorm.DB, user *User) error {
	tx.Where("id = ?", user.Id).Updates(user)
	return nil
}

func QueryUserById(id int64) (*User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
