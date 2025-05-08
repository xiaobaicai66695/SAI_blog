package router

import (
	"SAI_blog/controller"
	"SAI_blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	userRoutes := r.Group("/user")
	userRoutes.POST("/register", controller.Register)
	userRoutes.POST("/login", controller.Login)
	userRoutes.GET("/:id", controller.UserInfo)

	blogRoutes := r.Group("/blog")
	blogRoutes.POST("/upload", middleware.AutoMiddleware(), controller.UploadBlog)
	blogRoutes.GET("/:blogID", controller.BlogInfo)
	blogRoutes.GET("/list/follow", middleware.AutoMiddleware(), controller.FollowBlogList)
	blogRoutes.GET("/list", controller.BlogWithKey)

	r.POST("/follow/:lid", middleware.AutoMiddleware(), controller.Follow)
	r.DELETE("/follow/:lid", middleware.AutoMiddleware(), controller.UnFollow)
	r.POST("like/:lid", middleware.AutoMiddleware(), controller.Like)
	r.DELETE("like/:lid", middleware.AutoMiddleware(), controller.UnLike)
	r.GET("/ws", middleware.AutoMiddleware(), controller.WsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	return r
}
