package controller

import (
	"SAI_blog/repository"
	"SAI_blog/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type historyMsg struct {
	Response
	Msgs []*repository.Group
}

func WsHandler(c *gin.Context) {
	groupIdStr := c.Query("groupId")
	uid, ok := c.Get("uid")
	if !ok {
		fmt.Printf("用户未登录")
		return
	}
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	err = service.HandlerChatConnection(c.Writer, c.Request, groupId, uid.(int64))
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func PullHistoryMsg(c *gin.Context) {
	groupIdStr := c.Param("groupID")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	cursorStr := c.GetHeader("Cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	msgs, err := service.PullHistoryMsg(groupId, cursor)
	if err != nil {
		fmt.Printf(err.Error())
		c.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "请求失败",
		})
		return
	}
	c.JSON(http.StatusOK, historyMsg{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "请求成功",
		},
		Msgs: msgs,
	})
}
