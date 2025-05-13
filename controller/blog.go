package controller

import (
	"SAI_blog/repository"
	"SAI_blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UploadBlogResponse struct {
	Response
	*repository.Blog
}

type BlogInfoResponse struct {
	Response
	*service.BlogInfo
}

type BlogListResponse struct {
	Response
	Blogs []service.BlogInfo
}

type CommentData struct {
	Comment string `json:"comment"`
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
	c.JSON(http.StatusOK, UploadBlogResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "上传成功",
		},
		Blog: &blog,
	})
}

func BlogInfo(c *gin.Context) {
	blogIdStr := c.Param("blogID")
	blogId, err := strconv.ParseInt(blogIdStr, 10, 64)
	if err != nil {
		return
	}
	blogInfo := service.BlogInfoById(blogId)
	if blogInfo == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "没有这篇博客",
		})
		return
	}
	c.JSON(http.StatusOK, BlogInfoResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询成功",
		},
		BlogInfo: blogInfo,
	})
}

func FollowBlogList(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok || uid == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "获取参数失败",
		})
		return
	}
	blogInfos := service.QueryFollowBlogList(uid.(int64), page)
	c.JSON(http.StatusOK, BlogListResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询成功",
		},
		Blogs: blogInfos,
	})
}

func BlogWithKey(c *gin.Context) {
	key := c.Query("key")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "获取查询参数失败",
		})
		return
	}
	blogs := service.QueryBlogListWithKey(key, page)
	c.JSON(http.StatusOK, BlogListResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询成功",
		},
		Blogs: blogs,
	})
}

func BlogTotalCount(c *gin.Context) {
	count := repository.BlogTotalCount()
	c.JSON(http.StatusOK, gin.H{
		"status_code": 1,
		"status_msg":  count,
	})
}

func UploadComment(c *gin.Context) {
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
	var data CommentData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "评论失败",
		})
		return
	}

	service.UploadComment(blogId, uid.(int64), data.Comment)
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  "评论成功",
	})
}

func BlogOfUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	blogInfos := service.QueryBlogByUserId(id)
	c.JSON(http.StatusOK, BlogListResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "请求成功",
		},
		Blogs: *blogInfos,
	})
}

func BlogContentImagesUpload(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "请先登录",
		})
		return
	}
	file, err := c.FormFile("file")
	url, err := service.BlogContentImagesUpload(uid.(int64), file)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  url,
	})
}
