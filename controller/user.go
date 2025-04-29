package controller

import (
	"SAI_blog/middleware"
	"SAI_blog/repository"
	"SAI_blog/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	Response
	User  *service.UserVO `json:"user"`
	Token string          `json:"token"`
}

type UserInfoResponse struct {
	Response
	User *repository.UserInfo `json:"user"`
}

var user repository.User

func Register(c *gin.Context) {
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, Response{0, "系统异常"})
		return
	}
	name := user.Name
	account := user.Account
	password := user.Password
	id, err := service.UserRegister(name, account, password)
	if err != nil || id == 0 {
		fmt.Printf("%s", id)
		c.JSON(http.StatusOK, Response{0, "用户已存在"})
		return
	}
	c.JSON(200, Response{
		StatusCode: 1,
		StatusMsg:  "恭喜你，注册成功",
	})
}

func Login(c *gin.Context) {
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, Response{0, "登录失败，系统原因"})
		return
	}
	account := user.Account
	password := user.Password
	userVO, err := service.UserLogin(account, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{0, err.Error()})
		return
	}
	token, err := middleware.GetToken(userVO.ID)
	if err != nil {
		c.JSON(http.StatusOK, Response{0, err.Error()})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{1, "登录成功"},
		User:     userVO,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	userInfo := repository.QueryUserInfoById(id)
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询成功",
		},
		User: userInfo,
	})
}
