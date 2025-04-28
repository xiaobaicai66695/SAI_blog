package controller

import (
	"SAI_blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	_ "strconv"
)

func Follow(c *gin.Context) {
	uidStr, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	lidStr := c.Param("lid")
	uid := uidStr.(int64)
	lid, _ := strconv.ParseInt(lidStr, 10, 64)
	if err := service.Follow(uid, lid); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "关注成功",
	})
}

func UnFollow(c *gin.Context) {
	uidStr, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	lidStr := c.Param("lid")
	uid := uidStr.(int64)
	lid, _ := strconv.ParseInt(lidStr, 10, 64)
	if err := service.UnFollow(uid, lid); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "已取关",
	})
}
