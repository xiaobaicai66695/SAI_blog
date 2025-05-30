package service

import (
	"SAI_blog/common"
	"SAI_blog/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type CommentsInfo struct {
	Comment repository.BlogComment `json:"Comment"`
	UserVO  *repository.UserVO     `json:"user"`
}
type BlogInfo struct {
	BlogId      int64          `gorm:"column:blog_id;primary_key" json:"blog_id"`
	UID         int64          `gorm:"column:uid" json:"uid"`
	UserName    string         `gorm:"column:user_name" json:"user_name"`
	UserIcon    string         `gorm:"column:user_icon" json:"user_icon"`
	Title       string         `gorm:"column:title" json:"title"`
	Content     string         `gorm:"column:content;type:longtext" json:"content"`
	Likes       int64          `gorm:"column:likes" json:"likes"`
	Comments    int64          `gorm:"column:comments" json:"comments"`
	CommentInfo []CommentsInfo `gorm:"-" json:"comment_info"`
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
func BlogInfoById(blogId int64, uid int64) (*BlogInfo, bool) {
	if blogId == 0 {
		return nil, false
	}
	blog := repository.QueryBlogById(blogId)
	isLiked := repository.BlogIsLike(blogId, uid)
	if blog == nil {
		return nil, false
	}
	blogInfo := packingBlogToBlogInfo(blog)
	return blogInfo, isLiked
}

func QueryFollowBlogList(uid int64, page int) []BlogInfo {
	offset := (page - 1) * 10
	followBlogIds := repository.QueryBlogIdByUid(uid)
	var blogInfos []BlogInfo
	var ids []int64
	for _, blogId := range followBlogIds {
		ids = append(ids, blogId.BlogId)
	}
	blogs := repository.QueryBlogWithIds(ids, offset)
	for _, blog := range blogs {
		if blog.BlogId == 0 {
			continue
		}

		blogInfo := packingBlogToBlogInfo(blog)
		blogInfos = append(blogInfos, *blogInfo)
	}
	return blogInfos
}

func QueryBlogListWithKey(key string, page int) []BlogInfo {
	offset := (page - 1) * 10
	blogs := repository.QueryBlogByKey(key, offset)
	var blogInfos []BlogInfo
	for _, blog := range blogs {
		if blog.BlogId == 0 {
			continue
		}
		blogInfo := packingBlogToBlogInfo(blog)
		blogInfos = append(blogInfos, *blogInfo)
	}
	return blogInfos
}

func packingBlogToBlogInfo(blog *repository.Blog) *BlogInfo {
	comments := repository.QueryCommentsById(blog.BlogId)
	user, _ := repository.QueryUserById(blog.UID)
	var commentInfos []CommentsInfo
	for _, comment := range comments {
		commentUser, _ := repository.QueryUserById(comment.Cid)
		commentInfo := CommentsInfo{
			Comment: comment,
			UserVO: &repository.UserVO{
				ID:      commentUser.Id,
				Name:    commentUser.Name,
				IcoUrl:  commentUser.IcoUrl,
				Account: commentUser.Account,
			},
		}
		commentInfos = append(commentInfos, commentInfo)
	}
	userName := user.Name
	blogInfo := &BlogInfo{
		BlogId:   blog.BlogId,
		UID:      blog.UID,
		UserName: userName,
		UserIcon: user.IcoUrl,
		Title:    blog.Title,
		Content:  blog.Content,
		Comments: blog.Comments,
		//Comments:    int64(len(commentInfos)),
		Likes:       blog.Likes,
		CommentInfo: commentInfos,
	}
	return blogInfo
}

func UploadComment(blogId int64, uid int64, msg string) {
	msg = strings.Trim(msg, "")
	repository.UploadComment(blogId, uid, msg)
}

func QueryBlogByUserId(id int64) *[]BlogInfo {
	blogs := repository.QueryBlogByUserId(id)
	var blogInfos []BlogInfo
	for _, blog := range *blogs {
		blogInfo := packingBlogToBlogInfo(&blog)
		blogInfos = append(blogInfos, *blogInfo)
	}
	return &blogInfos
}

func BlogContentImagesUpload(uid int64, file *multipart.FileHeader) (string, error) {
	ext := path.Ext(file.Filename)
	if !(ext == "" || ext == ".png" || ext == ".jpg" || ext == ".gif" || ext == ".jpeg") {
		return "", fmt.Errorf("格式错误")
	}
	filePath := fmt.Sprintf("./static/blog_images/%s%d%s", time.Now().Format("20060102150405"), uid, ext)
	errChan := make(chan error)
	go func() {
		errChan <- SaveImages(file, filePath)
	}()
	err := <-errChan
	if err != nil {
		return "", fmt.Errorf("保存失败")
	}
	filePath = strings.TrimPrefix(filePath, ".")
	return filePath, nil
}
