package service

import (
	"SAI_blog/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type LoginFlow struct {
	account  string
	password string
}
type UserVO struct {
	ID      int64  `gorm:"primary_key" json:"ID"`
	Name    string `gorm:"column:name" json:"name"`
	Account string `gorm:"column:account" json:"account"`
}

var userVO *UserVO

func UserLogin(account string, password string) (*UserVO, error) {
	return NewLoginFlow(account, password).Do()
}

func NewLoginFlow(account string, password string) *LoginFlow {
	return &LoginFlow{account, password}
}

func (f *LoginFlow) Do() (*UserVO, error) {
	user, err := f.userLogin()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (f *LoginFlow) userLogin() (*UserVO, error) {
	user, _ := repository.NewUserDao().QueryByAccount(f.account)
	if user == nil {
		return nil, fmt.Errorf("用户未注册")
	}
	//数据库中的password
	userPassword := user.Password
	if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(f.password)) != nil {
		return nil, fmt.Errorf("密码错误")
	}
	userVO := &UserVO{
		ID:      user.Id,
		Name:    user.Name,
		Account: user.Account,
	}
	return userVO, nil
}
