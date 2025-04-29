package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const keyPrefix = "blog:"

var ctx = context.Background()
var blog = &Blog{}

const blogUploadTopic = "blog"

type Blog struct {
	BlogId   int64  `gorm:"column:blog_id;primary_key;autoIncrement" json:"blog_id"`
	UID      int64  `gorm:"column:uid" json:"uid"`
	Title    string `gorm:"column:title" json:"title"`
	Content  string `gorm:"column:content;type:longtext" json:"content"`
	Likes    int64  `gorm:"column:likes" json:"likes"`
	Comments int64  `gorm:"column:comments" json:"comments"`
}

type BlogMsg struct {
	UUID string
	Blog *Blog
}

// 保存博客的消费流程
type BlogConsumerGroupHandler struct {
}

// 消费者
func (BlogConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (BlogConsumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

// 上传博客时消费者逻辑
func (BlogConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Print("收到消息了\n")
		json.Unmarshal(msg.Value, &blog)
		db.Save(blog)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (Blog) TableName() string {
	return "blog"
}

func BlogToRedis(blog *Blog) (int64, error) {
	blogId := time.Now().Unix()
	blogId = blogId*10000 + blog.UID
	blog.BlogId = blogId
	blogIdStr := strconv.FormatInt(blogId, 10)
	blogJson, _ := json.Marshal(blog)
	key := fmt.Sprintf("%s%s", keyPrefix, blogIdStr)
	err := rdb1.Set(ctx, key, blogJson, time.Hour*10).Err()
	if err != nil {
		return 0, err
	}
	return blogId, err
}

// 给粉丝推送博客的流程
type PushToFollower struct{}

func (PushToFollower) Setup(sess sarama.ConsumerGroupSession) error   { return nil }
func (PushToFollower) Cleanup(sess sarama.ConsumerGroupSession) error { return nil }

// 具体推送逻辑
func (PushToFollower) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("准备推送消息了\n")
		json.Unmarshal(msg.Value, &blog)
		var uid = blog.UID
		var relationships []Relationship
		//根据作者的uid找出粉丝的id
		//fid是关注的人 uid是自己
		//要找自己是哪一些人所关注的人即条件是fid=作者的uid
		db.Find(&relationships, "fid = ?", uid)
		//为每个查出来的粉丝推送
		for _, relationship := range relationships {
			var followBlog = &FollowBlog{
				Uid:    relationship.Uid,
				BlogId: blog.BlogId,
			}
			db.Create(followBlog)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func ConsumeBlogFromKafka(groupID string, topic string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, groupID, config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT)
		<-signals
		cancel()
	}()
	for {
		err := consumerGroup.Consume(ctx, []string{topic}, handler)
		if err != nil {
			return err
		}
	}
}

//	func SaveBlogFromKafka() error {
//		config := sarama.NewConfig()
//		config.Consumer.Return.Errors = true
//		config.Consumer.Offsets.Initial = sarama.OffsetOldest
//		consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "blogGroup", config)
//		if err != nil {
//			return err
//		}
//		defer consumerGroup.Close()
//		ctx, cancel := context.WithCancel(ctx)
//		defer cancel()
//		go func() {
//			signals := make(chan os.Signal, 1)
//			signal.Notify(signals, syscall.SIGINT)
//			<-signals
//			cancel()
//		}()
//		handler := &BlogConsumerGroupHandler{}
//		for {
//			err := consumerGroup.Consume(ctx, []string{blogUploadTopic}, handler)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
