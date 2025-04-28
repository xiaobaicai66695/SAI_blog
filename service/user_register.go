package service

import (
	"SAI_blog/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterFlow struct {
	name     string
	account  string
	password string
}

func UserRegister(name string, account string, password string) (int64, error) {
	return NewRegisterFlow(name, account, password).Do()
}

func NewRegisterFlow(name string, account string, password string) *RegisterFlow {
	return &RegisterFlow{name, account, password}
}

func (f *RegisterFlow) Do() (int64, error) {
	uid, err := f.Register()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (f *RegisterFlow) Register() (int64, error) {
	//fmt.Printf("将要查询用户")
	user, _ := repository.NewUserDao().QueryByAccount(f.account)
	//fmt.Printf("已查询用户")
	//查询用户，若不存在才注册
	if user != nil {
		return 0, fmt.Errorf("用户已存在")
	}
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(f.password), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordBytes)
	if err != nil {
		return 0, fmt.Errorf("存储密码时出了问题")
	}
	user = &repository.User{
		Name:     f.name,
		Account:  f.account,
		Password: hashedPassword,
	}
	//fmt.Printf("用户已初始化")
	if err := repository.NewUserDao().Save(user); err != nil {
		return 0, fmt.Errorf("注册失败")
	}
	return user.Id, nil
}
