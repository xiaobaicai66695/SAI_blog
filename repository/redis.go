package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb1 *redis.Client //操作blog相关数据
var rdb2 *redis.Client //操作like相关数据
var rdb3 *redis.Client //记录聊天信息

func InitRedis() error {
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "210618",
		DB:       1,
	})
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "210618",
		DB:       2,
	})
	rdb3 = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "210618",
		DB:       3,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := rdb1.Ping(ctx).Result()
	if err != nil {
		return errors.New(err.Error())
	}
	fmt.Println(pong)
	fmt.Println("---------------------------------------")
	return nil
}
