package controller

import (
	"SAI_blog/repository"
	"SAI_blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BlogResponse struct {
	Response
	*repository.Blog
}

func UploadBlog(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok || uid == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	var blog repository.Blog
	err := c.ShouldBind(&blog)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	blog.UID = uid.(int64)
	err = service.UploadBlog(&blog)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	blog.Likes = 0
	blog.Comments = 0
	c.JSON(http.StatusOK, BlogResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "上传成功",
		},
		Blog: &blog,
	})
}
