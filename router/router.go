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
	userRoutes.POST("/bindmail", middleware.AutoMiddleware(), controller.BindEmail)
	userRoutes.POST("/reset/sendmail", controller.SendMail)
	userRoutes.POST("/reset/:account", controller.ResetPassword)
	userRoutes.GET("/:id", middleware.AutoMiddleware(), controller.UserInfo)
	userRoutes.POST("ico-update", middleware.AutoMiddleware(), controller.UpdateIco)

	blogRoutes := r.Group("/blog")
	blogRoutes.POST("/content-images-load", middleware.AutoMiddleware(), controller.BlogContentImagesUpload)
	blogRoutes.POST("/upload", middleware.AutoMiddleware(), controller.UploadBlog)
	blogRoutes.GET("/:blogID", controller.BlogInfo)
	blogRoutes.GET("/list/follow", middleware.AutoMiddleware(), controller.FollowBlogList)
	blogRoutes.GET("/list", controller.BlogWithKey)
	blogRoutes.GET("/count", controller.BlogTotalCount)
	blogRoutes.GET("userblog/:id", controller.BlogOfUser)
	blogRoutes.POST("/comment/:blogId", middleware.AutoMiddleware(), controller.UploadComment)

	r.POST("/follow/:lid", middleware.AutoMiddleware(), controller.Follow)
	r.DELETE("/follow/:lid", middleware.AutoMiddleware(), controller.UnFollow)
	r.POST("like/:blogId", middleware.AutoMiddleware(), controller.Like)
	r.DELETE("like/:blogId", middleware.AutoMiddleware(), controller.UnLike)

	r.GET("/ws", middleware.AutoMiddleware(), controller.WsHandler)
	r.GET("/history/:groupID", controller.PullHistoryMsg)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	return r
}
