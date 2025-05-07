package main

import (
	"SAI_blog/repository"
	"SAI_blog/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	if err := repository.InitDB(); err != nil {
		panic(err)
	}
	if err := repository.InitRedis(); err != nil {
		panic(err)
	}
	go func() {
		err := repository.ConsumeBlogFromKafka("save-blog", "blog", repository.BlogConsumerGroupHandler{})
		if err != nil {
			panic("未能成功保存博客:" + err.Error())
		}
	}()
	go func() {
		err := repository.ConsumeBlogFromKafka("push-follower-blog", "blog", repository.PushToFollower{})
		if err != nil {
			panic("粉丝推送博客失败:" + err.Error())
		}
	}()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			repository.UpdateLikeFromRedis()
		}
	}()
	r := gin.Default()
	r = router.InitRouter(r)
	r.Use(cors.Default())
	panic(r.Run("0.0.0.0:8081"))
}
