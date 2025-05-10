package controller

import (
	"SAI_blog/middleware"
	"SAI_blog/repository"
	"SAI_blog/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
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
	uid, ok := c.Get("uid")
	if !ok {
		uid = int64(0)
	}

	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	userInfo := repository.QueryUserInfoById(id, uid.(int64))
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询成功",
		},
		User: userInfo,
	})
}

func ResetPassword(c *gin.Context) {
	account := c.Param("account")
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "修改失败",
		})
	}
	repository.ResetPassword(account, user.Password)
	c.JSON(200, Response{
		StatusCode: 1,
		StatusMsg:  "修改密码成功",
	})
}

func SendMail(c *gin.Context) {
	_ = c.ShouldBindJSON(&user)
	emailAddress := repository.QueryAddressByAccount(user.Account)
	if emailAddress != user.Email {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "邮箱或账号错误",
		})
		return
	}
	from := "575121043@qq.com"
	password := "wbohjifkququbebd"
	to := []string{emailAddress}

	stmpHost := "smtp.qq.com"
	stmpPort := "587"

	message := []byte("From: 575121043@qq.com\r\n" +
		"To: " + emailAddress + "\r\n" +
		"Subject: 重置密码链接\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		//后续需要改的地方
		"点击这里重置密码: https://localhost:8081/reset/" + user.Account + "\r\n")

	auth := smtp.PlainAuth("", from, password, stmpHost)
	err := smtp.SendMail(stmpHost+":"+stmpPort, auth, from, to, message)
	if err != nil && !strings.Contains(err.Error(), "short response") {
		fmt.Println(err)
		c.JSON(http.StatusOK, Response{0, "发送失败"})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "发送成功",
	})
}

func BindEmail(c *gin.Context) {
	uid, ok := c.Get("uid")
	_ = c.ShouldBindJSON(&user)
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	repository.BindEmail(uid.(int64), user.Email)
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "绑定成功",
	})
}
