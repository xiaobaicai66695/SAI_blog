package service

import (
	"SAI_blog/common"
	"SAI_blog/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"os"
)

type BlogInfo struct {
	BlogId   int64  `gorm:"column:blog_id;primary_key" json:"blog_id"`
	UID      int64  `gorm:"column:uid" json:"uid"`
	Title    string `gorm:"column:title" json:"title"`
	Content  string `gorm:"column:content;type:longtext" json:"content"`
	Likes    int64  `gorm:"column:likes" json:"likes"`
	Comments int64  `gorm:"column:comments" json:"comments"`
	Lid      int64  `gorm:"column:lid" json:"lid"`
	Comment  string `gorm:"column:comment;type:text" json:"comment"`
	Cid      int64  `gorm:"column:cid" json:"cid"`
}

// 用消息队列要发送的消息
type BlogMsg struct {
	UUID string
	Blog *repository.Blog
}

// 定义常量
const blogUploadTopic = "blog"

// var blog *repository.Blog
var ctx = context.Background()

type ConsumerGroupHandler struct {
}

// 消费者
func (ConsumerGroupHandler) Setup() error {
	return nil
}

func (ConsumerGroupHandler) Cleanup() error {
	return nil
}

// 上传博客时消费者逻辑
func (ConsumerGroupHandler) UploadBlogConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}

func sendBlogUploadMessage(msg *repository.Blog) {
	//初始化生产者
	producer, err := common.InitProducer()
	if err != nil {
		fmt.Println("producer init fail")
		os.Exit(1)
	}
	defer producer.Close()
	msgJson, _ := json.Marshal(msg)
	msgBytes := []byte(msgJson)
	message := &sarama.ProducerMessage{
		Topic: blogUploadTopic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err1 := producer.SendMessage(message)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println("发送消息了")
	return
}

func UploadBlog(blog *repository.Blog) error {
	blogId, err := repository.BlogToRedis(blog)
	//errChan := make(chan error, 1)
	if blogId == 0 || err != nil {
		return fmt.Errorf("UploadBlog blog to redis fail")
	}
	blog.BlogId = blogId
	go sendBlogUploadMessage(blog)

	//go func() {
	//	err = repository.ConsumeBlogFromKafka("save-blog", blogUploadTopic, repository.BlogConsumerGroupHandler{})
	//}()
	//go func() {
	//	err = repository.ConsumeBlogFromKafka("push-follower-blog", blogUploadTopic, repository.PushToFollower{})
	//}()
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	return nil
}
