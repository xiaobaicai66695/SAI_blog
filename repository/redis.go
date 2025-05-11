package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var rdb1 *redis.Client //操作blog相关数据
var rdb2 *redis.Client //操作like相关数据
var rdb3 *redis.Client //记录聊天信息
var rdb4 *redis.Client //存储token
var config Config

func InitRedis() error {
	config, _ = LoadConfig()
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Password,
		DB:       1,
	})
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Password,
		DB:       2,
	})
	rdb3 = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Password,
		DB:       3,
	})
	rdb4 = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: config.Redis.Password,
		DB:       4,
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
