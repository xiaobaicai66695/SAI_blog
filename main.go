package main

import (
	"SAI_blog/repository"
	"SAI_blog/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := repository.InitDB(); err != nil {
		panic(err)
	}
	if err := repository.InitRedis(); err != nil {
		panic(err)
	}
	r := gin.Default()
	r = router.InitRouter(r)
	panic(r.Run(":8081"))
}
