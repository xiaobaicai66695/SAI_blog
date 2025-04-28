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

	/*r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})*/
	return r
}
