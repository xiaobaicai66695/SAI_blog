package controller

import (
	"SAI_blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Like(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	blogIdStr := c.Param("blogId")
	blogId, _ := strconv.ParseInt(blogIdStr, 10, 64)
	err := service.Like(blogId, uid.(int64))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "点赞成功",
	})
}

func UnLike(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	blogIdStr := c.Param("blogId")
	blogId, _ := strconv.ParseInt(blogIdStr, 10, 64)
	err := service.UnLike(blogId, uid.(int64))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "取消点赞成功",
	})
}
