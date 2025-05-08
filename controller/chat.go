package controller

import (
	"SAI_blog/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

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
